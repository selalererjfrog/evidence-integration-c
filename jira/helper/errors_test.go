package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitError(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		err           error
		expectedError string
	}{
		{
			name:          "GitError with simple error",
			operation:     "rev-parse",
			err:           errors.New("command failed"),
			expectedError: "git operation 'rev-parse' failed: command failed",
		},
		{
			name:          "GitError with complex operation",
			operation:     "log --pretty=format:%s",
			err:           errors.New("fatal: bad revision"),
			expectedError: "git operation 'log --pretty=format:%s' failed: fatal: bad revision",
		},
		{
			name:          "GitError with nil error",
			operation:     "status",
			err:           nil,
			expectedError: "git operation 'status' failed: <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitErr := &GitError{
				Operation: tt.operation,
				Err:       tt.err,
			}
			assert.Equal(t, tt.expectedError, gitErr.Error())
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		value         string
		err           error
		expectedError string
	}{
		{
			name:          "ValidationError with empty value",
			field:         "JIRA_API_TOKEN",
			value:         "",
			err:           errors.New("environment variable is required"),
			expectedError: "validation failed for JIRA_API_TOKEN='': environment variable is required",
		},
		{
			name:          "ValidationError with invalid value",
			field:         "commit",
			value:         "xyz123",
			err:           errors.New("invalid format"),
			expectedError: "validation failed for commit='xyz123': invalid format",
		},
		{
			name:          "ValidationError with complex field name",
			field:         "jira_id_regex",
			value:         "[A-Z+-[0-9]+",
			err:           errors.New("error parsing regexp"),
			expectedError: "validation failed for jira_id_regex='[A-Z+-[0-9]+': error parsing regexp",
		},
		{
			name:          "ValidationError with nil error",
			field:         "test_field",
			value:         "test_value",
			err:           nil,
			expectedError: "validation failed for test_field='test_value': <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validationErr := &ValidationError{
				Field: tt.field,
				Value: tt.value,
				Err:   tt.err,
			}
			assert.Equal(t, tt.expectedError, validationErr.Error())
		})
	}
}
