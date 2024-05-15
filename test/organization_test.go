package test

import (
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/google/go-github/v53/github"
	"testing"
	"time"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
)

type organizationMockConfiguration struct {
	config     map[string]interface{}
	ssoEnabled *bool
	name       string
	url        string
	secrets    []*githubcollected.OrganizationSecret
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

	return githubcollected.Organization{
		Organization: nil,
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
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newOrganizationMock(test.args),
			namespace.Organization, test.policyName, test.shouldBeViolated, scm_type.GitHub)
	}
}
