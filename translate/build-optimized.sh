#!/bin/bash

set -e

echo "🚀 Building Translation Service with Optimizations..."

# Configuration
IMAGE_NAME="ai-translate"
CACHE_IMAGE="${IMAGE_NAME}:cache"
LATEST_IMAGE="${IMAGE_NAME}:latest"
BUILD_NUMBER=${BUILD_NUMBER:-$(date +%s)}

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

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

    # Clean up any existing containers with the same name
    print_status "Cleaning up existing containers..."
    docker rm -f test-ai-translate-container 2>/dev/null || true

# Build with layer caching and Artifactory configuration
print_status "Building Docker image with layer caching and Artifactory configuration..."
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --cache-from type=local,src=/tmp/.buildx-cache \
    --cache-to type=local,dest=/tmp/.buildx-cache-new,mode=max \
    --build-arg HF_HUB_ETAG_TIMEOUT=86400 \
    --build-arg HF_HUB_DOWNLOAD_TIMEOUT=86400 \
    --build-arg HF_ENDPOINT=https://apptrustswampupc.jfrog.io/artifactory/api/huggingfaceml/dev-huggingfaceml-remote \
    --tag ${IMAGE_NAME}:${BUILD_NUMBER} \
    --tag ${LATEST_IMAGE} \
    --progress=plain \
    .

# Move cache
rm -rf /tmp/.buildx-cache
mv /tmp/.buildx-cache-new /tmp/.buildx-cache

print_status "Build completed successfully!"
print_status "Image tags: ${IMAGE_NAME}:${BUILD_NUMBER}, ${LATEST_IMAGE}"

# Optional: Test the built image
if [[ "${TEST_IMAGE}" == "true" ]]; then
    print_status "Testing the built image..."
    
    # Run container
    docker run -d --name test-ai-translate-container -p 8002:8002 ${LATEST_IMAGE}
    
    # Wait for container to start
    sleep 10
    
    # Test health endpoint
    if curl -f http://localhost:8002/health > /dev/null 2>&1; then
        print_status "✅ Health check passed"
    else
        print_error "❌ Health check failed"
        docker logs test-ai-translate-container
        exit 1
    fi
    
    # Test root endpoint
    if curl -f http://localhost:8002/ > /dev/null 2>&1; then
        print_status "✅ Root endpoint check passed"
    else
        print_error "❌ Root endpoint check failed"
        docker logs test-ai-translate-container
        exit 1
    fi
    
    # Clean up test container
    docker rm -f test-ai-translate-container
    print_status "✅ All tests passed!"
fi

echo "🎉 Build process completed successfully!"
echo "📦 Image: ${IMAGE_NAME}:${BUILD_NUMBER}"
echo "📦 Latest: ${LATEST_IMAGE}"
echo "💡 To test the image: TEST_IMAGE=true ./build-optimized.sh"
