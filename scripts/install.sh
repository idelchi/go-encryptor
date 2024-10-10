#!/bin/sh
set -e

# Default values
VERSION="latest"
OUTPUT_DIR="/usr/local/bin"

# Usage function
usage() {
    echo "Usage: $0 [-v version] [-o output_directory]"
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

case $ARCH in
    x86_64) ARCH="x86_64" ;;
    armv6*) ARCH="armv6" ;;
    armv7*) ARCH="armv7" ;;
    aarch64) ARCH="arm64" ;;
    i386) ARCH="i386" ;;
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
URL="${BASE_URL}/v${VERSION}/${BINARY_NAME}"

# Download and extract/install
echo "Downloading $BINARY_NAME from $URL"
curl -L -o /tmp/$BINARY_NAME $URL

if [ "$FORMAT" = "tar.gz" ]; then
    tar -C $OUTPUT_DIR -xzf /tmp/$BINARY_NAME
else
    unzip -d $OUTPUT_DIR /tmp/$BINARY_NAME
fi

# Cleanup
rm /tmp/$BINARY_NAME

echo "gocry installed to $OUTPUT_DIR"
