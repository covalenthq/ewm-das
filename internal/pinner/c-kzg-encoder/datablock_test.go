package ckzgencoder

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	ckzg4844 "github.com/ethereum/c-kzg-4844/v2/bindings/go"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// readAndDecodeNode reads a CBOR encoded file, decodes it using the DAG-CBOR decoder, and returns the decoded node.
func readAndDecodeNode(cid string) (datamodel.Node, error) {
	// Read the file
	data, err := os.ReadFile(fmt.Sprintf("../../../test/data/%s.bin", cid))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create a NodeAssembler to build the decoded node
	nb := basicnode.Prototype.Any.NewBuilder()

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	if err := dagcbor.Decode(nb, reader); err != nil {
		return nil, fmt.Errorf("failed to decode CBOR data: %w", err)
	}

	return nb.Build(), nil
}

// getCommitment extracts a specific commitment from the commitments node by index and returns it as a Bytes48.
func getCommitment(commitments datamodel.Node, index int) (*ckzg4844.Bytes48, error) {
	// Get the commitment node
	commitmentNode, err := commitments.LookupByIndex(int64(index))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup commitment: %w", err)
	}

	// Get the commitment bytes
	commitmentBytes, err := commitmentNode.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get commitment bytes: %w", err)
	}

	var commitment48 ckzg4844.Bytes48
	copy(commitment48[:], commitmentBytes)

	return &commitment48, nil
}

// getCommitments extracts the commitments node from the given node.
func getCommitments(node datamodel.Node) (datamodel.Node, error) {
	return node.LookupByString("commitments")
}

// getCommitmentAsBytes48 reads a CBOR encoded file, decodes it, and extracts a specific commitment by index as a Bytes48.
func getCommitmentAsBytes48(cid string, index int) (*ckzg4844.Bytes48, error) {
	node, err := readAndDecodeNode(cid)
	if err != nil {
		return nil, err
	}

	commitments, err := getCommitments(node)
	if err != nil {
		return nil, fmt.Errorf("getCommitments returned an error: %w", err)
	}

	return getCommitment(commitments, index)
}

// extractProofAndCell extracts the proof and cell values from a decoded node and returns them as Bytes48 and Cell respectively.
func extractProofAndCell(node datamodel.Node) (*ckzg4844.Bytes48, *ckzg4844.Cell, error) {
	// Extract proof
	proofNode, err := node.LookupByString("proof")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup proof: %w", err)
	}

	proofBytes, err := proofNode.AsBytes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get proof bytes: %w", err)
	}

	var proof48 ckzg4844.Bytes48
	copy(proof48[:], proofBytes)

	// Extract cell
	cellNode, err := node.LookupByString("cell")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup cell: %w", err)
	}

	cellBytes, err := cellNode.AsBytes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cell bytes: %w", err)
	}

	var cell2048 ckzg4844.Cell
	copy(cell2048[:], cellBytes)

	return &proof48, &cell2048, nil
}

// getProofAndCell reads a CBOR encoded file, decodes it, and extracts the proof and cell values.
func getProofAndCell(cid string) (*ckzg4844.Bytes48, *ckzg4844.Cell, error) {
	node, err := readAndDecodeNode(cid)
	if err != nil {
		return nil, nil, err
	}

	return extractProofAndCell(node)
}

// TestVerifyRootNode tests the extraction of commitments, proofs, and cells, and verifies them using VerifyCellKZGProofBatch.
func TestVerifyRootNode(t *testing.T) {
	rootCid := "bafyreiahay5quioczvzk5tdr7muuiyozmtsq6yizncwi6r6bst42v5jnqi"
	commitment, err := getCommitmentAsBytes48(rootCid, 2)
	if err != nil {
		t.Fatalf("getCommitmentAsBytes48 returned an error: %v", err)
	}

	cellCid := "bafyreiconqo5jeezyvxbu5xr5mcrk2xmdj76gtdwswuhvism25e5penkse"
	proof, cell, err := getProofAndCell(cellCid)
	if err != nil {
		t.Fatalf("getProofAndCell returned an error: %v", err)
	}

	commitments := [1]ckzg4844.Bytes48{*commitment}
	proofs := [1]ckzg4844.Bytes48{*proof}
	cells := [1]ckzg4844.Cell{*cell}
	indexes := [1]uint64{6}

	fmt.Printf("commitment: %v\n", commitment)

	ok, err := ckzg4844.VerifyCellKZGProofBatch(commitments[:], indexes[:], cells[:], proofs[:])
	if err != nil {
		t.Fatalf("VerifyCellKZGProofBatch returned an error: %v", err)
	}

	if !ok {
		t.Fatalf("VerifyCellKZGProofBatch failed")
	}
}
