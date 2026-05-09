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

sudo install -m 755 "$tmp/$BINARY" "$INSTALL_DIR/$BINARY"

echo "Installed skpm $version to $INSTALL_DIR/$BINARY"
