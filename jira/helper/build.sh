#!/bin/bash

# Build script for JIRA transition checker
# This script checks if the main binary exists and builds it if it doesn't

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first."
    exit 1
fi

# Check if main.go exists
if [ ! -f "main.go" ]; then
    print_error "main.go not found in current directory"
    exit 1
fi

# Check if main binary exists and is newer than main.go
if [ -f "main" ] && [ "main" -nt "main.go" ]; then
    print_status "Binary 'main' already exists and is up to date"
    print_status "Binary location: $(pwd)/main"
else
    print_status "Building main binary..."
    
    # Check if go.mod exists, if not initialize module
    if [ ! -f "go.mod" ]; then
        print_warning "go.mod not found, initializing Go module..."
        go mod init jira-helper
    fi
    
    # Download dependencies
    print_status "Downloading dependencies..."
    go mod tidy
    
    # Build the binary
    print_status "Running: go build -o main ."
    if go build -o main .; then
        print_status "Build successful!"
        
        # Make it executable
        chmod +x main
        print_status "Binary is now executable"
    else
        print_error "Build failed!"
        exit 1
    fi
fi

print_status "Build process completed successfully!"
print_status "You can now run: ./main [JIRA-ID1] [JIRA-ID2] ..."
print_status "Binary location: $(pwd)/main" 