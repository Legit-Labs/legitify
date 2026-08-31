package test

import (
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/google/go-github/v53/github"
	"testing"
	"time"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	gitlabcollected "github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	gitlab2 "github.com/xanzy/go-gitlab"
)

type organizationMockConfiguration struct {
	config     map[string]interface{}
	ssoEnabled *bool
	name       string
	url        string
	secrets    []*githubcollected.OrganizationSecret
	org        *github.Organization
}

func newOrganizationMock(config organizationMockConfiguration) githubcollected.Organization {
	samlEnabledMockResult := false
	if config.ssoEnabled != nil {
		samlEnabledMockResult = *config.ssoEnabled
	}
	var hooks []*github.Hook
	if config.config != nil {
		hook := github.Hook{
			Config: config.config,
			Name:   &config.name,
			URL:    &config.url,
		}
		hooks = append(hooks, &hook)
	}
	var orgSecrets []*githubcollected.OrganizationSecret = nil
	if config.secrets != nil {
		orgSecrets = append(orgSecrets, config.secrets...)
	}

	var org *githubcollected.ExtendedOrg
	if config.org != nil {
		extended := githubcollected.NewExtendedOrg(config.org, permissions.OrgRoleNone)
		org = &extended
	}

	return githubcollected.Organization{
		Organization: org,
		SamlEnabled:  &samlEnabledMockResult,
		Hooks:        hooks,
		OrgSecrets:   orgSecrets,
	}
}

func TestOrganization(t *testing.T) {
	boolTrue := true
	boolFalse := false

	tests := []struct {
		name             string
		policyName       string
		shouldBeViolated bool
		args             organizationMockConfiguration
	}{
		{
			name:             "webhook secured",
			policyName:       "organization_webhook_no_secret",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				config: map[string]interface{}{
					"insecure_ssl": "0",
					"secret":       "123",
				},
				name: "test",
				url:  "test",
			},
		},
		{
			name:             "Violate webhook policy no configuration",
			policyName:       "organization_webhook_no_secret",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				config: map[string]interface{}{},
				name:   "test",
				url:    "test",
			},
		},
		{
			name:             "Violate webhook policy no secret",
			policyName:       "organization_webhook_no_secret",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				config: map[string]interface{}{
					"insecure_ssl": "0", // no secret
				},
				name: "test",
				url:  "test",
			},
		},
		{
			name:             "Violate webhook policy insecure configuration",
			policyName:       "organization_webhook_doesnt_require_ssl",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				config: map[string]interface{}{
					"insecure_ssl": "1",
					"secret":       "123",
				},
				name: "test",
				url:  "test",
			},
		},
		// -- SSO tests
		{
			name:             "SSO should be disabled",
			policyName:       "organization_not_using_single_sign_on",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				ssoEnabled: &boolFalse,
			},
		},
		{
			name:             "SSO should be enabled",
			policyName:       "organization_not_using_single_sign_on",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				ssoEnabled: &boolTrue,
			},
		},
		{
			name:             "Organization has no stale secrets",
			policyName:       "organization_secret_is_stale",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				secrets: []*githubcollected.OrganizationSecret{
					{
						Name:      "test1",
						UpdatedAt: int(time.Now().UnixNano()) - 2628000000000000, // one month
					},
					{
						Name:      "test2",
						UpdatedAt: int(time.Now().UnixNano()) - (3 * 2628000000000000), // three months
					},
					{
						Name:      "test3",
						UpdatedAt: int(time.Now().UnixNano()) - (6 * 2628000000000000), // six month
					},
				},
			},
		},
		{
			name:             "Organization has stale secrets",
			policyName:       "organization_secret_is_stale",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				secrets: []*githubcollected.OrganizationSecret{
					{
						Name:      "test1",
						UpdatedAt: 1652020546000000000, //08.05.2022
					},
					{
						Name:      "test2",
						UpdatedAt: 957796546000000000, //08.05.2000
					},
					{
						Name:      "test3",
						UpdatedAt: int(time.Now().UnixNano()),
					},
				},
			},
		},
		{
			name:             "Organization has no secrets",
			policyName:       "organization_secret_is_stale",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				secrets: nil,
			},
		},
		{
			name:             "two factor authentication required for the organization",
			policyName:       "two_factor_authentication_not_required_for_org",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				org: &github.Organization{
					TwoFactorRequirementEnabled: github.Bool(true),
				},
			},
		},
		{
			name:             "two factor authentication not required for the organization",
			policyName:       "two_factor_authentication_not_required_for_org",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				org: &github.Organization{
					TwoFactorRequirementEnabled: github.Bool(false),
				},
			},
		},
		{
			name:             "default repository permission is none",
			policyName:       "default_repository_permission_is_not_none",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				org: &github.Organization{
					DefaultRepoPermission: github.String("none"),
				},
			},
		},
		{
			name:             "default repository permission grants read to every member",
			policyName:       "default_repository_permission_is_not_none",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				org: &github.Organization{
					DefaultRepoPermission: github.String("read"),
				},
			},
		},
		{
			name:             "only admins can create public repositories",
			policyName:       "non_admins_can_create_public_repositories",
			shouldBeViolated: false,
			args: organizationMockConfiguration{
				org: &github.Organization{
					MembersCanCreatePublicRepos: github.Bool(false),
				},
			},
		},
		{
			name:             "non admins can create public repositories",
			policyName:       "non_admins_can_create_public_repositories",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				org: &github.Organization{
					MembersCanCreatePublicRepos: github.Bool(true),
				},
			},
		},
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newOrganizationMock(test.args),
			namespace.Organization, test.policyName, test.shouldBeViolated, scm_type.GitHub)
	}
}

func gitlabGroupTestTemplate(t *testing.T, name string, mockData interface{}, testedPolicyName string, expectFailure bool) {
	PolicyTestTemplate(t, name, mockData, namespace.Organization, testedPolicyName, expectFailure, scm_type.GitLab)
}

func TestGitlabGroupTwoFactorRequired(t *testing.T) {
	name := "Group Should Require Two-Factor Authentication"
	testedPolicyName := "two_factor_authentication_not_required_for_group"

	makeMockData := func(flag bool) gitlabcollected.Organization {
		return gitlabcollected.Organization{Group: &gitlab2.Group{RequireTwoFactorAuth: flag}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		gitlabGroupTestTemplate(t, name, makeMockData(options[expectFailure]), testedPolicyName, expectFailure)
	}
}

func TestGitlabGroupForkingOutsideNamespace(t *testing.T) {
	name := "Group Should Prevent Forking To External Namespaces"
	testedPolicyName := "collaborators_can_fork_repositories_to_external_namespaces"

	makeMockData := func(flag bool) gitlabcollected.Organization {
		return gitlabcollected.Organization{Group: &gitlab2.Group{PreventForkingOutsideGroup: flag}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		gitlabGroupTestTemplate(t, name, makeMockData(options[expectFailure]), testedPolicyName, expectFailure)
	}
}

func TestGitlabGroupBranchProtection(t *testing.T) {
	name := "Group Should Enforce Branch Protection By Default"
	testedPolicyName := "group_does_not_enforce_branch_protection_by_default"

	makeMockData := func(protection int) gitlabcollected.Organization {
		return gitlabcollected.Organization{Group: &gitlab2.Group{DefaultBranchProtection: protection}}
	}

	options := map[bool]int{
		false: 2,
		true:  0,
	}

	for _, expectFailure := range bools {
		gitlabGroupTestTemplate(t, name, makeMockData(options[expectFailure]), testedPolicyName, expectFailure)
	}
}
