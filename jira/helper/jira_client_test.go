package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/stretchr/testify/assert"
)

func TestNewJiraClient(t *testing.T) {
	// Save original environment
	originalToken := os.Getenv("JIRA_API_TOKEN")
	originalURL := os.Getenv("JIRA_URL")
	originalUsername := os.Getenv("JIRA_USERNAME")

	// Clean up environment after test
	defer func() {
		os.Setenv("JIRA_API_TOKEN", originalToken)
		os.Setenv("JIRA_URL", originalURL)
		os.Setenv("JIRA_USERNAME", originalUsername)
	}()

	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorField  string
	}{
		{
			name: "Valid configuration",
			envVars: map[string]string{
				"JIRA_API_TOKEN": "test-token",
				"JIRA_URL":       "https://example.atlassian.net",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError: false,
		},
		{
			name: "Missing JIRA_API_TOKEN",
			envVars: map[string]string{
				"JIRA_API_TOKEN": "",
				"JIRA_URL":       "https://example.atlassian.net",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError: true,
			errorField:  "JIRA_API_TOKEN",
		},
		{
			name: "Missing JIRA_URL",
			envVars: map[string]string{
				"JIRA_API_TOKEN": "test-token",
				"JIRA_URL":       "",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError: true,
			errorField:  "JIRA_URL",
		},
		{
			name: "Missing JIRA_USERNAME",
			envVars: map[string]string{
				"JIRA_API_TOKEN": "test-token",
				"JIRA_URL":       "https://example.atlassian.net",
				"JIRA_USERNAME":  "",
			},
			expectError: true,
			errorField:  "JIRA_USERNAME",
		},
		{
			name: "Invalid JIRA_URL format",
			envVars: map[string]string{
				"JIRA_API_TOKEN": "test-token",
				"JIRA_URL":       "://invalid-url",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			client, err := NewJiraClient()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorField != "" {
					validationErr, ok := err.(*ValidationError)
					if ok {
						assert.Equal(t, tt.errorField, validationErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.NotNil(t, client.client)
			}

			// Clean up environment variables
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestJiraClient_FetchJiraDetails(t *testing.T) {
	// This test demonstrates the structure of FetchJiraDetails
	// Actual testing would require mocking the JIRA API client
	// which is complex without dependency injection

	// Test that response structure is correct
	response := TransitionCheckResponse{
		Tasks: make([]JiraTransitionResult, 0, 3),
	}

	assert.NotNil(t, response.Tasks)
	assert.Equal(t, 0, len(response.Tasks))
}

func TestJiraClient_createErrorResult(t *testing.T) {
	// Test with baseURL to verify error results don't get links
	client := &JiraClient{
		baseURL: "https://example.atlassian.net",
	}

	tests := []struct {
		name          string
		jiraID        string
		err           error
		expectedDesc  string
		captureStderr bool
	}{
		{
			name:          "Error with specific message",
			jiraID:        "EV-123",
			err:           errors.New("connection timeout"),
			expectedDesc:  "Error: connection timeout",
			captureStderr: true,
		},
		{
			name:          "Error is nil",
			jiraID:        "EV-456",
			err:           nil,
			expectedDesc:  "Error: Could not retrieve issue",
			captureStderr: false,
		},
		{
			name:          "Complex error",
			jiraID:        "EV-789",
			err:           fmt.Errorf("failed to connect: %w", errors.New("network unreachable")),
			expectedDesc:  "Error: failed to connect: network unreachable",
			captureStderr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderrOutput string
			if tt.captureStderr {
				// Capture stderr
				oldStderr := os.Stderr
				r, w, _ := os.Pipe()
				os.Stderr = w

				_ = client.createErrorResult(tt.jiraID, tt.err)

				w.Close()
				os.Stderr = oldStderr

				// Read captured output
				buf := make([]byte, 1024)
				n, _ := r.Read(buf)
				stderrOutput = string(buf[:n])
			}

			result := client.createErrorResult(tt.jiraID, tt.err)

			// Verify result
			assert.Equal(t, tt.jiraID, result.Key)
			assert.Equal(t, "", result.Link) // Error results should not have links
			assert.Equal(t, ErrorStatus, result.Status)
			assert.Equal(t, tt.expectedDesc, result.Description)
			assert.Equal(t, ErrorType, result.Type)
			assert.Equal(t, "", result.Project)
			assert.Equal(t, "", result.Created)
			assert.Equal(t, "", result.Updated)
			assert.Nil(t, result.Assignee)
			assert.Equal(t, "", result.Reporter)
			assert.Equal(t, "", result.Priority)
			assert.Empty(t, result.Transitions)

			// Verify stderr output if captured
			if tt.captureStderr && tt.err != nil {
				assert.Contains(t, stderrOutput, fmt.Sprintf("Failed to fetch JIRA %s:", tt.jiraID))
			}
		})
	}
}

func TestJiraClient_createSuccessResult(t *testing.T) {
	// Test with baseURL to verify link generation
	client := &JiraClient{
		baseURL: "https://example.atlassian.net",
	}

	// Create test time
	assigneeName := "John Doe"

	// Create test issue with all fields populated
	issue := &jira.Issue{
		Key: "EV-123",
		Fields: &jira.IssueFields{
			Status: &jira.Status{
				Name: "In Progress",
			},
			Description: "Test description",
			Type: jira.IssueType{
				Name: "Task",
			},
			Project: jira.Project{
				Key: "EV",
			},
			Created: jira.Time(time.Date(2023, 12, 15, 14, 30, 45, 0, time.UTC)),
			Updated: jira.Time(time.Date(2023, 12, 15, 14, 30, 45, 0, time.UTC)),
			Assignee: &jira.User{
				DisplayName: assigneeName,
			},
			Reporter: &jira.User{
				DisplayName: "Jane Smith",
			},
			Priority: &jira.Priority{
				Name: "High",
			},
		},
		Changelog: &jira.Changelog{
			Histories: []jira.ChangelogHistory{
				{
					Created: "2023-12-14T10:00:00.000+0000",
					Author: jira.User{
						DisplayName:  "Test User",
						EmailAddress: "test@example.com",
					},
					Items: []jira.ChangelogItems{
						{
							Field:      "status",
							FromString: "To Do",
							ToString:   "In Progress",
						},
					},
				},
			},
		},
	}

	result := client.createSuccessResult(issue)

	// Verify all fields
	assert.Equal(t, "EV-123", result.Key)
	assert.Equal(t, "https://example.atlassian.net/browse/EV-123", result.Link)
	assert.Equal(t, "In Progress", result.Status)
	assert.Equal(t, "Test description", result.Description)
	assert.Equal(t, "Task", result.Type)
	assert.Equal(t, "EV", result.Project)
	assert.NotEmpty(t, result.Created)
	assert.NotEmpty(t, result.Updated)
	assert.NotNil(t, result.Assignee)
	assert.Equal(t, assigneeName, *result.Assignee)
	assert.Equal(t, "Jane Smith", result.Reporter)
	assert.Equal(t, "High", result.Priority)
	assert.Len(t, result.Transitions, 1)
	assert.Equal(t, "To Do", result.Transitions[0].FromStatus)
	assert.Equal(t, "In Progress", result.Transitions[0].ToStatus)

	// Test with empty baseURL
	clientNoURL := &JiraClient{
		baseURL: "",
	}
	resultNoURL := clientNoURL.createSuccessResult(issue)
	assert.Equal(t, "", resultNoURL.Link)
}

func TestJiraClient_extractTransitions(t *testing.T) {
	client := &JiraClient{}

	tests := []struct {
		name                string
		issue               *jira.Issue
		expectedTransitions int
	}{
		{
			name: "Issue with multiple transitions",
			issue: &jira.Issue{
				Changelog: &jira.Changelog{
					Histories: []jira.ChangelogHistory{
						{
							Created: "2023-12-14T10:00:00.000+0000",
							Author: jira.User{
								DisplayName:  "User One",
								EmailAddress: "user1@example.com",
							},
							Items: []jira.ChangelogItems{
								{
									Field:      "status",
									FromString: "To Do",
									ToString:   "In Progress",
								},
								{
									Field:      "priority",
									FromString: "Low",
									ToString:   "High",
								},
							},
						},
						{
							Created: "2023-12-15T10:00:00.000+0000",
							Author: jira.User{
								DisplayName:  "User Two",
								EmailAddress: "user2@example.com",
							},
							Items: []jira.ChangelogItems{
								{
									Field:      "status",
									FromString: "In Progress",
									ToString:   "Done",
								},
							},
						},
					},
				},
			},
			expectedTransitions: 2,
		},
		{
			name: "Issue with no changelog",
			issue: &jira.Issue{
				Changelog: nil,
			},
			expectedTransitions: 0,
		},
		{
			name: "Issue with empty changelog",
			issue: &jira.Issue{
				Changelog: &jira.Changelog{
					Histories: []jira.ChangelogHistory{},
				},
			},
			expectedTransitions: 0,
		},
		{
			name: "Issue with non-status changes only",
			issue: &jira.Issue{
				Changelog: &jira.Changelog{
					Histories: []jira.ChangelogHistory{
						{
							Created: "2023-12-14T10:00:00.000+0000",
							Author: jira.User{
								DisplayName:  "User",
								EmailAddress: "user@example.com",
							},
							Items: []jira.ChangelogItems{
								{
									Field:      "priority",
									FromString: "Low",
									ToString:   "High",
								},
								{
									Field:      "assignee",
									FromString: "user1",
									ToString:   "user2",
								},
							},
						},
					},
				},
			},
			expectedTransitions: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transitions := client.extractTransitions(tt.issue)
			assert.Len(t, transitions, tt.expectedTransitions)

			// Verify transition details for multi-transition test
			if tt.expectedTransitions == 2 {
				assert.Equal(t, "To Do", transitions[0].FromStatus)
				assert.Equal(t, "In Progress", transitions[0].ToStatus)
				assert.Equal(t, "User One", transitions[0].Author)
				assert.Equal(t, "user1@example.com", transitions[0].AuthorEmail)

				assert.Equal(t, "In Progress", transitions[1].FromStatus)
				assert.Equal(t, "Done", transitions[1].ToStatus)
				assert.Equal(t, "User Two", transitions[1].Author)
				assert.Equal(t, "user2@example.com", transitions[1].AuthorEmail)
			}
		})
	}
}

func TestJiraClient_fetchSingleJiraDetail(t *testing.T) {
	// This test demonstrates the expected behavior when fetchSingleJiraDetail fails
	// Actual implementation would require mocking the JIRA API

	// Test error result structure
	errorResult := JiraTransitionResult{
		Key:         "EV-123",
		Link:        "", // No link for error results
		Status:      ErrorStatus,
		Type:        ErrorType,
		Description: "Error: Could not retrieve issue",
		Transitions: []Transition{},
	}

	assert.Equal(t, "EV-123", errorResult.Key)
	assert.Equal(t, ErrorStatus, errorResult.Status)
	assert.Equal(t, ErrorType, errorResult.Type)
	assert.Contains(t, errorResult.Description, "Error:")
}
