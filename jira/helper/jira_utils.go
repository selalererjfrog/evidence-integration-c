package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
)

// Field extractors for JIRA objects - all handle nil values safely

func getStatusName(status *jira.Status) string {
	if status == nil {
		return ""
	}
	return status.Name
}

func getIssueTypeName(issueType jira.IssueType) string {
	return issueType.Name
}

func getProjectKey(project jira.Project) string {
	return project.Key
}

func getReporterName(reporter *jira.User) string {
	if reporter == nil {
		return ""
	}
	return reporter.DisplayName
}

func getPriorityName(priority *jira.Priority) string {
	if priority == nil {
		return ""
	}
	return priority.Name
}

func getAssignee(assignee *jira.User) *string {
	if assignee == nil {
		return nil
	}
	return &assignee.DisplayName
}

// getTimeAsString converts various time representations to string format
func getTimeAsString(timeField interface{}) string {
	if timeField == nil {
		return ""
	}

	switch v := timeField.(type) {
	case string:
		return v
	case time.Time:
		return v.Format(JiraTimeFormat)
	case *time.Time:
		if v != nil {
			return v.Format(JiraTimeFormat)
		}
		return ""
	default:
		// Try JSON marshaling as last resort
		if jsonBytes, err := json.Marshal(timeField); err == nil {
			var timeStr string
			if json.Unmarshal(jsonBytes, &timeStr) == nil && timeStr != "" {
				return timeStr
			}
		}
		// Final fallback
		return fmt.Sprintf("%v", timeField)
	}
}

// getDescription extracts description text from JIRA description field
func getDescription(desc interface{}) string {
	if desc == nil {
		return ""
	}

	// Handle the Atlassian Document Format (ADF) structure
	descMap, ok := desc.(map[string]interface{})
	if !ok {
		// Fallback to string representation
		return fmt.Sprintf("%v", desc)
	}

	content, ok := descMap["content"].([]interface{})
	if !ok {
		return fmt.Sprintf("%v", desc)
	}

	var result strings.Builder
	for _, item := range content {
		text := extractTextFromADFNode(item)
		if text != "" {
			result.WriteString(text)
		}
	}

	if result.Len() == 0 {
		return fmt.Sprintf("%v", desc)
	}
	return result.String()
}

// extractTextFromADFNode extracts text from an ADF node (paragraph, text, etc.)
func extractTextFromADFNode(node interface{}) string {
	nodeMap, ok := node.(map[string]interface{})
	if !ok {
		return ""
	}

	nodeType, _ := nodeMap["type"].(string)

	switch nodeType {
	case "paragraph":
		// Extract text from paragraph's content
		content, ok := nodeMap["content"].([]interface{})
		if !ok {
			return ""
		}

		var texts []string
		for _, item := range content {
			if text := extractTextFromADFNode(item); text != "" {
				texts = append(texts, text)
			}
		}
		return strings.Join(texts, "")

	case "text":
		// Direct text node
		text, _ := nodeMap["text"].(string)
		return text

	default:
		// Handle other node types if needed
		return ""
	}
}
