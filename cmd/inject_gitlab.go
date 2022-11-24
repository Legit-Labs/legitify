//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	github2 "github.com/Legit-Labs/legitify/internal/collectors/github"
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
		provideOutputer,
		initializeAnalyzeExecutor,
		provideCollectors,
		analyzers.NewAnalyzer,
		skippers.NewSkipper,
		enricher.NewEnricherManager,
		collectors_manager.NewCollectorsManager)
	return nil, nil
}

func provideCollectors(ctx context.Context, client github.Client, analyzeArgs *args) []collectors.Collector {
	type newCollectorFunc func(ctx context.Context, client github.Client) collectors.Collector
	var collectorsMapping = map[namespace.Namespace]newCollectorFunc{
		namespace.Repository:   github2.NewRepositoryCollector,
		namespace.Organization: github2.NewOrganizationCollector,
		namespace.Member:       github2.NewMemberCollector,
		namespace.Actions:      github2.NewActionCollector,
		namespace.RunnerGroup:  github2.NewRunnersCollector,
	}

	var result []collectors.Collector
	for _, ns := range analyzeArgs.Namespaces {
		result = append(result, collectorsMapping[ns](ctx, client))
	}

	return result
}

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
