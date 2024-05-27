package skippers

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/analyzers/parsing_utils"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"log"
)

type Skipper interface {
	ShouldSkip(data collectors.CollectedData, violation opa_engine.QueryResult) bool
}

type IsPrerequisitesSatisfied func(data collectors.CollectedData) bool

func NewSkipper(ctx context.Context) Skipper {
	return &skipper{
		ctx:             ctx,
		ignoredPolicies: context_utils.GetIgnoredPolicies(ctx),
		prerequisitesCheckers: map[string]IsPrerequisitesSatisfied{
			"premium": func(data collectors.CollectedData) bool {
				return data.Context.Premium()
			},
			"scorecard_enabled": func(data collectors.CollectedData) bool {
				return context_utils.GetScorecardEnabled(ctx)
			},
			"has_branch_protection_permission": func(data collectors.CollectedData) bool {
				repositoryContext, ok := data.Context.(collectors.CollectedDataRepositoryContext)
				if !ok {
					log.Printf("invalid type %T", data.Context)
					return false
				}
				return repositoryContext.HasBranchProtectionPermission()
			},
			"enterprise": func(_ collectors.CollectedData) bool {
				return !context_utils.GetIsCloud(ctx)
			},
			"is_branch_protection_supported": func(data collectors.CollectedData) bool {
				repositoryContext, ok := data.Context.(collectors.CollectedDataRepositoryContext)
				if !ok {
					log.Printf("invalid type %T", data.Context)
					return false
				}
				return repositoryContext.IsBranchProtectionSupported()
			},
		},
	}
}

type skipper struct {
	ctx                   context.Context
	prerequisitesCheckers map[string]IsPrerequisitesSatisfied
	ignoredPolicies       []string
}

func (sm *skipper) ShouldSkip(data collectors.CollectedData, violation opa_engine.QueryResult) bool {
	if sm.ignoredPolicy(violation) {
		return true
	}

	prerequisites := parsing_utils.ResolveAnnotation(violation.Annotations.Custom["prerequisites"])

	sufficient, missingPrerequisite := sm.arePrerequisitesSatisfied(prerequisites, data)
	if !sufficient {
		errlog.AddSkipIssue(violation.PolicyName, data.Entity.Name(), errlog.NewPrerequisiteSkipReason(missingPrerequisite))
		return true
	}

	currentScopes := context_utils.GetTokenScopes(sm.ctx)
	scopes := parsing_utils.ResolveAnnotation(violation.Annotations.Custom["requiredScopes"])

	sufficient, missingScope := sufficientScopes(data.Context.Roles(), currentScopes, scopes)
	if !sufficient {
		errlog.AddSkipIssue(violation.PolicyName, data.Entity.Name(), errlog.NewPermissionSkipReason(missingScope))
		return true
	}

	return false
}

func (sm *skipper) ignoredPolicy(policy opa_engine.QueryResult) bool {
	for _, ignored := range sm.ignoredPolicies {
		if policy.PolicyName == ignored {
			return true
		}
	}

	return false
}

func (sm *skipper) arePrerequisitesSatisfied(pre []string, data collectors.CollectedData) (satisfied bool, prerequisite string) {
	for _, p := range pre {
		predicate, ok := sm.prerequisitesCheckers[p]
		if !ok || !predicate(data) {
			return false, p
		}
	}

	return true, ""
}

func sufficientScopes(roles []permissions.Role, availableScopes permissions.TokenScopes, requiredScopes []string) (sufficient bool, missing string) {
	for _, requiredScope := range requiredScopes {
		if !permissions.HasScope(requiredScope, availableScopes, roles) {
			return false, requiredScope
		}
	}

	return true, ""
}
