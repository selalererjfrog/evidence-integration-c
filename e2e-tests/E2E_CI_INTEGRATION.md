# E2E Test Integration with CI Pipeline

## Overview

This document describes the end-to-end test integration that has been added to the `promote-to-qa.yml` GitHub Actions workflow.

## What Was Added

### 1. New E2E Test Job

A new job called `e2e-test` has been added to the workflow that:
- Runs after the promotion to QA is successful
- Pulls Docker images from the QA repository (`evidencetrial.jfrog.io/commons-qa-docker-local/`)
- Runs the services without rebuilding them
- Executes comprehensive E2E tests
- Generates detailed JSON test results

### 2. Key Features

#### Docker Image Retrieval
- Automatically finds the latest Docker images in the QA repository
- Uses Artifactory AQL to search for the most recent images
- Falls back to using the QA version if specific tags are not found

#### Service Orchestration
- Runs both quote service and translation service containers
- Uses health checks to ensure services are ready before running tests
- Configures proper environment variables for Artifactory integration

#### Test Execution
- Uses a custom CI test runner (`run-ci-tests.py`)
- Generates JSON test reports using `pytest-json-report`
- Provides detailed test summaries and failure information

#### Results Generation
- Creates comprehensive JSON test results including:
  - Test run metadata (timestamp, version, environment)
  - Service health status
  - Detailed test results with pass/fail information
- Uploads results as GitHub Actions artifacts
- Provides test summary in workflow logs

## Required Configuration

### GitHub Secrets
- `JF_ACCESS_TOKEN`: JFrog Artifactory access token
- `HF_TOKEN`: Hugging Face token for model access

### GitHub Variables
- `JF_URL`: JFrog Artifactory URL
- `JF_USER`: JFrog Artifactory username
- `HF_ENDPOINT`: Hugging Face endpoint for Artifactory

## Workflow Structure

```
promote-to-qa.yml
├── promote-to-qa (existing job)
│   ├── Setup JFrog CLI
│   ├── Get latest DEV version
│   └── Promote to QA
└── e2e-test (new job)
    ├── Checkout code
    ├── Setup Docker Buildx
    ├── Setup JFrog CLI
    ├── Configure JFrog Artifactory
    ├── Verify required secrets
    ├── Get Docker images from QA repository
    ├── Create E2E test Docker Compose file
    ├── Run E2E tests
    ├── Generate test results JSON
    ├── Upload test results
    └── Cleanup containers
```

## Test Results

The workflow generates two types of test results:

### 1. Comprehensive JSON Report (`e2e-test-results.json`)
```json
{
  "test_run": {
    "timestamp": "2025-08-25T10:30:00Z",
    "workflow": "promote-to-qa",
    "qa_version": "1.0.0-SNAPSHOT",
    "quote_service_image": "evidencetrial.jfrog.io/commons-qa-docker-local/quote-of-day-service:1.0.0-SNAPSHOT",
    "translation_service_image": "evidencetrial.jfrog.io/commons-qa-docker-local/ai-translate:1.0.0-SNAPSHOT",
    "environment": "qa",
    "status": "completed"
  },
  "services": {
    "quote_service": {
      "status": "running",
      "health": "UP"
    },
    "translation_service": {
      "status": "running",
      "health": "healthy"
    }
  },
  "test_results": {
    "summary": {
      "total": 7,
      "passed": 7,
      "failed": 0
    },
    "tests": [...]
  }
}
```

### 2. Pytest JSON Report (`test-results.json`)
Standard pytest JSON report with detailed test information.

## Artifacts

The workflow uploads the following artifacts:
- `e2e-test-results-$VERSION/e2e-test-results.json`
- `e2e-test-results-$VERSION/test-results.json`

These are retained for 30 days and can be downloaded from the GitHub Actions run page.

## Usage

The E2E tests run automatically when:
1. The `promote-to-qa.yml` workflow is triggered
2. The promotion to QA is successful
3. All required secrets and variables are configured

## Troubleshooting

### Common Issues

1. **Missing Secrets**: Ensure `HF_TOKEN` and `HF_ENDPOINT` are configured
2. **Image Not Found**: Check that Docker images exist in the QA repository
3. **Service Health Checks**: Verify that services have proper health endpoints
4. **Network Issues**: Ensure containers can communicate on the test network

### Debugging

- Check the workflow logs for detailed error messages
- Download the test results artifacts for detailed failure information
- Verify service logs in the Docker Compose output

## Future Enhancements

Potential improvements:
- Add performance metrics collection
- Include load testing scenarios
- Add test result notifications (Slack, email)
- Implement test result trend analysis
- Add parallel test execution for faster feedback
