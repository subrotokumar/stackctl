#!/usr/bin/env bash
set -e

# 1. Detect OS
OS=$(uname -s)
case "$OS" in
  Linux*)  OS="Linux" ;;
  Darwin*) OS="Darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS="Windows" ;;
  *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

# 2. Detect Architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)        ARCH="x86_64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

# 3. Fetch latest tag
TAG=$(curl -s https://api.github.com/repos/subrotokumar/stackctl/releases/latest | grep '"tag_name":' | cut -d'"' -f4)

# 4. Download & Extract based on OS
echo "Downloading stackctl $TAG for ${OS}_${ARCH}..."

if [ "$OS" = "Windows" ]; then
  # Windows uses .zip archives
  TARGET_DIR="$HOME/AppData/Local/Programs/stackctl"
  mkdir -p "$TARGET_DIR"
  
  TMP_ZIP=$(mktemp)
  curl -sL "https://github.com/subrotokumar/stackctl/releases/download/$TAG/stackctl_${OS}_${ARCH}.zip" -o "$TMP_ZIP"
  unzip -o -q "$TMP_ZIP" -d "$TARGET_DIR"
  rm -f "$TMP_ZIP"
  
  # Ensure target directory is in PATH
  if [[ ":$PATH:" != *":$TARGET_DIR:"* ]]; then
    echo "export PATH=\"\$PATH:$TARGET_DIR\"" >> "$HOME/.bashrc"
    export PATH="$PATH:$TARGET_DIR"
  fi
  echo "Successfully installed stackctl to $TARGET_DIR/stackctl.exe"
else
  # Linux and macOS use .tar.gz archives
  SUDO=""
  if [ "$(id -u)" -ne 0 ]; then
    SUDO="sudo"
  fi
  curl -sL "https://github.com/subrotokumar/stackctl/releases/download/$TAG/stackctl_${OS}_${ARCH}.tar.gz" | $SUDO tar -xz -C /usr/local/bin stackctl
  echo "Successfully installed stackctl to /usr/local/bin/stackctl"
fi