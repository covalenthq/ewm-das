#!/bin/bash

# Define the source file and destination directory
SOURCE_FILE="test/data/trusted_setup.txt"
DEST_DIR="$HOME/.covalent"
DEST_FILE="$DEST_DIR/$(basename $SOURCE_FILE)"

# Check if the source file exists
if [ ! -f "$SOURCE_FILE" ]; then
    echo "Source file $SOURCE_FILE does not exist."
    exit 1
fi

# Create the destination directory if it doesn't exist
if [ ! -d "$DEST_DIR" ]; then
    echo "Creating destination directory $DEST_DIR..."
    mkdir -p "$DEST_DIR"
fi

# Copy the file to the destination directory
echo "Copying $SOURCE_FILE to $DEST_DIR..."
cp "$SOURCE_FILE" "$DEST_DIR/"

# Confirm the file has been copied
if [ -f "$DEST_FILE" ]; then
    echo "File copied successfully to $DEST_FILE"
else
    echo "Failed to copy the file"
    exit 1
fi