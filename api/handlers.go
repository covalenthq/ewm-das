package api

import (
	"fmt"
	"io"
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
			succ_str := fmt.Sprintf("{\"cid\": \"%s\"}", cid.String())
			if _, err := w.Write([]byte(succ_str)); err != nil {
				log.Errorf("error writing data to connection: %w", err)
			}
		}
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		handleError(w, "CID is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Extracting data for CID: %s", cid)
}
