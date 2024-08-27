#!/bin/bash

# Paths
COVALENT_DIR="$HOME/.covalenthq"
IPFS_PATH=$(which ipfs)
EXECUTABLE="light-client"
TRUSTED_SETUP="trusted_setup.txt"
GCP_CREDENTIALS="gcp-credentials.json"
WRAPPER_SCRIPT="run_light_client.sh"
PLIST_FILE="com.covalenthq.light-client.plist"
IPFS_PLIST_FILE="com.covalenthq.ipfs.plist"

# Check if the destination directory exists
mkdir -p "$COVALENT_DIR"

# Copy the executable and trusted setup to the destination directory
cp "$EXECUTABLE" "$COVALENT_DIR/"
cp "$TRUSTED_SETUP" "$COVALENT_DIR/"
cp "$GCP_CREDENTIALS" "$COVALENT_DIR/"

# Make the executable runnable
chmod +x "$COVALENT_DIR/$EXECUTABLE"

# Bypass Gatekeeper for the executable
spctl --add --label "Trusted" "$COVALENT_DIR/$EXECUTABLE"
spctl --enable --label "Trusted"

# Create the IPFS launchd plist file without a wrapper script
cat <<EOF > "$HOME/Library/LaunchAgents/$IPFS_PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.covalenthq.ipfs</string>
    <key>ProgramArguments</key>
    <array>
        <string>$IPFS_PATH</string>
        <string>daemon</string>
        <string>--enable-gc</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>ThrottleInterval</key>
    <integer>10</integer> <!-- Prevent rapid restarts -->
    <key>StandardOutPath</key>
    <string>$COVALENT_DIR/ipfs.log</string>
    <key>StandardErrorPath</key>
    <string>$COVALENT_DIR/ipfs.error.log</string>
</dict>
</plist>
EOF

# Load the IPFS daemon
launchctl load "$HOME/Library/LaunchAgents/$IPFS_PLIST_FILE"

# Create the light client wrapper script
cat <<EOF > "$COVALENT_DIR/$WRAPPER_SCRIPT"
#!/bin/bash

# Define the directory again in the wrapper script, it will use the value assigned during installation
COVALENT_DIR="$COVALENT_DIR"

# Ensure that only one instance of the service is running
SERVICE_NAME="$EXECUTABLE"
if pgrep -f "\$SERVICE_NAME" > /dev/null 2>&1; then
    echo "\$SERVICE_NAME is already running."
    exit 1
fi

# Wait for IPFS daemon to be fully available
until pgrep -f "ipfs daemon" > /dev/null 2>&1; do
    echo "Waiting for IPFS daemon to start..."
    sleep 5
done

# Run your service binary with all the arguments
"\$COVALENT_DIR/\$SERVICE_NAME" \\
    --loglevel debug \\
    --rpc-url wss://moonbase-alpha.blastapi.io/618fd77b-a090-457b-b08a-373398006a5e \\
    --contract 0x916B54696A70588a716F899bE1e8f2A5fFd5f135 \\
    --topic-id DAS-TO-BQ \\
    --gcp-creds-file "\$COVALENT_DIR/$GCP_CREDENTIALS" \\
    --client-id {YOUR_UNIQUE_ID}

EOF

# Make the wrapper script executable
chmod +x "$COVALENT_DIR/$WRAPPER_SCRIPT"

# Create the launchd plist file for the light client
cat <<EOF > "$HOME/Library/LaunchAgents/$PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.covalenthq.light-client</string>
    <key>ProgramArguments</key>
    <array>
        <string>$COVALENT_DIR/$WRAPPER_SCRIPT</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PINNER_DIR</key>
        <string>$COVALENT_DIR</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>ThrottleInterval</key>
    <integer>30</integer> <!-- Prevents rapid restarts -->
    <key>StandardOutPath</key>
    <string>$COVALENT_DIR/light-client.log</string>
    <key>StandardErrorPath</key>
    <string>$COVALENT_DIR/light-client.log</string>
</dict>
</plist>
EOF

# Load the light client daemon
launchctl load "$HOME/Library/LaunchAgents/$PLIST_FILE"

echo "Installation completed. The IPFS daemon and the light client daemon are now running."