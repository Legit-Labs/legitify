package test

import (
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/google/go-github/v49/github"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
)

type organizationMockConfiguration struct {
	config     map[string]interface{}
	ssoEnabled *bool
	name       string
	url        string
}

func newOrganizationMock(config organizationMockConfiguration) githubcollected.Organization {
	samlEnabledMockResult := false
	visibilityChangePolicyMockResult := false
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

	return githubcollected.Organization{
		Organization:                             nil,
		SamlEnabled:                              &samlEnabledMockResult,
		Hooks:                                    hooks,
		EnterpriseVisibilityChangePolicyDisabled: &visibilityChangePolicyMockResult,
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
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newOrganizationMock(test.args),
			namespace.Organization, test.policyName, test.shouldBeViolated, scm_type.GitHub)
	}
}
