package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockGitCommand creates a mock git command function for testing
func createMockGitCommand(responses map[string]struct {
	output string
	err    error
}) func(args ...string) (string, error) {
	return func(args ...string) (string, error) {
		key := fmt.Sprintf("%v", args)
		if response, ok := responses[key]; ok {
			return response.output, response.err
		}
		return "", fmt.Errorf("unexpected git command: %v", args)
	}
}

func TestNewGitServiceComplete(t *testing.T) {
	service := NewGitService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.execCommand)
}

func TestDefaultGitCommandComplete(t *testing.T) {
	// This test verifies the real git command execution
	// It will only pass if git is installed
	output, err := defaultGitCommand("--version")

	if err == nil {
		assert.Contains(t, output, "git version")
	} else {
		t.Skip("Git not installed, skipping real command test")
	}
}

func TestGitService_GetBranchInfoComplete(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses map[string]struct {
			output string
			err    error
		}
		expectedBranch string
		expectedCommit string
		expectedJiraID string
		expectError    bool
		errorContains  string
	}{
		{
			name: "Successful branch info retrieval",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "feature/EV-123-test", err: nil},
				"[log -1 --format=%H%n%s]": {output: "abc123def456\nEV-123: Fix bug in feature", err: nil},
			},
			expectedBranch: "feature/EV-123-test",
			expectedCommit: "abc123def456",
			expectedJiraID: "EV-123",
			expectError:    false,
		},
		{
			name: "Branch command fails",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]": {output: "", err: errors.New("not a git repository")},
			},
			expectError:   true,
			errorContains: "not a git repository",
		},
		{
			name: "Log command fails",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "main", err: nil},
				"[log -1 --format=%H%n%s]": {output: "", err: errors.New("bad revision")},
			},
			expectError:   true,
			errorContains: "bad revision",
		},
		{
			name: "Log output has unexpected format - single line",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "main", err: nil},
				"[log -1 --format=%H%n%s]": {output: "onlyoneline", err: nil},
			},
			expectError:   true,
			errorContains: "unexpected output format",
		},
		{
			name: "Log output has unexpected format - empty",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "main", err: nil},
				"[log -1 --format=%H%n%s]": {output: "", err: nil},
			},
			expectError:   true,
			errorContains: "unexpected output format",
		},
		{
			name: "No JIRA ID in commit message",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "feature/no-jira", err: nil},
				"[log -1 --format=%H%n%s]": {output: "abc123def456\nGeneral cleanup", err: nil},
			},
			expectedBranch: "feature/no-jira",
			expectedCommit: "abc123def456",
			expectedJiraID: "",
			expectError:    false,
		},
		{
			name: "Multiple JIRA IDs - returns first",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "feature/multi", err: nil},
				"[log -1 --format=%H%n%s]": {output: "abc123\nEV-123, EV-456: Fix multiple issues", err: nil},
			},
			expectedBranch: "feature/multi",
			expectedCommit: "abc123",
			expectedJiraID: "EV-123",
			expectError:    false,
		},
		{
			name: "Branch with JIRA in name",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[branch --show-current]":  {output: "feature/EV-999-test", err: nil},
				"[log -1 --format=%H%n%s]": {output: "def456\nEV-789: Different ID in commit", err: nil},
			},
			expectedBranch: "feature/EV-999-test",
			expectedCommit: "def456",
			expectedJiraID: "EV-789",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &GitService{
				execCommand: createMockGitCommand(tt.mockResponses),
			}

			branch, commit, jiraID, err := git.GetBranchInfo()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBranch, branch)
				assert.Equal(t, tt.expectedCommit, commit)
				assert.Equal(t, tt.expectedJiraID, jiraID)
			}
		})
	}
}

