package main

import (
	"context"
	"fmt"
	"os"

	jira "github.com/andygrunwald/go-jira/v2/cloud"
)

// JiraClient wraps the JIRA client and provides methods for JIRA operations
type JiraClient struct {
	client  *jira.Client
	baseURL string
}

// NewJiraClient creates a new JIRA client with authentication
func NewJiraClient() (*JiraClient, error) {
	jiraToken := os.Getenv("JIRA_API_TOKEN")
	if jiraToken == "" {
		return nil, &ValidationError{Field: "JIRA_API_TOKEN", Value: "", Err: fmt.Errorf("environment variable not found")}
	}

	jiraURL := os.Getenv("JIRA_URL")
	if jiraURL == "" {
		return nil, &ValidationError{Field: "JIRA_URL", Value: "", Err: fmt.Errorf("environment variable not found")}
	}

	jiraUsername := os.Getenv("JIRA_USERNAME")
	if jiraUsername == "" {
		return nil, &ValidationError{Field: "JIRA_USERNAME", Value: "", Err: fmt.Errorf("environment variable not found")}
	}

	// Create JIRA client with basic auth transport
	tp := jira.BasicAuthTransport{
		Username: jiraUsername,
		APIToken: jiraToken,
	}

	client, err := jira.NewClient(jiraURL, tp.Client())
	if err != nil {
		return nil, fmt.Errorf("failed to create JIRA client: %w", err)
	}

	return &JiraClient{
		client:  client,
		baseURL: jiraURL,
	}, nil
}

// FetchJiraDetails fetches JIRA details sequentially
func (jc *JiraClient) FetchJiraDetails(jiraIDs []string) TransitionCheckResponse {
	response := TransitionCheckResponse{
		Tasks: make([]JiraTransitionResult, 0, len(jiraIDs)),
	}

	for _, jiraID := range jiraIDs {
		result := jc.fetchSingleJiraDetail(jiraID)
		response.Tasks = append(response.Tasks, result)
	}

	return response
}

// fetchSingleJiraDetail fetches details for a single JIRA ID
func (jc *JiraClient) fetchSingleJiraDetail(jiraID string) JiraTransitionResult {
	issue, _, err := jc.client.Issue.Get(context.Background(), jiraID, &jira.GetQueryOptions{Expand: "changelog"})

	if err != nil || issue == nil || issue.Fields == nil {
		return jc.createErrorResult(jiraID, err)
	}

	return jc.createSuccessResult(issue)
}

// createErrorResult creates an error result for a failed JIRA fetch
func (jc *JiraClient) createErrorResult(jiraID string, err error) JiraTransitionResult {
	errorMsg := "Error: Could not retrieve issue"
	if err != nil {
		errorMsg = fmt.Sprintf("Error: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to fetch JIRA %s: %v\n", jiraID, err)
	}

	return JiraTransitionResult{
		Key:         jiraID,
		Link:        "", // No link for error results
		Status:      ErrorStatus,
		Description: errorMsg,
		Type:        ErrorType,
		Project:     "",
		Created:     "",
		Updated:     "",
		Assignee:    nil,
		Reporter:    "",
		Priority:    "",
		Transitions: []Transition{},
	}
}

// createSuccessResult creates a result from a successfully fetched JIRA issue
func (jc *JiraClient) createSuccessResult(issue *jira.Issue) JiraTransitionResult {
	// Create the JIRA link
	link := ""
	if jc.baseURL != "" {
		link = fmt.Sprintf("%s/browse/%s", jc.baseURL, issue.Key)
	}

	result := JiraTransitionResult{
		Key:         issue.Key,
		Link:        link,
		Status:      getStatusName(issue.Fields.Status),
		Description: getDescription(issue.Fields.Description),
		Type:        getIssueTypeName(issue.Fields.Type),
		Project:     getProjectKey(issue.Fields.Project),
		Created:     getTimeAsString(issue.Fields.Created),
		Updated:     getTimeAsString(issue.Fields.Updated),
		Assignee:    getAssignee(issue.Fields.Assignee),
		Reporter:    getReporterName(issue.Fields.Reporter),
		Priority:    getPriorityName(issue.Fields.Priority),
		Transitions: jc.extractTransitions(issue),
	}

	return result
}

// extractTransitions extracts status transitions from issue changelog
func (jc *JiraClient) extractTransitions(issue *jira.Issue) []Transition {
	var transitions []Transition

	if issue.Changelog == nil || len(issue.Changelog.Histories) == 0 {
		return transitions
	}

	for _, history := range issue.Changelog.Histories {
		for _, item := range history.Items {
			if item.Field == "status" {
				transition := Transition{
					FromStatus:     item.FromString,
					ToStatus:       item.ToString,
					Author:         history.Author.DisplayName,
					AuthorEmail:    history.Author.EmailAddress,
					TransitionTime: history.Created,
				}
				transitions = append(transitions, transition)
			}
		}
	}

	return transitions
}
