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

type AnalysisAndSecurityPolicies struct {
	AdvancedSecurityEnabledForNewRepositories      bool   `json:"advanced_security_enabled_for_new_repositories"`
	DependabotAlertsEnabledForNewRepositories      bool   `json:"dependabot_alerts_enabled_for_new_repositories"`
	SecretScanningEnabledForNewRepositories        bool   `json:"secret_scanning_enabled_for_new_repositories"`
	SecretScanningPushProtectionEnabledForNewRepos bool   `json:"secret_scanning_push_protection_enabled_for_new_repositories"`
	SecretScanningPushProtectionCustomLink         string `json:"secret_scanning_push_protection_custom_link"`
} 
