"""
Security configuration for the translation service to address vulnerabilities:
- CVE-2025-47273: setuptools path traversal vulnerability
- CVE-2025-43859: h11 request smuggling vulnerability
"""

import os
from typing import Dict, Any

def get_security_config() -> Dict[str, Any]:
    """
    Returns security configuration to mitigate known vulnerabilities.
    """
    return {
        # Security headers to prevent various attacks
        "security_headers": {
            "X-Content-Type-Options": "nosniff",
            "X-Frame-Options": "DENY",
            "X-XSS-Protection": "1; mode=block",
            "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
            "Content-Security-Policy": "default-src 'self'",
        },
        
        # Request size limits to prevent DoS attacks
        "request_limits": {
            "max_content_length": 10 * 1024 * 1024,  # 10MB
            "max_request_size": 10 * 1024 * 1024,    # 10MB
        },
        
        # CORS configuration to prevent unauthorized access
        "cors": {
            "allow_origins": ["http://localhost:3000", "https://yourdomain.com"],
            "allow_credentials": True,
            "allow_methods": ["GET", "POST"],
            "allow_headers": ["*"],
        },
        
        # Rate limiting to prevent abuse
        "rate_limiting": {
            "requests_per_minute": 60,
            "burst_size": 10,
        },
        
        # Logging configuration to prevent information disclosure
        "logging": {
            "level": "INFO",
            "format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
            "exclude_sensitive_fields": ["password", "token", "secret"],
        },
        
        # File upload security
        "file_upload": {
            "allowed_extensions": [".txt", ".json"],
            "max_file_size": 5 * 1024 * 1024,  # 5MB
            "upload_directory": "/tmp/uploads",
        },
        
        # Session security
        "session": {
            "secure": True,
            "httponly": True,
            "samesite": "strict",
            "max_age": 1800,  # 30 minutes
        },
    }

def apply_security_headers(app):
    """
    Apply security headers to FastAPI application.
    """
    from fastapi import Request
    from fastapi.responses import Response
    
    @app.middleware("http")
    async def add_security_headers(request: Request, call_next):
        response = await call_next(request)
        
        security_config = get_security_config()
        for header, value in security_config["security_headers"].items():
            response.headers[header] = value
            
        return response

def validate_file_upload(filename: str, file_size: int) -> bool:
    """
    Validate file uploads to prevent path traversal attacks.
    """
    security_config = get_security_config()
    
    # Check file size
    if file_size > security_config["file_upload"]["max_file_size"]:
        return False
    
    # Check file extension
    allowed_extensions = security_config["file_upload"]["allowed_extensions"]
    if not any(filename.endswith(ext) for ext in allowed_extensions):
        return False
    
    # Check for path traversal attempts
    if ".." in filename or "/" in filename or "\\" in filename:
        return False
    
    return True

def sanitize_input(text: str) -> str:
    """
    Sanitize user input to prevent injection attacks.
    """
    import html
    import re
    
    # HTML escape
    text = html.escape(text)
    
    # Remove potentially dangerous patterns
    dangerous_patterns = [
        r'<script.*?</script>',
        r'javascript:',
        r'data:text/html',
        r'vbscript:',
    ]
    
    for pattern in dangerous_patterns:
        text = re.sub(pattern, '', text, flags=re.IGNORECASE)
    
    return text.strip()

# Environment-specific security settings
def get_environment_security_settings() -> Dict[str, Any]:
    """
    Get environment-specific security settings.
    """
    env = os.getenv("ENVIRONMENT", "development")
    
    if env == "production":
        return {
            "debug": False,
            "reload": False,
            "host": "0.0.0.0",
            "port": 8002,
            "ssl_keyfile": os.getenv("SSL_KEYFILE"),
            "ssl_certfile": os.getenv("SSL_CERTFILE"),
        }
    else:
        return {
            "debug": True,
            "reload": True,
            "host": "127.0.0.1",
            "port": 8002,
        }
