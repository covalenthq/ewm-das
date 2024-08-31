#!/bin/bash

# Define the directory again in the wrapper script, it will use the value assigned during installation
COVALENT_DIR="${COVALENT_DIR:-$HOME/.covalent}"  # Default to ~/.covalent if not set
SERVICE_NAME="${EXECUTABLE:-light-client}"        # Default to light-client if not set
GCP_CREDENTIALS="${GCP_CREDENTIALS:-gcp-credentials.json}"  # Default to gcp-credentials.json if not set

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
    --rpc-url wss://moonbase-alpha.blastapi.io/618fd77b-a090-457b-b08a-373398006a5e \
    --contract 0x916B54696A70588a716F899bE1e8f2A5fFd5f135 \
    --topic-id DAS-TO-BQ \
    --gcp-creds-file "$COVALENT_DIR/$GCP_CREDENTIALS" \
    --client-id "$CLIENT_ID"