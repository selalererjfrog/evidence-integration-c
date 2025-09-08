#!/bin/bash

echo "🧹 Cleaning up old test reports..."

# Remove old report files
if [ -f "test-report.md" ]; then
    rm test-report.md
    echo "✅ Removed test-report.md"
fi

if [ -f "test-report.json" ]; then
    rm test-report.json
    echo "✅ Removed test-report.json"
fi

# Remove Cypress results directory
if [ -d "cypress/results" ]; then
    rm -rf cypress/results
    echo "✅ Removed cypress/results directory"
fi

echo "🎉 Cleanup completed!"
