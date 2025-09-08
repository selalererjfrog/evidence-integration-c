#!/bin/bash

set -e

echo "🧪 Testing Artifactory Integration with Hugging Face Models..."

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

# Check if Python is available
if ! command -v python3 &> /dev/null; then
    print_error "Python3 is not installed"
    exit 1
fi

# Check if required packages are installed
print_status "Checking required packages..."
python3 -c "import transformers, torch, huggingface_hub" 2>/dev/null || {
    print_error "Required packages not found. Please install: transformers torch huggingface_hub"
    exit 1
}

# Set environment variables
export HF_HUB_ETAG_TIMEOUT=86400
export HF_HUB_DOWNLOAD_TIMEOUT=86400
export HF_ENDPOINT=https://apptrustswampupc.jfrog.io/artifactory/api/huggingfaceml/dev-huggingfaceml-remote
export HF_TOKEN=
export TRANSFORMERS_CACHE=./test_models
export HF_HOME=./test_models

# Create test directory
mkdir -p ./test_models

print_status "Testing model download from Artifactory..."

# Test script
python3 -c "
import os
import logging
from transformers import MarianTokenizer, MarianMTModel
import huggingface_hub

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Configure Hugging Face Hub to use Artifactory
if os.getenv('HF_ENDPOINT'):
    print(f'Using Artifactory endpoint: {os.getenv(\"HF_ENDPOINT\")}')
    huggingface_hub.set_http_backend(os.getenv('HF_ENDPOINT'))

model_name = 'Helsinki-NLP/opus-mt-en-fr'
cache_dir = os.getenv('TRANSFORMERS_CACHE', './test_models')
token = os.getenv('HF_TOKEN')

print(f'Downloading model: {model_name}')
print(f'Cache directory: {cache_dir}')

try:
    # Load tokenizer
    print('Loading tokenizer...')
    tokenizer = MarianTokenizer.from_pretrained(
        model_name, 
        cache_dir=cache_dir,
        token=token,
        trust_remote_code=True
    )
    print('✅ Tokenizer loaded successfully')
    
    # Load model
    print('Loading model...')
    model = MarianMTModel.from_pretrained(
        model_name, 
        cache_dir=cache_dir,
        token=token,
        trust_remote_code=True
    )
    print('✅ Model loaded successfully')
    
    # Test translation
    print('Testing translation...')
    text = 'Hello, how are you?'
    inputs = tokenizer(text, return_tensors='pt', padding=True, truncation=True, max_length=512)
    
    with torch.no_grad():
        translated = model.generate(**inputs)
    
    translated_text = tokenizer.decode(translated[0], skip_special_tokens=True)
    print(f'Original: {text}')
    print(f'Translated: {translated_text}')
    print('✅ Translation test successful')
    
except Exception as e:
    print(f'❌ Error: {str(e)}')
    exit(1)
"

if [ $? -eq 0 ]; then
    print_status "✅ Artifactory integration test completed successfully!"
    print_status "Model downloaded and cached in: ./test_models"
else
    print_error "❌ Artifactory integration test failed"
    exit 1
fi

# Cleanup
print_status "Cleaning up test files..."
rm -rf ./test_models

echo "🎉 Artifactory integration test completed!"
