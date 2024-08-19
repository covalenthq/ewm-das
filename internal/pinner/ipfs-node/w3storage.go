package ipfsnode

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-ucanto/did"
	"github.com/web3-storage/go-ucanto/principal/ed25519/signer"
)

// W3Storage struct encapsulates the web3.storage client.
type W3Storage struct {
	agentKey            string
	agentDID            did.DID
	delegationProofPath string
	w3Temp              string
}

// Binary name for the `w3` command-line tool.
var binName = "w3"

// NewW3Storage initializes a new W3Storage instance.
func NewW3Storage(agentKey, delegationProofPath string) (*W3Storage, error) {
	// Ensure the `w3` binary is available in the PATH.
	if _, err := exec.LookPath(binName); err != nil {
		return nil, fmt.Errorf("%s binary not found in PATH", binName)
	}

	// Create a temporary directory for `w3` usage.
	tmpDir, err := os.MkdirTemp("", "w3up_config")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Parse the agent key using the ed25519 signer.
	agentSigner, err := signer.Parse(agentKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent key: %w", err)
	}

	// Verify that the delegation proof path exists.
	if _, err := os.Stat(delegationProofPath); err != nil {
		return nil, fmt.Errorf("delegation proof path does not exist: %w", err)
	}

	return &W3Storage{
		agentKey:            agentKey,
		agentDID:            agentSigner.DID().DID(),
		delegationProofPath: delegationProofPath,
		w3Temp:              tmpDir,
	}, nil
}

// Close cleans up resources used by W3Storage.
func (w3 *W3Storage) Close() {
	_ = os.RemoveAll(w3.w3Temp) // Ignoring error as it's not critical
}

// Initialize initializes the W3Storage instance by adding the agent and space.
func (w3 *W3Storage) Initialize() error {
	i, err := w3.whoAmI()
	if err != nil {
		return err
	}
	log.Infof("Initialized W3Storage with agent DID: %s", i)

	spaceDID, err := w3.addSpace()
	if err != nil {
		return err
	}

	log.Infof("Added space with DID: %s", spaceDID)

	return nil
}

// whoAmI returns the DID of the current W3Storage instance.
func (w3 *W3Storage) whoAmI() (did.DID, error) {
	command := fmt.Sprintf("W3_STORE_NAME=%s W3_PRINCIPAL=\"%s\" %s whoami", w3.w3Temp, w3.agentKey, binName)
	output, err := runCommand(command)
	if err != nil {
		return did.Undef, err
	}

	pdid, err := did.Parse(strings.TrimSpace(output))
	if err != nil {
		return did.Undef, fmt.Errorf("failed to parse DID from whoami output: %w", err)
	}

	return pdid, nil
}

// addSpace adds a space to the W3Storage instance using the delegation proof.
func (w3 *W3Storage) addSpace() (did.DID, error) {
	command := fmt.Sprintf("W3_STORE_NAME=%s %s space add %s", w3.w3Temp, binName, w3.delegationProofPath)
	output, err := runCommand(command)
	if err != nil {
		return did.Undef, err
	}

	pdid, err := did.Parse(strings.TrimSpace(output))
	if err != nil {
		return did.Undef, fmt.Errorf("failed to parse DID from add space output: %w", err)
	}

	return pdid, nil
}

// Pin uploads a CAR file and pins it, returning the resulting CID.
func (w3 *W3Storage) Pin(carFile *os.File) (cid.Cid, error) {
	command := fmt.Sprintf("W3_STORE_NAME=%s %s up %s --json --no-wrap --car", w3.w3Temp, binName, carFile.Name())
	output, err := runCommand(command)
	if err != nil {
		return cid.Undef, err
	}

	log.Infof("w3 up output: %s", output)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return cid.Undef, fmt.Errorf("failed to unmarshal w3 up output: %w", err)
	}

	root, ok := result["root"].(map[string]interface{})
	if !ok {
		return cid.Undef, fmt.Errorf("unexpected format: root key missing in w3 up output")
	}

	rcid, ok := root["/"].(string)
	if !ok {
		return cid.Undef, fmt.Errorf("unexpected format: CID key missing in root object")
	}

	return cid.Parse(rcid)
}

// runCommand is a helper function to execute shell commands and return the output.
func runCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, error: %w", string(output), err)
	}
	return strings.TrimSpace(string(output)), nil
}
