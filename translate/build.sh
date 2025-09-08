#!/bin/bash

echo "🚀 Building Translation Service Python Package..."

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf build/ dist/ *.egg-info/

# Install build dependencies
echo "📦 Installing build dependencies..."
pip install build twine

# Run tests
echo "🧪 Running tests..."
python -m pytest tests/ -v

# Build the package
echo "🔨 Building package..."
python -m build

# Check the package
echo "✅ Checking package..."
twine check dist/*

echo "🎉 Build completed successfully!"
echo "📦 Package files created in dist/ directory"
echo "📋 To install locally: pip install dist/*.whl"
echo "📋 To upload to PyPI: twine upload dist/*"
