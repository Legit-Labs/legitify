package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/types"

	"github.com/Legit-Labs/legitify/internal/common/namespace"

	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newAnalyzeCommand())
}

const (
	argOrg          = "org"
	argRepository   = "repo"
	argPoliciesPath = "policies-path"
	argNamespace    = "namespace"
	argOutputFormat = "output-format"
	argOutputScheme = "output-scheme"
	argColor        = "color"
	argScorecard    = "scorecard"
	argFailedOnly   = "failed-only"
)

func toOptionsString(options []string) string {
	return "[" + strings.Join(options, "/") + "]"
}

var analyzeArgs args
var parsedRepositories []types.RepositoryWithOwner

func newAnalyzeCommand() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:          "analyze",
		Short:        `Analyze GitHub organizations associated with a PAT to find security issues`,
		RunE:         executeAnalyzeCommand,
		SilenceUsage: true,
	}

	formats := toOptionsString(formatter.OutputFormats())
	schemeTypes := toOptionsString(converter.SchemeTypes())
	colorWhens := toOptionsString(ColorOptions())
	scorecardWhens := toOptionsString(scorecardOptions())

	viper.AutomaticEnv()
	flags := analyzeCmd.Flags()
	analyzeArgs.AddCommonOptions(flags)

	flags.StringSliceVarP(&analyzeArgs.Organizations, argOrg, "", nil, "specific organizations to collect")
	flags.StringSliceVarP(&analyzeArgs.Repositories, argRepository, "", nil, "specific repositories to collect (--repo owner/repo_name (e.g. ossf/scorecard)")
	flags.StringSliceVarP(&analyzeArgs.PoliciesPath, argPoliciesPath, "p", []string{}, "directory containing opa policies")
	flags.StringSliceVarP(&analyzeArgs.Namespaces, argNamespace, "n", namespace.All, "which namespace to run")
	flags.StringVarP(&analyzeArgs.OutputFormat, argOutputFormat, "f", formatter.Human, "output format "+formats)
	flags.StringVarP(&analyzeArgs.OutputScheme, argOutputScheme, "", converter.DefaultScheme, "output scheme "+schemeTypes)
	flags.StringVarP(&analyzeArgs.ColorWhen, argColor, "", DefaultColorOption, "when to use coloring "+colorWhens)
	flags.StringVarP(&analyzeArgs.ScorecardWhen, argScorecard, "", DefaultScOption, "Whether to run additional scorecard checks "+scorecardWhens)
	flags.BoolVarP(&analyzeArgs.FailedOnly, argFailedOnly, "", false, "Only show violated policied (do not show succeeded/skipped)")

	return analyzeCmd
}

func validateAnalyzeArgs() error {
	if err := namespace.ValidateNamespaces(analyzeArgs.Namespaces); err != nil {
		return err
	}

	if err := converter.ValidateOutputScheme(analyzeArgs.OutputScheme); err != nil {
		return err
	}

	if err := formatter.ValidateOutputFormat(analyzeArgs.OutputFormat, analyzeArgs.OutputScheme); err != nil {
		return err
	}

	if err := ValidateScorecardOption(analyzeArgs.ScorecardWhen); err != nil {
		return err
	}

	if len(analyzeArgs.Organizations) != 0 && len(analyzeArgs.Repositories) != 0 {
		return fmt.Errorf("cannot use --org & --repo options together")
	}

	return nil
}

func buildContext() (context.Context, error) {
	var ctx context.Context
	if len(analyzeArgs.Organizations) != 0 && len(analyzeArgs.Repositories) != 0 {
		return nil, fmt.Errorf("cannot use --org & --repo options together")
	} else if len(analyzeArgs.Organizations) != 0 {
		ctx = context_utils.NewContextWithOrg(analyzeArgs.Organizations)
	} else if len(analyzeArgs.Repositories) != 0 {
		validated, err := validateRepositories(analyzeArgs.Repositories)
		if err != nil {
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

	return ctx, nil
}

func executeAnalyzeCommand(cmd *cobra.Command, _args []string) error {
	analyzeArgs.ApplyEnvVars()

	// to make sure scorecard works
	if err := os.Setenv("GITHUB_AUTH_TOKEN", analyzeArgs.Token); err != nil {
		return err
	}

	err := validateAnalyzeArgs()
	if err != nil {
		return err
	}

	if err = setErrorFile(analyzeArgs.ErrorFile); err != nil {
		return err
	}

	err = setOutputFile(analyzeArgs.OutputFile)
	if err != nil {
		return err
	}

	err = InitColorPackage(analyzeArgs.ColorWhen)
	if err != nil {
		return err
	}

	stdErrLog := log.New(os.Stderr, "", 0)

	executor, err := setupGitHub(&analyzeArgs, stdErrLog)
	if err != nil {
		return err
	}

	return executor.Run()
}
