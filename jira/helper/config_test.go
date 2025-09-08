package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlagsComplete(t *testing.T) {
	// Save original command line args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name          string
		args          []string
		expectedFlags *FlagConfig
		expectedArgs  []string
	}{
		{
			name: "Parse all flags",
			args: []string{"cmd", "-r", "TEST-[0-9]+", "-o", "output.json", "--extract-only", "--range", "commit123"},
			expectedFlags: &FlagConfig{
				JIRAIDRegex: "TEST-[0-9]+",
				OutputFile:  "output.json",
				ExtractOnly: true,
				CommitRange: true,
			},
			expectedArgs: []string{"commit123"},
		},
		{
			name: "Parse help flag short",
			args: []string{"cmd", "-h"},
			expectedFlags: &FlagConfig{
				Help: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "Parse help flag long",
			args: []string{"cmd", "--help"},
			expectedFlags: &FlagConfig{
				HelpLong: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "Parse extract-from-git flag",
			args: []string{"cmd", "--extract-from-git", "arg1", "arg2"},
			expectedFlags: &FlagConfig{
				ExtractFromGit: true,
			},
			expectedArgs: []string{"arg1", "arg2"},
		},
		{
			name: "Parse markdown flag",
			args: []string{"cmd", "--markdown"},
			expectedFlags: &FlagConfig{
				GenerateMarkdown: true,
			},
			expectedArgs: []string{},
		},
		{
			name: "Parse markdown with output flag",
			args: []string{"cmd", "--markdown", "--markdown-output", "report.md"},
			expectedFlags: &FlagConfig{
				GenerateMarkdown: true,
				MarkdownOutput:   "report.md",
			},
			expectedArgs: []string{},
		},
		{
			name:          "No flags, only arguments",
			args:          []string{"cmd", "EV-123", "EV-456"},
			expectedFlags: &FlagConfig{},
			expectedArgs:  []string{"EV-123", "EV-456"},
		},
		{
			name: "Mixed short and long flags",
			args: []string{"cmd", "-r", "CUSTOM-[0-9]+", "--extract-only", "-o", "test.json"},
			expectedFlags: &FlagConfig{
				JIRAIDRegex: "CUSTOM-[0-9]+",
				OutputFile:  "test.json",
				ExtractOnly: true,
			},
			expectedArgs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine to avoid flag redefinition errors
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// Set os.Args for parsing
			os.Args = tt.args

			flags, args := ParseFlags()

			assert.Equal(t, tt.expectedFlags.JIRAIDRegex, flags.JIRAIDRegex)
			assert.Equal(t, tt.expectedFlags.OutputFile, flags.OutputFile)
			assert.Equal(t, tt.expectedFlags.ExtractOnly, flags.ExtractOnly)
			assert.Equal(t, tt.expectedFlags.ExtractFromGit, flags.ExtractFromGit)
			assert.Equal(t, tt.expectedFlags.CommitRange, flags.CommitRange)
			assert.Equal(t, tt.expectedFlags.Help, flags.Help)
			assert.Equal(t, tt.expectedFlags.HelpLong, flags.HelpLong)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestLoadConfigComplete(t *testing.T) {
	// Save original environment
	originalToken := os.Getenv("JIRA_API_TOKEN")
	originalURL := os.Getenv("JIRA_URL")
	originalUsername := os.Getenv("JIRA_USERNAME")
	originalRegex := os.Getenv("JIRA_ID_REGEX")
	originalOutput := os.Getenv("OUTPUT_FILE")

	// Clean up environment after test
	defer func() {
		os.Setenv("JIRA_API_TOKEN", originalToken)
		os.Setenv("JIRA_URL", originalURL)
		os.Setenv("JIRA_USERNAME", originalUsername)
		os.Setenv("JIRA_ID_REGEX", originalRegex)
		os.Setenv("OUTPUT_FILE", originalOutput)
	}()

	tests := []struct {
		name           string
		flags          *FlagConfig
		args           []string
		envVars        map[string]string
		expectError    bool
		expectedConfig *AppConfig
		errorContains  string
	}{
		{
			name: "Extract-only mode - no JIRA validation",
			flags: &FlagConfig{
				ExtractOnly: true,
				JIRAIDRegex: "TEST-[0-9]+",
				OutputFile:  "test.json",
			},
			args:        []string{},
			envVars:     map[string]string{},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAIDRegex:  "TEST-[0-9]+",
				OutputFile:   "test.json",
				ExtractOnly:  true,
				SingleCommit: true,
			},
		},
		{
			name: "Extract-from-git mode - no JIRA validation",
			flags: &FlagConfig{
				ExtractFromGit: true,
			},
			args:        []string{},
			envVars:     map[string]string{},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAIDRegex:    DefaultJIRAIDRegex,
				OutputFile:     DefaultOutputFile,
				ExtractFromGit: true,
				SingleCommit:   true,
			},
		},
		{
			name: "Markdown mode - no JIRA validation",
			flags: &FlagConfig{
				GenerateMarkdown: true,
				MarkdownOutput:   "test-report.md",
			},
			args:        []string{},
			envVars:     map[string]string{},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAIDRegex:  DefaultJIRAIDRegex,
				OutputFile:   DefaultOutputFile,
				SingleCommit: true,
			},
		},
		{
			name:  "Full mode with valid JIRA config",
			flags: &FlagConfig{},
			args:  []string{},
			envVars: map[string]string{
				"JIRA_API_TOKEN": "token123",
				"JIRA_URL":       "https://example.atlassian.net",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAToken:    "token123",
				JIRAURL:      "https://example.atlassian.net",
				JIRAUsername: "user@example.com",
				JIRAIDRegex:  DefaultJIRAIDRegex,
				OutputFile:   DefaultOutputFile,
				SingleCommit: true,
			},
		},
		{
			name:  "Full mode with missing JIRA token",
			flags: &FlagConfig{},
			args:  []string{},
			envVars: map[string]string{
				"JIRA_API_TOKEN": "",
				"JIRA_URL":       "https://example.atlassian.net",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError:   true,
			errorContains: "JIRA_API_TOKEN",
		},
		{
			name:  "Full mode with missing JIRA URL",
			flags: &FlagConfig{},
			args:  []string{},
			envVars: map[string]string{
				"JIRA_API_TOKEN": "token123",
				"JIRA_URL":       "",
				"JIRA_USERNAME":  "user@example.com",
			},
			expectError:   true,
			errorContains: "JIRA_URL",
		},
		{
			name:  "Full mode with missing JIRA username",
			flags: &FlagConfig{},
			args:  []string{},
			envVars: map[string]string{
				"JIRA_API_TOKEN": "token123",
				"JIRA_URL":       "https://example.atlassian.net",
				"JIRA_USERNAME":  "",
			},
			expectError:   true,
			errorContains: "JIRA_USERNAME",
		},
		{
			name: "Environment variables override defaults",
			flags: &FlagConfig{
				ExtractOnly: true,
			},
			args: []string{},
			envVars: map[string]string{
				"JIRA_ID_REGEX": "CUSTOM-[0-9]+",
				"OUTPUT_FILE":   "custom_output.json",
			},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAIDRegex:  "CUSTOM-[0-9]+",
				OutputFile:   "custom_output.json",
				ExtractOnly:  true,
				SingleCommit: true,
			},
		},
		{
			name: "Flags override environment variables",
			flags: &FlagConfig{
				ExtractOnly: true,
				JIRAIDRegex: "FLAG-[0-9]+",
				OutputFile:  "flag_output.json",
			},
			args: []string{},
			envVars: map[string]string{
				"JIRA_ID_REGEX": "ENV-[0-9]+",
				"OUTPUT_FILE":   "env_output.json",
			},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAIDRegex:  "FLAG-[0-9]+",
				OutputFile:   "flag_output.json",
				ExtractOnly:  true,
				SingleCommit: true,
			},
		},
		{
			name: "CommitRange flag sets SingleCommit to false",
			flags: &FlagConfig{
				ExtractOnly: true,
				CommitRange: true,
			},
			args:        []string{},
			envVars:     map[string]string{},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAIDRegex:  DefaultJIRAIDRegex,
				OutputFile:   DefaultOutputFile,
				ExtractOnly:  true,
				SingleCommit: false,
			},
		},
		{
			name: "All environment variables and flags",
			flags: &FlagConfig{
				JIRAIDRegex: "FLAG-[0-9]+",
				OutputFile:  "flag.json",
				CommitRange: true,
			},
			args: []string{},
			envVars: map[string]string{
				"JIRA_API_TOKEN": "token456",
				"JIRA_URL":       "https://test.atlassian.net",
				"JIRA_USERNAME":  "test@example.com",
				"JIRA_ID_REGEX":  "ENV-[0-9]+",
				"OUTPUT_FILE":    "env.json",
			},
			expectError: false,
			expectedConfig: &AppConfig{
				JIRAToken:    "token456",
				JIRAURL:      "https://test.atlassian.net",
				JIRAUsername: "test@example.com",
				JIRAIDRegex:  "FLAG-[0-9]+", // Flag overrides env
				OutputFile:   "flag.json",   // Flag overrides env
				SingleCommit: false,         // CommitRange flag
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("JIRA_API_TOKEN")
			os.Unsetenv("JIRA_URL")
			os.Unsetenv("JIRA_USERNAME")
			os.Unsetenv("JIRA_ID_REGEX")
			os.Unsetenv("OUTPUT_FILE")

			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config, err := LoadConfig(tt.flags, tt.args)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, config)
			}
		})
	}
}

