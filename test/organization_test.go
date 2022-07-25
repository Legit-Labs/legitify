package test

import (
	"github.com/google/go-github/v44/github"
	"testing"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
)

type organizationMockConfiguration struct {
	config     map[string]interface{}
	ssoEnabled *bool
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
		}
		hooks = append(hooks, &hook)
	}

	return githubcollected.Organization{
		Organization: nil,
		SamlEnabled:  &samlEnabledMockResult,
		Hooks:        hooks,
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
			},
		},
		{
			name:             "Violate webhook policy no configuration",
			policyName:       "organization_webhook_no_secret",
			shouldBeViolated: true,
			args: organizationMockConfiguration{
				config: map[string]interface{}{},
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
			namespace.Organization, test.policyName, test.shouldBeViolated)
	}
}
