from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import logging
import asyncio
from transformers import MarianMTModel, MarianTokenizer

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Translation Service (Test Version)",
    description="A simple test REST API service for English to French translation",
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

# Global variables for the models
models = {}
tokenizers = {}
models_loaded = False

async def load_models():
    """Load the Hugging Face models asynchronously"""
    global models, tokenizers, models_loaded
    try:
        logger.info("Loading translation models...")
        
        # Run model loading in a thread pool to avoid blocking
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(None, _load_models_sync)
        
        models_loaded = True
        logger.info("All models loaded successfully!")
        
    except Exception as e:
        logger.error(f"Failed to load models: {str(e)}")
        raise

def _load_models_sync():
    """Load the models synchronously"""
    global models, tokenizers
    
    # Define the models to load
    model_configs = {
        "fr": "Helsinki-NLP/opus-mt-en-fr",
        "he": "Helsinki-NLP/opus-mt-en-he"
    }
    
    for lang_code, model_name in model_configs.items():
        logger.info(f"Loading {model_name} for {lang_code}...")
        
        # Load tokenizer and model
        tokenizers[lang_code] = MarianTokenizer.from_pretrained(model_name)
        models[lang_code] = MarianMTModel.from_pretrained(model_name)
        
        # Set model to evaluation mode
        models[lang_code].eval()
        
        logger.info(f"Model for {lang_code} loaded successfully!")

async def translate_text_ai(text: str, source_lang: str = "en", target_lang: str = "fr") -> str:
    """Real AI translation using Helsinki-NLP models"""
    if not models_loaded:
        raise RuntimeError("Models not loaded")
    
    if source_lang != "en":
        raise ValueError("Currently only English source language is supported")
    
    if target_lang not in ["fr", "he"]:
        raise ValueError("Currently only French (fr) and Hebrew (he) target languages are supported")
    
    try:
        # Run translation in a thread pool to avoid blocking
        loop = asyncio.get_event_loop()
        translated_text = await loop.run_in_executor(None, _translate_sync, text, target_lang)
        return translated_text
        
    except Exception as e:
        logger.error(f"Translation error: {str(e)}")
        raise

def _translate_sync(text: str, target_lang: str) -> str:
    """Perform the actual translation synchronously"""
    try:
        # Get the appropriate model and tokenizer for the target language
        model = models[target_lang]
        tokenizer = tokenizers[target_lang]
        
        # Prepare the input for the model using the method you specified
        tokenized_text = tokenizer.prepare_seq2seq_batch(text, return_tensors="pt")
        
        # Generate the translation using the model with max_new_tokens to avoid deprecation warning
        translated_tokens = model.generate(**tokenized_text, max_new_tokens=512)
        
        # Decode the translated tokens back into human-readable text
        translated_text = tokenizer.batch_decode(translated_tokens, skip_special_tokens=True)[0]
        
        return translated_text
        
    except Exception as e:
        logger.error(f"Error during translation: {str(e)}")
        raise

@app.on_event("startup")
async def startup_event():
    """Initialize the translation models on startup"""
    await load_models()

@app.get("/")
async def root():
    """Health check endpoint"""
    return {
        "message": "Translation Service with Real AI Models is running",
        "version": "1.0.0",
        "models": ["Helsinki-NLP/opus-mt-en-fr", "Helsinki-NLP/opus-mt-en-he"]
    }

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "models_loaded": models_loaded,
        "models": ["Helsinki-NLP/opus-mt-en-fr", "Helsinki-NLP/opus-mt-en-he"]
    }

@app.post("/translate", response_model=TranslationResponse)
async def translate_text(request: TranslationRequest):
    """Translate a single text from English to French or Hebrew"""
    try:
        if not models_loaded:
            raise HTTPException(status_code=503, detail="Translation service not ready")
        
        translated_text = await translate_text_ai(
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

@app.get("/translate")
async def translate_text_get(
    text: str,
    source_lang: str = "en",
    target_lang: str = "fr"
):
    """Translate text using GET request with query parameters"""
    try:
        if not text:
            raise HTTPException(status_code=400, detail="Text parameter is required")
        
        if not models_loaded:
            raise HTTPException(status_code=503, detail="Translation service not ready")
        
        translated_text = await translate_text_ai(text, source_lang, target_lang)
        
        return TranslationResponse(
            original_text=text,
            translated_text=translated_text,
            source_lang=source_lang,
            target_lang=target_lang
        )
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"Translation error: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Translation failed: {str(e)}")

@app.post("/translate/batch", response_model=BatchTranslationResponse)
async def translate_batch(request: BatchTranslationRequest):
    """Translate multiple texts from English to French or Hebrew"""
    try:
        if not models_loaded:
            raise HTTPException(status_code=503, detail="Translation service not ready")
        
        translations = []
        for text in request.texts:
            translated_text = await translate_text_ai(
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
            "target": ["fr", "he"]
        },
        "models": ["Helsinki-NLP/opus-mt-en-fr", "Helsinki-NLP/opus-mt-en-he"]
    }

@app.get("/translate/quick")
async def quick_translate():
    """Quick translation with default text for testing"""
    if not models_loaded:
        raise HTTPException(status_code=503, detail="Translation service not ready")
    
    default_text = "Hello, how are you today?"
    translated_text = await translate_text_ai(default_text, "en", "fr")
    
    return {
        "message": "Quick translation test with real AI model",
        "original_text": default_text,
        "translated_text": translated_text,
        "source_lang": "en",
        "target_lang": "fr"
    }

@app.get("/translate/quick/hebrew")
async def quick_translate_hebrew():
    """Quick Hebrew translation with default text for testing"""
    if not models_loaded:
        raise HTTPException(status_code=503, detail="Translation service not ready")
    
    default_text = "Hello, how are you today?"
    translated_text = await translate_text_ai(default_text, "en", "he")
    
    return {
        "message": "Quick Hebrew translation test with real AI model",
        "original_text": default_text,
        "translated_text": translated_text,
        "source_lang": "en",
        "target_lang": "he"
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8002)
