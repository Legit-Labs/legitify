package test

import (
	"github.com/Legit-Labs/legitify/internal/clients/github/types"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/go-github/v53/github"
)

type organizationActionsMockConfiguration struct {
	allowedActions         *string
	enabledRepositories    *string
	tokenDefaultPermission string
	workflowsCanApprovePRs bool
}

func newOrganizationActionsMock(config organizationActionsMockConfiguration) githubcollected.OrganizationActions {
	return githubcollected.OrganizationActions{
		Organization: defaultOrg,
		ActionsPermissions: &github.ActionsPermissions{
			EnabledRepositories: config.enabledRepositories,
			AllowedActions:      config.allowedActions,
		},
		TokenPermissions: &types.TokenPermissions{
			DefaultWorkflowPermissions:   &config.tokenDefaultPermission,
			CanApprovePullRequestReviews: &config.workflowsCanApprovePRs,
		},
	}
}

func TestActions(t *testing.T) {
	all := "all"
	selected := "selected"
	tests := []struct {
		name             string
		policyName       string
		scmType          scm_type.ScmType
		shouldBeViolated bool
		args             organizationActionsMockConfiguration
	}{
		{
			name:             "all github actions are allowed to run",
			policyName:       "all_github_actions_are_allowed",
			shouldBeViolated: true,
			args: organizationActionsMockConfiguration{
				allowedActions: &all,
			},
		},
		{
			name:             "not all github actions are allowed to run",
			policyName:       "all_github_actions_are_allowed",
			shouldBeViolated: false,
			args: organizationActionsMockConfiguration{
				allowedActions: &selected,
			},
		},
		{
			name:             "all repositories can run GitHub actions",
			policyName:       "all_repositories_can_run_github_actions",
			shouldBeViolated: true,
			args: organizationActionsMockConfiguration{
				enabledRepositories: &all,
			},
		},
		{
			name:             "not all repositories can run GitHub actions",
			policyName:       "all_repositories_can_run_github_actions",
			shouldBeViolated: false,
			args: organizationActionsMockConfiguration{
				enabledRepositories: &selected,
			},
		},
		{
			name:             "actions can approve pull requests",
			policyName:       "actions_can_approve_pull_requests",
			shouldBeViolated: true,
			args: organizationActionsMockConfiguration{
				enabledRepositories:    &selected,
				workflowsCanApprovePRs: true,
			},
		},
		{
			name:             "actions can not approve pull requests",
			policyName:       "actions_can_approve_pull_requests",
			shouldBeViolated: false,
			args: organizationActionsMockConfiguration{
				enabledRepositories:    &selected,
				workflowsCanApprovePRs: false,
			},
		},
		{
			name:             "workflow token default permissions is not set to read only",
			policyName:       "token_default_permissions_is_read_write",
			shouldBeViolated: true,
			args: organizationActionsMockConfiguration{
				enabledRepositories:    &selected,
				tokenDefaultPermission: "write",
			},
		},
		{
			name:             "workflow token default permissions is set to read only",
			policyName:       "token_default_permissions_is_read_write",
			shouldBeViolated: false,
			args: organizationActionsMockConfiguration{
				enabledRepositories:    &selected,
				tokenDefaultPermission: "read",
			},
		},
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newOrganizationActionsMock(test.args),
			namespace.Actions, test.policyName, test.shouldBeViolated, scm_type.GitHub)
	}
}
