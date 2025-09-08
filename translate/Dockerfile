# Simple single-stage build for AMD64
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# No system dependencies needed - all Python packages have pre-compiled wheels

# Copy requirements first for better caching
COPY requirements-optimized.txt requirements.txt

# Install Python dependencies
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt && \
    pip cache purge

# Copy application code
COPY . .

# Accept build arguments for model caching and Artifactory
ARG HF_HUB_ETAG_TIMEOUT
ARG HF_HUB_DOWNLOAD_TIMEOUT
ARG HF_ENDPOINT
ARG HF_TOKEN

# Set environment variables for model caching and Artifactory
ENV TRANSFORMERS_CACHE="/app/models"
ENV HF_HOME="/app/models"
ENV HF_HUB_ETAG_TIMEOUT=${HF_HUB_ETAG_TIMEOUT}
ENV HF_HUB_DOWNLOAD_TIMEOUT=${HF_HUB_DOWNLOAD_TIMEOUT}
ENV HF_ENDPOINT=${HF_ENDPOINT}
ENV HF_TOKEN=${HF_TOKEN}

# Create model cache directory
RUN mkdir -p /app/models

# Create a non-root user
RUN useradd --create-home --shell /bin/bash app \
    && chown -R app:app /app
USER app

# Expose port
EXPOSE 8002

# Run the application
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8002"]
