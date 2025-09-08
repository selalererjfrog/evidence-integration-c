package main

import "fmt"

// GitError represents errors from git operations
type GitError struct {
	Operation string
	Err       error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git operation '%s' failed: %v", e.Operation, e.Err)
}

// ValidationError represents validation errors
type ValidationError struct {
	Field string
	Value string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s='%s': %v", e.Field, e.Value, e.Err)
}
