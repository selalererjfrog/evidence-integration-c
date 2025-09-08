package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMarkdownFromJSON(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "markdown-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		jsonContent string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid JSON with tasks",
			jsonContent: `{
				"tasks": [
					{
						"key": "TEST-123",
						"status": "Done",
						"description": "Test task",
						"type": "Task",
						"project": "TEST",
						"created": "2025-01-01T10:00:00.000+0300",
						"updated": "2025-01-02T15:00:00.000+0300",
						"assignee": "John Doe",
						"reporter": "Jane Smith",
						"priority": "High",
						"transitions": []
					}
				]
			}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			jsonContent: `{invalid json}`,
			expectError: true,
			errorMsg:    "error parsing JSON",
		},
		{
			name:        "Empty file",
			jsonContent: "",
			expectError: true,
			errorMsg:    "error parsing JSON",
		},
		{
			name: "Valid JSON with empty tasks",
			jsonContent: `{
				"tasks": []
			}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create input JSON file
			inputFile := filepath.Join(tempDir, "input.json")
			outputFile := filepath.Join(tempDir, "output.md")

			err := os.WriteFile(inputFile, []byte(tt.jsonContent), 0644)
			require.NoError(t, err)

			// Test GenerateMarkdownFromJSON
			err = GenerateMarkdownFromJSON(inputFile, outputFile)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				// Verify output file exists
				_, err = os.Stat(outputFile)
				assert.NoError(t, err)
			}

			// Clean up
			os.Remove(inputFile)
			os.Remove(outputFile)
		})
	}

	// Test with non-existent file
	t.Run("Non-existent input file", func(t *testing.T) {
		err := GenerateMarkdownFromJSON("/non/existent/file.json", "/tmp/output.md")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error reading JSON file")
	})

	// Test with invalid output path
	t.Run("Invalid output path", func(t *testing.T) {
		inputFile := filepath.Join(tempDir, "valid.json")
		err := os.WriteFile(inputFile, []byte(`{"tasks": []}`), 0644)
		require.NoError(t, err)

		err = GenerateMarkdownFromJSON(inputFile, "/root/invalid/path/output.md")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error writing markdown file")
	})
}

