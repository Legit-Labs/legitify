package skippers

import (
	"context"
	"log"

	"github.com/Legit-Labs/legitify/internal/analyzers/parsing_utils"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
)

type Skipper interface {
	ShouldSkip(data collectors.CollectedData, violation opa_engine.QueryResult) bool
}

type IsPrerequisitesSatisfied func(data collectors.CollectedData) bool

func NewSkipper(ctx context.Context) Skipper {
	return &skipper{
		ctx: ctx,
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
			"enterprise_slug": func(_ collectors.CollectedData) bool {
				_, ok := context_utils.GetEnterprise(ctx)
				return ok
			},
		},
	}
}

type skipper struct {
	ctx                   context.Context
	prerequisitesCheckers map[string]IsPrerequisitesSatisfied
}

func (sm *skipper) ShouldSkip(data collectors.CollectedData, violation opa_engine.QueryResult) bool {
	prerequisites := parsing_utils.ResolveAnnotation(violation.Annotations.Custom["prerequisites"])

	sufficient, missingPrerequisite := sm.arePrerequisitesSatisfied(prerequisites, data)
	if !sufficient {
		errlog.PrereqIssueF("Skipping policy: %s, missing prerequisite: %s\n", violation.PolicyName, missingPrerequisite)
		return true
	}

	currentScopes := context_utils.GetTokenScopes(sm.ctx)
	scopes := parsing_utils.ResolveAnnotation(violation.Annotations.Custom["requiredScopes"])

	sufficient, missingScope := sufficientScopes(data.Context.Roles(), currentScopes, scopes)
	if !sufficient {
		errlog.PrereqIssueF("Skipping policy: %s, missing scope: %s\n", violation.PolicyName, missingScope)
		return true
	}

	return false
}

func (sm *skipper) arePrerequisitesSatisfied(pre []string, data collectors.CollectedData) (satisfied bool, predicate string) {
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
