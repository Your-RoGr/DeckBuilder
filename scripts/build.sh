#!/bin/bash

cd "$(dirname $0)"
cd ..

set -e

# Path to the source file
SRC="src/main.go"

# Name of the resulting executable
BIN_NAME="DeckBuilder"

# Installation directory
INSTALL_DIR="/usr/local/bin"

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed! Please install Go and try again."
  exit 1
fi

# Build the project
echo "Building the $BIN_NAME..."
go build -o "$BIN_NAME" "$SRC"

# Install the binary (requires sudo privileges)
echo "Installing $BIN_NAME to $INSTALL_DIR (you may be prompted for your password)..."
sudo mv "$BIN_NAME" "$INSTALL_DIR/"

# Make sure it's executable
sudo chmod +x "$INSTALL_DIR/$BIN_NAME"

echo "Done! You can now run your app using: $BIN_NAME"
