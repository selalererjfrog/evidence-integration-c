#!/bin/bash
set -e

echo "🧪 Testing E2E Workflow Locally"
echo "=================================="

# Simulate the workflow steps
echo "1. 📦 Installing dependencies..."
npm ci

echo "2. 🚀 Starting services..."
# Start quoteofday service (assuming it's already running)
echo "   - Quoteofday service should be running on port 8001"
# Start UI service
echo "   - Starting Quotopia UI service..."
cd ../quotopia-ui
docker build -t quotopia-ui:test .
docker run -d --name quotopia-ui-test -p 8081:80 quotopia-ui:test
cd ../e2e-tests

echo "3. ⏳ Waiting for services to be ready..."
sleep 5

echo "4. 🧪 Running E2E tests with report generation..."
npm run test:report

echo "5. 🔍 Verifying test results..."
if [ ! -f "test-report.md" ]; then
    echo "❌ Markdown report not found"
    exit 1
fi

if [ ! -f "test-report.json" ]; then
    echo "❌ JSON report not found"
    exit 1
fi

# Check test results
TOTAL_TESTS=$(jq -r '.summary.totalTests' test-report.json)
PASSING_TESTS=$(jq -r '.summary.passing' test-report.json)
FAILING_TESTS=$(jq -r '.summary.failing' test-report.json)

echo "📊 Test Results:"
echo "- Total Tests: $TOTAL_TESTS"
echo "- Passing: $PASSING_TESTS"
echo "- Failing: $FAILING_TESTS"

if [ "$FAILING_TESTS" -gt 0 ]; then
    echo "❌ Some tests failed"
    exit 1
fi

echo "6. 📋 Displaying test summary..."
echo "📄 Markdown Report Preview:"
head -20 test-report.md

echo "📊 JSON Report Summary:"
jq '.summary' test-report.json

echo "7. 🧹 Cleaning up..."
docker stop quotopia-ui-test 2>/dev/null || true
docker rm quotopia-ui-test 2>/dev/null || true

echo "✅ Workflow test completed successfully!"
