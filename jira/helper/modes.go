package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// runExtractOnlyMode runs the tool in extract-only mode
func runExtractOnlyMode(config *AppConfig) error {
	git := NewGitService()

	fmt.Println("=== JIRA ID Extraction (Extract Only Mode) ===")
	if config.SingleCommit {
		fmt.Printf("Commit: %s\n", config.StartCommit)
	} else {
		fmt.Printf("Start Commit: %s\n", config.StartCommit)
	}
	fmt.Printf("JIRA ID Regex: %s\n", config.JIRAIDRegex)
	fmt.Println("")

	// Get branch info
	branchName, commitHash, currentJiraID, err := git.GetBranchInfo()
	if err != nil {
		return fmt.Errorf("failed to get branch info: %w", err)
	}

	fmt.Printf("Branch: %s\n", branchName)
	fmt.Printf("Latest Commit: %s\n", commitHash)

	// Validate HEAD
	if err := git.ValidateHEAD(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		return nil // Exit gracefully
	}

	// Extract JIRA IDs
	jiraIDs, err := git.ExtractJiraIDs(config.StartCommit, config.JIRAIDRegex, currentJiraID, config.SingleCommit)
	if err != nil {
		return fmt.Errorf("failed to extract JIRA IDs: %w", err)
	}

	if len(jiraIDs) == 0 {
		fmt.Println("No JIRA IDs found")
		return nil
	}

	// Output comma-separated JIRA IDs
	fmt.Println(strings.Join(jiraIDs, ","))
	return nil
}

// runLegacyExtractFromGit runs the legacy extract-from-git mode
func runLegacyExtractFromGit(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: ./main --extract-from-git <start_commit> <jira_id_regex>")
		return fmt.Errorf("insufficient arguments")
	}

	git := NewGitService()
	startCommit := args[0]
	regex := args[1]

	// Get branch info
	branchName, commitHash, currentJiraID, err := git.GetBranchInfo()
	if err != nil {
		return fmt.Errorf("error getting branch info: %v", err)
	}

	fmt.Printf("BRANCH_NAME: %s\n", branchName)
	fmt.Printf("JIRA ID: %s\n", currentJiraID)
	fmt.Printf("START_COMMIT: %s\n", commitHash)

	// Validate HEAD
	if err := git.ValidateHEAD(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		return nil // Exit gracefully
	}

	// Validate commit
	if err := git.ValidateCommit(startCommit); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		return nil // Exit gracefully
	}

	// Extract JIRA IDs
	jiraIDs, err := git.ExtractJiraIDs(startCommit, regex, currentJiraID, false)
	if err != nil {
		return fmt.Errorf("error extracting JIRA IDs: %v", err)
	}

	if len(jiraIDs) == 0 {
		fmt.Println("No JIRA IDs found")
		return nil
	}

	// Print comma-separated JIRA IDs
	fmt.Println(strings.Join(jiraIDs, ","))
	return nil
}

// runFullMode runs the complete JIRA evidence gathering process
func runFullMode(config *AppConfig) error {
	git := NewGitService()

	fmt.Println("=== JIRA Details Fetching Process ===")
	if config.SingleCommit {
		fmt.Printf("Commit: %s\n", config.StartCommit)
	} else {
		fmt.Printf("Start Commit: %s\n", config.StartCommit)
	}
	fmt.Printf("JIRA ID Regex: %s\n", config.JIRAIDRegex)
	fmt.Printf("Output File: %s\n", config.OutputFile)
	fmt.Println("")

	// Step 1: Extract JIRA IDs from git commits
	if config.SingleCommit {
		fmt.Println("Step 1: Extracting JIRA IDs from commit...")
	} else {
		fmt.Println("Step 1: Extracting JIRA IDs from git commits...")
	}

	// Get branch info
	branchName, commitHash, currentJiraID, err := git.GetBranchInfo()
	if err != nil {
		return fmt.Errorf("error getting branch info: %v", err)
	}

	// Display branch information
	fmt.Printf("Branch: %s\n", branchName)
	fmt.Printf("Latest Commit: %s\n", commitHash)

	// Validate HEAD
	if err := git.ValidateHEAD(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		return nil // Exit gracefully
	}

	// Extract JIRA IDs
	jiraIDs, err := git.ExtractJiraIDs(config.StartCommit, config.JIRAIDRegex, currentJiraID, config.SingleCommit)
	if err != nil {
		return fmt.Errorf("error extracting JIRA IDs: %v", err)
	}

	if len(jiraIDs) == 0 {
		fmt.Println("No JIRA IDs found in commit range")
		return nil
	}

	fmt.Printf("Found JIRA IDs: %s\n", strings.Join(jiraIDs, ", "))
	config.JIRAIDs = jiraIDs

	// Step 2: Fetch JIRA details
	fmt.Println("")
	fmt.Println("Step 2: Fetching JIRA details...")

	// Create JIRA client
	jiraClient, err := NewJiraClient()
	if err != nil {
		return fmt.Errorf("error creating JIRA client: %v", err)
	}

	// Process JIRA IDs and get results
	response := jiraClient.FetchJiraDetails(config.JIRAIDs)

	// Step 3: Write results to file
	fmt.Println("")
	fmt.Println("Step 3: Writing results...")

	if err := saveJiraResults(response, config); err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("=== Process completed successfully ===")
	return nil
}

