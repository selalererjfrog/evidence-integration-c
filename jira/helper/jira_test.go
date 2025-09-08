package main

import (
	"context"
	"testing"
	"time"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJiraClient is a mock implementation of the JIRA client
type MockJiraClient struct {
	mock.Mock
}

// MockIssueService mocks the JIRA issue service
type MockIssueService struct {
	mock.Mock
}

func (m *MockIssueService) Get(ctx context.Context, issueID string, options *jira.GetQueryOptions) (*jira.Issue, *jira.Response, error) {
	args := m.Called(ctx, issueID, options)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*jira.Issue), args.Get(1).(*jira.Response), args.Error(2)
}

// Test getDescription function
func TestGetDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "string input",
			input:    "simple string",
			expected: "simple string",
		},
		{
			name: "ADF format with text",
			input: map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Test description",
							},
						},
					},
				},
			},
			expected: "Test description",
		},
		{
			name: "ADF format with multiple paragraphs",
			input: map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "First paragraph",
							},
						},
					},
					map[string]interface{}{
						"type": "paragraph",
						"content": []interface{}{
							map[string]interface{}{
								"type": "text",
								"text": "Second paragraph",
							},
						},
					},
				},
			},
			expected: "First paragraphSecond paragraph",
		},
		{
			name:     "invalid ADF format",
			input:    map[string]interface{}{"invalid": "format"},
			expected: "map[invalid:format]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDescription(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test extractTextFromADFNode function
func TestExtractTextFromADFNode(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "non-map input",
			input:    "string",
			expected: "",
		},
		{
			name: "text node",
			input: map[string]interface{}{
				"type": "text",
				"text": "Hello World",
			},
			expected: "Hello World",
		},
		{
			name: "paragraph node",
			input: map[string]interface{}{
				"type": "paragraph",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": "Paragraph text",
					},
				},
			},
			expected: "Paragraph text",
		},
		{
			name: "unknown node type",
			input: map[string]interface{}{
				"type": "unknown",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTextFromADFNode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test getTimeAsString function
func TestGetTimeAsString(t *testing.T) {
	now := time.Now()
	nowPtr := &now
	expectedFormat := now.Format("2006-01-02T15:04:05.000-0700")

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "string input",
			input:    "2023-01-01T12:00:00.000+0000",
			expected: "2023-01-01T12:00:00.000+0000",
		},
		{
			name:     "time.Time input",
			input:    now,
			expected: expectedFormat,
		},
		{
			name:     "*time.Time input",
			input:    nowPtr,
			expected: expectedFormat,
		},
		{
			name:     "unknown type",
			input:    123,
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTimeAsString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test field extractors
func TestFieldExtractors(t *testing.T) {
	t.Run("getStatusName", func(t *testing.T) {
		// Test nil status
		assert.Equal(t, "", getStatusName(nil))

		// Test valid status
		status := &jira.Status{Name: "In Progress"}
		assert.Equal(t, "In Progress", getStatusName(status))
	})

	t.Run("getIssueTypeName", func(t *testing.T) {
		issueType := jira.IssueType{Name: "Task"}
		assert.Equal(t, "Task", getIssueTypeName(issueType))
	})

	t.Run("getProjectKey", func(t *testing.T) {
		project := jira.Project{Key: "PROJ"}
		assert.Equal(t, "PROJ", getProjectKey(project))
	})

	t.Run("getReporterName", func(t *testing.T) {
		// Test nil reporter
		assert.Equal(t, "", getReporterName(nil))

		// Test valid reporter
		reporter := &jira.User{DisplayName: "John Doe"}
		assert.Equal(t, "John Doe", getReporterName(reporter))
	})

	t.Run("getPriorityName", func(t *testing.T) {
		// Test nil priority
		assert.Equal(t, "", getPriorityName(nil))

		// Test valid priority
		priority := &jira.Priority{Name: "High"}
		assert.Equal(t, "High", getPriorityName(priority))
	})

	t.Run("getAssignee", func(t *testing.T) {
		// Test nil assignee
		assert.Nil(t, getAssignee(nil))

		// Test valid assignee
		assignee := &jira.User{DisplayName: "Jane Doe"}
		result := getAssignee(assignee)
		assert.NotNil(t, result)
		assert.Equal(t, "Jane Doe", *result)
	})
}
