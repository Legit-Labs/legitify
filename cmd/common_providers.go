package cmd

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/opa"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/Legit-Labs/legitify/internal/outputer"
	"log"
)

func provideOutputer(ctx context.Context, analyzeArgs *args) outputer.Outputer {
	return outputer.NewOutputer(ctx, analyzeArgs.OutputFormat, analyzeArgs.OutputScheme, analyzeArgs.FailedOnly)
}

func provideOpa(analyzeArgs *args) (opa_engine.Enginer, error) {
	opaEngine, err := opa.Load(analyzeArgs.PoliciesPath)
	if err != nil {
		return nil, err
	}
	return opaEngine, nil
}

func provideContext(client Client, logger *log.Logger) (context.Context, error) {
	var ctx context.Context
	if len(analyzeArgs.Organizations) != 0 {
		ctx = context_utils.NewContextWithOrg(analyzeArgs.Organizations)
	} else if len(analyzeArgs.Repositories) != 0 {
		validated, err := validateRepositories(analyzeArgs.Repositories)
		if err != nil {
			return nil, err
		}
		if err = repositoriesAnalyzable(ctx, client, validated); err != nil {
			return nil, err
		}
		ctx = context_utils.NewContextWithRepos(validated)
		analyzeArgs.Namespaces = []namespace.Namespace{namespace.Repository}
	} else {
		ctx = context.Background()
	}

	ctx = context_utils.NewContextWithScorecard(ctx,
		IsScorecardEnabled(analyzeArgs.ScorecardWhen),
		IsScorecardVerbose(analyzeArgs.ScorecardWhen))

	if !IsScorecardEnabled(analyzeArgs.ScorecardWhen) {
		logger.Printf("Note: to get the OpenSSF scorecard results for the organization repositories use the --scorecard option\n\n")
	}

	return context_utils.NewContextWithTokenScopes(ctx, client.Scopes()), nil
}
