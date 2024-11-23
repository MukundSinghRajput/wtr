#!/bin/bash
APP_NAME="wtr"
OUTPUT_DIR="binaries"
mkdir -p $OUTPUT_DIR
PLATFORMS=("linux/amd64" "darwin/amd64" "windows/amd64")
echo "Building binaries..."
for PLATFORM in "${PLATFORMS[@]}"; do
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)
    OUTPUT_NAME="$OUTPUT_DIR/$APP_NAME-$OS-$ARCH"
    if [ "$OS" == "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi
    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT_NAME .
    if [ $? -ne 0 ]; then
        echo "Failed to build for $OS/$ARCH"
    else
        echo "Binary created: $OUTPUT_NAME"
    fi
done
echo "All binaries built successfully."