func TestValidateCommitHashComplete(t *testing.T) {
	tests := []struct {
		name         string
		hash         string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid full SHA",
			hash:        "abc123def456789012345678901234567890abcd",
			expectError: false,
		},
		{
			name:        "Valid short SHA",
			hash:        "abc123d",
			expectError: false,
		},
		{
			name:        "Valid single character",
			hash:        "a",
			expectError: false,
		},
		{
			name:         "Empty hash",
			hash:         "",
			expectError:  true,
			errorMessage: "cannot be empty",
		},
		{
			name:         "Invalid characters - contains g",
			hash:         "abc123g",
			expectError:  true,
			errorMessage: "invalid format",
		},
		{
			name:         "Invalid characters - special chars",
			hash:         "abc-123",
			expectError:  true,
			errorMessage: "invalid format",
		},
		{
			name:         "Invalid characters - spaces",
			hash:         "abc 123",
			expectError:  true,
			errorMessage: "invalid format",
		},
		{
			name:        "Mixed case (valid)",
			hash:        "AbC123DeF",
			expectError: false,
		},
		{
			name:        "All uppercase",
			hash:        "ABCDEF123",
			expectError: false,
		},
		{
			name:        "All lowercase",
			hash:        "abcdef123",
			expectError: false,
		},
		{
			name:         "Contains non-alphanumeric",
			hash:         "abc@123",
			expectError:  true,
			errorMessage: "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommitHash(tt.hash)
			if tt.expectError {
				assert.Error(t, err)
				validationErr, ok := err.(*ValidationError)
				assert.True(t, ok)
				assert.Equal(t, "commit", validationErr.Field)
				assert.Contains(t, validationErr.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitService_ValidateCommitComplete(t *testing.T) {
	tests := []struct {
		name          string
		commit        string
		mockResponses map[string]struct {
			output string
			err    error
		}
		expectError bool
		errorType   string
	}{
		{
			name:   "Valid commit exists",
			commit: "abc123",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]": {output: "abc123def456789", err: nil},
			},
			expectError: false,
		},
		{
			name:   "Commit does not exist",
			commit: "abc123def456",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123def456]": {output: "", err: errors.New("fatal: bad revision")},
			},
			expectError: true,
			errorType:   "GitError",
		},
		{
			name:        "Invalid commit format",
			commit:      "xyz@123",
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name:        "Empty commit",
			commit:      "",
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name:   "Valid short commit",
			commit: "a1b2c3",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify a1b2c3]": {output: "a1b2c3d4e5f6", err: nil},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &GitService{
				execCommand: createMockGitCommand(tt.mockResponses),
			}

			err := git.ValidateCommit(tt.commit)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType == "GitError" {
					// The error might be wrapped, check the error message
					assert.Contains(t, err.Error(), "git operation")
				} else if tt.errorType == "ValidationError" {
					_, ok := err.(*ValidationError)
					assert.True(t, ok, "Expected ValidationError type")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitService_ValidateHEADComplete(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses map[string]struct {
			output string
			err    error
		}
		expectError  bool
		errorMessage string
	}{
		{
			name: "HEAD exists",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify HEAD]": {output: "abc123def456", err: nil},
			},
			expectError: false,
		},
		{
			name: "HEAD does not exist - empty repo",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify HEAD]": {output: "", err: errors.New("fatal: bad revision 'HEAD'")},
			},
			expectError:  true,
			errorMessage: "repository may be empty or corrupted",
		},
		{
			name: "Git command error",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify HEAD]": {output: "", err: errors.New("fatal: not a git repository")},
			},
			expectError:  true,
			errorMessage: "repository may be empty or corrupted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &GitService{
				execCommand: createMockGitCommand(tt.mockResponses),
			}

			err := git.ValidateHEAD()
			if tt.expectError {
				assert.Error(t, err)
				gitErr, ok := err.(*GitError)
				assert.True(t, ok)
				assert.Contains(t, gitErr.Error(), "rev-parse --verify HEAD")
				assert.Contains(t, gitErr.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitService_CheckRepositoryComplete(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses map[string]struct {
			output string
			err    error
		}
		expectError  bool
		errorMessage string
	}{
		{
			name: "Valid git repository",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --git-dir]": {output: ".git", err: nil},
			},
			expectError: false,
		},
		{
			name: "Not a git repository",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --git-dir]": {output: "", err: errors.New("fatal: not a git repository")},
			},
			expectError:  true,
			errorMessage: "not in a git repository",
		},
		{
			name: "Git command fails",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --git-dir]": {output: "", err: errors.New("command not found")},
			},
			expectError:  true,
			errorMessage: "not in a git repository",
		},
		{
			name: "Submodule or worktree",
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --git-dir]": {output: ".git/modules/submodule", err: nil},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git := &GitService{
				execCommand: createMockGitCommand(tt.mockResponses),
			}

			err := git.CheckRepository()
			if tt.expectError {
				assert.Error(t, err)
				gitErr, ok := err.(*GitError)
				assert.True(t, ok)
				assert.Contains(t, gitErr.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractFirstJIRAIDComplete(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected string
	}{
		{
			name:     "Extract first JIRA ID",
			text:     "EV-123: Fix bug, also fixes EV-456",
			pattern:  "[A-Z]+-[0-9]+",
			expected: "EV-123",
		},
		{
			name:     "No JIRA ID found",
			text:     "General cleanup",
			pattern:  "[A-Z]+-[0-9]+",
			expected: "",
		},
		{
			name:     "Invalid pattern",
			text:     "EV-123: Fix bug",
			pattern:  "[A-Z(+[0-9]+)", // Invalid regex - unmatched bracket
			expected: "",              // Returns empty on regex compile error
		},
		{
			name:     "Multiple patterns with different prefixes",
			text:     "TEST-999 and EV-123 and BUG-456",
			pattern:  "[A-Z]+-[0-9]+",
			expected: "TEST-999",
		},
		{
			name:     "JIRA ID at end of string",
			text:     "This fixes issue EV-789",
			pattern:  "[A-Z]+-[0-9]+",
			expected: "EV-789",
		},
		{
			name:     "JIRA ID with single digit",
			text:     "Fix: EV-1",
			pattern:  "[A-Z]+-[0-9]+",
			expected: "EV-1",
		},
		{
			name:     "Custom pattern - lowercase",
			text:     "fix: bug-123 and BUG-456",
			pattern:  "[a-z]+-[0-9]+",
			expected: "bug-123",
		},
		{
			name:     "Empty text",
			text:     "",
			pattern:  "[A-Z]+-[0-9]+",
			expected: "",
		},
		{
			name:     "Empty pattern",
			text:     "EV-123: Fix bug",
			pattern:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFirstJIRAID(tt.text, tt.pattern)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractUniqueJIRAIDsComplete(t *testing.T) {
	regex, _ := regexp.Compile("[A-Z]+-[0-9]+")

	tests := []struct {
		name           string
		commitMessages string
		currentJiraID  string
		regex          *regexp.Regexp
		expected       []string
	}{
		{
			name:           "Extract unique IDs from multiple commits",
			commitMessages: "EV-123: Fix bug\nEV-456: Add feature\nEV-123: Update docs",
			currentJiraID:  "EV-789",
			regex:          regex,
			expected:       []string{"EV-789", "EV-123", "EV-456"},
		},
		{
			name:           "No duplicate current JIRA ID",
			commitMessages: "EV-123: Fix bug\nEV-456: Add feature",
			currentJiraID:  "EV-123",
			regex:          regex,
			expected:       []string{"EV-123", "EV-456"},
		},
		{
			name:           "Empty commit messages",
			commitMessages: "",
			currentJiraID:  "EV-100",
			regex:          regex,
			expected:       []string{"EV-100"},
		},
		{
			name:           "Current JIRA ID doesn't match pattern",
			commitMessages: "EV-123: Fix bug",
			currentJiraID:  "invalid-id",
			regex:          regex,
			expected:       []string{"EV-123"},
		},
		{
			name:           "Multiple IDs in single line",
			commitMessages: "EV-123, EV-456: Fix bugs related to EV-789",
			currentJiraID:  "",
			regex:          regex,
			expected:       []string{"EV-123", "EV-456", "EV-789"},
		},
		{
			name:           "No matches at all",
			commitMessages: "General cleanup\nRefactoring\nUpdate docs",
			currentJiraID:  "",
			regex:          regex,
			expected:       []string{},
		},
		{
			name:           "Duplicates across multiple lines",
			commitMessages: "EV-100: Start\nEV-100: Continue\nEV-100: Finish",
			currentJiraID:  "EV-100",
			regex:          regex,
			expected:       []string{"EV-100"},
		},
		{
			name:           "Mixed case IDs",
			commitMessages: "EV-123: Fix\nev-456: Update\nEv-789: Add",
			currentJiraID:  "",
			regex:          regex,
			expected:       []string{"EV-123"}, // Only uppercase matches
		},
		{
			name:           "IDs with varying digit counts",
			commitMessages: "EV-1: Fix\nEV-12: Update\nEV-123: Add\nEV-1234: Test",
			currentJiraID:  "",
			regex:          regex,
			expected:       []string{"EV-1", "EV-12", "EV-123", "EV-1234"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractUniqueJIRAIDs(tt.commitMessages, tt.currentJiraID, tt.regex)
			// Sort results for consistent comparison
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestGitService_ExtractJiraIDsComplete(t *testing.T) {
	tests := []struct {
		name          string
		startCommit   string
		jiraIDRegex   string
		currentJiraID string
		singleCommit  bool
		mockResponses map[string]struct {
			output string
			err    error
		}
		expectedIDs   []string
		expectError   bool
		expectWarning bool
	}{
		{
			name:          "Single commit mode with one JIRA ID",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "EV-100",
			singleCommit:  true,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":        {output: "abc123def", err: nil},
				"[log -1 --pretty=format:%s abc123]": {output: "EV-123: Fix critical bug", err: nil},
			},
			expectedIDs: []string{"EV-123"},
			expectError: false,
		},
		{
			name:          "Single commit mode with no JIRA ID",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "",
			singleCommit:  true,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":        {output: "abc123def", err: nil},
				"[log -1 --pretty=format:%s abc123]": {output: "General cleanup", err: nil},
			},
			expectedIDs:   []string{},
			expectError:   false,
			expectWarning: true,
		},
		{
			name:          "Range mode with multiple commits",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "EV-100",
			singleCommit:  false,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":           {output: "abc123def", err: nil},
				"[log --pretty=format:%s abc123..HEAD]": {output: "EV-123: Fix bug\nEV-456: Add feature\nGeneral cleanup\nEV-123: Update docs", err: nil},
			},
			expectedIDs: []string{"EV-100", "EV-123", "EV-456"},
			expectError: false,
		},
		{
			name:          "Invalid commit hash",
			startCommit:   "invalid@commit",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "",
			singleCommit:  true,
			expectError:   true,
		},
		{
			name:          "Commit not found",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "",
			singleCommit:  true,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]": {output: "", err: errors.New("bad revision")},
			},
			expectError: true,
		},
		{
			name:          "Invalid regex pattern",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z(+[0-9]+)", // Invalid regex - unmatched bracket
			currentJiraID: "",
			singleCommit:  true,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":        {output: "abc123def", err: nil},
				"[log -1 --pretty=format:%s abc123]": {output: "EV-123: Fix bug", err: nil},
			},
			expectError: true,
		},
		{
			name:          "Log command fails",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "",
			singleCommit:  true,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":        {output: "abc123def", err: nil},
				"[log -1 --pretty=format:%s abc123]": {output: "", err: errors.New("log failed")},
			},
			expectError: true,
		},
		{
			name:          "Range mode with empty range",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "EV-100",
			singleCommit:  false,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":           {output: "abc123def", err: nil},
				"[log --pretty=format:%s abc123..HEAD]": {output: "", err: nil},
			},
			expectedIDs:   []string{"EV-100"},
			expectError:   false,
			expectWarning: false, // Warning only shows when no IDs found at all
		},
		{
			name:          "Multiple IDs in single commit",
			startCommit:   "abc123",
			jiraIDRegex:   "[A-Z]+-[0-9]+",
			currentJiraID: "",
			singleCommit:  true,
			mockResponses: map[string]struct {
				output string
				err    error
			}{
				"[rev-parse --verify abc123]":        {output: "abc123def", err: nil},
				"[log -1 --pretty=format:%s abc123]": {output: "EV-123, EV-456: Fix multiple issues", err: nil},
			},
			expectedIDs: []string{"EV-123", "EV-456"},
			expectError: false,
		},
	}

	// Capture stderr for warning messages
	oldStderr := os.Stderr

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr for this test
			r, w, _ := os.Pipe()
			os.Stderr = w

			git := &GitService{
				execCommand: createMockGitCommand(tt.mockResponses),
			}

			ids, err := git.ExtractJiraIDs(tt.startCommit, tt.jiraIDRegex, tt.currentJiraID, tt.singleCommit)

			// Close and restore stderr
			w.Close()
			os.Stderr = oldStderr

			// Read stderr output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			stderrOutput := string(buf[:n])

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedIDs, ids)
			}

			if tt.expectWarning {
				assert.Contains(t, stderrOutput, "⚠️")
			}
		})
	}
}
