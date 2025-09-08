#!/bin/bash

echo "🔍 Verifying Environment Variables for Artifactory Integration..."

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

# Check environment variables
echo "Checking required environment variables:"

# HF_HUB_ETAG_TIMEOUT
if [ -n "${HF_HUB_ETAG_TIMEOUT}" ]; then
    print_status "✅ HF_HUB_ETAG_TIMEOUT: ${HF_HUB_ETAG_TIMEOUT}"
else
    print_warning "⚠️  HF_HUB_ETAG_TIMEOUT not set"
fi

# HF_HUB_DOWNLOAD_TIMEOUT
if [ -n "${HF_HUB_DOWNLOAD_TIMEOUT}" ]; then
    print_status "✅ HF_HUB_DOWNLOAD_TIMEOUT: ${HF_HUB_DOWNLOAD_TIMEOUT}"
else
    print_warning "⚠️  HF_HUB_DOWNLOAD_TIMEOUT not set"
fi

# HF_ENDPOINT
if [ -n "${HF_ENDPOINT}" ]; then
    print_status "✅ HF_ENDPOINT: ${HF_ENDPOINT}"
else
    print_warning "⚠️  HF_ENDPOINT not set"
fi

# HF_TOKEN
if [ -n "${HF_TOKEN}" ]; then
    print_status "✅ HF_TOKEN: [HIDDEN]"
else
    print_warning "⚠️  HF_TOKEN not set"
fi

# TRANSFORMERS_CACHE
if [ -n "${TRANSFORMERS_CACHE}" ]; then
    print_status "✅ TRANSFORMERS_CACHE: ${TRANSFORMERS_CACHE}"
else
    print_warning "⚠️  TRANSFORMERS_CACHE not set"
fi

# HF_HOME
if [ -n "${HF_HOME}" ]; then
    print_status "✅ HF_HOME: ${HF_HOME}"
else
    print_warning "⚠️  HF_HOME not set"
fi

echo ""
echo "📋 Summary:"
echo "- All environment variables should be set for optimal Artifactory integration"
echo "- Missing variables will use default values or may cause issues"
echo "- Run 'source .env' to load variables from .env file"
