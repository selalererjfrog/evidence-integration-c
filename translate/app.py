from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import logging
from translation_service import TranslationService

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Translation Service",
    description="A REST API service for English to French translation using Hugging Face models",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize translation service
translation_service = TranslationService()

class TranslationRequest(BaseModel):
    text: str
    source_lang: str = "en"
    target_lang: str = "fr"

class TranslationResponse(BaseModel):
    original_text: str
    translated_text: str
    source_lang: str
    target_lang: str
    confidence: Optional[float] = None

class BatchTranslationRequest(BaseModel):
    texts: List[str]
    source_lang: str = "en"
    target_lang: str = "fr"

class BatchTranslationResponse(BaseModel):
    translations: List[TranslationResponse]

@app.on_event("startup")
async def startup_event():
    """Initialize the translation model on startup"""
    logger.info("Starting translation service...")
    await translation_service.initialize()
    logger.info("Translation service started successfully!")

@app.get("/")
async def root():
    """Health check endpoint"""
    return {
        "message": "Translation Service is running",
        "version": "1.0.0",
        "model": "Helsinki-NLP/opus-mt-en-fr"
    }

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "model_loaded": translation_service.is_initialized()
    }

@app.post("/translate", response_model=TranslationResponse)
async def translate_text(request: TranslationRequest):
    """Translate a single text from English to French"""
    try:
        if not translation_service.is_initialized():
            raise HTTPException(status_code=503, detail="Translation service not ready")
        
        translated_text = await translation_service.translate(
            request.text, 
            request.source_lang, 
            request.target_lang
        )
        
        return TranslationResponse(
            original_text=request.text,
            translated_text=translated_text,
            source_lang=request.source_lang,
            target_lang=request.target_lang
        )
    except Exception as e:
        logger.error(f"Translation error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Translation failed: {str(e)}")

@app.post("/translate/batch", response_model=BatchTranslationResponse)
async def translate_batch(request: BatchTranslationRequest):
    """Translate multiple texts from English to French"""
    try:
        if not translation_service.is_initialized():
            raise HTTPException(status_code=503, detail="Translation service not ready")
        
        translations = []
        for text in request.texts:
            translated_text = await translation_service.translate(
                text, 
                request.source_lang, 
                request.target_lang
            )
            translations.append(TranslationResponse(
                original_text=text,
                translated_text=translated_text,
                source_lang=request.source_lang,
                target_lang=request.target_lang
            ))
        
        return BatchTranslationResponse(translations=translations)
    except Exception as e:
        logger.error(f"Batch translation error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Batch translation failed: {str(e)}")

@app.get("/languages")
async def get_supported_languages():
    """Get supported languages"""
    return {
        "supported_languages": {
            "source": ["en"],
            "target": ["fr"]
        },
        "model": "Helsinki-NLP/opus-mt-en-fr"
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8002)
