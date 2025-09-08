# Disk Space Optimization Guide for AI Translate Service

## 🚨 Problem Analysis

The GitHub Actions build is failing with "No space left on device" error. This is caused by:

1. **Large dependencies**: `torch` (~2GB) and `transformers` (~500MB)
2. **Multi-platform builds**: Building for both AMD64 and ARM64 simultaneously
3. **No cleanup**: Accumulated cache and temporary files
4. **Development dependencies**: Installing unnecessary dev packages in CI

## 🔧 Optimizations Implemented

### 1. Split Build Process into Stages

**Before**: Single job building both platforms simultaneously
**After**: Separate jobs for testing, AMD64 build, ARM64 build, and publishing

```yaml
jobs:
  test:           # Lightweight testing only
  build-amd64:    # Build AMD64 image only
  build-arm64:    # Build ARM64 image only  
  publish:        # Create manifests and publish
```

**Benefits**:
- Reduces peak disk usage by ~50%
- Parallel execution possible
- Better error isolation

### 2. Optimized Dependencies

**Created `requirements-optimized.txt`** with production-only dependencies:

```txt
# Production dependencies only
fastapi==0.104.1
uvicorn[standard]==0.24.0
transformers==4.53.0
torch==2.8.0
numpy<2.0.0
sentencepiece
pydantic==2.5.0
python-multipart==0.0.18
aiohttp==3.12.14
```

**Removed from production**:
- `pytest` and `pytest-asyncio` (moved to dev dependencies)
- `httpx` (not used in production)

**Potential savings**: ~50-100MB

### 3. Aggressive Cleanup

**Added cleanup steps**:
```yaml
- name: Cleanup pip cache
  run: |
    pip cache purge
    rm -rf ~/.cache/pip

- name: Cleanup Docker
  run: |
    docker system prune -f
    docker builder prune -f
```

**Dockerfile optimizations**:
```dockerfile
RUN pip install --no-cache-dir -r requirements.txt && \
    pip cache purge

RUN apt-get update && apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean && \
    apt-get autoremove -y
```

### 4. Separate Cache Keys

**Platform-specific caching**:
```yaml
# AMD64 cache
--cache-from type=registry,ref=${{ env.DOCKER_REGISTRY }}/ai-translate:cache-amd64
--cache-to type=registry,ref=${{ env.DOCKER_REGISTRY }}/ai-translate:cache-amd64

# ARM64 cache  
--cache-from type=registry,ref=${{ env.DOCKER_REGISTRY }}/ai-translate:cache-arm64
--cache-to type=registry,ref=${{ env.DOCKER_REGISTRY }}/ai-translate:cache-arm64
```

## 📊 Expected Space Savings

| Component | Before | After | Savings |
|-----------|--------|-------|---------|
| Multi-platform build | ~4GB | ~2GB | 50% |
| Dependencies | ~2.5GB | ~2.4GB | 4% |
| Cache accumulation | ~1GB | ~0.1GB | 90% |
| **Total** | **~7.5GB** | **~4.5GB** | **40%** |

## 🛠️ Additional Optimization Options

### 1. Use CPU-only PyTorch

Replace `torch==2.8.0` with `torch==2.8.0+cpu`:
```txt
# In requirements-optimized.txt
torch==2.8.0+cpu
```

**Potential savings**: 300-500MB (removes CUDA dependencies)

### 2. Use Smaller Base Images

Consider using `python:3.11-alpine` instead of `python:3.11-slim`:
```dockerfile
FROM python:3.11-alpine AS builder
```

**Potential savings**: 100-200MB

### 3. Model Quantization

Use quantized models for smaller footprint:
```python
# In translation_service.py
model = MarianMTModel.from_pretrained(
    model_name, 
    cache_dir=cache_dir,
    token=token,
    trust_remote_code=True,
    torch_dtype=torch.float16  # Use half precision
)
```

**Potential savings**: 50-100MB

### 4. Conditional Model Loading

Only load models when needed:
```python
# Lazy loading
def _load_model_if_needed(self):
    if not self.model:
        self._load_model()
```

## 🔍 Monitoring and Analysis

### Dependency Analysis Script

Run the analysis script to identify optimization opportunities:
```bash
cd translate
python analyze-dependencies.py
```

This will:
- Show package sizes
- Identify large dependencies
- Suggest optimizations
- Calculate potential savings

### CI Monitoring

Add disk space monitoring to CI:
```yaml
- name: Check disk space
  run: |
    df -h
    du -sh /home/runner/*
```

## 🚀 Implementation Steps

1. **Immediate** (already implemented):
   - ✅ Split build jobs
   - ✅ Optimized requirements
   - ✅ Added cleanup steps

2. **Short-term** (recommended):
   - Use CPU-only PyTorch
   - Implement model quantization
   - Add disk space monitoring

3. **Long-term** (optional):
   - Alpine base image
   - Model compression
   - CDN for model distribution

## 📝 Best Practices

### For Future Dependencies

1. **Always specify versions** to avoid unexpected updates
2. **Use `--no-cache-dir`** in CI builds
3. **Separate dev and prod dependencies**
4. **Regular cleanup** of caches and temporary files
5. **Monitor package sizes** before adding new dependencies

### For CI/CD

1. **Split large jobs** into smaller, focused tasks
2. **Use platform-specific caching**
3. **Implement cleanup steps** after each stage
4. **Monitor disk usage** during builds
5. **Use parallel jobs** when possible

## 🔧 Troubleshooting

### If Build Still Fails

1. **Check disk space**:
   ```bash
   df -h
   ```

2. **Clean Docker**:
   ```bash
   docker system prune -a -f
   ```

3. **Reduce parallelism**:
   ```yaml
   strategy:
     matrix:
       platform: [linux/amd64]  # Build one platform at a time
   ```

4. **Use external runners** with more disk space

### Performance Monitoring

Track build times and success rates:
- Monitor CI job durations
- Track cache hit rates
- Measure disk usage patterns
- Analyze failure patterns

## 📈 Results Tracking

After implementing these optimizations:

1. **Monitor build success rate** (should improve from ~70% to >95%)
2. **Track build times** (should remain similar or improve)
3. **Measure disk usage** (should reduce by 40-50%)
4. **Monitor cache effectiveness** (should improve hit rates)

The optimized workflow should resolve the "No space left on device" errors while maintaining build performance and reliability.
