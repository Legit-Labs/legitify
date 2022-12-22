package cmd

import (
	"context"
	"fmt"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/opa"
	"github.com/Legit-Labs/legitify/internal/opa/opa_engine"
	"github.com/Legit-Labs/legitify/internal/outputer"
	"github.com/Legit-Labs/legitify/internal/screen"
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

func provideContext(client Client) (context.Context, error) {
	var ctx = context.Background()

	if len(analyzeArgs.Organizations) != 0 {
		ctx = context_utils.NewContextWithOrg(analyzeArgs.Organizations)
	} else if len(analyzeArgs.Repositories) != 0 {
		validated, err := validateRepositories(analyzeArgs.Repositories)
		if err != nil {
			return nil, err
		}
		if err = repositoriesAnalyzable(client, validated); err != nil {
			return nil, err
		}
		ctx = context_utils.NewContextWithRepos(validated)
		analyzeArgs.Namespaces = []namespace.Namespace{namespace.Repository}
	}

	ctx = context_utils.NewContextWithScorecard(ctx,
		IsScorecardEnabled(analyzeArgs.ScorecardWhen),
		IsScorecardVerbose(analyzeArgs.ScorecardWhen))

	if !IsScorecardEnabled(analyzeArgs.ScorecardWhen) {
		screen.Printf("Note: to get the OpenSSF scorecard results for the organization repositories use the --scorecard option\n\n")
	}

	return context_utils.NewContextWithTokenScopes(ctx, client.Scopes()), nil
}
