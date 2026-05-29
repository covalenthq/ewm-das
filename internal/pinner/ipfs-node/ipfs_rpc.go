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
	"strings"
	"time"

	"github.com/ipfs/go-cid"
)

const ipfsRPCUploadTimeout = 5 * time.Minute

type IPFSRPCStorage struct {
	rpcURL     string
	rpcToken   string
	httpClient *http.Client
}

type IPFSRPCConfig struct {
	RPCURL   string // IPFS_RPC_URL
	RPCToken string // IPFS_RPC_TOKEN (optional)
}

func NewIPFSRPCStorage(cfg IPFSRPCConfig) (*IPFSRPCStorage, error) {
	if cfg.RPCURL == "" {
		return nil, fmt.Errorf("ipfs-rpc: RPC URL is required (set IPFS_RPC_URL, e.g. http://127.0.0.1:5001/api/v0)")
	}
	return &IPFSRPCStorage{
		rpcURL:     strings.TrimRight(cfg.RPCURL, "/"),
		rpcToken:   cfg.RPCToken,
		httpClient: &http.Client{Timeout: ipfsRPCUploadTimeout},
	}, nil
}

func (s *IPFSRPCStorage) Initialize() error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, s.rpcURL+"/version", nil)
	if err != nil {
		return fmt.Errorf("ipfs-rpc: build version request: %w", err)
	}
	s.setAuth(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ipfs-rpc: version request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("ipfs-rpc: version returned status %d: %s",
			resp.StatusCode, truncate(string(body), 256))
	}
	log.Infof("IPFS RPC reachable at %s", s.rpcURL)
	return nil
}

func (s *IPFSRPCStorage) Pin(carFile *os.File, expectedRoots []cid.Cid) error {
	if len(expectedRoots) == 0 {
		return fmt.Errorf("ipfs-rpc: expectedRoots is empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), ipfsRPCUploadTimeout)
	defer cancel()

	if _, err := carFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("ipfs-rpc: rewind CAR file: %w", err)
	}

	body, contentType, err := s.buildMultipart(carFile)
	if err != nil {
		return err
	}

	url := s.rpcURL + "/dag/import?pin-roots=true"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("ipfs-rpc: build upload request: %w", err)
	}
	s.setAuth(req)
	req.Header.Set("Content-Type", contentType)

	resp, err := s.doWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("ipfs-rpc: dag/import returned status %d: %s",
			resp.StatusCode, truncate(string(respBody), 256))
	}

	pinned := make(map[string]struct{}, len(expectedRoots))
	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		var entry dagImportRoot
		if err := dec.Decode(&entry); err != nil {
			return fmt.Errorf("ipfs-rpc: decode dag/import response: %w", err)
		}
		c := entry.Root.Cid.Slash
		if c == "" {
			return fmt.Errorf("ipfs-rpc: dag/import response entry has empty Root.Cid")
		}
		if entry.Root.PinErrorMsg != "" {
			return fmt.Errorf("ipfs-rpc: pin failed for %s: %s", c, entry.Root.PinErrorMsg)
		}
		pinned[c] = struct{}{}
	}

	for _, c := range expectedRoots {
		if _, ok := pinned[c.String()]; !ok {
			return fmt.Errorf("ipfs-rpc: expected root %s missing from dag/import response (%d roots acknowledged)",
				c, len(pinned))
		}
	}

	log.Infof("IPFS RPC pinned %d blocks (DAS root=%s)", len(pinned), expectedRoots[0])
	return nil
}

func (s *IPFSRPCStorage) setAuth(req *http.Request) {
	if s.rpcToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.rpcToken)
	}
}

func (s *IPFSRPCStorage) buildMultipart(carFile *os.File) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("file", filepath.Base(carFile.Name()))
	if err != nil {
		return nil, "", fmt.Errorf("ipfs-rpc: multipart create file: %w", err)
	}
	if _, err := io.Copy(part, carFile); err != nil {
		return nil, "", fmt.Errorf("ipfs-rpc: multipart copy file: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}

// doWithRetry backs off on 5xx; 4xx is returned immediately so bad creds/URLs surface fast.
func (s *IPFSRPCStorage) doWithRetry(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("ipfs-rpc: buffer request body: %w", err)
		}
		bodyBytes = b
	}

	backoff := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for attempt := 0; attempt <= len(backoff); attempt++ {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		resp, err := s.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < len(backoff) {
				time.Sleep(backoff[attempt])
				continue
			}
			return nil, fmt.Errorf("ipfs-rpc: HTTP request failed after retries: %w", err)
		}
		if resp.StatusCode < 500 {
			return resp, nil
		}
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()
		lastErr = fmt.Errorf("ipfs-rpc: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 256))
		if attempt < len(backoff) {
			time.Sleep(backoff[attempt])
			continue
		}
		return nil, lastErr
	}
	return nil, lastErr
}
