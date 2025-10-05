#!/bin/bash

# 🛡️ Praetorian Go - Testing Setup Script
# Installs mutation testing tools and sets up testing environment

set -euo pipefail

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Guard clause: Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

log_info "Setting up mutation testing environment..."

# Install go-mutesting
log_info "Installing go-mutesting..."
if go install github.com/zimmski/go-mutesting/...@latest; then
    log_success "go-mutesting installed successfully"
else
    log_warning "Failed to install go-mutesting. You may need to accept Xcode license: sudo xcodebuild -license"
    exit 1
fi

# Verify installation
if command -v go-mutesting &> /dev/null; then
    log_success "go-mutesting is ready to use"
else
    log_warning "go-mutesting installation may have failed. Check your GOPATH/bin"
fi

log_info "Testing environment setup complete!"
log_info "You can now run: ./scripts/generate-reports.sh"
