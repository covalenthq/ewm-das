package ckzgencoder

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	ckzg4844 "github.com/ethereum/c-kzg-4844/bindings/go"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func getCommitments(node datamodel.Node) (datamodel.Node, error) {
	// Get the commitments node
	commitments, err := node.LookupByString("commitments")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup commitments: %w", err)
	}

	return commitments, nil
}

func getCommitment(commitments datamodel.Node, index int) (*ckzg4844.Bytes48, error) {
	// Get the commitment node
	commitment, err := commitments.LookupByIndex(int64(index))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup commitment: %w", err)
	}

	// Get the commitment bytes
	commitmentBytes, err := commitment.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get commitment bytes: %w", err)
	}

	var commitment48 ckzg4844.Bytes48
	copy(commitment48[:], commitmentBytes)

	return &commitment48, nil
}

func getCommitmentAsBytes48(cid string, index int) (*ckzg4844.Bytes48, error) {
	// Read the file
	data, err := os.ReadFile(fmt.Sprintf("../../test/data/%s.bin", cid))
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

	commitments, err := getCommitments(nb.Build())
	if err != nil {
		return nil, fmt.Errorf("getCommitments returned an error: %w", err)
	}

	return getCommitment(commitments, index)
}

func getProofAndCell(cid string) (*ckzg4844.Bytes48, *ckzg4844.Cell, error) {
	// Read the file
	data, err := os.ReadFile(fmt.Sprintf("../../test/data/%s.bin", cid))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Create a NodeAssembler to build the decoded node
	nb := basicnode.Prototype.Any.NewBuilder()

	// Wrap the data in a bytes.Reader to satisfy io.Reader interface
	reader := bytes.NewReader(data)

	// Decode the data using the DAG-CBOR decoder
	err = dagcbor.Decode(nb, reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode CBOR data: %w", err)
	}

	// node is a map, we need to extract the proof and cell
	proof, err := nb.Build().LookupByString("proof")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup proof: %w", err)
	}

	proofBytes, err := proof.AsBytes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get proof bytes: %w", err)
	}

	cell, err := nb.Build().LookupByString("cell")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup cell: %w", err)
	}

	cellBytes, err := cell.AsBytes()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cell bytes: %w", err)
	}

	var proof48 ckzg4844.Bytes48
	copy(proof48[:], proofBytes)

	var cell2048 ckzg4844.Cell
	copy(cell2048[:], cellBytes)

	return &proof48, &cell2048, nil
}

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
