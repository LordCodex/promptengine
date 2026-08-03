#!/bin/bash
set -e

# PromptEngine POSIX installer script
# Downloads and installs the correct binary for your OS and Architecture.

OWNER="LordCodex"
REPO="promptengine"
BINARY_NAME="promptengine"
DEFAULT_INSTALL_DIR="/usr/local/bin"

# Detect OS
OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS_TYPE" in
  darwin)  OS="darwin" ;;
  linux)   OS="linux" ;;
  *)       echo "Unsupported operating system: $OS_TYPE"; exit 1 ;;
esac

# Detect Architecture
ARCH_TYPE=$(uname -m)
case "$ARCH_TYPE" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)            echo "Unsupported CPU architecture: $ARCH_TYPE"; exit 1 ;;
esac

# Determine installation directory
INSTALL_DIR="${PROMPTENGINE_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"

# Fetch latest release version
echo "Fetching latest version tag..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$LATEST_TAG" ]; then
  LATEST_TAG="v1.0.0"
fi

VERSION=${LATEST_TAG#v}

echo "Installing $BINARY_NAME $LATEST_TAG for $OS-$ARCH..."

# Build download url
DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/download/$LATEST_TAG/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
TEMP_DIR=$(mktemp -d)

clean_up() {
  rm -rf "$TEMP_DIR"
}
trap clean_up EXIT

echo "Downloading from: $DOWNLOAD_URL"
if ! curl -L -s -f -o "$TEMP_DIR/archive.tar.gz" "$DOWNLOAD_URL"; then
  echo "Error: Failed to download release archive."
  exit 1
fi

# Extract and install
tar -xzf "$TEMP_DIR/archive.tar.gz" -C "$TEMP_DIR"

if [ -f "$TEMP_DIR/$BINARY_NAME" ]; then
  echo "Installing binary to $INSTALL_DIR..."
  if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
  else
    echo "Requires sudo privileges to write to $INSTALL_DIR"
    sudo mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
  fi
  echo "✓ Successfully installed $BINARY_NAME version $VERSION to $INSTALL_DIR!"
else
  echo "Error: Binary not found in extracted archive."
  exit 1
fi
