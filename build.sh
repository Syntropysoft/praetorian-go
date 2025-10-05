#!/bin/bash

# Build script for Praetorian Go - Multi-platform binaries
echo "🔨 Building Praetorian binaries for multiple platforms..."

# Create dist directory
mkdir -p dist

# Build for different platforms
echo "📦 Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/praetorian-linux-amd64 ./cmd/praetorian

echo "📦 Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/praetorian-windows-amd64.exe ./cmd/praetorian

echo "📦 Building for macOS AMD64 (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/praetorian-darwin-amd64 ./cmd/praetorian

echo "📦 Building for macOS ARM64 (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/praetorian-darwin-arm64 ./cmd/praetorian

echo "📦 Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/praetorian-linux-arm64 ./cmd/praetorian

echo "📦 Building for Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dist/praetorian-windows-arm64.exe ./cmd/praetorian

# Show results
echo ""
echo "✅ Build completed! Binaries created in dist/ directory:"
ls -la dist/

echo ""
echo "📊 Binary sizes:"
for file in dist/*; do
    size=$(du -h "$file" | cut -f1)
    echo "  $(basename "$file"): $size"
done 