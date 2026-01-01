#!/bin/bash
# mcinit installation script for Unix-like systems (macOS, Linux)

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    *)          echo -e "${RED}Unsupported operating system: $OS${NC}"; exit 1;;
esac

case "$ARCH" in
    x86_64)     ARCH="amd64";;
    aarch64)    ARCH="arm64";;
    arm64)      ARCH="arm64";;
    *)          echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1;;
esac

# Get latest release version
echo -e "${GREEN}Fetching latest release...${NC}"
LATEST_VERSION=$(curl -s https://api.github.com/repos/jackh54/mcinit/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Failed to fetch latest version${NC}"
    exit 1
fi

echo -e "${GREEN}Latest version: $LATEST_VERSION${NC}"

# Download URL
DOWNLOAD_URL="https://github.com/jackh54/mcinit/releases/download/${LATEST_VERSION}/mcinit_${LATEST_VERSION#v}_${OS}_${ARCH}.tar.gz"

echo -e "${GREEN}Downloading mcinit...${NC}"
curl -L -o /tmp/mcinit.tar.gz "$DOWNLOAD_URL"

# Extract
echo -e "${GREEN}Extracting...${NC}"
tar -xzf /tmp/mcinit.tar.gz -C /tmp

# Install
INSTALL_DIR="/usr/local/bin"
if [ -w "$INSTALL_DIR" ]; then
    echo -e "${GREEN}Installing to $INSTALL_DIR...${NC}"
    mv /tmp/mcinit "$INSTALL_DIR/mcinit"
    chmod +x "$INSTALL_DIR/mcinit"
else
    echo -e "${YELLOW}No write permission for $INSTALL_DIR, using sudo...${NC}"
    sudo mv /tmp/mcinit "$INSTALL_DIR/mcinit"
    sudo chmod +x "$INSTALL_DIR/mcinit"
fi

# Cleanup
rm -f /tmp/mcinit.tar.gz

echo -e "${GREEN}mcinit installed successfully!${NC}"
echo -e "${GREEN}Run 'mcinit --help' to get started${NC}"

