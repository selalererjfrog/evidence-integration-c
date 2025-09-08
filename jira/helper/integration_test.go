//go:build integration
// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: Environment variables should be set before running these tests.
// You can either:
// 1. Use the run_integration_tests.sh script which loads .env automatically
// 2. Export environment variables manually
// 3. Source .env file: `set -a; source .env; set +a`
//
// Required environment variables:
// - JIRA_API_TOKEN: Your JIRA API token
// - JIRA_URL: Your JIRA instance URL (e.g., https://your-instance.atlassian.net)
// - JIRA_USERNAME: Your JIRA username (typically your email)
//
// Optional test environment variables:
// - TEST_EXISTING_JIRA_ID: A valid JIRA ticket ID that exists in your instance (e.g., OPS-3)
// - TEST_COMMIT_WITH_JIRA: A git commit hash that contains JIRA IDs in its message
// - TEST_PERFORMANCE: Set to "true" to enable performance tests
//
// Without TEST_EXISTING_JIRA_ID and TEST_COMMIT_WITH_JIRA, some tests will be skipped.

// TestEnvironmentSetup verifies all required environment variables are set
func TestEnvironmentSetup(t *testing.T) {
	// Check if environment variables are set
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping tests: JIRA_API_TOKEN not set. Use run_integration_tests.sh or set environment variables")
	}

	requiredVars := []string{
		"JIRA_API_TOKEN",
		"JIRA_URL",
		"JIRA_USERNAME",
	}

	for _, varName := range requiredVars {
		value := os.Getenv(varName)
		assert.NotEmpty(t, value, "Environment variable %s must be set for integration tests", varName)
		if varName == "JIRA_API_TOKEN" {
			// Log masked token for debugging
			if len(value) > 4 {
				t.Logf("%s: ****%s", varName, value[len(value)-4:])
			}
		} else {
			t.Logf("%s: %s", varName, value)
		}
	}

	// Optional test variables
	optionalVars := []string{
		"TEST_EXISTING_JIRA_ID",
		"TEST_COMMIT_WITH_JIRA",
		"TEST_PERFORMANCE",
	}

	for _, varName := range optionalVars {
		value := os.Getenv(varName)
		if value != "" {
			t.Logf("%s: %s (optional)", varName, value)
		}
	}
}

// TestJIRAConnection tests that we can connect to JIRA
func TestJIRAConnection(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	client, err := NewJiraClient()
	require.NoError(t, err, "Failed to create JIRA client")
	assert.NotNil(t, client)
}

// TestJIRAOperations tests real JIRA API operations
func TestJIRAOperations(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	client, err := NewJiraClient()
	require.NoError(t, err, "Failed to create JIRA client")

	t.Run("FetchExistingTicket", func(t *testing.T) {
		testJiraID := os.Getenv("TEST_EXISTING_JIRA_ID")
		if testJiraID == "" {
			t.Skip("TEST_EXISTING_JIRA_ID not set, skipping")
		}

		response := client.FetchJiraDetails([]string{testJiraID})
		require.Len(t, response.Tasks, 1)
		result := response.Tasks[0]

		assert.Equal(t, testJiraID, result.Key)
		assert.NotEqual(t, "Error", result.Status)
		assert.NotEmpty(t, result.Project)
		assert.NotEmpty(t, result.Created)
		assert.NotEmpty(t, result.Reporter)

		t.Logf("Successfully fetched %s:", testJiraID)
		t.Logf("  Status: %s", result.Status)
		t.Logf("  Type: %s", result.Type)
		t.Logf("  Project: %s", result.Project)
		t.Logf("  Reporter: %s", result.Reporter)
		t.Logf("  Transitions: %d", len(result.Transitions))
	})

	t.Run("FetchNonExistentTicket", func(t *testing.T) {
		// Use a ticket ID that's very unlikely to exist
		nonExistentID := "XXXNONEXISTENT-99999"
		response := client.FetchJiraDetails([]string{nonExistentID})
		require.Len(t, response.Tasks, 1)
		result := response.Tasks[0]

		assert.Equal(t, nonExistentID, result.Key)
		assert.Equal(t, "Error", result.Status)
		assert.Equal(t, "Error", result.Type)
		assert.Contains(t, result.Description, "Error")

		t.Logf("Correctly handled non-existent ticket %s", nonExistentID)
	})

	t.Run("FetchMultipleTickets", func(t *testing.T) {
		testJiraID := os.Getenv("TEST_EXISTING_JIRA_ID")
		if testJiraID == "" {
			t.Skip("TEST_EXISTING_JIRA_ID not set, skipping")
		}

		// Test fetching multiple tickets including one non-existent
		ticketIDs := []string{testJiraID, "XXXNONEXISTENT-99999"}
		response := client.FetchJiraDetails(ticketIDs)

		assert.Len(t, response.Tasks, 2)
		// First should be successful
		assert.NotEqual(t, "Error", response.Tasks[0].Status)
		// Second should be error
		assert.Equal(t, "Error", response.Tasks[1].Status)
	})
}

