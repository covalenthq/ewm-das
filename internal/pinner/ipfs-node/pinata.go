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

// pinataPersistencePollSchedule controls how long Pin waits for Pinata's async
// CAR validator to confirm the file before failing. The schedule is generous
// because Pinata's `is_duplicate=false` upload response only confirms "bytes
// received", not "pinned" — a CAR that fails async validation is dropped
// silently. Empirically a valid wrapped CAR appears in the account index in
// well under 10s; we wait up to ~30s before declaring rejection.
var pinataPersistencePollSchedule = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	4 * time.Second,
	8 * time.Second,
	15 * time.Second,
}

// PinataStorage uploads CAR files to Pinata via the v3 Files API.
type PinataStorage struct {
	jwt           string
	groupID       string
	network       string
	uploadsHost   string
	uploadsScheme string // "https" in production; tests override to "http"
	apiHost       string
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
		apiHost:       defaultPinataAPIHost,
		apiScheme:     "https",
		httpClient:    &http.Client{Timeout: pinataUploadTimeout},
	}
	return ps, nil
}

// Initialize verifies the JWT against Pinata's testAuthentication endpoint.
func (p *PinataStorage) Initialize() error {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s://%s/data/testAuthentication", p.apiScheme, p.apiHost), nil)
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

// Pin uploads carFile to Pinata and blocks until Pinata's async CAR validator
// confirms the file is in the account index. expectedRoot is the CAR's root
// CID (the dag-pb wrapper CID built by pin.go::writeWrappedCAR), used to
// verify Pinata's upload response and to spot drift.
//
// Returning nil means Pinata has actually persisted the upload, not merely
// that the HTTP POST returned 200. Returning an error means the file was
// either not accepted or silently dropped by the async validator — surface
// this to the caller; never log success on uncertainty.
func (p *PinataStorage) Pin(carFile *os.File, expectedRoot cid.Cid) error {
	ctx, cancel := context.WithTimeout(context.Background(), pinataUploadTimeout)
	defer cancel()

	if _, err := carFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("pinata: rewind CAR file: %w", err)
	}

	body, contentType, err := p.buildMultipart(carFile)
	if err != nil {
		return err
	}

	uploadURL := fmt.Sprintf("%s://%s/v3/files", p.uploadsScheme, p.uploadsHost)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, body)
	if err != nil {
		return fmt.Errorf("pinata: build upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.jwt)
	req.Header.Set("Content-Type", contentType)

	resp, err := p.doWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("pinata: upload returned status %d: %s",
			resp.StatusCode, truncate(string(respBody), 256))
	}

	var parsed pinataUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return fmt.Errorf("pinata: decode upload response: %w", err)
	}
	if parsed.Data.CID == "" {
		return fmt.Errorf("pinata: upload response had empty data.cid")
	}
	if parsed.Data.ID == "" {
		return fmt.Errorf("pinata: upload response had empty data.id")
	}
	if parsed.Data.CID != expectedRoot.String() {
		return fmt.Errorf("pinata: upload returned cid %q, expected wrapper %q",
			parsed.Data.CID, expectedRoot.String())
	}

	if err := p.waitForPersistence(ctx, parsed.Data.ID); err != nil {
		return err
	}

	log.Infof("Pinata pinned upload: cid=%s id=%s size=%d",
		parsed.Data.CID, parsed.Data.ID, parsed.Data.Size)
	return nil
}

// waitForPersistence polls GET /v3/files/{network}/{id} until the file appears
// in Pinata's account index or the schedule is exhausted. Pinata's async CAR
// validator can drop uploads silently (most notably any CAR whose root block
// has a dag-cbor codec); this poll is the only way to distinguish a real pin
// from "bytes received, then discarded".
func (p *PinataStorage) waitForPersistence(ctx context.Context, fileID string) error {
	url := fmt.Sprintf("%s://%s/v3/files/%s/%s", p.apiScheme, p.apiHost, p.network, fileID)
	for attempt, wait := range pinataPersistencePollSchedule {
		select {
		case <-ctx.Done():
			return fmt.Errorf("pinata: persistence check canceled after %d attempts: %w", attempt, ctx.Err())
		case <-time.After(wait):
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("pinata: build persistence request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+p.jwt)
		resp, err := p.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("pinata: persistence request failed: %w", err)
		}
		sc := resp.StatusCode
		resp.Body.Close()
		if sc == http.StatusOK {
			log.Debugf("Pinata persistence confirmed for %s after %d attempts", fileID, attempt+1)
			return nil
		}
	}
	return fmt.Errorf("pinata: file %s did not appear in account index — "+
		"Pinata likely rejected the upload silently (CAR root codec must be dag-pb)", fileID)
}

// buildMultipart constructs the multipart form payload for /v3/files.
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
