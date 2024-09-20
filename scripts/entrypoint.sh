#!/bin/sh

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

# Start IPFS daemon in the background
ipfs daemon --enable-gc=true &

# Wait for IPFS daemon to start
until netstat -an | grep 'LISTEN' | grep ':5001'; do
  printf '.'
  sleep 1
done

# Start light-client with the provided and hardcoded arguments
light-client \
  --loglevel debug \
  --rpc-url ws://34.42.69.93:8080/rpc \
  --topic-id DAS-TO-BQ \
  --gcp-creds-file /gcp-credentials.json \
  --private-key "$PRIVATE_KEY"