package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test main function behavior
// Note: The main() function itself is typically not unit tested directly
// as it's the entry point. Integration tests cover its behavior.

// Test LoadConfigCommitRange to ensure compatibility with existing tests
func TestLoadConfigCommitRange(t *testing.T) {
	// Save and restore environment variables
	oldJiraIDRegex := os.Getenv("JIRA_ID_REGEX")
	oldOutputFile := os.Getenv("OUTPUT_FILE")
	defer func() {
		os.Setenv("JIRA_ID_REGEX", oldJiraIDRegex)
		os.Setenv("OUTPUT_FILE", oldOutputFile)
	}()

	tests := []struct {
		name           string
		flags          *FlagConfig
		expectedSingle bool
	}{
		{
			name: "default behavior (single commit)",
			flags: &FlagConfig{
				CommitRange: false,
			},
			expectedSingle: true,
		},
		{
			name: "range mode enabled",
			flags: &FlagConfig{
				CommitRange: true,
			},
			expectedSingle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set extract-only to skip JIRA config validation
			tt.flags.ExtractOnly = true

			config, err := LoadConfig(tt.flags, []string{})
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSingle, config.SingleCommit)
		})
	}
}
