# Makefile for JIRA Helper

.PHONY: all build test test-unit test-integration test-coverage test-integration-script clean run-example show-jira-ids create-env help

# Default target
all: build

# Build the binary
build:
	go build -o main .

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "=== Running Unit Tests ==="
	go test -v ./...

# Run integration tests
test-integration:
	@echo "=== Running Integration Tests ==="
	@if [ -z "$$JIRA_API_TOKEN" ]; then \
		echo "❌ JIRA_API_TOKEN not set. Skipping integration tests."; \
		echo "   Set environment variables or create a .env file."; \
		exit 0; \
	fi
	go test -v -tags=integration ./...



# Run tests with coverage
test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests with the helper script
test-integration-script:
	./run_integration_tests.sh

# Clean build artifacts and test files
clean:
	rm -f main test_main
	rm -f transformed_jira_data.json transformed_jira_data.md
	rm -f coverage.out coverage.html coverage_report.html
	rm -f *.test

# Run the binary with example
run-example:
	@if [ ! -f main ]; then \
		echo "Building binary first..."; \
		make build; \
	fi
	@COMMIT=$$(git rev-parse HEAD 2>/dev/null || echo "HEAD"); \
	if [ -z "$$JIRA_API_TOKEN" ]; then \
		echo "Running in extract-only mode (no JIRA API token)"; \
		./main --extract-only $$COMMIT; \
	else \
		./main $$COMMIT; \
	fi

# Show available JIRA IDs in recent commits
show-jira-ids:
	@echo "=== JIRA IDs in recent commits ==="
	@git log --oneline -10 | grep -E '[A-Z]+-[0-9]+' || echo "No JIRA IDs found in last 10 commits"

# Create a sample .env file
create-env:
	@if [ -f .env ]; then \
		echo ".env file already exists"; \
	else \
		echo "# JIRA Configuration" > .env; \
		echo "JIRA_API_TOKEN=your-jira-api-token-here" >> .env; \
		echo "JIRA_URL=https://your-domain.atlassian.net" >> .env; \
		echo "JIRA_USERNAME=your-email@example.com" >> .env; \
		echo "" >> .env; \
		echo "# Optional Test Configuration" >> .env; \
		echo "TEST_EXISTING_JIRA_ID=TEST-123" >> .env; \
		echo "TEST_COMMIT_WITH_JIRA=abc123def456" >> .env; \
		echo "" >> .env; \
		echo "✅ Created .env file. Please update with your actual values."; \
	fi

# Display help
help:
	@echo "JIRA Helper Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all                   Build the binary (default)"
	@echo "  build                 Build the binary"
	@echo "  test                  Run all tests (unit + integration)"
	@echo "  test-unit             Run unit tests only"
	@echo "  test-integration      Run integration tests (requires JIRA credentials)"
	@echo "  test-integration-script Run integration tests using the shell script"
	@echo "  test-coverage         Run tests with coverage report"
	@echo "  clean                 Remove build artifacts and test files"
	@echo "  run-example           Run the binary with HEAD commit"
	@echo "  show-jira-ids         Show JIRA IDs in recent commits"
	@echo "  create-env            Create a sample .env file"
	@echo "  help                  Display this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  JIRA_API_TOKEN        JIRA API token (required for integration tests)"
	@echo "  JIRA_URL             JIRA instance URL"
	@echo "  JIRA_USERNAME        JIRA username (email)"
	@echo ""
	@echo "Examples:"
	@echo "  make build            # Build the binary"
	@echo "  make test-unit        # Run unit tests"
	@echo "  make test             # Run all tests"
	@echo "  make clean            # Clean up"
