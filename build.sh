#!/bin/bash

# Create executables directory if it doesn't exist
mkdir -p executables

# Build for all platforms
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o executables/GridironGo-linux

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o executables/GridironGo-mac

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o executables/GridironGo-windows.exe

echo "Build complete. Executables are in the 'executables' directory."