// processDirectJiraIDs handles direct JIRA ID processing (no git operations)
func processDirectJiraIDs(config *AppConfig) error {
	fmt.Printf("Processing JIRA IDs: %s\n", strings.Join(config.JIRAIDs, ", "))

	// Create a new Jira client
	jiraClient, err := NewJiraClient()
	if err != nil {
		return fmt.Errorf("error creating JIRA client: %v", err)
	}

	// Get response
	response := jiraClient.FetchJiraDetails(config.JIRAIDs)

	// Save results to file using the same method as other modes
	if err := saveJiraResults(response, config); err != nil {
		return err
	}

	return nil
}

// saveJiraResults saves JIRA results to JSON
func saveJiraResults(response TransitionCheckResponse, config *AppConfig) error {
	// Save JSON
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	if err := writeToFile(config.OutputFile, jsonBytes); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Printf("JIRA data saved to: %s\n", config.OutputFile)

	return nil
}

// determineExecutionMode determines which mode to run based on flags and arguments
func determineExecutionMode(flags *FlagConfig, args []string, config *AppConfig) error {
	// Handle markdown generation mode
	if flags.GenerateMarkdown {
		return runMarkdownMode(flags)
	}

	// Handle legacy extract-from-git mode
	if flags.ExtractFromGit {
		return runLegacyExtractFromGit(args)
	}

	// Check if we have required arguments
	if len(args) == 0 {
		return fmt.Errorf("missing required arguments")
	}

	// Check if this is direct JIRA ID processing mode
	if !flags.ExtractOnly && len(args) > 0 {
		// Check if all arguments match JIRA ID pattern
		regex, err := regexp.Compile(config.JIRAIDRegex)
		if err == nil && allArgsMatchPattern(args, regex) {
			// All arguments are JIRA IDs - process them directly
			config.JIRAIDs = args
			return processDirectJiraIDs(config)
		}
	}

	// Otherwise, we're in git-based mode
	config.StartCommit = args[0]

	// Check if we're in a git repository
	git := NewGitService()
	if err := git.CheckRepository(); err != nil {
		return err
	}

	// Run the appropriate mode
	if config.ExtractOnly {
		return runExtractOnlyMode(config)
	}
	return runFullMode(config)
}

// allArgsMatchPattern checks if all arguments match the given regex pattern
func allArgsMatchPattern(args []string, regex *regexp.Regexp) bool {
	for _, arg := range args {
		if !regex.MatchString(arg) {
			return false
		}
	}
	return true
}

// runMarkdownMode runs the markdown generation mode
func runMarkdownMode(flags *FlagConfig) error {
	// Determine input and output files
	inputFile := getOrDefault(flags.OutputFile, os.Getenv("OUTPUT_FILE"), DefaultOutputFile)
	outputFile := getOrDefault(flags.MarkdownOutput, "transformed_jira_data.md")

	fmt.Println("=== Markdown Generation Mode ===")
	fmt.Printf("Input JSON file: %s\n", inputFile)
	fmt.Printf("Output Markdown file: %s\n", outputFile)
	fmt.Println("")

	// Generate markdown from JSON
	if err := GenerateMarkdownFromJSON(inputFile, outputFile); err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("=== Markdown generation completed successfully ===")
	return nil
}
