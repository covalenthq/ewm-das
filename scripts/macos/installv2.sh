#!/bin/bash

# Check if the OS is macOS (Darwin)
check_os() {
  if [[ "$(uname)" != "Darwin" ]]; then
    echo "This script is only compatible with macOS."
    exit 1
  fi
}

# Validate the private key
validate_private_key() {
  if [ -z "$1" ]; then
    echo "Private Key is required."
    echo "Usage: $0 <private-key>"
    exit 1
  fi

  if ! [[ "$1" =~ ^[0-9a-fA-F]{64}$ ]]; then
    echo "Error: private key is not a valid 64-character hexadecimal number."
    exit 1
  fi
}

# Define variables and paths
define_paths() {
  COVALENT_DIR="$HOME/.covalent"
  PLIST_FILE="$HOME/Library/LaunchAgents/com.covalent.light-client.plist"
  IPFS_PLIST_FILE="$HOME/Library/LaunchAgents/com.covalent.ipfs.plist"
  IPFS_PATH=$(which ipfs)
  WRAPPER_SCRIPT="$COVALENT_DIR/run.sh"
  IPFS_REPO_DIR="$HOME/.ipfs"

  EXECUTABLE_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.7.0/macos/light-client"
  TRUSTED_SETUP_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.7.0/macos/trusted_setup.txt"
  UNINSTALL_SCRIPT_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.7.0/macos/uninstall.sh"
  RUN_SCRIPT_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.7.0/macos/run.sh"
}

# Uninstall previous versions
uninstall_previous() {
  if [ -f uninstall.sh ]; then
    bash uninstall.sh
  fi
}

# Download and install files
download_files() {
  mkdir -p "$COVALENT_DIR"
  
  curl -o "$COVALENT_DIR/light-client" "$EXECUTABLE_URL"
  curl -o "$COVALENT_DIR/trusted_setup.txt" "$TRUSTED_SETUP_URL"
  curl -o "$COVALENT_DIR/uninstall.sh" "$UNINSTALL_SCRIPT_URL"
  curl -o "$WRAPPER_SCRIPT" "$RUN_SCRIPT_URL"
  
  chmod +x "$COVALENT_DIR/light-client"
  chmod +x "$COVALENT_DIR/uninstall.sh"
  chmod +x "$WRAPPER_SCRIPT"

  # Bypass Gatekeeper for the executable
  spctl --add --label "Trusted" "$COVALENT_DIR/light-client"
  spctl --enable --label "Trusted"
}

# Create and configure IPFS plist
create_ipfs_plist() {
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
  <integer>10</integer>
  <key>StandardOutPath</key>
  <string>$COVALENT_DIR/ipfs.log</string>
  <key>StandardErrorPath</key>
  <string>$COVALENT_DIR/ipfs.error.log</string>
</dict>
</plist>
EOF

  plutil -lint "$IPFS_PLIST_FILE"
  launchctl unload "$IPFS_PLIST_FILE"
  launchctl load "$IPFS_PLIST_FILE"
}

# Create and configure Light Client plist
create_light_client_plist() {
  cat <<EOF > "$PLIST_FILE"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.covalent.light-client</string>
  <key>ProgramArguments</key>
  <array>
    <string>$WRAPPER_SCRIPT</string>
  </array>
  <key>EnvironmentVariables</key>
  <dict>
    <key>PINNER_DIR</key>
    <string>$COVALENT_DIR</string>
    <key>PRIVATE_KEY</key>
    <string>$1</string>
  </dict>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>ThrottleInterval</key>
  <integer>30</integer>
  <key>StandardOutPath</key>
  <string>$COVALENT_DIR/light-client.log</string>
  <key>StandardErrorPath</key>
  <string>$COVALENT_DIR/light-client.log</string>
</dict>
</plist>
EOF

  plutil -lint "$PLIST_FILE"
  launchctl unload "$PLIST_FILE"
  launchctl load "$PLIST_FILE"
}

# Main installation function
install() {
  check_os
  validate_private_key "$1"
  define_paths
  uninstall_previous
  download_files
  create_ipfs_plist
  create_light_client_plist "$1"

  echo "Installation completed. The IPFS daemon and light client are now running."
}

# Execute the installation
install "$1"