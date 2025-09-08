# E2E Tests for Quoteofday Service and Quotopia UI

Comprehensive end-to-end tests for the quoteofday service and Quotopia UI using Cypress with automatic test report generation.

## Prerequisites

- Node.js and npm installed
- Quoteofday service running on `http://localhost:8001`
- Quotopia UI service running on `http://localhost:8080`

## Setup

1. Install dependencies:
   ```bash
   npm ci
   ```

2. Start the quoteofday service:
   ```bash
   cd ../quoteofday
   ./mvnw spring-boot:run
   ```

3. Start the Quotopia UI service (via Docker):
   ```bash
   cd ../quotopia-ui
   docker run -p 8080:80 quotopia-ui:latest
   ```

## Running Tests

### Quoteofday Service Tests Only
```bash
npm test
# or
./run-cypress.sh
```

### Comprehensive E2E Tests (Service + UI)
```bash
npm run test:e2e
# or
./run-cypress.sh e2e
```

### All Tests (Service + E2E UI)
```bash
npm run test:all
# or
./run-cypress.sh all
```

### Run tests with UI
```bash
npm run cypress:open
```

## 📊 Automatic Test Report Generation

The test suite includes automatic generation of comprehensive test reports in both Markdown and JSON formats.

### Generate Reports with Tests

#### All Tests with Report
```bash
npm run test:report
# or
./run-cypress.sh report
```

#### Service Tests with Report
```bash
npm run test:service:report
# or
./run-cypress.sh service-report
```

#### E2E Tests with Report
```bash
npm run test:e2e:report
# or
./run-cypress.sh e2e-report
```

#### Generate Report Only (from existing results)
```bash
npm run report
```

### Generated Reports

After running tests with report generation, two files are automatically created:

- **`test-report.md`** - Comprehensive markdown report with:
  - Executive summary with test statistics
  - Detailed test suite breakdown
  - Performance metrics and analysis
  - Test categories and coverage
  - Recommendations and next steps

- **`test-report.json`** - Structured JSON report with:
  - Metadata (test runner, browser, environment)
  - Summary statistics
  - Performance metrics
  - Test categorization
  - Detailed test results
  - Service information

### Report Features

#### 📈 Performance Metrics
- Fastest and slowest test identification
- Average test execution times
- Total suite performance analysis
- Service response time tracking

#### 🎯 Test Categorization
- **UI/UX Tests**: Interface and user experience validation
- **API Tests**: Service endpoints and functionality
- **Integration Tests**: Service-to-service communication

#### 📋 Detailed Analysis
- Test success/failure rates
- Service health status
- Component validation results
- Error handling verification

## 🚀 GitHub Actions Workflow

The E2E tests are integrated into a GitHub Actions workflow that:

### Workflow Features
- **Automatic Triggering**: Runs on pushes to `e2e-tests/**` or workflow dispatch
- **Docker Service Management**: Pulls latest images from JFrog Artifactory
- **Service Orchestration**: Starts both quoteofday and quotopia-ui services
- **Test Execution**: Runs comprehensive E2E tests
- **Report Generation**: Automatically generates test reports
- **Evidence Creation**: Creates JFrog AppTrust evidence from test results
- **Artifact Upload**: Uploads test reports as GitHub artifacts

### Workflow Steps
1. **Setup**: Checkout code, setup JFrog CLI and Node.js
2. **Docker Login**: Authenticate with JFrog Artifactory
3. **Image Discovery**: Get latest Docker image versions from Artifactory
4. **Service Deployment**: Pull and start Docker containers
5. **Test Execution**: Run E2E tests with automatic report generation
6. **Evidence Creation**: Create JFrog AppTrust evidence
7. **Artifact Upload**: Upload test reports
8. **Cleanup**: Stop and remove Docker containers

### Local Workflow Testing
Test the workflow locally before committing:
```bash
./test-workflow.sh
```

This script simulates the GitHub Actions workflow steps locally.

## Test Coverage

