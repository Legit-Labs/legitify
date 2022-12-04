//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	glclient "github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/collectors/gitlab"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/google/wire"
	"log"
)

func setupGitLab(analyzeArgs *args, log *log.Logger) (*analyzeExecutor, error) {
	wire.Build(
		wire.Bind(new(Client), new(*glclient.Client)),
		analyzeProviderSet,
		provideGitLabClient,
		provideGitLabCollectors,
	)
	return nil, nil
}

func provideGitLabCollectors(ctx context.Context, client *glclient.Client, analyzeArgs *args) []collectors.Collector {
	var collectorsMapping = map[namespace.Namespace]func(ctx context.Context, client *glclient.Client) collectors.Collector{
		namespace.Organization: gitlab.NewGroupCollector,
	}

	var result []collectors.Collector
	for _, ns := range analyzeArgs.Namespaces {
		if creator, ok := collectorsMapping[ns]; ok {
			result = append(result, creator(ctx, client))
		}
	}

	return result
}

func provideGitLabClient(analyzeArgs *args) (*glclient.Client, error) {
	return glclient.NewClient(context.Background(), analyzeArgs.Token, analyzeArgs.Endpoint, analyzeArgs.Organizations, false)
}
