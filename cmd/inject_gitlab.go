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

func setupGitlab(analyzeArgs *args, log *log.Logger) (*analyzeExecutor, error) {
	wire.Build(
		wire.Bind(new(Client), new(*glclient.Client)),
		analyzeProviderSet,
		provideGitlabClient,
		provideGitlabCollectors,
	)
	return nil, nil
}

func provideGitlabCollectors(ctx context.Context, client *glclient.Client, analyzeArgs *args) []collectors.Collector {
	type newCollectorFunc func(ctx context.Context, client glclient.Client) collectors.Collector
	var collectorsMapping = map[namespace.Namespace]newCollectorFunc{
		namespace.Organization: gitlab.NewGroupCollector,
	}

	var result []collectors.Collector
	for _, ns := range analyzeArgs.Namespaces {
		result = append(result, collectorsMapping[ns](ctx, *client))
	}

	return result
}

func provideGitlabClient(analyzeArgs *args) (*glclient.Client, error) {
	return glclient.NewClient(context.Background(), analyzeArgs.Token, analyzeArgs.Endpoint, false)
}
