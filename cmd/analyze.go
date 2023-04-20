package cmd

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/screen"
	"os"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newAnalyzeCommand())
}

const (
	argOrg                        = "org"
	argRepository                 = "repo"
	argEnterprises                = "enterprise"
	argPoliciesPath               = "policies-path"
	argNamespace                  = "namespace"
	argOutputFormat               = "output-format"
	argOutputScheme               = "output-scheme"
	argColor                      = "color"
	argScorecard                  = "scorecard"
	argFailedOnly                 = "failed-only"
	argSimulateSecondaryRateLimit = "simulate-secondary-rate-limit"
)

func toOptionsString(options []string) string {
	return "[" + strings.Join(options, "/") + "]"
}

var analyzeArgs args

func newAnalyzeCommand() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:          "analyze",
		Short:        `Analyze GitHub/GitLab organizations associated with a PAT to find security issues`,
		RunE:         executeAnalyzeCommand,
		SilenceUsage: true,
	}

	scorecardWhens := toOptionsString(scorecardOptions())

	viper.AutomaticEnv()
	flags := analyzeCmd.Flags()
	analyzeArgs.addSchemeOutputOptions(flags)
	analyzeArgs.addCommonCollectionOptions(flags)

	flags.StringSliceVarP(&analyzeArgs.Organizations, argOrg, "", nil, "specific organizations to collect")
	flags.StringSliceVarP(&analyzeArgs.Repositories, argRepository, "", nil, "specific repositories to collect (--repo owner/repo_name (e.g. ossf/scorecard)")
	flags.StringSliceVarP(&analyzeArgs.Enterprises, argEnterprises, "", nil, "specific enterprises to collect (--enterprise your_enterprise_slug) this flag must be provided with a value")
	flags.StringSliceVarP(&analyzeArgs.PoliciesPath, argPoliciesPath, "p", []string{}, "directory containing opa policies")
	flags.StringSliceVarP(&analyzeArgs.Namespaces, argNamespace, "n", namespace.All, "which namespace to run")
	flags.StringVarP(&analyzeArgs.ScorecardWhen, argScorecard, "", DefaultScOption, "Whether to run additional scorecard checks "+scorecardWhens)
	flags.BoolVarP(&analyzeArgs.SimulateSecondaryRateLimit, argSimulateSecondaryRateLimit, "", false, "Simulate secondary rate limits (for testing purposes)")
	_ = flags.MarkHidden(argSimulateSecondaryRateLimit)

	return analyzeCmd
}

func validateAnalyzeArgs() error {
	if err := namespace.ValidateNamespaces(analyzeArgs.Namespaces); err != nil {
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

func setupExecutor(analyzeArgs *args) (*analyzeExecutor, error) {
	switch analyzeArgs.ScmType {
	case scm_type.GitHub:
		return setupGitHub(analyzeArgs)
	case scm_type.GitLab:
		return setupGitLab(analyzeArgs)
	default:
		// shouldn't happen since scm type is validated before
		return nil, fmt.Errorf("invalid scm type %s", analyzeArgs.ScmType)
	}
}

func executeAnalyzeCommand(cmd *cobra.Command, _args []string) error {
	if err := analyzeArgs.applyCommonCollectionOptions(); err != nil {
		return err
	}

	if preExit, err := analyzeArgs.applySchemeOutputOptions(); err != nil {
		return err
	} else {
		defer preExit()
	}

	if err := validateAnalyzeArgs(); err != nil {
		return err
	}

	// to make sure scorecard works
	if err := os.Setenv("GITHUB_AUTH_TOKEN", analyzeArgs.Token); err != nil {
		return err
	}

	executor, err := setupExecutor(&analyzeArgs)
	if err != nil {
		return err
	}

	screen.Printf("Note: to get the OpenSSF scorecard results for the organization repositories use the --scorecard option\n\n")

	return executor.Run()
}
