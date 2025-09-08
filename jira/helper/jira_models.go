package main

// Constants for JIRA operations
const (
	JiraTimeFormat = "2006-01-02T15:04:05.000-0700"
	ErrorStatus    = "Error"
	ErrorType      = "Error"
)

/*
    TransitionCheckResponse is the json formatted predicate that will be returned to the calling build process for creating an evidence
    its structure should be:

    {
        "tasks": [
            {
                "key": "EV-1",
                "status": "QA in Progress",
                "description": "<description text>",
                "type": "Task",
                "project": "EV",
                "created": "2020-01-01T12:11:56.063+0530",
                "updated": "2020-01-01T12:12:01.876+0530",
                "assignee": "<assignee name>",
                "reporter": "<reporter name>",
                "priority": "Medium",
                "transitions": [
                    {
                        "from_status": "To Do",
                        "to_status": "In Progress",
                        "author": "<>author name>",
                        "author_user_name": "<author email>",
                        "transition_time": "2020-07-28T16:39:54.620+0530"
                    }
                ]
            },
            {
                "key": "EV-2",
                "status": "Error",
                "description": "Error: Could not retrieve issue",
                "type": "Error",
                "project": "",
                "created": "",
                "updated": "",
                "assignee": null,
                "reporter": "",
                "priority": "",
                "transitions": []
            }
        ]
    }

   notice that the calling client should first check that return value was 0 before using the response JSON,
   otherwise the response is an error message which cannot be parsed
*/

type TransitionCheckResponse struct {
	Tasks []JiraTransitionResult `json:"tasks"`
}

type JiraTransitionResult struct {
	Key         string       `json:"key"`
	Link        string       `json:"link,omitempty"`
	Status      string       `json:"status"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Project     string       `json:"project"`
	Created     string       `json:"created"`
	Updated     string       `json:"updated"`
	Assignee    *string      `json:"assignee"`
	Reporter    string       `json:"reporter"`
	Priority    string       `json:"priority"`
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	FromStatus     string `json:"from_status"`
	ToStatus       string `json:"to_status"`
	Author         string `json:"author"`
	AuthorEmail    string `json:"author_user_name"`
	TransitionTime string `json:"transition_time"`
}
