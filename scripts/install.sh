#!/bin/bash

set -e

get_os() {
    case "$(uname -s)" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        *)          echo "Unsupported operating system"; exit 1;;
    esac
    echo "$OS"
}

get_arch() {
    case "$(uname -m)" in
        x86_64) ARCH=amd64;;
        arm64)  ARCH=arm64;;
        *)      echo "Unsupported architecture"; exit 1;;
    esac
    echo "$ARCH"
}

get_latest_version() {
    curl -s "https://api.github.com/repos/$OWNER/$NAME/releases/latest" | grep -o '"tag_name": ".*"' | sed 's/"tag_name": "//;s/"//'
}

OWNER="jahvon"
NAME="flow"
BINARY="flow"

OS=$(get_os)
ARCH=$(get_arch)
if [ -z "$VERSION" ]; then
    VERSION=$(get_latest_version)
fi

DOWNLOAD_URL="https://github.com/${OWNER}/${NAME}/releases/download/${VERSION}/${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
TMP_DIR=$(mktemp -d)
DOWNLOAD_PATH="${TMP_DIR}/${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"

echo "Downloading $BINARY $VERSION for $OS/$ARCH..."
wget -q "$DOWNLOAD_URL" -O "$DOWNLOAD_PATH"
if [ $? -ne 0 ]; then
    echo "Failed to download $DOWNLOAD_URL"
    exit 1
fi

INSTALL_DIR="/usr/local/bin"
echo "Installing $BINARY $VERSION to $INSTALL_DIR..."
tar -xzf "$DOWNLOAD_PATH" -C "$TMP_DIR"

chmod +x "$TMP_DIR/$BINARY"
sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"

echo "$BINARY was installed successfully to $INSTALL_DIR/$BINARY"
if command -v $BINARY --version >/dev/null; then
    echo "Run '$BINARY --help' to get started"
else
    echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
    echo "  export PATH=$INSTALL_DIR:\$PATH"
fi

exit 0