func TestGenerateMarkdown(t *testing.T) {
	// Test various task scenarios
	tests := []struct {
		name     string
		response TransitionCheckResponse
		checks   []string // Strings that should appear in the output
	}{
		{
			name: "Single task with transitions and link",
			response: TransitionCheckResponse{
				Tasks: []JiraTransitionResult{
					{
						Key:         "EV-123",
						Link:        "https://example.atlassian.net/browse/EV-123",
						Status:      "Done",
						Description: "Test description",
						Type:        "Task",
						Project:     "EV",
						Created:     "2025-01-01T10:00:00.000+0300",
						Updated:     "2025-01-05T15:00:00.000+0300",
						Assignee:    strPtr("John Doe"),
						Reporter:    "Jane Smith",
						Priority:    "High",
						Transitions: []Transition{
							{
								FromStatus:     "To Do",
								ToStatus:       "In Progress",
								Author:         "John Doe",
								AuthorEmail:    "john@example.com",
								TransitionTime: "2025-01-02T09:00:00.000+0300",
							},
						},
					},
				},
			},
			checks: []string{
				"# JIRA Tasks Report",
				"Total tasks: 1",
				"[EV-123](https://example.atlassian.net/browse/EV-123)",
				"### 1. [EV-123](https://example.atlassian.net/browse/EV-123)",
				"**Status:** Done",
				"**Type:** Task",
				"**Assignee:** John Doe",
				"**Reporter:** Jane Smith",
				"**Description:**",
				"> Test description",
				"**Transition History:**",
				"| To Do | In Progress | John Doe | 2025-01-02 09:00:00 |",
				"## Status Distribution",
				"| Done | 1 |",
			},
		},
		{
			name: "Task with null assignee and no transitions",
			response: TransitionCheckResponse{
				Tasks: []JiraTransitionResult{
					{
						Key:         "BUG-456",
						Status:      "Open",
						Description: "",
						Type:        "Bug",
						Project:     "PROJ",
						Created:     "2025-02-01T12:00:00.000+0300",
						Updated:     "2025-02-01T12:00:00.000+0300",
						Assignee:    nil,
						Reporter:    "Bob Builder",
						Priority:    "Medium",
						Transitions: []Transition{},
					},
				},
			},
			checks: []string{
				"| BUG-456 | Open | Bug | Medium | Unassigned |",
				"**Assignee:** Unassigned",
				"**Created:** 2025-02-01 12:00:00",
			},
		},
		{
			name: "Error task",
			response: TransitionCheckResponse{
				Tasks: []JiraTransitionResult{
					{
						Key:         "ERR-789",
						Status:      "Error",
						Description: "Error: Could not retrieve issue",
						Type:        "Error",
						Project:     "",
						Created:     "",
						Updated:     "",
						Assignee:    nil,
						Reporter:    "",
						Priority:    "",
						Transitions: []Transition{},
					},
				},
			},
			checks: []string{
				"| ERR-789 | Error | Error |  | Unassigned |", // No link for error tasks
				"### 1. ERR-789", // No link in header for error tasks
				"**Created:** N/A",
				"**Updated:** N/A",
				"> Error: Could not retrieve issue",
			},
		},
		{
			name: "Multiple tasks with different statuses",
			response: TransitionCheckResponse{
				Tasks: []JiraTransitionResult{
					{
						Key:      "T1",
						Status:   "Done",
						Type:     "Task",
						Project:  "P1",
						Created:  "2025-01-01T10:00:00.000+0300",
						Updated:  "2025-01-01T10:00:00.000+0300",
						Assignee: strPtr("Alice"),
						Reporter: "Bob",
						Priority: "High",
					},
					{
						Key:      "T2",
						Status:   "In Progress",
						Type:     "Bug",
						Project:  "P1",
						Created:  "2025-01-01T10:00:00.000+0300",
						Updated:  "2025-01-01T10:00:00.000+0300",
						Assignee: strPtr("Charlie"),
						Reporter: "Dave",
						Priority: "Low",
					},
					{
						Key:      "T3",
						Status:   "Done",
						Type:     "Story",
						Project:  "P2",
						Created:  "2025-01-01T10:00:00.000+0300",
						Updated:  "2025-01-01T10:00:00.000+0300",
						Assignee: nil,
						Reporter: "Eve",
						Priority: "Medium",
					},
				},
			},
			checks: []string{
				"Total tasks: 3",
				"| T1 | Done | Task | High | Alice |",
				"| T2 | In Progress | Bug | Low | Charlie |",
				"| T3 | Done | Story | Medium | Unassigned |",
				"| Done | 2 |",
				"| In Progress | 1 |",
			},
		},
		{
			name:     "Empty task list",
			response: TransitionCheckResponse{Tasks: []JiraTransitionResult{}},
			checks: []string{
				"Total tasks: 0",
				"## Summary",
				"## Task Details",
				"## Status Distribution",
			},
		},
		{
			name: "Task with multiline description and link",
			response: TransitionCheckResponse{
				Tasks: []JiraTransitionResult{
					{
						Key:         "MULTI-123",
						Link:        "https://test.atlassian.net/browse/MULTI-123",
						Status:      "Open",
						Description: "Line 1\nLine 2\nLine 3",
						Type:        "Task",
						Project:     "TEST",
						Created:     "2025-01-01T10:00:00.000+0300",
						Updated:     "2025-01-01T10:00:00.000+0300",
						Assignee:    strPtr("Test User"),
						Reporter:    "Reporter",
						Priority:    "Low",
						Transitions: []Transition{},
					},
				},
			},
			checks: []string{
				"> Line 1\n> Line 2\n> Line 3",
				"[MULTI-123](https://test.atlassian.net/browse/MULTI-123)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown := generateMarkdown(tt.response)

			// Check that all expected strings are present
			for _, check := range tt.checks {
				assert.Contains(t, markdown, check, "Expected to find: %s", check)
			}

			// Basic structure checks
			assert.True(t, strings.HasPrefix(markdown, "# JIRA Tasks Report"))
			assert.Contains(t, markdown, "Generated on:")
			assert.Contains(t, markdown, "## Summary")
			assert.Contains(t, markdown, "## Task Details")
			assert.Contains(t, markdown, "## Status Distribution")
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid JIRA date",
			input:    "2025-01-01T10:30:45.123+0300",
			expected: "2025-01-01 10:30:45",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "N/A",
		},
		{
			name:     "Invalid date format",
			input:    "not-a-date",
			expected: "not-a-date",
		},
		{
			name:     "Different timezone",
			input:    "2025-12-31T23:59:59.999-0500",
			expected: "2025-12-31 23:59:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}

// Test that generated markdown includes current time
func TestGenerateMarkdownIncludesCurrentTime(t *testing.T) {
	response := TransitionCheckResponse{
		Tasks: []JiraTransitionResult{},
	}

	markdown := generateMarkdown(response)

	// Check that it includes a date in the expected format
	now := time.Now()
	yearStr := now.Format("2006")
	assert.Contains(t, markdown, yearStr)
	assert.Contains(t, markdown, "Generated on:")
}

// Test markdown formatting edge cases
func TestMarkdownFormattingEdgeCases(t *testing.T) {
	response := TransitionCheckResponse{
		Tasks: []JiraTransitionResult{
			{
				Key:         "EDGE-1",
				Status:      "Status|With|Pipes",
				Description: "Description with | pipes and **markdown**",
				Type:        "Type*With*Stars",
				Project:     "PROJ",
				Created:     "2025-01-01T10:00:00.000+0300",
				Updated:     "2025-01-01T10:00:00.000+0300",
				Assignee:    strPtr("User|With|Pipes"),
				Reporter:    "Reporter**Bold**",
				Priority:    "Priority`Code`",
				Transitions: []Transition{
					{
						FromStatus:     "From|Status",
						ToStatus:       "To|Status",
						Author:         "Author|Name",
						AuthorEmail:    "email@with|pipe.com",
						TransitionTime: "2025-01-01T11:00:00.000+0300",
					},
				},
			},
		},
	}

	markdown := generateMarkdown(response)

	// Verify the special characters are preserved in the output
	assert.Contains(t, markdown, "Status|With|Pipes")
	assert.Contains(t, markdown, "Description with | pipes and **markdown**")
	assert.Contains(t, markdown, "User|With|Pipes")
	assert.Contains(t, markdown, "From|Status")
	assert.Contains(t, markdown, "To|Status")
}
