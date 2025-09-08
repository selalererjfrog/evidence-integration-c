package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GitService handles all git operations
type GitService struct {
	execCommand func(args ...string) (string, error)
}

// NewGitService creates a new git service
func NewGitService() *GitService {
	return &GitService{
		execCommand: defaultGitCommand,
	}
}

// defaultGitCommand executes a git command and returns the output
func defaultGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", &GitError{Operation: strings.Join(args, " "), Err: err}
	}
	return strings.TrimSpace(string(output)), nil
}

// GetBranchInfo returns current branch name, latest commit hash, and JIRA ID from latest commit
func (g *GitService) GetBranchInfo() (string, string, string, error) {
	// Get current branch
	branchName, err := g.execCommand("branch", "--show-current")
	if err != nil {
		return "", "", "", err
	}

	// Get commit hash and subject in one command to reduce git calls
	commitOutput, err := g.execCommand("log", "-1", "--format=%H%n%s")
	if err != nil {
		return "", "", "", err
	}

	lines := strings.Split(commitOutput, "\n")
	if len(lines) < 2 {
		return "", "", "", &GitError{Operation: "log -1", Err: fmt.Errorf("unexpected output format")}
	}

	commitHash := lines[0]
	subject := lines[1]

	// Extract JIRA ID using default pattern
	jiraID := extractFirstJIRAID(subject, DefaultJIRAIDRegex)

	return branchName, commitHash, jiraID, nil
}

// ValidateCommit checks if a commit exists in the repository
func (g *GitService) ValidateCommit(commit string) error {
	// First validate the commit hash format
	if err := validateCommitHash(commit); err != nil {
		return err
	}

	if _, err := g.execCommand("rev-parse", "--verify", commit); err != nil {
		return &GitError{Operation: "rev-parse --verify", Err: fmt.Errorf("commit '%s' not found", commit)}
	}
	return nil
}

// ValidateHEAD checks if HEAD commit exists in the repository
func (g *GitService) ValidateHEAD() error {
	if _, err := g.execCommand("rev-parse", "--verify", "HEAD"); err != nil {
		return &GitError{Operation: "rev-parse --verify HEAD", Err: fmt.Errorf("repository may be empty or corrupted")}
	}
	return nil
}

// ExtractJiraIDs extracts JIRA IDs from git commit messages
func (g *GitService) ExtractJiraIDs(startCommit, jiraIDRegex, currentJiraID string, singleCommit bool) ([]string, error) {
	// Validate commit first
	if err := g.ValidateCommit(startCommit); err != nil {
		return nil, err
	}

	var output string
	var err error

	if singleCommit {
		// Get only the specified commit message
		output, err = g.execCommand("log", "-1", "--pretty=format:%s", startCommit)
		if err != nil {
			return nil, err
		}
	} else {
		// Get commit messages from startCommit to HEAD (original behavior)
		output, err = g.execCommand("log", "--pretty=format:%s", startCommit+"..HEAD")
		if err != nil {
			return nil, err
		}
	}

	// Parse regex
	regex, err := regexp.Compile(jiraIDRegex)
	if err != nil {
		return nil, &ValidationError{Field: "jira_id_regex", Value: jiraIDRegex, Err: err}
	}

	// Extract unique JIRA IDs
	// In single commit mode, don't add currentJiraID from branch
	jiraIDToAdd := currentJiraID
	if singleCommit {
		jiraIDToAdd = ""
	}
	uniqueIDs := extractUniqueJIRAIDs(output, jiraIDToAdd, regex)

	if len(uniqueIDs) == 0 {
		if singleCommit {
			fmt.Fprintf(os.Stderr, "⚠️  No JIRA IDs found in commit %s\n", startCommit)
		} else {
			fmt.Fprintf(os.Stderr, "⚠️  No JIRA IDs found in commit range %s..HEAD\n", startCommit)
		}
	}

	return uniqueIDs, nil
}

// CheckRepository checks if we're in a git repository
func (g *GitService) CheckRepository() error {
	if _, err := g.execCommand("rev-parse", "--git-dir"); err != nil {
		return &GitError{Operation: "rev-parse --git-dir", Err: fmt.Errorf("not in a git repository")}
	}
	return nil
}

// validateCommitHash validates that a commit hash looks valid
func validateCommitHash(hash string) error {
	if hash == "" {
		return &ValidationError{Field: "commit", Value: hash, Err: fmt.Errorf("cannot be empty")}
	}

	// Basic validation - should be hex characters (allowing short hashes)
	validHex := regexp.MustCompile("^[a-fA-F0-9]+$")
	if !validHex.MatchString(hash) {
		return &ValidationError{Field: "commit", Value: hash, Err: fmt.Errorf("invalid format")}
	}

	return nil
}

// extractFirstJIRAID extracts the first JIRA ID from a string
func extractFirstJIRAID(text, pattern string) string {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}

	matches := regex.FindAllString(text, -1)
	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}

// extractUniqueJIRAIDs extracts unique JIRA IDs from commit messages
func extractUniqueJIRAIDs(commitMessages, currentJiraID string, regex *regexp.Regexp) []string {
	jiraIDs := make(map[string]bool)

	// Add current JIRA ID if it matches the pattern
	if currentJiraID != "" && regex.MatchString(currentJiraID) {
		jiraIDs[currentJiraID] = true
	}

	// Extract from commit messages
	lines := strings.Split(commitMessages, "\n")
	for _, line := range lines {
		matches := regex.FindAllString(line, -1)
		for _, match := range matches {
			jiraIDs[match] = true
		}
	}

	// Convert map to slice
	var result []string
	for jiraID := range jiraIDs {
		if jiraID != "" {
			result = append(result, jiraID)
		}
	}

	return result
}
