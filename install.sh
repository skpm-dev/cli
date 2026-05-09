#!/bin/sh
set -e

REPO="skpm-dev/cli"
BINARY="skpm"
INSTALL_DIR="/usr/local/bin"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

case "$arch" in
  x86_64) arch="amd64" ;;
  arm64 | aarch64) arch="arm64" ;;
  *)
    echo "Unsupported architecture: $arch"
    exit 1
    ;;
esac

case "$os" in
  linux | darwin) ;;
  *)
    echo "Unsupported OS: $os"
    echo "On Windows, download the binary manually from https://github.com/$REPO/releases"
    exit 1
    ;;
esac

version=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$version" ]; then
  echo "Could not determine latest version"
  exit 1
fi

archive="skpm_${os}_${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$version/$archive"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading skpm $version..."
curl -fsSL "$url" -o "$tmp/$archive"
tar -xzf "$tmp/$archive" -C "$tmp"

binary=$(find "$tmp" -name "$BINARY" -type f | head -1)

if [ -z "$binary" ]; then
  echo "Could not find $BINARY in downloaded archive"
  exit 1
fi

chmod +x "$binary"

if [ -w "$INSTALL_DIR" ]; then
  cp "$binary" "$INSTALL_DIR/$BINARY"
else
  echo "Installing to $INSTALL_DIR (you may be prompted for your password)..."
  sudo cp "$binary" "$INSTALL_DIR/$BINARY"
fi

echo "Installed skpm $version to $INSTALL_DIR/$BINARY"
