package ipldencoder

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// expectedKeys is a list of expected keys in the CBOR encoded data.
var expectedKeys = []string{"version", "links", "size", "commitments", "length"}

// contains checks if a string is in a slice of strings.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// readAndDecodeCbor reads a CBOR encoded file and decodes it using the DAG-CBOR decoder.
func readAndDecodeCbor(filePath string) (datamodel.Node, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create a NodeAssembler to build the decoded node
	nb := basicnode.Prototype.Any.NewBuilder()

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	err = dagcbor.Decode(nb, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CBOR data: %w", err)
	}

	// Return the decoded node
	return nb.Build(), nil
}

// TestReadAndDecodeCbor tests the ReadAndDecodeCbor function.
func TestReadAndDecodeCbor(t *testing.T) {
	// Call the function with the temporary file path
	rootCid := "bafyreiahay5quioczvzk5tdr7muuiyozmtsq6yizncwi6r6bst42v5jnqi"
	root, err := readAndDecodeCbor(fmt.Sprintf("../../test/data/%s.bin", rootCid))
	if err != nil {
		t.Fatalf("ReadAndDecodeCbor returned an error: %v", err)
	}

	// extract key and value from the root node
	itr := root.MapIterator()
	for !itr.Done() {
		k, v, err := itr.Next()
		if err != nil {
			t.Fatalf("Failed to extract key and value: %v", err)
		}

		key, err := k.AsString()
		if err != nil {
			t.Fatalf("Failed to extract key: %v", err)
		}

		// raise an error if the key is not in the expected keys
		if !contains(expectedKeys, key) {
			t.Fatalf("Key not found: %v", key)
		}

		// extract links from the root node
		if key == "links" {
			rowLink, err := v.LookupByIndex(2)
			if err != nil {
				t.Fatalf("Failed to extract link: %v", err)
			}

			// link to cid
			cid, err := rowLink.AsLink()
			if err != nil {
				t.Fatalf("Failed to extract link: %v", err)
			}

			cols, err := readAndDecodeCbor(fmt.Sprintf("../../test/data/%s.bin", cid.String()))
			if err != nil {
				t.Fatalf("ReadAndDecodeCbor returned an error: %v", err)
			}

			// extract key and value from the cols node
			colLink, err := cols.LookupByIndex(6)
			if err != nil {
				t.Fatalf("Failed to extract key and value: %v", err)
			}

			// link to cid
			cid, err = colLink.AsLink()
			if err != nil {
				t.Fatalf("Failed to extract link: %v", err)
			}

			// read the data from the cid
			cell, err := readAndDecodeCbor(fmt.Sprintf("../../test/data/%s.bin", cid.String()))
			if err != nil {
				t.Fatalf("ReadAndDecodeCbor returned an error: %v", err)
			}

			// extract key and value from the cell node iteartor
			itr := cell.MapIterator()
			for !itr.Done() {
				k, v, err := itr.Next()
				if err != nil {
					t.Fatalf("Failed to extract key and value: %v", err)
				}
				key, err := k.AsString()
				if err != nil {
					t.Fatalf("Failed to extract key: %v", err)
				}
				// verify that "proof" value is a byte slice of length 48
				if key == "proof" {
					proof, err := v.AsBytes()
					if err != nil {
						t.Fatalf("Failed to extract proof: %v", err)
					}
					if len(proof) != 48 {
						t.Fatalf("Proof is not of length 48: %v", err)
					}
				} else if key == "cell" {
					// verify that "cell" value is a byte slice of length 48
					cell, err := v.AsBytes()
					if err != nil {
						t.Fatalf("Failed to extract cell: %v", err)
					}
					if len(cell) != 2048 {
						t.Fatalf("Cell is not of length 48: %v", err)
					}
				}
			}
		}
	}
}
