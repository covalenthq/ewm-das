package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/covalenthq/das-ipfs-pinner/internal/pinner/das"
	ipfsnode "github.com/covalenthq/das-ipfs-pinner/internal/pinner/ipfs-node"
)

func parseMultipartFormData(r *http.Request, maxMemory int64) (map[string][]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	files := make(map[string][]byte)
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			files[fileHeader.Filename] = data
		}
	}
	return files, nil
}

func createUploadHandler(ipfsNode *ipfsnode.IPFSNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			handleError(w, "Content-Type must be multipart/form-data", http.StatusUnsupportedMediaType)
			return
		}

		files, err := parseMultipartFormData(r, MaxMultipartMemory)
		if err != nil {
			log.Errorf("Failed to parse multipart form: %w", err)
			handleError(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		for filename, data := range files {
			log.Debugf("Received %d bytes of data from file: %s", len(data), filename)

			block, err := das.Encode(data)
			if err != nil {
				log.Errorf("Failed to encode data: %w", err)
				handleError(w, "Failed to encode data", http.StatusInternalServerError)
				return
			}

			cid, err := ipfsNode.PublishBlock(block, true)
			if err != nil {
				log.Errorf("Failed to upload data to IPFS: %w", err)
				handleError(w, "Failed to upload data to IPFS", http.StatusInternalServerError)
				return
			}

			log.Infof("Data upload successfully with CID: %s", cid)
			succStr := fmt.Sprintf("{\"cid\": \"%s\"}", cid.String())
			if _, err := w.Write([]byte(succStr)); err != nil {
				log.Errorf("error writing data to connection: %w", err)
			}
		}
	}
}

func createDownloadHandler(ipfsNode *ipfsnode.IPFSNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET method
		if r.Method != http.MethodGet {
			handleError(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the CID from the query parameters
		cid := r.URL.Query().Get("cid")
		if cid == "" {
			handleError(w, "CID query parameter is required", http.StatusBadRequest)
			return
		}

		// Extract the block from IPFS
		_, err := ipfsNode.ExtractBlock(r.Context(), cid)
		if err != nil {
			handleError(w, "Failed to extract data from IPFS", http.StatusInternalServerError)
			return
		}
	}
}

func createCalculateCIDHandler(ipfsNode *ipfsnode.IPFSNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			handleError(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			handleError(w, "Content-Type must be multipart/form-data", http.StatusUnsupportedMediaType)
			return
		}

		files, err := parseMultipartFormData(r, MaxMultipartMemory)
		if err != nil {
			log.Errorf("Failed to parse multipart form: %w", err)
			handleError(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		for filename, data := range files {
			log.Debugf("Received %d bytes of data from file: %s", len(data), filename)

			block, err := das.Encode(data)
			if err != nil {
				log.Errorf("Failed to encode data: %w", err)
				handleError(w, "Failed to encode data", http.StatusInternalServerError)
				return
			}

			cid, err := ipfsNode.PublishBlock(block, false)
			if err != nil {
				log.Errorf("Failed to upload data to IPFS: %w", err)
				handleError(w, "Failed to upload data to IPFS", http.StatusInternalServerError)
				return
			}

			log.Infof("Data locally cached with CID: %s", cid)
			succStr := fmt.Sprintf("{\"cid\": \"%s\"}", cid.String())
			if _, err := w.Write([]byte(succStr)); err != nil {
				log.Errorf("error writing data to connection: %w", err)
			}
		}
	}
}

func deprecatedHandler(originalHandler http.HandlerFunc, newEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the use of a deprecated endpoint
		log.Warnf("Deprecated endpoint accessed: %s", r.URL.Path)

		// Add a deprecation notice to the headers
		w.Header().Set("Warning", `199 - "Deprecated API: Please use `+newEndpoint+` instead"`)

		// Call the original handler to preserve the original response body
		originalHandler(w, r)
	}
}
