# Translation Service

A Python-based REST API service for English to French translation using the Helsinki-NLP/opus-mt-en-fr model from Hugging Face.

## Features

- **FastAPI-based REST API** with automatic OpenAPI documentation
- **Local AI Model**: Uses the Helsinki-NLP/opus-mt-en-fr model for high-quality translations
- **Async Processing**: Non-blocking translation operations
- **Batch Translation**: Support for translating multiple texts at once
- **Health Monitoring**: Built-in health check endpoints
- **Docker Support**: Containerized deployment
- **GPU Support**: Automatic GPU acceleration when available

## Model Information

This service uses the [Helsinki-NLP/opus-mt-en-fr](https://huggingface.co/Helsinki-NLP/opus-mt-en-fr) model, which is:
- A MarianMT transformer model
- Trained on OPUS datasets
- Optimized for English to French translation
- Provides high-quality translations with BLEU scores around 30-40

## API Endpoints

### Health Check
- `GET /` - Service status
- `GET /health` - Detailed health information

### Translation
- `POST /translate` - Translate single text
- `POST /translate/batch` - Translate multiple texts
- `GET /languages` - Get supported languages

## Quick Start

### Local Development

1. **Install dependencies:**
   ```bash
   cd translate
   pip install -r requirements.txt
   ```

2. **Run the service:**
   ```bash
   python app.py
   ```

3. **Access the API:**
   - API: http://localhost:8002
- Documentation: http://localhost:8002/docs
- Health check: http://localhost:8002/health

### Docker Deployment

1. **Build the image:**
   ```bash
   docker build -t translation-service .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8002:8002 translation-service
   ```

## API Usage Examples

### Single Translation

```bash
curl -X POST "http://localhost:8002/translate" \
     -H "Content-Type: application/json" \
     -d '{
       "text": "Hello, how are you?",
       "source_lang": "en",
       "target_lang": "fr"
     }'
```

**Response:**
```json
{
  "original_text": "Hello, how are you?",
  "translated_text": "Bonjour, comment allez-vous?",
  "source_lang": "en",
  "target_lang": "fr",
  "confidence": null
}
```

### Batch Translation

```bash
curl -X POST "http://localhost:8002/translate/batch" \
     -H "Content-Type: application/json" \
     -d '{
       "texts": [
         "Hello, how are you?",
         "The weather is nice today",
         "I love this application"
       ],
       "source_lang": "en",
       "target_lang": "fr"
     }'
```

### Health Check

```bash
curl http://localhost:8002/health
```

**Response:**
```json
{
  "status": "healthy",
  "model_loaded": true
}
```

## Configuration

### Environment Variables

- `MODEL_NAME`: Hugging Face model name (default: "Helsinki-NLP/opus-mt-en-fr")
- `PORT`: Service port (default: 8002)
- `HOST`: Service host (default: "0.0.0.0")

### Model Loading

The model is automatically downloaded from Hugging Face on first startup. The model files are cached locally for subsequent runs.

## Performance

- **First Request**: May take 1-2 seconds due to model loading
- **Subsequent Requests**: Typically 100-500ms per translation
- **GPU Acceleration**: Significantly faster with CUDA-enabled GPU
- **Batch Processing**: More efficient for multiple translations

## Development

### Project Structure

```
translate/
├── app.py                 # FastAPI application
├── translation_service.py # Translation logic
├── requirements.txt       # Python dependencies
├── Dockerfile            # Container configuration
└── README.md             # This file
```

### Adding New Languages

To support additional language pairs:

1. Update the model name in `translation_service.py`
2. Modify the language validation logic
3. Update the `/languages` endpoint response

### Testing

```bash
# Test the service
curl http://localhost:8002/health

# Test translation
curl -X POST "http://localhost:8002/translate" \
     -H "Content-Type: application/json" \
     -d '{"text": "Hello world"}'
```

## Troubleshooting

### Common Issues

1. **Model Download Fails**: Check internet connection and Hugging Face access
2. **Out of Memory**: Reduce batch size or use smaller model
3. **Slow Performance**: Ensure GPU drivers are installed for CUDA support

### Logs

The service provides detailed logging for debugging:
- Model loading status
- Translation requests and responses
- Error details

## License

This project uses the Helsinki-NLP/opus-mt-en-fr model which is licensed under Apache 2.0.
