#!/bin/bash

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