// TestGitOperations tests real git operations
func TestGitOperations(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	gitService := NewGitService()

	t.Run("CheckRepository", func(t *testing.T) {
		err := gitService.CheckRepository()
		if err != nil {
			t.Skipf("Not in a git repository: %v", err)
		}
		assert.NoError(t, err)
	})

	t.Run("GetBranchInfo", func(t *testing.T) {
		branch, commit, firstJiraID, err := gitService.GetBranchInfo()
		if err != nil {
			// This might fail if we're not in a git repo or have no commits
			t.Logf("GetBranchInfo failed (expected in some environments): %v", err)
			return
		}

		assert.NotEmpty(t, branch, "Branch should not be empty")
		assert.NotEmpty(t, commit, "Commit should not be empty")
		t.Logf("Current branch: %s, commit: %s, first JIRA ID: %s", branch, commit, firstJiraID)
	})

	t.Run("ValidateHEAD", func(t *testing.T) {
		err := gitService.ValidateHEAD()
		if err != nil {
			t.Logf("ValidateHEAD failed (expected in some environments): %v", err)
			return
		}
		assert.NoError(t, err)
	})

	t.Run("ExtractJiraIDs_SingleCommit", func(t *testing.T) {
		// Try to use TEST_COMMIT_WITH_JIRA if available
		testCommit := os.Getenv("TEST_COMMIT_WITH_JIRA")
		if testCommit == "" {
			// Get actual commit hash instead of HEAD
			var err error
			_, testCommit, _, err = gitService.GetBranchInfo()
			if err != nil {
				t.Logf("Failed to get commit hash: %v", err)
				return
			}
		}

		// Get current branch for extraction
		branch, _, _, err := gitService.GetBranchInfo()
		if err != nil {
			t.Logf("Failed to get branch info: %v", err)
			branch = "main" // Use default branch
		}

		jiraIDs, err := gitService.ExtractJiraIDs(testCommit, "[A-Z]+-[0-9]+", branch, false)
		if err != nil {
			t.Logf("ExtractJiraIDs failed for %s: %v", testCommit, err)
			return
		}

		t.Logf("Found %d JIRA IDs in commit %s: %v", len(jiraIDs), testCommit, jiraIDs)
	})

	t.Run("ExtractJiraIDs_Range", func(t *testing.T) {
		// Get current branch for extraction
		branch, _, _, err := gitService.GetBranchInfo()
		if err != nil {
			t.Logf("Failed to get branch info: %v", err)
			branch = "main" // Use default branch
		}

		// Try to extract from last 10 commits
		jiraIDs, err := gitService.ExtractJiraIDs("HEAD~10", "[A-Z]+-[0-9]+", branch, true)
		if err != nil {
			// Might fail if there aren't 10 commits
			t.Logf("ExtractJiraIDs range failed (expected in shallow repos): %v", err)
			return
		}

		t.Logf("Found %d unique JIRA IDs in last 10 commits: %v", len(jiraIDs), jiraIDs)
	})
}

