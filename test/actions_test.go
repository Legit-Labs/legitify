package test

import (
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/go-github/v44/github"
)

type organizationActionsMockConfiguration struct {
	allowedActions      *string
	enabledRepositories *string
}

func newOrganizationActionsMock(config organizationActionsMockConfiguration) githubcollected.OrganizationActions {
	return githubcollected.OrganizationActions{
		Organization: defaultOrg,
		ActionsPermissions: &github.ActionsPermissions{
			EnabledRepositories: config.enabledRepositories,
			AllowedActions:      config.allowedActions,
		},
	}
}

func TestActions(t *testing.T) {
	all := "all"
	selected := "selected"
	tests := []struct {
		name             string
		policyName       string
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
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newOrganizationActionsMock(test.args),
			namespace.Actions, test.policyName, test.shouldBeViolated)
	}
}
