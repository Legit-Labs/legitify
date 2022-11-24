//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/opa"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/Legit-Labs/legitify/internal/outputer"
	"github.com/google/wire"
	"log"
)

func setupGitHub(analyzeArgs *args, log *log.Logger) (*analyzeExecutor, error) {
	wire.Build(provideOpa,
		provideContext,
		provideGitHubClient,
		provideCollectorsManager,
		provideOutputer,
		initializeAnalyzeExecutor,
		analyzers.NewAnalyzer,
		skippers.NewSkipper,
		enricher.NewEnricherManager)
	return nil, nil
}

func provideOutputer(ctx context.Context, analyzeArgs *args) outputer.Outputer {
	return outputer.NewOutputer(ctx, analyzeArgs.OutputFormat, analyzeArgs.OutputScheme, analyzeArgs.FailedOnly)
}

func provideCollectorsManager(ctx context.Context, analyzeArgs *args, client github.Client) collectors_manager.CollectorManager {
	return collectors_manager.NewCollectorsManager(ctx, analyzeArgs.Namespaces, client)
}

func provideOpa(analyzeArgs *args) (opa_engine.Enginer, error) {
	opaEngine, err := opa.Load(analyzeArgs.PoliciesPath)
	if err != nil {
		return nil, err
	}
	return opaEngine, nil
}

func provideGitHubClient(analyzeArgs *args) (github.Client, error) {
	return github.NewClient(context.Background(), analyzeArgs.Token, analyzeArgs.Endpoint,
		analyzeArgs.Organizations, false)
}

func provideContext(client github.Client, logger *log.Logger) (context.Context, error) {
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
		parsedRepositories = validated
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
