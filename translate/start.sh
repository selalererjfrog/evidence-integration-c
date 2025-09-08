#!/bin/bash

echo "🚀 Starting Translation Service..."
echo "📦 Installing dependencies..."

# Install dependencies
pip install -r requirements.txt

echo "🌐 Starting the service on http://localhost:8002"
echo "📚 API Documentation will be available at http://localhost:8002/docs"
echo "🏥 Health check at http://localhost:8002/health"
echo ""
echo "Press Ctrl+C to stop the service"
echo ""

# Start the service
python app.py
