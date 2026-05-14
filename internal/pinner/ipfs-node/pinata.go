package ipfsnode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ipfs/go-cid"
)

const (
	defaultPinataUploadsHost = "uploads.pinata.cloud"
	defaultPinataAPIHost     = "api.pinata.cloud"
	defaultPinataNetwork     = "public"
	pinataUploadTimeout      = 5 * time.Minute
)

// PinataStorage uploads CAR files to Pinata via the v3 Files API.
type PinataStorage struct {
	jwt           string
	groupID       string
	network       string
	uploadsHost   string
	uploadsScheme string // "https" in production; tests override to "http"
	apiScheme     string
	httpClient    *http.Client
}

// PinataConfig is the value-object passed in from main / api.ServerConfig.
type PinataConfig struct {
	JWT     string // required
	GroupID string // optional
	Network string // optional; defaults to "public"
}

// NewPinataStorage validates the JWT eagerly by calling /data/testAuthentication.
func NewPinataStorage(cfg PinataConfig) (*PinataStorage, error) {
	if cfg.JWT == "" {
		return nil, fmt.Errorf("pinata: JWT is required (set PINATA_JWT)")
	}
	network := cfg.Network
	if network == "" {
		network = defaultPinataNetwork
	}
	if network != "public" && network != "private" {
		return nil, fmt.Errorf("pinata: network must be \"public\" or \"private\", got %q", network)
	}

	ps := &PinataStorage{
		jwt:           cfg.JWT,
		groupID:       cfg.GroupID,
		network:       network,
		uploadsHost:   defaultPinataUploadsHost,
		uploadsScheme: "https",
		apiScheme:     "https",
		httpClient:    &http.Client{Timeout: pinataUploadTimeout},
	}
	return ps, nil
}

// Initialize verifies the JWT against Pinata's testAuthentication endpoint.
// Called once at startup; non-fatal warm-up failures are returned so caller can decide.
func (p *PinataStorage) Initialize() error {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s://%s/data/testAuthentication", p.apiScheme, defaultPinataAPIHost), nil)
	if err != nil {
		return fmt.Errorf("pinata: build auth-test request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.jwt)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("pinata: auth-test request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pinata: auth-test returned status %d: %s",
			resp.StatusCode, truncate(string(body), 256))
	}
	log.Info("Pinata authentication verified")
	return nil
}

// pinataUploadResponse mirrors the JSON shape from /v3/files.
type pinataUploadResponse struct {
	Data struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		CID         string `json:"cid"`
		Size        int64  `json:"size"`
		CreatedAt   string `json:"created_at"`
		MimeType    string `json:"mime_type"`
		GroupID     string `json:"group_id"`
		Network     string `json:"network"`
		IsDuplicate bool   `json:"is_duplicate"`
	} `json:"data"`
}

// Pin uploads the given CAR file to Pinata and returns the CID Pinata reports.
//
// The caller is responsible for verifying that the returned CID matches the
// locally-computed root CID. The check lives at publish.go:66.
func (p *PinataStorage) Pin(carFile *os.File) (cid.Cid, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pinataUploadTimeout)
	defer cancel()

	if _, err := carFile.Seek(0, io.SeekStart); err != nil {
		return cid.Undef, fmt.Errorf("pinata: rewind CAR file: %w", err)
	}

	body, contentType, err := p.buildMultipart(carFile)
	if err != nil {
		return cid.Undef, err
	}

	url := fmt.Sprintf("%s://%s/v3/files", p.uploadsScheme, p.uploadsHost)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return cid.Undef, fmt.Errorf("pinata: build upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.jwt)
	req.Header.Set("Content-Type", contentType)

	resp, err := p.doWithRetry(req)
	if err != nil {
		return cid.Undef, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return cid.Undef, fmt.Errorf("pinata: upload returned status %d: %s",
			resp.StatusCode, truncate(string(respBody), 256))
	}

	var parsed pinataUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return cid.Undef, fmt.Errorf("pinata: decode upload response: %w", err)
	}
	if parsed.Data.CID == "" {
		return cid.Undef, fmt.Errorf("pinata: upload response had empty data.cid")
	}

	parsedCID, err := cid.Parse(parsed.Data.CID)
	if err != nil {
		return cid.Undef, fmt.Errorf("pinata: parse returned CID %q: %w", parsed.Data.CID, err)
	}

	log.Infof("Pinata accepted upload: cid=%s id=%s size=%d duplicate=%t",
		parsed.Data.CID, parsed.Data.ID, parsed.Data.Size, parsed.Data.IsDuplicate)
	return parsedCID, nil
}

// buildMultipart constructs the multipart form payload for /v3/files.
//
// Fields: file, network, name, car=true, plus optional group_id.
func (p *PinataStorage) buildMultipart(carFile *os.File) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("file", filepath.Base(carFile.Name()))
	if err != nil {
		return nil, "", fmt.Errorf("pinata: multipart create file: %w", err)
	}
	if _, err := io.Copy(part, carFile); err != nil {
		return nil, "", fmt.Errorf("pinata: multipart copy file: %w", err)
	}

	if err := w.WriteField("network", p.network); err != nil {
		return nil, "", err
	}
	if err := w.WriteField("name", filepath.Base(carFile.Name())); err != nil {
		return nil, "", err
	}
	// car=true tells Pinata to process this upload as a CAR file and index its
	// inner DAG, returning the CAR's root CID rather than a wrapper.
	if err := w.WriteField("car", "true"); err != nil {
		return nil, "", err
	}
	if p.groupID != "" {
		if err := w.WriteField("group_id", p.groupID); err != nil {
			return nil, "", err
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}

// doWithRetry retries on 5xx and transient network errors. Never on 4xx.
func (p *PinataStorage) doWithRetry(req *http.Request) (*http.Response, error) {
	// Cache body so we can re-send on retry.
	var bodyBytes []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("pinata: buffer request body: %w", err)
		}
		bodyBytes = b
	}

	backoff := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for attempt := 0; attempt <= len(backoff); attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		resp, err := p.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < len(backoff) {
				time.Sleep(backoff[attempt])
				continue
			}
			return nil, fmt.Errorf("pinata: HTTP request failed after retries: %w", err)
		}
		if resp.StatusCode < 500 {
			return resp, nil
		}
		// 5xx — close and retry.
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		lastErr = fmt.Errorf("pinata: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 256))
		if attempt < len(backoff) {
			time.Sleep(backoff[attempt])
			continue
		}
		return nil, lastErr
	}
	return nil, lastErr
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
