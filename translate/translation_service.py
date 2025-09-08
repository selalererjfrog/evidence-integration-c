import asyncio
import logging
import os
from transformers import MarianMTModel, MarianTokenizer
import torch
from typing import Optional

logger = logging.getLogger(__name__)

class TranslationService:
    def __init__(self):
        self.model: Optional[MarianMTModel] = None
        self.tokenizer: Optional[MarianTokenizer] = None
        self.model_name = "Helsinki-NLP/opus-mt-en-fr"
        self._initialized = False
        
    async def initialize(self):
        """Initialize the translation model asynchronously"""
        try:
            logger.info(f"Loading model: {self.model_name}")
            
            # Run model loading in a thread pool to avoid blocking
            loop = asyncio.get_event_loop()
            await loop.run_in_executor(None, self._load_model)
            
            self._initialized = True
            logger.info("Model loaded successfully!")
            
        except Exception as e:
            logger.error(f"Failed to load model: {str(e)}")
            raise
    
    def _load_model(self):
        """Load the Hugging Face model and tokenizer from Artifactory"""
        try:
            # Load tokenizer and model with caching and Artifactory configuration
            cache_dir = "/app/models" if os.path.exists("/app/models") else None
            
            # Configure Hugging Face Hub to use Artifactory
            import huggingface_hub
            hf_endpoint = os.getenv('HF_ENDPOINT')
            hf_token = os.getenv('HF_TOKEN')
            
            if hf_endpoint and hf_token:
                try:
                    # Set the endpoint using the newer API
                    if hasattr(huggingface_hub, 'set_http_backend'):
                        huggingface_hub.set_http_backend(hf_endpoint)
                        logger.info(f"Configured HF endpoint: {hf_endpoint}")
                    elif hasattr(huggingface_hub, 'HfApi'):
                        # Alternative approach using HfApi
                        api = huggingface_hub.HfApi(endpoint=hf_endpoint)
                        logger.info(f"Configured HF endpoint via HfApi: {hf_endpoint}")
                    else:
                        logger.warning("Could not configure HF endpoint - API not available")
                except Exception as e:
                    logger.warning(f"Could not configure HF endpoint: {e}")
            else:
                logger.info("Using public Hugging Face Hub")
            
            # Load with token if available
            token = os.getenv('HF_TOKEN')
            if token:
                self.tokenizer = MarianTokenizer.from_pretrained(
                    self.model_name, 
                    cache_dir=cache_dir,
                    token=token,
                    trust_remote_code=True
                )
                self.model = MarianMTModel.from_pretrained(
                    self.model_name, 
                    cache_dir=cache_dir,
                    token=token,
                    trust_remote_code=True
                )
            else:
                self.tokenizer = MarianTokenizer.from_pretrained(
                    self.model_name, 
                    cache_dir=cache_dir,
                    trust_remote_code=True
                )
                self.model = MarianMTModel.from_pretrained(
                    self.model_name, 
                    cache_dir=cache_dir,
                    trust_remote_code=True
                )
            
            # Set model to evaluation mode
            self.model.eval()
            
            # Move to GPU if available
            if torch.cuda.is_available():
                self.model = self.model.to('cuda')
                logger.info("Model moved to GPU")
            else:
                logger.info("Using CPU for inference")
                
        except Exception as e:
            logger.error(f"Error loading model: {str(e)}")
            raise
    
    def is_initialized(self) -> bool:
        """Check if the model is initialized"""
        return self._initialized and self.model is not None and self.tokenizer is not None
    
    async def translate(self, text: str, source_lang: str = "en", target_lang: str = "fr") -> str:
        """Translate text from source language to target language"""
        if not self.is_initialized():
            raise RuntimeError("Translation service not initialized")
        
        if source_lang != "en" or target_lang != "fr":
            raise ValueError("Currently only English to French translation is supported")
        
        try:
            # Run translation in a thread pool to avoid blocking
            loop = asyncio.get_event_loop()
            translated_text = await loop.run_in_executor(None, self._translate_text, text)
            return translated_text
            
        except Exception as e:
            logger.error(f"Translation error: {str(e)}")
            raise
    
    def _translate_text(self, text: str) -> str:
        """Perform the actual translation"""
        try:
            # Tokenize the input text
            inputs = self.tokenizer(text, return_tensors="pt", padding=True, truncation=True, max_length=512)
            
            # Move inputs to the same device as the model
            if torch.cuda.is_available():
                inputs = {k: v.to('cuda') for k, v in inputs.items()}
            
            # Generate translation
            with torch.no_grad():
                translated = self.model.generate(**inputs)
            
            # Decode the translation
            translated_text = self.tokenizer.decode(translated[0], skip_special_tokens=True)
            
            return translated_text
            
        except Exception as e:
            logger.error(f"Error during translation: {str(e)}")
            raise
    
    async def translate_batch(self, texts: list, source_lang: str = "en", target_lang: str = "fr") -> list:
        """Translate multiple texts"""
        if not self.is_initialized():
            raise RuntimeError("Translation service not initialized")
        
        if source_lang != "en" or target_lang != "fr":
            raise ValueError("Currently only English to French translation is supported")
        
        try:
            # Run batch translation in a thread pool
            loop = asyncio.get_event_loop()
            translated_texts = await loop.run_in_executor(None, self._translate_batch_texts, texts)
            return translated_texts
            
        except Exception as e:
            logger.error(f"Batch translation error: {str(e)}")
            raise
    
    def _translate_batch_texts(self, texts: list) -> list:
        """Perform batch translation"""
        try:
            # Tokenize all texts
            inputs = self.tokenizer(texts, return_tensors="pt", padding=True, truncation=True, max_length=512)
            
            # Move inputs to the same device as the model
            if torch.cuda.is_available():
                inputs = {k: v.to('cuda') for k, v in inputs.items()}
            
            # Generate translations
            with torch.no_grad():
                translated = self.model.generate(**inputs)
            
            # Decode all translations
            translated_texts = [self.tokenizer.decode(t, skip_special_tokens=True) for t in translated]
            
            return translated_texts
            
        except Exception as e:
            logger.error(f"Error during batch translation: {str(e)}")
            raise
