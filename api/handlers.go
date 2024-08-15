package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/covalenthq/das-ipfs-pinner/internal/das"
	ipfsnode "github.com/covalenthq/das-ipfs-pinner/internal/ipfs-node"
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

func createStoreHandler(ipfsNode *ipfsnode.IPFSNode) http.HandlerFunc {
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
			handleError(w, "Failed to parse multipart form", http.StatusBadRequest)
			return
		}

		for filename, data := range files {
			log.Printf("Received %d bytes of data from file: %s\n", len(data), filename)

			block, err := das.Encode(data)
			if err != nil {
				handleError(w, "Failed to encode data", http.StatusInternalServerError)
				return
			}

			cid, err := ipfsNode.PublishBlock(block, true)
			if err != nil {
				handleError(w, "Failed to store data to IPFS", http.StatusInternalServerError)
				return
			}

			log.Printf("Data stored successfully with CID: %s\n", cid)
			fmt.Fprintf(w, "File %s stored successfully with CID: %s\n", filename, cid)
		}
	}
}

func extractHandler(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		handleError(w, "CID is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Extracting data for CID: %s\n", cid)
}
