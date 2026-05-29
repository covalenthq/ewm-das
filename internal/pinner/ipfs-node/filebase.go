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
	defaultFilebaseRPCURL = "https://rpc.filebase.io/api/v0"
	filebaseUploadTimeout = 5 * time.Minute
)

// FilebaseStorage uploads CAR files to Filebase via the Kubo-shaped IPFS RPC
// API and verifies every block listed as a CAR root is pinned.
//
// The migration relies on a quirk of /dag/import with pin-roots=true: it pins
// every CID declared as a root in the CAR header recursively. Our CAR build in
// pin.go therefore declares every block of the DAG as a root, which causes
// Filebase to pin and DHT-advertise each dag-cbor block individually. This is
// load-bearing — without it Filebase pins only the CAR's nominal root, and
// inner block CIDs return CONTENT_NOT_HOSTED from the dedicated gateway.
type FilebaseStorage struct {
	rpcToken   string
	baseURL    string
	httpClient *http.Client
}

// FilebaseConfig is the value-object passed in from main / api.ServerConfig.
type FilebaseConfig struct {
	RPCToken string // required; an IPFS RPC API token from the Filebase console
}

// NewFilebaseStorage constructs a client. Auth is verified separately by
// Initialize so callers can decide how to react to startup failures.
func NewFilebaseStorage(cfg FilebaseConfig) (*FilebaseStorage, error) {
	if cfg.RPCToken == "" {
		return nil, fmt.Errorf("filebase: RPC token is required (set FILEBASE_RPC_TOKEN)")
	}
	return &FilebaseStorage{
		rpcToken:   cfg.RPCToken,
		baseURL:    defaultFilebaseRPCURL,
		httpClient: &http.Client{Timeout: filebaseUploadTimeout},
	}, nil
}

// Initialize verifies the RPC token by hitting the Kubo-shaped /version
// endpoint. Returns an error if the token is missing or invalid so the pinner
// fails to start instead of surfacing auth errors at first upload.
func (f *FilebaseStorage) Initialize() error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, f.baseURL+"/version", nil)
	if err != nil {
		return fmt.Errorf("filebase: build version request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+f.rpcToken)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("filebase: version request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("filebase: version returned status %d: %s",
			resp.StatusCode, truncate(string(body), 256))
	}
	log.Info("Filebase RPC authentication verified")
	return nil
}

// Pin uploads carFile to Filebase via POST /dag/import?pin-roots=true and
// verifies every CID in expectedRoots appears in the response with an empty
// PinErrorMsg. Returning nil means every expected block was pinned.
//
// The caller is responsible for ensuring carFile's header lists every CID it
// wants pinned as a root (pin.go::writeMultiRootCAR handles this for the DAS
// pipeline). Filebase only pins blocks named in the CAR header.
func (f *FilebaseStorage) Pin(carFile *os.File, expectedRoots []cid.Cid) error {
	if len(expectedRoots) == 0 {
		return fmt.Errorf("filebase: expectedRoots is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), filebaseUploadTimeout)
	defer cancel()

	if _, err := carFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("filebase: rewind CAR file: %w", err)
	}

	body, contentType, err := f.buildMultipart(carFile)
	if err != nil {
		return err
	}

	url := f.baseURL + "/dag/import?pin-roots=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("filebase: build upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+f.rpcToken)
	req.Header.Set("Content-Type", contentType)

	resp, err := f.doWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("filebase: dag/import returned status %d: %s",
			resp.StatusCode, truncate(string(respBody), 256))
	}

	pinned := make(map[string]struct{}, len(expectedRoots))
	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		var entry dagImportRoot
		if err := dec.Decode(&entry); err != nil {
			return fmt.Errorf("filebase: decode dag/import response: %w", err)
		}
		c := entry.Root.Cid.Slash
		if c == "" {
			return fmt.Errorf("filebase: dag/import response entry has empty Root.Cid")
		}
		if entry.Root.PinErrorMsg != "" {
			return fmt.Errorf("filebase: pin failed for %s: %s", c, entry.Root.PinErrorMsg)
		}
		pinned[c] = struct{}{}
	}

	for _, c := range expectedRoots {
		if _, ok := pinned[c.String()]; !ok {
			return fmt.Errorf("filebase: expected root %s missing from dag/import response (%d roots acknowledged)",
				c, len(pinned))
		}
	}

	log.Infof("Filebase pinned %d blocks (DAS root=%s)", len(pinned), expectedRoots[0])
	return nil
}

// buildMultipart wraps the CAR file in the multipart/form-data shape /dag/import
// expects (single "file" field, body is the raw CAR bytes).
func (f *FilebaseStorage) buildMultipart(carFile *os.File) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("file", filepath.Base(carFile.Name()))
	if err != nil {
		return nil, "", fmt.Errorf("filebase: multipart create file: %w", err)
	}
	if _, err := io.Copy(part, carFile); err != nil {
		return nil, "", fmt.Errorf("filebase: multipart copy file: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}

// doWithRetry retries on 5xx and transient network errors. Never on 4xx so
// misconfigured tokens surface immediately rather than burning the backoff
// budget on something a retry cannot fix.
func (f *FilebaseStorage) doWithRetry(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("filebase: buffer request body: %w", err)
		}
		bodyBytes = b
	}

	backoff := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for attempt := 0; attempt <= len(backoff); attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		resp, err := f.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < len(backoff) {
				time.Sleep(backoff[attempt])
				continue
			}
			return nil, fmt.Errorf("filebase: HTTP request failed after retries: %w", err)
		}
		if resp.StatusCode < 500 {
			return resp, nil
		}
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		lastErr = fmt.Errorf("filebase: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 256))
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
