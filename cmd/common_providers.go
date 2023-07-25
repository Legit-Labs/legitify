package cmd

import (
	"context"
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/gpt"
	"github.com/Legit-Labs/legitify/internal/opa"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/Legit-Labs/legitify/internal/outputer"
)

func provideGenericClient(args *args) (Client, error) {
	switch args.ScmType {
	case scm_type.GitHub:
		return provideGitHubClient(args)
	case scm_type.GitLab:
		return provideGitLabClient(args)
	default:
		return nil, fmt.Errorf("invalid scm type")
	}
}

func provideOutputer(ctx context.Context, analyzeArgs *args) outputer.Outputer {
	return outputer.NewOutputer(ctx, analyzeArgs.OutputFormat, analyzeArgs.OutputScheme, analyzeArgs.FailedOnly)
}

func provideOpa(analyzeArgs *args) (opa_engine.Enginer, error) {
	opaEngine, err := opa.Load(analyzeArgs.PoliciesPath, analyzeArgs.ScmType)
	if err != nil {
		return nil, err
	}
	return opaEngine, nil
}

func provideContext(client Client, args *args) (context.Context, error) {
	ctx := context.Background()

	if len(args.Organizations) != 0 {
		ctx = context_utils.NewContextWithOrg(args.Organizations)
	} else if len(args.Repositories) != 0 {
		validated, err := validateRepositories(args.Repositories)
		if err != nil {
			return nil, err
		}
		if err = repositoriesAnalyzable(client, validated); err != nil {
			return nil, err
		}
		ctx = context_utils.NewContextWithRepos(validated)
		args.Namespaces = []namespace.Namespace{namespace.Repository}
	}

	ctx = context_utils.NewContextWithScorecard(ctx,
		IsScorecardEnabled(args.ScorecardWhen),
		IsScorecardVerbose(args.ScorecardWhen))

	ctx = context_utils.NewContextWithIsCloud(ctx, args.Endpoint == "")
	ctx = context_utils.NewContextWithIgnoredPolicies(ctx, args.IgnoredPolicies)

	return context_utils.NewContextWithTokenScopes(ctx, client.Scopes()), nil
}

func provideGPTAnalyzer(context context.Context, args *args) *gpt.Analyzer {
	return gpt.NewAnalyzer(context, args.OpenAIToken)
}
