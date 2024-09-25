#!/bin/bash

if [ -z "$1" ]; then
    echo "Private Key is required."
    echo "Usage: $0 <private-key>"
    exit 1
fi

# Check if private key is a valid 64-character hexadecimal number
if ! [[ "$1" =~ ^[0-9a-fA-F]{64}$ ]]; then
  echo "Error: private key is not a valid 64-character hexadecimal number."
  exit 1
fi

# Paths
COVALENT_DIR="$HOME/.covalent"
IPFS_PATH=$(which ipfs)  # Get the actual path of the IPFS binary
EXECUTABLE="light-client"
TRUSTED_SETUP="trusted_setup.txt"
WRAPPER_SCRIPT="$COVALENT_DIR/run.sh"
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
cp "uninstall.sh" "$COVALENT_DIR/"
cp "run.sh" "$WRAPPER_SCRIPT"

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

# Validate the IPFS plist file
plutil -lint "$IPFS_PLIST_FILE"

# Load the IPFS daemon
launchctl unload "$IPFS_PLIST_FILE" # Unload first to ensure no conflicts
launchctl load "$IPFS_PLIST_FILE" || {
    echo "Failed to load IPFS plist. Check the plist and system logs for more details."
    exit 1
}

# Create the light client launchd plist file with the correct executable path
cat <<EOF > "$PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.covalent.light-client</string>
    <key>ProgramArguments</key>
    <array>
        <string>$HOME/.covalent/run.sh</string>
    </array>
    <key>EnvironmentVariables</key>
    <dict>
        <key>COVALENT_DIR</key>
        <string>$HOME/.covalent</string>
        <key>PRIVATE_KEY</key>
        <string>$1</string>
    </dict>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>ThrottleInterval</key>
    <integer>30</integer> <!-- Prevents rapid restarts -->
    <key>StandardOutPath</key>
    <string>$HOME/.covalent/light-client.log</string>
    <key>StandardErrorPath</key>
    <string>$HOME/.covalent/light-client.log</string>
</dict>
</plist>
EOF

# Validate the light client plist file
plutil -lint "$PLIST_FILE"

# Load the light client daemon
launchctl unload "$PLIST_FILE" # Unload first to ensure no conflicts
launchctl load "$PLIST_FILE" || {
    echo "Failed to load Light Client plist. Check the plist and system logs for more details."
    exit 1
}

echo "Installation completed. The IPFS daemon and the light client daemon are now running."