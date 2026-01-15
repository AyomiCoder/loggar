#!/bin/bash
# Build script for Loggar CLI

set -e

VERSION="0.1.0"
APP_NAME="loggar"
BUILD_DIR="dist"

echo "Building Loggar CLI v${VERSION}..."

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build for macOS (Intel)
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/${APP_NAME}_darwin_amd64 ./cmd/loggar

# Build for macOS (Apple Silicon)
echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o ${BUILD_DIR}/${APP_NAME}_darwin_arm64 ./cmd/loggar

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${APP_NAME}_linux_amd64 ./cmd/loggar

echo "✓ Build complete! Binaries are in ${BUILD_DIR}/"
ls -lh ${BUILD_DIR}/
