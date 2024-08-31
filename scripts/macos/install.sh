#!/bin/bash

# Paths
COVALENT_DIR="$HOME/.covalent"
IPFS_PATH=$(which ipfs)  # Get the actual path of the IPFS binary
EXECUTABLE="light-client"
TRUSTED_SETUP="trusted_setup.txt"
GCP_CREDENTIALS="gcp-credentials.json"
WRAPPER_SCRIPT="$COVALENT_DIR/run_light_client.sh"
PLIST_FILE="$HOME/Library/LaunchAgents/com.covalent.light-client.plist"
IPFS_PLIST_FILE="$HOME/Library/LaunchAgents/com.covalent.ipfs.plist"
IPFS_REPO_DIR="$HOME/.ipfs"

# Uninstallation step (run the uninstallation script)
bash uninstall.sh

# Check if the destination directory exists
mkdir -p "$COVALENT_DIR"

# Copy the executable, trusted setup, and wrapper script to the destination directory
cp "$EXECUTABLE" "$COVALENT_DIR/"
cp "$TRUSTED_SETUP" "$COVALENT_DIR/"
cp "$GCP_CREDENTIALS" "$COVALENT_DIR/"
cp "uninstall.sh" "$COVALENT_DIR/"
cp "run_light_client.sh" "$WRAPPER_SCRIPT"

# Make the executable and wrapper script runnable
chmod +x "$COVALENT_DIR/$EXECUTABLE"
chmod +x "$WRAPPER_SCRIPT"

# Bypass Gatekeeper for the executable
spctl --add --label "Trusted" "$COVALENT_DIR/$EXECUTABLE"
spctl --enable --label "Trusted"

# Create the IPFS launchd plist file with the correct IPFS path
cat <<EOF > "$IPFS_PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.covalent.ipfs</string>
    <key>ProgramArguments</key>
    <array>
        <string>$IPFS_PATH</string>
        <string>daemon</string>
        <string>--repo-dir</string>
        <string>$IPFS_REPO_DIR</string>
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
launchctl load "$IPFS_PLIST_FILE"

# Set CLIENT_ID environment variable for the light client
export CLIENT_ID="$1"

# Copy the light client launchd plist file
cp "light_client.plist" "$PLIST_FILE"

# Load the light client daemon
launchctl load "$PLIST_FILE"

echo "Installation completed. The IPFS daemon and the light client daemon are now running."