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

  EXECUTABLE_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.8.0/macos/light-client"
  TRUSTED_SETUP_URL="https://storage.googleapis.com/ewm-release-artefacts/v0.8.0/macos/trusted_setup.txt"
}

# Uninstall previous versions
uninstall_previous() {
  # Paths for the setup
  COVALENT_DIR="$HOME/.covalent"
  PLIST_FILE_NAME="com.covalent.light-client.plist"
  IPFS_PLIST_FILE_NAME="com.covalent.ipfs.plist"

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
  remove_plist "$PLIST_FILE_NAME"
  remove_plist "$IPFS_PLIST_FILE_NAME"

  # Remove the .covalent and .covalenthq directories and their contents
  remove_directory "$COVALENT_DIR"

  echo "Uninstallation completed. The light client and IPFS daemons have been removed."
}

# Create uninstall script in the covalent directory
create_uninstall_script() {
  cat <<EOF > "$COVALENT_DIR/uninstall.sh"
#!/bin/bash

# Paths for the new setup
COVALENT_DIR="\$HOME/.covalent"
PLIST_FILE_NAME="com.covalent.light-client.plist"
IPFS_PLIST_FILE_NAME="com.covalent.ipfs.plist"

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
remove_plist "\$PLIST_FILE_NAME"
remove_plist "\$IPFS_PLIST_FILE_NAME"

# Remove the .covalent and .covalenthq directories and their contents
remove_directory "\$COVALENT_DIR"

echo "Uninstallation completed. The light client and IPFS daemons have been removed."
EOF

  chmod +x "$COVALENT_DIR/uninstall.sh"
}

# Create the run script inside the installation
create_run_script() {
  cat <<EOF > "$WRAPPER_SCRIPT"
#!/bin/bash

# Define the directory again in the wrapper script, it will use the value assigned during installation
COVALENT_DIR="\${COVALENT_DIR:-\$HOME/.covalent}"  # Default to ~/.covalent if not set
SERVICE_NAME="\${EXECUTABLE:-light-client}"        # Default to light-client if not set

# Check if PRIVATE_KEY is set; if not, exit with an error
if [ -z "\$PRIVATE_KEY" ]; then
  echo "Error: PRIVATE_KEY environment variable is not set."
  exit 1
fi

# Check if PRIVATE_KEY is a valid 64-character hexadecimal number
if ! [[ "\$PRIVATE_KEY" =~ ^[0-9a-fA-F]{64}$ ]]; then
  echo "Error: PRIVATE_KEY is not a valid 64-character hexadecimal number."
  exit 1
fi

# Ensure that only one instance of the service is running
if pgrep -f "\$SERVICE_NAME" > /dev/null 2>&1; then
    echo "\$SERVICE_NAME is already running."
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
"\$COVALENT_DIR/\$SERVICE_NAME" \\
    --loglevel debug \\
    --rpc-url ws://34.42.69.93:8080/rpc \\
    --collect-url https://ewm-light-clients-v2-838505730421.us-central1.run.app \\
    --private-key "\$PRIVATE_KEY"
EOF

  chmod +x "$WRAPPER_SCRIPT"
}

# Download and install files
download_files() {
  mkdir -p "$COVALENT_DIR"
  
  curl -o "$COVALENT_DIR/light-client" "$EXECUTABLE_URL"
  curl -o "$COVALENT_DIR/trusted_setup.txt" "$TRUSTED_SETUP_URL"
  
  chmod +x "$COVALENT_DIR/light-client"

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
    <key>COVALENT_DIR</key>
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

cleanup() {
  CWD=$(pwd)  # Get the current working directory
  rm -f "$CWD/$PLIST_FILE"
  rm -f "$CWD/$IPFS_PLIST_FILE"
}

# Main installation function
install() {
  check_os
  validate_private_key "$1"
  define_paths
  uninstall_previous
  download_files
  create_ipfs_plist
  create_run_script
  create_light_client_plist "$1"
  create_uninstall_script
  cleanup

  echo "Installation completed. The IPFS daemon and light client are now running."
}

# Execute the installation
install "$1"