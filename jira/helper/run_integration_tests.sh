#!/bin/bash

# Integration Test Runner for JIRA Helper
# This script runs integration tests against a real JIRA instance

set -e

echo "=== JIRA Helper Integration Tests ==="
echo ""

# Check if .env file exists
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    set -a
    source .env
    set +a
else
    echo "No .env file found. Using environment variables."
fi

# Check required environment variables
REQUIRED_VARS=(
    "JIRA_API_TOKEN"
    "JIRA_URL"
    "JIRA_USERNAME"
)

MISSING_VARS=()
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        MISSING_VARS+=("$var")
    fi
done

if [ ${#MISSING_VARS[@]} -ne 0 ]; then
    echo "❌ ERROR: The following required environment variables are not set:"
    printf '   - %s\n' "${MISSING_VARS[@]}"
    echo ""
    echo "Please set these variables or create a .env file with:"
    echo ""
    echo "JIRA_API_TOKEN=your-jira-api-token"
    echo "JIRA_URL=https://your-domain.atlassian.net"
    echo "JIRA_USERNAME=your-email@example.com"
    echo ""
    echo "Optional test variables:"
    echo "TEST_EXISTING_JIRA_ID=PROJ-123  # An existing JIRA ticket ID for testing"
    echo "TEST_COMMIT_WITH_JIRA=abc123    # A commit hash that contains JIRA IDs"
    echo "TEST_PERFORMANCE=true           # Enable performance tests"
    echo ""
    exit 1
fi

# Display configuration
echo "Configuration:"
echo "  JIRA_URL: $JIRA_URL"
echo "  JIRA_USERNAME: $JIRA_USERNAME"
echo "  JIRA_API_TOKEN: ****${JIRA_API_TOKEN: -4}"
echo ""

# Optional: Set test-specific variables if not already set
if [ -z "$TEST_EXISTING_JIRA_ID" ]; then
    echo "ℹ️  TEST_EXISTING_JIRA_ID not set. Some tests will be skipped."
fi

# Run integration tests
echo "Running integration tests..."
echo ""

# Use -v for verbose output, -tags=integration to run only integration tests
go test -v -tags=integration ./...

# Check exit code
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All integration tests passed!"
else
    echo ""
    echo "❌ Some integration tests failed!"
    exit 1
fi
