package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// GenerateMarkdownFromJSON reads a JSON file and generates markdown
func GenerateMarkdownFromJSON(inputFile string, outputFile string) error {
	// Read JSON file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	// Parse JSON into struct
	var response TransitionCheckResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Generate markdown
	markdown := generateMarkdown(response)

	// Write markdown to file
	err = os.WriteFile(outputFile, []byte(markdown), 0644)
	if err != nil {
		return fmt.Errorf("error writing markdown file: %v", err)
	}

	fmt.Printf("Markdown file generated: %s\n", outputFile)
	return nil
}

// generateMarkdown creates markdown content from JIRA data
func generateMarkdown(response TransitionCheckResponse) string {
	var sb strings.Builder

	// Header
	sb.WriteString("# JIRA Tasks Report\n\n")
	sb.WriteString(fmt.Sprintf("Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Total tasks: %d\n\n", len(response.Tasks)))

	// Summary table
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Key | Status | Type | Priority | Assignee |\n")
	sb.WriteString("|-----|--------|------|----------|----------|\n")

	for _, task := range response.Tasks {
		assignee := "Unassigned"
		if task.Assignee != nil && *task.Assignee != "" {
			assignee = *task.Assignee
		}
		// Use link from JSON if available
		keyDisplay := task.Key
		if task.Link != "" {
			keyDisplay = fmt.Sprintf("[%s](%s)", task.Key, task.Link)
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			keyDisplay, task.Status, task.Type, task.Priority, assignee))
	}
	sb.WriteString("\n")

	// Detailed task information
	sb.WriteString("## Task Details\n\n")

	for i, task := range response.Tasks {
		// Use link from JSON if available
		keyDisplay := task.Key
		if task.Link != "" {
			keyDisplay = fmt.Sprintf("[%s](%s)", task.Key, task.Link)
		}
		sb.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, keyDisplay))

		// Basic information
		sb.WriteString("**Basic Information:**\n")
		sb.WriteString(fmt.Sprintf("- **Status:** %s\n", task.Status))
		sb.WriteString(fmt.Sprintf("- **Type:** %s\n", task.Type))
		sb.WriteString(fmt.Sprintf("- **Project:** %s\n", task.Project))
		sb.WriteString(fmt.Sprintf("- **Priority:** %s\n", task.Priority))

		// People
		sb.WriteString("\n**People:**\n")
		assignee := "Unassigned"
		if task.Assignee != nil && *task.Assignee != "" {
			assignee = *task.Assignee
		}
		sb.WriteString(fmt.Sprintf("- **Assignee:** %s\n", assignee))
		sb.WriteString(fmt.Sprintf("- **Reporter:** %s\n", task.Reporter))

		// Dates
		sb.WriteString("\n**Dates:**\n")
		sb.WriteString(fmt.Sprintf("- **Created:** %s\n", formatDate(task.Created)))
		sb.WriteString(fmt.Sprintf("- **Updated:** %s\n", formatDate(task.Updated)))

		// Description
		if task.Description != "" {
			sb.WriteString("\n**Description:**\n")
			sb.WriteString(fmt.Sprintf("> %s\n", strings.ReplaceAll(task.Description, "\n", "\n> ")))
		}

		// Transitions
		if len(task.Transitions) > 0 {
			sb.WriteString("\n**Transition History:**\n\n")
			sb.WriteString("| From Status | To Status | Author | Date |\n")
			sb.WriteString("|-------------|-----------|--------|------|\n")

			for _, transition := range task.Transitions {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					transition.FromStatus,
					transition.ToStatus,
					transition.Author,
					formatDate(transition.TransitionTime)))
			}
		}

		sb.WriteString("\n---\n\n")
	}

	// Status distribution
	statusCount := make(map[string]int)
	for _, task := range response.Tasks {
		statusCount[task.Status]++
	}

	sb.WriteString("## Status Distribution\n\n")
	sb.WriteString("| Status | Count |\n")
	sb.WriteString("|--------|-------|\n")
	for status, count := range statusCount {
		sb.WriteString(fmt.Sprintf("| %s | %d |\n", status, count))
	}

	return sb.String()
}

// formatDate formats a JIRA date string to a more readable format
func formatDate(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}

	// Try to parse the JIRA date format
	t, err := time.Parse(JiraTimeFormat, dateStr)
	if err != nil {
		// If parsing fails, return the original string
		return dateStr
	}

	// Return in a more readable format
	return t.Format("2006-01-02 15:04:05")
}
