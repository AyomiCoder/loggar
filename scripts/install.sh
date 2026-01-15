#!/bin/bash
# Installation script for Loggar CLI

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="loggar"
GITHUB_REPO="AyomiCoder/loggar"
VERSION="latest"

echo "Installing Loggar CLI..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
    echo "Error: Unsupported operating system: $OS"
    exit 1
fi

# Map architecture names
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo "Error: Unsupported architecture: $ARCH"
    exit 1
fi

BINARY_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${BINARY_NAME}_${OS}_${ARCH}"

echo "Downloading ${BINARY_NAME} for ${OS}/${ARCH}..."

# Download binary
TMP_FILE=$(mktemp)
if command -v curl > /dev/null; then
    curl -fsSL "$BINARY_URL" -o "$TMP_FILE"
elif command -v wget > /dev/null; then
    wget -q "$BINARY_URL" -O "$TMP_FILE"
else
    echo "Error: curl or wget is required"
    exit 1
fi

# Make executable
chmod +x "$TMP_FILE"

# Install to /usr/local/bin
echo "Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
else
    sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "✓ Loggar CLI installed successfully!"
echo "Run 'loggar --help' to get started"
