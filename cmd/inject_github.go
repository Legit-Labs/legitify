//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/collectors"
	github2 "github.com/Legit-Labs/legitify/internal/collectors/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/google/wire"
)

func setupGitHub(analyzeArgs *args) (*analyzeExecutor, error) {
	wire.Build(
		wire.Bind(new(Client), new(*github.Client)),
		analyzeProviderSet,
		provideGitHubClient,
		provideGitHubCollectors,
	)
	return nil, nil
}

func setupGitHubGPTExecutor(analyzeArgs *args) (*analyzeGPTExecutor, error) {
	wire.Build(
		wire.Bind(new(Client), new(*github.Client)),
		analyzeProviderSet,
		provideGitHubClient,
		provideGitHubCollectors,
	)
	return nil, nil
}

func provideGitHubCollectors(ctx context.Context, client *github.Client, analyzeArgs *args) []collectors.Collector {
	type newCollectorFunc func(ctx context.Context, client *github.Client) collectors.Collector
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

func provideGitHubClient(analyzeArgs *args) (*github.Client, error) {
	ctx := context_utils.NewContextWithSimulatedSecondaryRateLimit(context.Background(), analyzeArgs.SimulateSecondaryRateLimit)
	return github.NewClient(ctx, analyzeArgs.Token, analyzeArgs.Endpoint,
		analyzeArgs.Organizations)
}