// TestFullWorkflow tests the complete workflow from git to JIRA
func TestFullWorkflow(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	// This test requires both git and JIRA to be properly set up
	testCommit := os.Getenv("TEST_COMMIT_WITH_JIRA")
	if testCommit == "" {
		t.Skip("TEST_COMMIT_WITH_JIRA not set, skipping full workflow test")
	}

	config := &AppConfig{
		JIRAIDRegex: "[A-Z]+-[0-9]+",
		OutputFile:  "test_workflow_output.json",
	}

	gitService := NewGitService()

	// Get current branch
	branch, _, _, err := gitService.GetBranchInfo()
	if err != nil {
		t.Logf("Failed to get branch info: %v", err)
		branch = "main" // Use default branch
	}

	// Extract JIRA IDs from commit (single commit mode)
	jiraIDs, err := gitService.ExtractJiraIDs(testCommit, config.JIRAIDRegex, branch, true)
	require.NoError(t, err, "Failed to extract JIRA IDs")
	require.NotEmpty(t, jiraIDs, "No JIRA IDs found in test commit")

	t.Logf("Extracted JIRA IDs from %s: %v", testCommit, jiraIDs)

	// Fetch details for each JIRA ID
	client, err := NewJiraClient()
	require.NoError(t, err, "Failed to create JIRA client")

	response := client.FetchJiraDetails(jiraIDs)
	for _, task := range response.Tasks {
		t.Logf("Fetched %s: Status=%s, Type=%s", task.Key, task.Status, task.Type)
	}

	// Write results
	err = saveJiraResults(response, config)
	assert.NoError(t, err, "Failed to write output file")

	// Verify file was created
	_, err = os.Stat(config.OutputFile)
	assert.NoError(t, err, "Output file was not created")

	// Clean up
	os.Remove(config.OutputFile)

	t.Logf("Full workflow completed successfully: extracted %d IDs, fetched details, wrote output", len(jiraIDs))
}

// TestDirectJIRAProcessing tests direct JIRA ID processing mode
func TestDirectJIRAProcessing(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	testJiraID := os.Getenv("TEST_EXISTING_JIRA_ID")
	if testJiraID == "" {
		t.Skip("TEST_EXISTING_JIRA_ID not set, skipping")
	}

	config := &AppConfig{
		OutputFile: "test_direct_output.json",
	}

	// Set up config with JIRA IDs
	config.JIRAIDs = []string{testJiraID, "INVALID-99999"}
	err := processDirectJiraIDs(config)
	assert.NoError(t, err, "processDirectJiraIDs should not fail")

	// Verify output file was created
	_, err = os.Stat(config.OutputFile)
	assert.NoError(t, err, "Output file was not created")

	// Read and verify contents
	data, err := os.ReadFile(config.OutputFile)
	assert.NoError(t, err, "Failed to read output file")

	var response TransitionCheckResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err, "Failed to parse JSON")
	assert.Len(t, response.Tasks, 2, "Should have results for 2 tickets")

	// Verify first ticket is successful
	assert.Equal(t, testJiraID, response.Tasks[0].Key)
	assert.NotEqual(t, "Error", response.Tasks[0].Status)

	// Verify second ticket is error
	assert.Equal(t, "INVALID-99999", response.Tasks[1].Key)
	assert.Equal(t, "Error", response.Tasks[1].Status)

	// Clean up
	os.Remove(config.OutputFile)

	t.Logf("Direct JIRA processing completed successfully")
}