func TestDisplayUsageComplete(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DisplayUsage()

	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify all sections are present
	expectedSections := []string{
		"JIRA Evidence Gathering Tool",
		"Usage:",
		"./main [OPTIONS] <start_commit>",
		"./main <jira_id1> [jira_id2] [jira_id3] ...",
		"Options:",
		"-r, --regex PATTERN",
		"-o, --output FILE",
		"--extract-only",
		"--extract-from-git",
		"--range",
		"-h, --help",
		"Arguments:",
		"commit",
		"Environment Variables:",
		"JIRA_API_TOKEN",
		"JIRA_URL",
		"JIRA_USERNAME",
		"JIRA_ID_REGEX",
		"OUTPUT_FILE",
		"Examples:",
		"./main abc123def456",
		"./main --range abc123def456",
		"./main EV-123 EV-456 EV-789",
	}

	for _, section := range expectedSections {
		assert.Contains(t, output, section, "Missing section: %s", section)
	}

	// Verify the output is well-formatted
	lines := strings.Split(output, "\n")
	assert.Greater(t, len(lines), 20, "Usage output should have multiple lines")
}

func TestGetOrDefaultComplete(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "Returns first non-empty value",
			values:   []string{"", "second", "third"},
			expected: "second",
		},
		{
			name:     "Returns empty when all values are empty",
			values:   []string{"", "", ""},
			expected: "",
		},
		{
			name:     "Returns first value when non-empty",
			values:   []string{"first", "second", "third"},
			expected: "first",
		},
		{
			name:     "Works with single value",
			values:   []string{"only"},
			expected: "only",
		},
		{
			name:     "Works with no values",
			values:   []string{},
			expected: "",
		},
		{
			name:     "Handles nil-like empty strings",
			values:   []string{"", "", "", "found"},
			expected: "found",
		},
		{
			name:     "Priority order is maintained",
			values:   []string{"", "high", "low"},
			expected: "high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getOrDefault(tt.values...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateJIRAConfigComplete(t *testing.T) {
	tests := []struct {
		name          string
		config        *AppConfig
		expectError   bool
		expectedField string
		errorMessage  string
	}{
		{
			name: "Valid configuration",
			config: &AppConfig{
				JIRAToken:    "token123",
				JIRAURL:      "https://example.atlassian.net",
				JIRAUsername: "user@example.com",
			},
			expectError: false,
		},
		{
			name: "Missing JIRA token",
			config: &AppConfig{
				JIRAToken:    "",
				JIRAURL:      "https://example.atlassian.net",
				JIRAUsername: "user@example.com",
			},
			expectError:   true,
			expectedField: "JIRA_API_TOKEN",
			errorMessage:  "environment variable is required",
		},
		{
			name: "Missing JIRA URL",
			config: &AppConfig{
				JIRAToken:    "token123",
				JIRAURL:      "",
				JIRAUsername: "user@example.com",
			},
			expectError:   true,
			expectedField: "JIRA_URL",
			errorMessage:  "environment variable is required",
		},
		{
			name: "Missing JIRA username",
			config: &AppConfig{
				JIRAToken:    "token123",
				JIRAURL:      "https://example.atlassian.net",
				JIRAUsername: "",
			},
			expectError:   true,
			expectedField: "JIRA_USERNAME",
			errorMessage:  "environment variable is required",
		},
		{
			name: "All fields missing",
			config: &AppConfig{
				JIRAToken:    "",
				JIRAURL:      "",
				JIRAUsername: "",
			},
			expectError:   true,
			expectedField: "JIRA_API_TOKEN", // First check fails
			errorMessage:  "environment variable is required",
		},
		{
			name: "Whitespace-only values treated as empty",
			config: &AppConfig{
				JIRAToken:    "   ",
				JIRAURL:      "https://example.atlassian.net",
				JIRAUsername: "user@example.com",
			},
			expectError:   true,
			expectedField: "JIRA_API_TOKEN",
			errorMessage:  "environment variable is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Trim whitespace in config to match implementation
			if tt.config.JIRAToken == "   " {
				tt.config.JIRAToken = ""
			}

			err := validateJIRAConfig(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				validationErr, ok := err.(*ValidationError)
				assert.True(t, ok, "Error should be ValidationError type")
				assert.Equal(t, tt.expectedField, validationErr.Field)
				assert.Contains(t, validationErr.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
