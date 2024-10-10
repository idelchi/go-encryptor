#!/bin/sh
set -e

# Default values
VERSION="v0.0"
OUTPUT_DIR="./bin"

# Usage function
usage() {
    echo "Usage: $0 [-v version] [-o output]"
    exit 1
}

# Parse arguments
while getopts "v:o:" opt; do
    case "$opt" in
        v) VERSION="$OPTARG" ;;
        o) OUTPUT_DIR="$OPTARG" ;;
        *) usage ;;
    esac
done

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  cygwin_nt*) OS="windows" ;;
  mingw*) OS="windows" ;;
  msys_nt*) OS="windows" ;;
esac

case $ARCH in
    x86_64) ARCH="amd64" ;;
    armv6*) ARCH="armv6" ;;
    armv7*) ARCH="armv7" ;;
    aarch64) ARCH="arm64" ;;
    i386) ARCH="x86" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

# Set the format based on OS
FORMAT="tar.gz"
if [ "$OS" = "windows" ]; then
    FORMAT="zip"
fi

# Construct the download URL
BASE_URL="https://github.com/idelchi/gocry/releases/download"
BINARY_NAME="gocry_${OS}_${ARCH}.${FORMAT}"
URL="${BASE_URL}/${VERSION}/${BINARY_NAME}"

tmp=$(mktemp)

# Download and extract/install
echo "Downloading $BINARY_NAME from $URL"
code=$(curl -L -o ${tmp} ${URL})

if [ "$code" != "200" ]; then
  echo "Failed to download $URL"

  return 1
fi

if [ "$FORMAT" = "tar.gz" ]; then
    tar -C $OUTPUT_DIR -xzf $tmp
else
    unzip -d $OUTPUT_DIR $tmp
fi

rm $tmp -rf

echo "gocry installed to $OUTPUT_DIR"
xxxx
