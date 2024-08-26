#!/bin/bash

# Paths
COVALENT_DIR="$HOME/.covalenthq"
PLIST_FILE="com.covalenthq.light-client.plist"

# Unload the daemon
launchctl unload "$HOME/Library/LaunchAgents/$PLIST_FILE"

# Remove the plist file
rm "$HOME/Library/LaunchAgents/$PLIST_FILE"

# Remove the .covalent directory and its contents
rm -rf "$COVALENT_DIR"

echo "Uninstallation completed. The light client daemon has been removed."