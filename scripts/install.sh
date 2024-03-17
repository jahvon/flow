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
    curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep -o '"tag_name": ".*"' | sed 's/"tag_name": "//;s/"//'
}

REPO_OWNER="jahvon"
REPO_NAME="flow"
CLI_BINARY="flow"

OS=$(get_os)
ARCH=$(get_arch)
if [ -z "$VERSION" ]; then
    VERSION=$(get_latest_version)
fi

DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/$CLI_BINARY-$VERSION-$OS-$ARCH.tar.gz"
TMP_DIR=$(mktemp -d)

echo "Downloading $CLI_BINARY $VERSION for $OS/$ARCH..."
wget -q "$DOWNLOAD_URL" -O "$TMP_DIR/$CLI_BINARY-$VERSION-$OS-$ARCH.tar.gz"
if [ $? -ne 0 ]; then
    echo "Failed to download $DOWNLOAD_URL"
    exit 1
fi

INSTALL_DIR="/usr/local/bin"
echo "Installing $CLI_BINARY $VERSION to $INSTALL_DIR..."
tar -xzf "$TMP_DIR/$CLI_BINARY-$VERSION-$OS-$ARCH.tar.gz" -C "$TMP_DIR"

mv "$TMP_DIR/$CLI_BINARY" "$INSTALL_DIR/$CLI_BINARY"
chmod +x "$INSTALL_DIR/$CLI_BINARY"

echo "$CLI_BINARY was installed successfully to $INSTALL_DIR/$CLI_BINARY"
if command -v $CLI_BINARY --version >/dev/null; then
    echo "Run '$CLI_BINARY --help' to get started"
else
    echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
    echo "  export PATH=$INSTALL_DIR:\$PATH"
fi

exit 0
