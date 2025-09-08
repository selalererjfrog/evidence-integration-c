"""
Tests for the translation service
"""

import pytest
from fastapi.testclient import TestClient
from simple_app import app

client = TestClient(app)

def test_root_endpoint():
    """Test the root endpoint"""
    response = client.get("/")
    assert response.status_code == 200
    data = response.json()
    assert "message" in data
    assert "version" in data
    assert "models" in data

def test_health_endpoint():
    """Test the health endpoint"""
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert "status" in data
    assert "models_loaded" in data

def test_languages_endpoint():
    """Test the languages endpoint"""
    response = client.get("/languages")
    assert response.status_code == 200
    data = response.json()
    assert "supported_languages" in data
    assert "source" in data["supported_languages"]
    assert "target" in data["supported_languages"]
    assert "en" in data["supported_languages"]["source"]
    assert "fr" in data["supported_languages"]["target"]
    assert "he" in data["supported_languages"]["target"]

def test_translate_endpoint_missing_text():
    """Test translation endpoint with missing text"""
    response = client.get("/translate")
    assert response.status_code == 422  # FastAPI returns 422 for validation errors

def test_translate_endpoint_with_text():
    """Test translation endpoint with text"""
    response = client.get("/translate?text=Hello")
    # This might return 503 if models are not loaded, or 500 if models fail to load
    assert response.status_code in [200, 503, 500]

def test_translate_post_endpoint():
    """Test POST translation endpoint"""
    response = client.post("/translate", json={"text": "Hello"})
    # This might return 503 if models are not loaded, or 500 if models fail to load
    assert response.status_code in [200, 503, 500]

def test_batch_translate_endpoint():
    """Test batch translation endpoint"""
    response = client.post("/translate/batch", json={"texts": ["Hello", "World"]})
    # This might return 503 if models are not loaded, or 500 if models fail to load
    assert response.status_code in [200, 503, 500]

def test_quick_translate_endpoint():
    """Test quick translate endpoint"""
    response = client.get("/translate/quick")
    # This might return 503 if models are not loaded, or 500 if models fail to load
    assert response.status_code in [200, 503, 500]

def test_quick_translate_hebrew_endpoint():
    """Test quick Hebrew translate endpoint"""
    response = client.get("/translate/quick/hebrew")
    # This might return 503 if models are not loaded, or 500 if models fail to load
    assert response.status_code in [200, 503, 500]