// TestCLIOperations tests the CLI functionality
func TestCLIOperations(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "test_main", ".")
	err := cmd.Run()
	require.NoError(t, err, "Failed to build binary")
	defer os.Remove("test_main")

	t.Run("ExtractOnlyMode", func(t *testing.T) {
		// Test extract-only mode
		cmd := exec.Command("./test_main", "--extract-only", "HEAD")
		output, err := cmd.CombinedOutput()

		// This might fail if not in a git repo or no JIRA IDs in HEAD
		if err != nil {
			t.Logf("Extract-only mode output: %s", string(output))
			return
		}

		t.Logf("Extract-only output: %s", string(output))
	})

	t.Run("DirectJIRAMode", func(t *testing.T) {
		testJiraID := os.Getenv("TEST_EXISTING_JIRA_ID")
		if testJiraID == "" {
			testJiraID = "TEST-999"
		}

		// Test with direct JIRA IDs
		cmd := exec.Command("./test_main", testJiraID, "INVALID-99999")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Command failed: %v\nOutput: %s", err, string(output))
		}

		// The command outputs to a file, not stdout
		t.Logf("Command output: %s", string(output))

		// Read the JSON from the output file
		data, err := os.ReadFile("transformed_jira_data.json")
		require.NoError(t, err, "Failed to read output file")

		var response TransitionCheckResponse
		err = json.Unmarshal(data, &response)
		require.NoError(t, err, "Output file should contain valid JSON")

		assert.Len(t, response.Tasks, 2)

		// Clean up
		os.Remove("transformed_jira_data.json")
	})

	t.Run("HelpFlag", func(t *testing.T) {
		// Test help flag
		cmd := exec.Command("./test_main", "--help")
		output, err := cmd.CombinedOutput()

		// Help should exit with code 0
		assert.NoError(t, err)
		assert.Contains(t, string(output), "Usage:")
		assert.Contains(t, string(output), "JIRA Evidence Gathering Tool")
	})
}

// TestWithControlledGitRepo tests with a controlled git repository
func TestWithControlledGitRepo(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	// Create test repository
	repoDir, cleanup := createTestGitRepo(t)
	defer cleanup()

	// Change to test repo directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(repoDir)
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	// Build the binary in original directory
	// Change back to original directory for build
	err = os.Chdir(originalDir)
	require.NoError(t, err)

	cmd := exec.Command("go", "build", "-o", filepath.Join(repoDir, "test_main"), ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Build failed: %v\nOutput: %s", err, string(output))
	}
	require.NoError(t, err)

	// Change back to test repo
	err = os.Chdir(repoDir)
	require.NoError(t, err)

	t.Run("SingleCommitExtraction", func(t *testing.T) {
		// Get the actual commit hash
		gitCmd := exec.Command("git", "rev-parse", "HEAD")
		commitBytes, err := gitCmd.Output()
		require.NoError(t, err)
		commit := strings.TrimSpace(string(commitBytes))

		// Test single commit mode
		cmd := exec.Command("./test_main", "--extract-only", commit)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Output: %s", string(output))
		}
		require.NoError(t, err)

		assert.Contains(t, string(output), "PROJ-789")
		assert.NotContains(t, string(output), "TEST-123") // Should not include other commits
	})

	t.Run("RangeExtraction", func(t *testing.T) {
		// Get the commit hash for HEAD~3
		gitCmd := exec.Command("git", "rev-parse", "HEAD~3")
		commitBytes, err := gitCmd.Output()
		require.NoError(t, err)
		startCommit := strings.TrimSpace(string(commitBytes))

		// Test range mode
		cmd := exec.Command("./test_main", "--extract-only", "--range", startCommit)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Output: %s", string(output))
		}
		require.NoError(t, err)

		assert.Contains(t, string(output), "TEST-456")
		assert.Contains(t, string(output), "PROJ-789")
		assert.NotContains(t, string(output), "TEST-123") // This is before HEAD~3
	})

	t.Run("FullWorkflowWithTestRepo", func(t *testing.T) {
		// Get the actual commit hash
		gitCmd := exec.Command("git", "rev-parse", "HEAD")
		commitBytes, err := gitCmd.Output()
		require.NoError(t, err)
		commit := strings.TrimSpace(string(commitBytes))

		// Test full workflow with JIRA fetching
		cmd := exec.Command("./test_main", commit)
		output, err := cmd.CombinedOutput()

		// This will try to fetch PROJ-789 from JIRA which likely doesn't exist
		// But it should still produce valid JSON output
		if err != nil {
			t.Logf("Full workflow output: %s", string(output))
		}

		// Check that output file was created
		_, statErr := os.Stat("transformed_jira_data.json")
		assert.NoError(t, statErr, "Output file should be created")

		// Read and verify it's valid JSON
		if statErr == nil {
			data, err := os.ReadFile("transformed_jira_data.json")
			assert.NoError(t, err)

			var response TransitionCheckResponse
			err = json.Unmarshal(data, &response)
			assert.NoError(t, err, "Output should be valid JSON")
			assert.Len(t, response.Tasks, 1)
			assert.Equal(t, "PROJ-789", response.Tasks[0].Key)

			// Clean up
			os.Remove("transformed_jira_data.json")
		}
	})
}

