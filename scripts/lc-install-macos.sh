#!/bin/bash

# Paths
COVALENT_DIR="$HOME/.covalenthq"
EXECUTABLE="light-client"
TRUSTED_SETUP="trusted_setup.txt"
PLIST_FILE="com.covalenthq.light-client.plist"
GCP_CREDENTIALS="lc-gcp-credentials.json"
WRAPPER_SCRIPT="run_light_client.sh"

# Check if the destination directory exists
mkdir -p "$COVALENT_DIR"

# Copy the executable and trusted setup to the destination directory
cp "$EXECUTABLE" "$COVALENT_DIR/"
cp "$TRUSTED_SETUP" "$COVALENT_DIR/"

# Create the wrapper script
cat <<EOF > "$COVALENT_DIR/$WRAPPER_SCRIPT"
#!/bin/bash

# Ensure that only one instance of the service is running
SERVICE_NAME="$EXECUTABLE"
if pgrep -f "\$SERVICE_NAME" > /dev/null 2>&1; then
    echo "\$SERVICE_NAME is already running."
    exit 1
fi

# Start the IPFS daemon with garbage collection if it's not already running
if ! pgrep -f "ipfs daemon" > /dev/null 2>&1; then
    echo "Starting IPFS daemon with garbage collection..."
    ipfs daemon --enable-gc &
    sleep 10 # Give IPFS some time to start
else
    echo "IPFS daemon is already running."
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
    --client-id "{YOUR_UNIQUE_ID}"

EOF

# Make the wrapper script executable
chmod +x "$COVALENT_DIR/$WRAPPER_SCRIPT"

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
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>ThrottleInterval</key>
    <integer>30</integer> <!-- Prevents rapid restarts -->
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