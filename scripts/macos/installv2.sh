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
  RUN_SCRIPT_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.7.0/macos/run.sh"
}

# Uninstall previous versions
uninstall_previous() {
  # Paths for the old setup
  OLD_COVALENT_DIR="$HOME/.covalenthq"
  OLD_PLIST_FILE="com.covalenthq.light-client.plist"
  OLD_IPFS_PLIST_FILE="com.covalenthq.ipfs.plist"

  # Paths for the new setup
  COVALENT_DIR="$HOME/.covalent"
  PLIST_FILE="com.covalent.light-client.plist"
  IPFS_PLIST_FILE="com.covalent.ipfs.plist"

  # Function to unload and remove plist files
  remove_plist() {
    local plist_file="$1"
    if [ -f "$HOME/Library/LaunchAgents/$plist_file" ]; then
      launchctl unload "$HOME/Library/LaunchAgents/$plist_file" || echo "Failed to unload $plist_file"
      rm "$HOME/Library/LaunchAgents/$plist_file" || echo "Failed to remove $plist_file"
    fi
  }

  # Function to remove directories
  remove_directory() {
    local dir="$1"
    if [ -d "$dir" ]; then
      rm -rf "$dir" || echo "Failed to remove directory $dir"
    fi
  }

  # Unload and remove plist files for both old and new setups
  remove_plist "$PLIST_FILE"
  remove_plist "$IPFS_PLIST_FILE"
  remove_plist "$OLD_PLIST_FILE"
  remove_plist "$OLD_IPFS_PLIST_FILE"

  # Remove the .covalent and .covalenthq directories and their contents
  remove_directory "$COVALENT_DIR"
  remove_directory "$OLD_COVALENT_DIR"

  echo "Uninstallation completed. The light client and IPFS daemons for both old and new versions have been removed."
}

# Create uninstall script in the covalent directory
create_uninstall_script() {
  cat <<EOF > "$COVALENT_DIR/uninstall.sh"
#!/bin/bash

# Paths for the old setup
OLD_COVALENT_DIR="\$HOME/.covalenthq"
OLD_PLIST_FILE="com.covalenthq.light-client.plist"
OLD_IPFS_PLIST_FILE="com.covalenthq.ipfs.plist"

# Paths for the new setup
COVALENT_DIR="\$HOME/.covalent"
PLIST_FILE="com.covalent.light-client.plist"
IPFS_PLIST_FILE="com.covalent.ipfs.plist"

# Function to unload and remove plist files
remove_plist() {
  local plist_file="\$1"
  if [ -f "\$HOME/Library/LaunchAgents/\$plist_file" ]; then
    launchctl unload "\$HOME/Library/LaunchAgents/\$plist_file" || echo "Failed to unload \$plist_file"
    rm "\$HOME/Library/LaunchAgents/\$plist_file" || echo "Failed to remove \$plist_file"
  fi
}

# Function to remove directories
remove_directory() {
  local dir="\$1"
  if [ -d "\$dir" ]; then
    rm -rf "\$dir" || echo "Failed to remove directory \$dir"
  fi
}

# Unload and remove plist files for both old and new setups
remove_plist "\$PLIST_FILE"
remove_plist "\$IPFS_PLIST_FILE"
remove_plist "\$OLD_PLIST_FILE"
remove_plist "\$OLD_IPFS_PLIST_FILE"

# Remove the .covalent and .covalenthq directories and their contents
remove_directory "\$COVALENT_DIR"
remove_directory "\$OLD_COVALENT_DIR"

echo "Uninstallation completed. The light client and IPFS daemons for both old and new versions have been removed."
EOF

  chmod +x "$COVALENT_DIR/uninstall.sh"
}

# Download and install files
download_files() {
  mkdir -p "$COVALENT_DIR"
  
  curl -o "$COVALENT_DIR/light-client" "$EXECUTABLE_URL"
  curl -o "$COVALENT_DIR/trusted_setup.txt" "$TRUSTED_SETUP_URL"
  curl -o "$WRAPPER_SCRIPT" "$RUN_SCRIPT_URL"
  
  chmod +x "$COVALENT_DIR/light-client"
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

  # Unload IPFS service using bootout (modern equivalent of unload)
  launchctl bootout gui/"$(id -u)" "$IPFS_PLIST_FILE" || echo "Failed to unload $IPFS_PLIST_FILE"

  # Load IPFS service using bootstrap (modern equivalent of load)
  launchctl bootstrap gui/"$(id -u)" "$IPFS_PLIST_FILE" || echo "Failed to load $IPFS_PLIST_FILE"
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

  # Unload Light Client service using bootout
  launchctl bootout gui/"$(id -u)" "$PLIST_FILE" || echo "Failed to unload $PLIST_FILE"

  # Load Light Client service using bootstrap
  launchctl bootstrap gui/"$(id -u)" "$PLIST_FILE" || echo "Failed to load $PLIST_FILE"
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
  create_uninstall_script

  echo "Installation completed. The IPFS daemon and light client are now running."
}

# Execute the installation
install "$1"