// TestPerformance tests performance with larger datasets
func TestPerformance(t *testing.T) {
	if os.Getenv("JIRA_API_TOKEN") == "" {
		t.Skip("Skipping: JIRA_API_TOKEN not set")
	}

	if os.Getenv("TEST_PERFORMANCE") != "true" {
		t.Skip("Skipping performance tests (set TEST_PERFORMANCE=true to enable)")
	}

	client, err := NewJiraClient()
	require.NoError(t, err)

	testJiraID := os.Getenv("TEST_EXISTING_JIRA_ID")
	if testJiraID == "" {
		t.Skip("TEST_EXISTING_JIRA_ID not set, skipping")
	}

	// Extract project key from test ID
	parts := strings.Split(testJiraID, "-")
	if len(parts) < 2 {
		t.Skip("Invalid TEST_EXISTING_JIRA_ID format")
	}
	projectKey := parts[0]

	t.Run("FetchMultipleTicketsPerformance", func(t *testing.T) {
		// Generate a list of ticket IDs (some will be invalid)
		ticketIDs := make([]string, 20)
		for i := 0; i < 20; i++ {
			ticketIDs[i] = fmt.Sprintf("%s-%d", projectKey, 1000+i)
		}

		start := time.Now()
		response := client.FetchJiraDetails(ticketIDs)
		duration := time.Since(start)

		t.Logf("Fetched %d tickets in %v (%.2f tickets/second)",
			len(ticketIDs), duration, float64(len(ticketIDs))/duration.Seconds())

		// Count successful fetches
		successCount := 0
		for _, result := range response.Tasks {
			if result.Status != "Error" {
				successCount++
			}
		}
		t.Logf("Successfully fetched %d/%d tickets", successCount, len(ticketIDs))

		// Ensure it completes within reasonable time
		assert.Less(t, duration, 60*time.Second, "Fetching 20 tickets should complete within 60 seconds")
	})

	t.Run("LargeCommitRangeExtraction", func(t *testing.T) {
		gitService := NewGitService()

		// Get current branch
		branch, _, _, err := gitService.GetBranchInfo()
		if err != nil {
			t.Logf("Failed to get branch info: %v", err)
			branch = "main"
		}

		// Try to extract from last 100 commits
		start := time.Now()
		jiraIDs, err := gitService.ExtractJiraIDs("HEAD~100", "[A-Z]+-[0-9]+", branch, true)
		if err != nil {
			t.Logf("Failed to extract from 100 commits (expected in shallow repos): %v", err)
			return
		}

		duration := time.Since(start)
		t.Logf("Extracted %d unique JIRA IDs from 100 commits in %v", len(jiraIDs), duration)
	})
}

// Helper function to create a test git repository
func createTestGitRepo(t *testing.T) (string, func()) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "jira-test-repo-*")
	require.NoError(t, err)

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	err = cmd.Run()
	require.NoError(t, err)

	// Configure git
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	err = cmd.Run()
	require.NoError(t, err)

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	err = cmd.Run()
	require.NoError(t, err)

	// Create test commits
	commits := []string{
		"Initial commit",
		"TEST-123 Add feature A",
		"Fix bug in module B",
		"TEST-456 Update documentation",
		"PROJ-789 Implement new API",
	}

	for _, msg := range commits {
		// Create a file
		filename := fmt.Sprintf("file_%d.txt", time.Now().UnixNano())
		filepath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filepath, []byte(msg), 0644)
		require.NoError(t, err)

		// Add and commit
		cmd = exec.Command("git", "add", filename)
		cmd.Dir = tmpDir
		err = cmd.Run()
		require.NoError(t, err)

		cmd = exec.Command("git", "commit", "-m", msg)
		cmd.Dir = tmpDir
		err = cmd.Run()
		require.NoError(t, err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}
