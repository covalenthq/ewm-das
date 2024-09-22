#!/bin/bash

# Define the directory again in the wrapper script, it will use the value assigned during installation
COVALENT_DIR="${COVALENT_DIR:-$HOME/.covalent}"  # Default to ~/.covalent if not set
SERVICE_NAME="${EXECUTABLE:-light-client}"        # Default to light-client if not set
GCP_CREDENTIALS="${GCP_CREDENTIALS:-gcp-credentials.json}"  # Default to gcp-credentials.json if not set

# Check if PRIVATE_KEY is set; if not, exit with an error
if [ -z "$PRIVATE_KEY" ]; then
  echo "Error: PRIVATE_KEY environment variable is not set."
  exit 1
fi

# Check if PRIVATE_KEY is a valid 64-character hexadecimal number
if ! [[ "$PRIVATE_KEY" =~ ^[0-9a-fA-F]{64}$ ]]; then
  echo "Error: PRIVATE_KEY is not a valid 64-character hexadecimal number."
  exit 1
fi

# Ensure that only one instance of the service is running
if pgrep -f "$SERVICE_NAME" > /dev/null 2>&1; then
    echo "$SERVICE_NAME is already running."
    exit 1
fi

# Wait for IPFS daemon to start by checking if it is listening on port 5001
echo "Waiting for IPFS daemon to start on port 5001..."
until lsof -i :5001 | grep LISTEN > /dev/null; do
  printf '.'
  sleep 1
done
echo "IPFS daemon has started."

# Run your service binary with all the arguments
"$COVALENT_DIR/$SERVICE_NAME" \
    --loglevel debug \
    --rpc-url ws://34.42.69.93:8080/rpc \
    --collect-url https://ewm-light-clients-v2-838505730421.us-central1.run.app
    --private-key "$PRIVATE_KEY"