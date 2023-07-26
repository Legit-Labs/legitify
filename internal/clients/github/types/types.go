package types

import "encoding/json"

type TokenPermissions struct {
	DefaultWorkflowPermissions   *string `json:"default_workflow_permissions,omitempty"`
	CanApprovePullRequestReviews *bool   `json:"can_approve_pull_request_reviews,omitempty"`
}

type RepositoryRule struct {
	Type       string           `json:"type"`
	Parameters *json.RawMessage `json:"parameters,omitempty"`
}
