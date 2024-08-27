#!/bin/bash

# Paths
COVALENT_DIR="$HOME/.covalenthq"
PLIST_FILE="com.covalenthq.light-client.plist"
IPFS_PLIST_FILE="com.covalenthq.ipfs.plist"

# Unload the light client daemon and ipfs daemon
launchctl unload "$HOME/Library/LaunchAgents/$PLIST_FILE"
launchctl unload "$HOME/Library/LaunchAgents/$IPFS_PLIST_FILE"

# Remove the plist file
rm "$HOME/Library/LaunchAgents/$PLIST_FILE"

# Remove the .covalent directory and its contents
rm -rf "$COVALENT_DIR"

echo "Uninstallation completed. The light client & ipfs daemons has been removed."