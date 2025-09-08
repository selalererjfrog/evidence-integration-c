"""
Translation Service Package

A REST API service for English to French and Hebrew translation using Hugging Face models.
"""

__version__ = "1.0.0"
__author__ = "JFrog Evidence Demo"
__email__ = "demo@jfrog.com"

from .simple_app import app

__all__ = ["app"]
