# Build Optimization Guide

This document explains the optimizations implemented to reduce the AI Translate service build time.

## Optimizations Implemented

### 1. Multi-Stage Docker Build
- **Builder stage**: Installs build dependencies and pre-downloads models
- **Production stage**: Only includes runtime dependencies and cached models
- **Benefit**: Reduces final image size and build time

### 2. Model Caching
- Hugging Face models are pre-downloaded during build
- Models are cached in `/app/models` directory
- Environment variables set for model caching:
  - `TRANSFORMERS_CACHE="/app/models"`
  - `HF_HOME="/app/models"`
- **Benefit**: Eliminates model download time on container startup

### 3. Docker Layer Caching
- CI workflow uses registry-based layer caching
- Local builds use local layer caching
- **Benefit**: Reuses cached layers for faster rebuilds

### 4. Dependency Caching
- Python pip packages cached in CI
- Hugging Face models cached in CI
- **Benefit**: Avoids reinstalling dependencies on every build

### 5. Optimized Requirements
- Pinned dependency versions for better caching
- Separated development and production dependencies
- **Benefit**: More predictable builds and better cache hits

### 6. Build Context Optimization
- `.dockerignore` file excludes unnecessary files
- **Benefit**: Smaller build context, faster Docker builds

## Build Commands

### Local Development
```bash
# Build with optimizations
./build-optimized.sh

# Build and test
TEST_IMAGE=true ./build-optimized.sh
```

### CI/CD
The GitHub Actions workflow automatically uses all optimizations:
- Layer caching from registry
- Model and dependency caching
- Multi-stage builds

## Expected Performance Improvements

- **First build**: ~10-15 minutes (model download)
- **Subsequent builds**: ~2-3 minutes (cached layers and models)
- **CI builds**: ~5-8 minutes (registry caching)

## Troubleshooting

### Cache Issues
If you experience cache-related issues:

1. **Clear local cache**:
   ```bash
   rm -rf /tmp/.buildx-cache
   docker system prune -f
   ```

2. **Clear model cache**:
   ```bash
   rm -rf ~/.cache/huggingface
   ```

3. **Rebuild without cache**:
   ```bash
   docker buildx build --no-cache .
   ```

### Model Loading Issues
If models fail to load:

1. Check if `/app/models` directory exists in container
2. Verify environment variables are set correctly
3. Check network connectivity for model downloads

## Monitoring Build Performance

To monitor build performance:

1. **Build time**: Check CI workflow execution time
2. **Cache hits**: Monitor cache usage in build logs
3. **Image size**: Compare before/after optimization

## Artifactory Integration

The service is configured to use Hugging Face models from JFrog Artifactory:

### Environment Variables
- `HF_HUB_ETAG_TIMEOUT=86400`: Timeout for ETag requests
- `HF_HUB_DOWNLOAD_TIMEOUT=86400`: Timeout for model downloads
- `HF_ENDPOINT`: Artifactory Hugging Face endpoint
- `HF_TOKEN`: Authentication token for Artifactory

### Testing
Run the Artifactory integration test:
```bash
./test-artifactory.sh
```

### Benefits
- **Faster downloads**: Models served from Artifactory instead of Hugging Face Hub
- **Better reliability**: Reduced dependency on external services
- **Enterprise security**: Models served through your Artifactory instance

## Future Optimizations

Consider these additional optimizations:

1. **Model quantization**: Use smaller, quantized models
2. **Distributed caching**: Use shared cache across runners
3. **Parallel builds**: Build multiple architectures in parallel
4. **Incremental builds**: Only rebuild changed components
5. **Model versioning**: Implement model version management in Artifactory
