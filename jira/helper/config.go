package main

import (
	"flag"
	"fmt"
	"os"
)

// Constants for default values
const (
	DefaultJIRAIDRegex = "[A-Z]+-[0-9]+"
	DefaultOutputFile  = "transformed_jira_data.json"
)

// AppConfig holds all configuration for the application
type AppConfig struct {
	// JIRA Configuration
	JIRAToken    string
	JIRAURL      string
	JIRAUsername string
	JIRAIDRegex  string

	// Output Configuration
	OutputFile string

	// Runtime Configuration
	ExtractOnly    bool
	ExtractFromGit bool
	SingleCommit   bool
	StartCommit    string
	JIRAIDs        []string
}

// FlagConfig holds command line flags
type FlagConfig struct {
	JIRAIDRegex      string
	OutputFile       string
	ExtractOnly      bool
	ExtractFromGit   bool
	CommitRange      bool
	Help             bool
	HelpLong         bool
	GenerateMarkdown bool
	MarkdownOutput   string
}

// ParseFlags parses command line flags
func ParseFlags() (*FlagConfig, []string) {
	flags := &FlagConfig{}
	flag.StringVar(&flags.JIRAIDRegex, "r", "", "JIRA ID regex pattern")
	flag.StringVar(&flags.OutputFile, "o", "", "Output file for JIRA data")
	flag.BoolVar(&flags.ExtractOnly, "extract-only", false, "Only extract JIRA IDs, don't fetch details")
	flag.BoolVar(&flags.ExtractFromGit, "extract-from-git", false, "Extract JIRA IDs from git commits (legacy mode)")
	flag.BoolVar(&flags.CommitRange, "range", false, "Process commits from the specified commit to HEAD (instead of single commit)")
	flag.BoolVar(&flags.Help, "h", false, "Display help message")
	flag.BoolVar(&flags.HelpLong, "help", false, "Display help message")
	flag.BoolVar(&flags.GenerateMarkdown, "markdown", false, "Generate markdown from existing JSON file")
	flag.StringVar(&flags.MarkdownOutput, "markdown-output", "", "Output file for markdown (default: transformed_jira_data.md)")
	flag.Parse()

	return flags, flag.Args()
}

// LoadConfig loads configuration from flags and environment variables
func LoadConfig(flags *FlagConfig, args []string) (*AppConfig, error) {
	config := &AppConfig{
		JIRAIDRegex:    getOrDefault(flags.JIRAIDRegex, os.Getenv("JIRA_ID_REGEX"), DefaultJIRAIDRegex),
		OutputFile:     getOrDefault(flags.OutputFile, os.Getenv("OUTPUT_FILE"), DefaultOutputFile),
		ExtractOnly:    flags.ExtractOnly,
		ExtractFromGit: flags.ExtractFromGit,
		SingleCommit:   !flags.CommitRange, // Default to single commit unless --range is specified
	}

	// Load JIRA credentials only if not in extract-only mode or markdown mode
	if !config.ExtractOnly && !config.ExtractFromGit && !flags.GenerateMarkdown {
		config.JIRAToken = os.Getenv("JIRA_API_TOKEN")
		config.JIRAURL = os.Getenv("JIRA_URL")
		config.JIRAUsername = os.Getenv("JIRA_USERNAME")

		// Validate JIRA configuration
		if err := validateJIRAConfig(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// validateJIRAConfig validates JIRA-related configuration
func validateJIRAConfig(config *AppConfig) error {
	if config.JIRAToken == "" {
		return &ValidationError{Field: "JIRA_API_TOKEN", Value: "", Err: fmt.Errorf("environment variable is required")}
	}
	if config.JIRAURL == "" {
		return &ValidationError{Field: "JIRA_URL", Value: "", Err: fmt.Errorf("environment variable is required")}
	}
	if config.JIRAUsername == "" {
		return &ValidationError{Field: "JIRA_USERNAME", Value: "", Err: fmt.Errorf("environment variable is required")}
	}
	return nil
}

// getOrDefault gets value with defaults
func getOrDefault(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// DisplayUsage shows the usage information
func DisplayUsage() {
	fmt.Println("JIRA Evidence Gathering Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  ./main [OPTIONS] <start_commit>")
	fmt.Println("  ./main <jira_id1> [jira_id2] [jira_id3] ...")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -r, --regex PATTERN    JIRA ID regex pattern (default: '[A-Z]+-[0-9]+')")
	fmt.Println("  -o, --output FILE      Output file for JIRA data (default: transformed_jira_data.json)")
	fmt.Println("  --extract-only         Only extract JIRA IDs, don't fetch details")
	fmt.Println("  --extract-from-git     Extract JIRA IDs from git commits (legacy mode)")
	fmt.Println("  --range                Process commits from the specified commit to HEAD (instead of single commit)")
	fmt.Println("  --markdown             Generate markdown from existing JSON file")
	fmt.Println("  --markdown-output FILE Output file for markdown (default: transformed_jira_data.md)")
	fmt.Println("  -h, --help             Display this help message")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("  commit                 The commit to process (default: process only this commit)")
	fmt.Println("                         With --range: Starting commit hash (excluded from evidence filter)")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  JIRA_API_TOKEN         JIRA API token")
	fmt.Println("  JIRA_URL              JIRA instance URL")
	fmt.Println("  JIRA_USERNAME         JIRA username")
	fmt.Println("  JIRA_ID_REGEX         JIRA ID regex pattern (can be overridden with -r)")
	fmt.Println("  OUTPUT_FILE           Output file path (can be overridden with -o)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  ./main abc123def456                   # Process only commit abc123def456")
	fmt.Println("  ./main --range abc123def456           # Process commits from abc123def456 to HEAD")
	fmt.Println("  ./main -r 'EV-\\d+' -o jira_results.json abc123def456")
	fmt.Println("  ./main --extract-only abc123def456")
	fmt.Println("  ./main EV-123 EV-456 EV-789         # Direct JIRA ticket processing")
	fmt.Println("  ./main --markdown                    # Generate markdown from transformed_jira_data.json")
	fmt.Println("  ./main --markdown --markdown-output report.md  # Generate markdown with custom output file")
}
