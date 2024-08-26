#!/bin/bash

# Paths
COVALENT_DIR="$HOME/.covalenthq"
EXECUTABLE="light-client"
TRUSTED_SETUP="trusted_setup.txt"
PLIST_FILE="com.covalenthq.light-client.plist"
GCP_CREDENTIALS="gcp-credentials.json"

# Check if the destination directory exists
mkdir -p "$COVALENT_DIR"

# Copy the executable to the destination directory
cp "$EXECUTABLE" "$COVALENT_DIR/"
cp "$TRUSTED_SETUP" "$COVALENT_DIR/"
cp "$GCP_CREDENTIALS" "$COVALENT_DIR/"

# Create the launchd plist file
cat <<EOF > "$HOME/Library/LaunchAgents/$PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.covalenthq.light-client</string>
    <key>ProgramArguments</key>
    <array>
        <string>$COVALENT_DIR/$EXECUTABLE</string>
        <string>--loglevel</string>
        <string>debug</string>
        <string>--rpc-url</string>
        <string>wss://moonbase-alpha.blastapi.io/618fd77b-a090-457b-b08a-373398006a5e</string>
        <string>--contract</string>
        <string>0x916B54696A70588a716F899bE1e8f2A5fFd5f135</string>
        <string>--topic-id</string>
        <string>DAS-TO-BQ</string>
        <string>--gcp-creds-file</string>
        <string>$COVALENT_DIR/$GCP_CREDENTIALS</string>
        <string>--client-id</string>
        <string>4e440d76-4a26-4117-b317-b5b407b0cd54</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$COVALENT_DIR/light-client.log</string>
    <key>StandardErrorPath</key>
    <string>$COVALENT_DIR/light-client.error.log</string>
</dict>
</plist>
EOF

# Load the daemon
launchctl load "$HOME/Library/LaunchAgents/$PLIST_FILE"
echo "Installation completed. The light client daemon is now running."