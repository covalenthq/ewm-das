#!/bin/sh

# Check if CLIENT_ID is set; if not, exit with an error
if [ -z "$CLIENT_ID" ]; then
  echo "Error: CLIENT_ID environment variable is not set."
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
  --rpc-url wss://moonbase-alpha.blastapi.io/618fd77b-a090-457b-b08a-373398006a5e \
  --contract 0x916B54696A70588a716F899bE1e8f2A5fFd5f135 \
  --topic-id DAS-TO-BQ \
  --gcp-creds-file /gcp-credentials.json \
  --client-id "$CLIENT_ID"