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

// TestReadAndDecodeCbor tests the readAndDecodeCbor function.
func TestReadAndDecodeCbor(t *testing.T) {
	rootCid := "bafyreiahay5quioczvzk5tdr7muuiyozmtsq6yizncwi6r6bst42v5jnqi"
	root, err := readAndDecodeCbor(fmt.Sprintf("../../test/data/%s.bin", rootCid))
	if err != nil {
		t.Fatalf("readAndDecodeCbor returned an error: %v", err)
	}

	verifyRootNode(t, root)
}

// verifyRootNode verifies the keys in the root node and processes the "links" key.
func verifyRootNode(t *testing.T, root datamodel.Node) {
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

		if !contains(expectedKeys, key) {
			t.Fatalf("Unexpected key found: %v", key)
		}

		if key == "links" {
			processLinks(t, v)
		}
	}
}

// processLinks processes the "links" key in the root node.
func processLinks(t *testing.T, linksNode datamodel.Node) {
	rowLink, err := linksNode.LookupByIndex(2)
	if err != nil {
		t.Fatalf("Failed to extract link: %v", err)
	}

	cid, err := rowLink.AsLink()
	if err != nil {
		t.Fatalf("Failed to extract CID: %v", err)
	}

	cols, err := readAndDecodeCbor(fmt.Sprintf("../../test/data/%s.bin", cid.String()))
	if err != nil {
		t.Fatalf("readAndDecodeCbor returned an error: %v", err)
	}

	colLink, err := cols.LookupByIndex(6)
	if err != nil {
		t.Fatalf("Failed to extract key and value: %v", err)
	}

	cid, err = colLink.AsLink()
	if err != nil {
		t.Fatalf("Failed to extract CID: %v", err)
	}

	cell, err := readAndDecodeCbor(fmt.Sprintf("../../test/data/%s.bin", cid.String()))
	if err != nil {
		t.Fatalf("readAndDecodeCbor returned an error: %v", err)
	}

	verifyCellNode(t, cell)
}

// verifyCellNode verifies the contents of the "cell" node.
func verifyCellNode(t *testing.T, cell datamodel.Node) {
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

		switch key {
		case "proof":
			verifyProof(t, v)
		case "cell":
			verifyCell(t, v)
		}
	}
}

// verifyProof checks that the "proof" value is a byte slice of length 48.
func verifyProof(t *testing.T, proofNode datamodel.Node) {
	proof, err := proofNode.AsBytes()
	if err != nil {
		t.Fatalf("Failed to extract proof: %v", err)
	}
	if len(proof) != 48 {
		t.Fatalf("Proof is not of length 48")
	}
}

// verifyCell checks that the "cell" value is a byte slice of length 2048.
func verifyCell(t *testing.T, cellNode datamodel.Node) {
	cell, err := cellNode.AsBytes()
	if err != nil {
		t.Fatalf("Failed to extract cell: %v", err)
	}
	if len(cell) != 2048 {
		t.Fatalf("Cell is not of length 2048")
	}
}
