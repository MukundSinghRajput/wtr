#!/bin/bash
APP_NAME="wtr"
BINARIES_DIR="./binaries"
INSTALL_DIR=""
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
fi
if [ "$OS" == "darwin" ]; then
    OS="darwin"
elif [[ "$OS" == *"mingw"* || "$OS" == *"msys"* ]]; then
    OS="windows"
fi

BINARY="$BINARIES_DIR/$APP_NAME-$OS-$ARCH"
if [ "$OS" == "windows" ]; then
    BINARY="$BINARY.exe"
fi

if [ ! -f "$BINARY" ]; then
    echo "Error: Binary for $OS/$ARCH not found in $BINARIES_DIR"
    exit 1
fi

if [ "$OS" == "windows" ]; then
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
else
    INSTALL_DIR="/usr/local/bin"
fi

echo "Installing $APP_NAME for $OS/$ARCH..."
cp "$BINARY" "$INSTALL_DIR/$APP_NAME"

if [ "$OS" != "windows" ]; then
    chmod +x "$INSTALL_DIR/$APP_NAME"
fi

if [ "$OS" == "windows" ]; then
    powershell.exe -Command "[System.Environment]::SetEnvironmentVariable('PATH', \$([System.Environment]::GetEnvironmentVariable('PATH', 'User') + ';$HOME\\bin'), 'User')"
    echo "Added $HOME/bin to PATH. Restart your terminal or shell to apply changes."
else
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo "Ensure $INSTALL_DIR is in your PATH to run $APP_NAME directly."
    fi
fi

echo "$APP_NAME installed successfully! Run '$APP_NAME' to start."