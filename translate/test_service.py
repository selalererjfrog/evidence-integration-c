#!/usr/bin/env python3
"""
Test script for the translation service
"""

import asyncio
import aiohttp
import json
import time

async def test_translation_service():
    """Test the translation service endpoints"""
    
    base_url = "http://localhost:8002"
    
    async with aiohttp.ClientSession() as session:
        
        # Test 1: Health check
        print("🔍 Testing health check...")
        async with session.get(f"{base_url}/health") as response:
            if response.status == 200:
                health_data = await response.json()
                print(f"✅ Health check passed: {health_data}")
            else:
                print(f"❌ Health check failed: {response.status}")
                return
        
        # Test 2: Single translation
        print("\n🌐 Testing single translation...")
        translation_data = {
            "text": "Hello, how are you today?",
            "source_lang": "en",
            "target_lang": "fr"
        }
        
        async with session.post(
            f"{base_url}/translate",
            json=translation_data
        ) as response:
            if response.status == 200:
                result = await response.json()
                print(f"✅ Translation successful:")
                print(f"   Original: {result['original_text']}")
                print(f"   Translated: {result['translated_text']}")
            else:
                print(f"❌ Translation failed: {response.status}")
                error_text = await response.text()
                print(f"   Error: {error_text}")
        
        # Test 3: Batch translation
        print("\n📦 Testing batch translation...")
        batch_data = {
            "texts": [
                "The weather is beautiful today",
                "I love this application",
                "Thank you for your help"
            ],
            "source_lang": "en",
            "target_lang": "fr"
        }
        
        async with session.post(
            f"{base_url}/translate/batch",
            json=batch_data
        ) as response:
            if response.status == 200:
                result = await response.json()
                print(f"✅ Batch translation successful:")
                for i, translation in enumerate(result['translations']):
                    print(f"   {i+1}. '{translation['original_text']}' → '{translation['translated_text']}'")
            else:
                print(f"❌ Batch translation failed: {response.status}")
                error_text = await response.text()
                print(f"   Error: {error_text}")
        
        # Test 4: Get supported languages
        print("\n🌍 Testing languages endpoint...")
        async with session.get(f"{base_url}/languages") as response:
            if response.status == 200:
                languages = await response.json()
                print(f"✅ Languages: {languages}")
            else:
                print(f"❌ Languages endpoint failed: {response.status}")

def main():
    """Main function"""
    print("🚀 Starting translation service tests...")
    print("Make sure the service is running on http://localhost:8002")
    print("=" * 50)
    
    try:
        asyncio.run(test_translation_service())
    except aiohttp.ClientConnectorError:
        print("❌ Could not connect to the service. Make sure it's running on http://localhost:8002")
    except Exception as e:
        print(f"❌ Test failed with error: {e}")
    
    print("\n" + "=" * 50)
    print("🏁 Tests completed!")

if __name__ == "__main__":
    main()
