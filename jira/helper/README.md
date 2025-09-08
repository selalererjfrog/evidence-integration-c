# JIRA Helper

A Go tool that extracts JIRA ticket IDs from git commits and fetches their details from JIRA API.

## Quick Start

```bash
# Build
go build -o main .

# Extract JIRA IDs from a commit and fetch details
./main <commit_hash>

# Process specific JIRA tickets
./main EV-123 EV-456

# Extract IDs only (no JIRA API calls)
./main --extract-only <commit_hash>

# Process commit range
./main --range <start_commit>

# Generate markdown from JSON output
./main --markdown
```

## Prerequisites

- Go 1.21+
- Git repository (for commit extraction)
- JIRA Cloud API access

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `JIRA_API_TOKEN` | JIRA API token | Yes¹ |
| `JIRA_URL` | JIRA instance URL | Yes¹ |
| `JIRA_USERNAME` | JIRA username (email) | Yes¹ |
| `JIRA_ID_REGEX` | Pattern for JIRA IDs | No (default: `[A-Z]+-[0-9]+`) |
| `OUTPUT_FILE` | Output file path | No (default: `transformed_jira_data.json`) |

¹ Only required when fetching JIRA details (not for `--extract-only` mode)

### Using .env Files

```bash
# Create .env file with template
make create-env

# Edit with your values
JIRA_API_TOKEN=your-actual-token
JIRA_URL=https://your-instance.atlassian.net
JIRA_USERNAME=your-email@example.com
```

## Usage Modes

### 1. Git-based Mode (Default)
Extracts JIRA IDs from git commits and fetches their details.

```bash
# Single commit
./main abc123def456

# Commit range (from commit to HEAD)
./main --range abc123def456
```

### 2. Direct JIRA Mode
Process specific JIRA tickets directly.

```bash
./main EV-123 EV-456 EV-789
```

### 3. Extract Only Mode
Extract JIRA IDs without fetching details (useful for debugging).

```bash
./main --extract-only abc123def456
```

### 4. Markdown Generation Mode
Generate a markdown report from the JSON output file.

```bash
# Generate markdown from default JSON file (transformed_jira_data.json)
./main --markdown

# Generate markdown with custom output file
./main --markdown --markdown-output report.md

# Use different input JSON file
./main --markdown -o custom_data.json --markdown-output custom_report.md
```

## Command Line Options

- `-r, --regex PATTERN` - JIRA ID regex pattern
- `-o, --output FILE` - Output file path
- `--extract-only` - Only extract IDs, don't fetch from JIRA
- `--range` - Process commit range instead of single commit
- `--markdown` - Generate markdown from existing JSON file
- `--markdown-output FILE` - Output file for markdown (default: transformed_jira_data.md)
- `-h, --help` - Show help

## Output Format

The tool outputs a JSON file with JIRA ticket details and their transition history:

```json
{
  "tasks": [
    {
      "key": "EV-123",
      "link": "https://example.atlassian.net/browse/EV-123",
      "status": "In Progress",
      "description": "Task description",
      "type": "Task",
      "project": "EV",
      "created": "2020-01-01T12:11:56.063+0530",
      "updated": "2020-01-01T12:12:01.876+0530",
      "assignee": "John Doe",
      "reporter": "Jane Smith",
      "priority": "Medium",
      "transitions": [
        {
          "from_status": "To Do",
          "to_status": "In Progress",
          "author": "John Doe",
          "author_user_name": "john.doe@company.com",
          "transition_time": "2020-07-28T16:39:54.620+0530"
        }
      ]
    }
  ]
}
```

### Error Response

When a JIRA ticket cannot be fetched:

```json
{
  "key": "EV-123",
  "status": "Error",
  "description": "Error: Could not retrieve issue",
  "type": "Error"
}
```

### Markdown Output Format

The markdown generation feature creates a comprehensive report with:

- **Summary Table** - Overview of all tasks with key information
- **Task Details** - Complete information for each task including:
  - Basic information (status, type, project, priority)
  - People (assignee, reporter)
  - Dates (created, updated)
  - Description
  - Transition history
- **Status Distribution** - Summary of task counts by status
- **Clickable JIRA Links** - When JIRA URLs are included in the JSON data, ticket keys become clickable links

Example markdown output structure:
```markdown
# JIRA Tasks Report

Generated on: 2025-09-03 14:08:52

Total tasks: 1

## Summary

| Key | Status | Type | Priority | Assignee |
|-----|--------|------|----------|----------|
| OPS-12 | In Progress | Task | Medium | Sela Lerer |

## Task Details

### 1. OPS-12

**Basic Information:**
- **Status:** In Progress
- **Type:** Task
...
```

## Development

### Building

```bash
# Using Make
make build

# Manual build
go build -o main .
```

### Testing

```bash
# Unit tests
make test-unit

# Integration tests (requires JIRA credentials)
make test-integration

# All tests
make test

# Coverage report
make test-coverage
```

### Integration Test Setup

For integration tests, set these additional environment variables:

- `TEST_EXISTING_JIRA_ID` - A valid JIRA ticket ID in your instance (e.g., `OPS-3`)
- `TEST_COMMIT_WITH_JIRA` - A git commit hash containing JIRA IDs

Example `.env` for testing:
```bash
JIRA_API_TOKEN=your-token-here
JIRA_URL=https://your-instance.atlassian.net
JIRA_USERNAME=your-email@example.com
TEST_EXISTING_JIRA_ID=OPS-3
TEST_COMMIT_WITH_JIRA=d54597b0ea5e6e2d026c4611a8185a60b8d03e80
```

### Makefile Targets

- `make build` - Build the binary
- `make test` - Run all tests
- `make clean` - Remove build artifacts
- `make run-example` - Run with current commit
- `make show-jira-ids` - Show JIRA IDs in recent commits
- `make create-env` - Create sample .env file
- `make help` - Show all targets

## Project Structure

```
├── main.go              # Entry point
├── config.go            # Configuration and CLI parsing
├── modes.go             # Execution modes
├── git.go               # Git operations
├── jira_client.go       # JIRA API client
├── jira_models.go       # Data structures
├── jira_utils.go        # JIRA utilities
├── markdown_generator.go # Markdown generation
├── errors.go            # Error types
├── utils.go             # File I/O
└── *_test.go            # Test files
```

## Troubleshooting

### Common Issues

**Git Repository Not Found**
```
Error: not in a git repository
```
→ Run from a git repository

**JIRA Authentication Failed**
```
JIRA token not found
```
→ Set required environment variables

**Commit Not Found**
```
❌ commit 'abc123' not found
```
→ Verify commit exists in the repository

### Debug Commands

```bash
# Test without JIRA API
./main --extract-only HEAD

# Check recent JIRA IDs
make show-jira-ids

# Verbose test output
go test -v -run TestName
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Fetch JIRA details
  run: |
    cd jira/helper
    ./main "${{ github.sha }}"
```

## License

Part of the Evidence-Examples repository.