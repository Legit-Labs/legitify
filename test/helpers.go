package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/opa"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/google/go-github/v44/github"
)

var defaultOrg githubcollected.ExtendedOrg = githubcollected.NewExtendedOrg(&github.Organization{}, permissions.OrgRoleNone)

func FindPolicy(opaResults []opa_engine.QueryResult, policyName string) (*opa_engine.QueryResult, error) {
	for _, s := range opaResults {
		if policyName == s.PolicyName {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("policy doesn't exists %s", policyName)
}

func BuildFailedString(expected bool, policy string) string {
	if expected {
		return fmt.Sprintf("expected policy %s to exist in results", policy)
	}
	return fmt.Sprintf("did not expect expect policy %s to exist in results", policy)
}

func AssertQueryResult(opaResults []opa_engine.QueryResult, policyName string, shouldBeViolated bool, t *testing.T) {
	require.NotNil(t, opaResults, "query result == nil")
	policy, err := FindPolicy(opaResults, policyName)
	require.NoErrorf(t, err, "failed to find policy")

	if policy.IsViolation != shouldBeViolated {
		t.Error(BuildFailedString(shouldBeViolated, policyName))
	}
}

func PolicyTestTemplate(t *testing.T, name string, mockData interface{}, ns namespace.Namespace, testedPolicyName string, expectFailure bool) {
	t.Run(name, func(t *testing.T) {
		engine, err := opa.Load([]string{})
		require.Nil(t, err, "failed initializing opa client")
		ctx := context.Background()
		result, err := engine.Query(ctx, ns, mockData)
		require.Nil(t, err, "failed query")
		AssertQueryResult(result, testedPolicyName, expectFailure, t)
	})
}