### Quoteofday Service Tests (`quote-service.cy.js`)
- Service health endpoint (`/actuator/health`)
- API health endpoint (`/api/quotes/health`)
- Today's quote endpoint (`/api/quotes/today`)

### E2E UI Tests (`end-to-end-quotopia.cy.js`)
- **Service Management**: Automatically starts services if not running
- **UI Structure**: Verifies all UI components are present and visible
- **Quote Loading**: Tests quote loading from API and display
- **Interactive Features**: Tests quote card click to refresh
- **Error Handling**: Tests graceful error handling when API is down
- **Date Selector**: Verifies date selector functionality
- **Network Integration**: Verifies API calls are made correctly
- **Content Validation**: Checks footer and other UI content

## Service Management

The E2E tests include automatic service management:
- **Service Health Check**: Verifies services are running before tests
- **Auto-start**: Starts quoteofday service if not running
- **Error Testing**: Temporarily stops services to test error handling
- **Service Restart**: Restarts services after error testing

## Configuration

- **Quote Service URL**: `http://localhost:8001`
- **UI Service URL**: `http://localhost:8081`
- **Test Files**: 
  - `cypress/e2e/quote-service.cy.js` (Service tests)
  - `cypress/e2e/end-to-end-quotopia.cy.js` (E2E UI tests)
- **Configuration**: `cypress.config.js`
- **Report Generator**: `generate-test-report.js`

## Test Data Attributes

The UI includes `data-testid` attributes for reliable testing:
- `quotopia-container`, `quotopia-header`, `quotopia-main`, `quotopia-footer`
- `quote-container`, `quote-card`, `quote-text`, `quote-author`, `quote-date`
- `date-selector-container`, `date-selector`, `date-label`
- `footer-text`, `quotopia-banner`

## Available Scripts

### Test Execution
- `npm test` - Service tests only
- `npm run test:e2e` - E2E tests only
- `npm run test:all` - All tests
- `npm run cypress:open` - Open Cypress UI

### Report Generation
- `npm run test:report` - All tests with report
- `npm run test:service:report` - Service tests with report
- `npm run test:e2e:report` - E2E tests with report
- `npm run report` - Generate report only

### Shell Scripts
- `./run-cypress.sh` - Service tests
- `./run-cypress.sh e2e` - E2E tests
- `./run-cypress.sh all` - All tests
- `./run-cypress.sh report` - All tests with report
- `./run-cypress.sh e2e-report` - E2E tests with report
- `./run-cypress.sh service-report` - Service tests with report
- `./start-ui-service.sh` - Start UI service
- `./test-workflow.sh` - Test workflow locally
- `./clean-reports.sh` - Clean up old reports

## Report Example

After running tests, you'll see output like:
```
✅ Test reports generated successfully!
📄 test-report.md - Markdown report
📄 test-report.json - JSON report
📊 Summary: 10/10 tests passed (100%)
```

The reports provide comprehensive analysis of test results, performance metrics, and actionable recommendations for improving test coverage and reliability.

## GitHub Actions Integration

### Workflow File
- **Location**: `.github/workflows/end2end-tests.yml`
- **Trigger**: Push to `e2e-tests/**` or manual dispatch
- **Environment**: Ubuntu latest with Docker support

### Required Variables
- `JF_URL`: JFrog Artifactory URL
- `JF_USER`: JFrog username
- `JF_ACCESS_TOKEN`: JFrog access token
- `DOCKER_REGISTRY`: Docker registry URL
- `JFROG_CLI_KEY_ALIAS`: JFrog CLI key alias
- `JFROG_CLI_SIGNING_KEY`: JFrog CLI signing key

### Evidence Integration
The workflow creates JFrog AppTrust evidence using:
- **Build Name**: `e2e-tests`
- **Build Number**: GitHub run number
- **Predicate**: `test-report.json`
- **Markdown**: `test-report.md`
- **Provider ID**: `cypress-e2e`
- **Predicate Type**: `https://cypress.io/evidence/e2e/v1`
