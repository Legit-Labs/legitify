package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/screen"

	"github.com/Legit-Labs/legitify/internal/common/namespace"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newAnalyzeCommand())
}

const (
	argOrg                        = "org"
	argRepository                 = "repo"
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
	analyzeArgs.addCommonOptions(flags)

	flags.StringSliceVarP(&analyzeArgs.Organizations, argOrg, "", nil, "specific organizations to collect")
	flags.StringSliceVarP(&analyzeArgs.Repositories, argRepository, "", nil, "specific repositories to collect (--repo owner/repo_name (e.g. ossf/scorecard)")
	flags.StringSliceVarP(&analyzeArgs.PoliciesPath, argPoliciesPath, "p", []string{}, "directory containing opa policies")
	flags.StringSliceVarP(&analyzeArgs.Namespaces, argNamespace, "n", namespace.All, "which namespace to run")
	flags.StringVarP(&analyzeArgs.OutputFormat, argOutputFormat, "f", formatter.Human, "output format "+formats)
	flags.StringVarP(&analyzeArgs.OutputScheme, argOutputScheme, "", converter.DefaultScheme, "output scheme "+schemeTypes)
	flags.StringVarP(&analyzeArgs.ColorWhen, argColor, "", DefaultColorOption, "when to use coloring "+colorWhens)
	flags.StringVarP(&analyzeArgs.ScorecardWhen, argScorecard, "", DefaultScOption, "Whether to run additional scorecard checks "+scorecardWhens)
	flags.BoolVarP(&analyzeArgs.FailedOnly, argFailedOnly, "", false, "Only show violated policied (do not show succeeded/skipped)")
	flags.BoolVarP(&analyzeArgs.SimulateSecondaryRateLimit, argSimulateSecondaryRateLimit, "", false, "Simulate secondary rate limits (for testing purposes)")
	_ = flags.MarkHidden(argSimulateSecondaryRateLimit)

	return analyzeCmd
}

func validateAnalyzeArgs() error {
	if err := analyzeArgs.validateCommonOptions(); err != nil {
		return err
	}

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
	analyzeArgs.ApplyEnvVars()

	// to make sure scorecard works
	if err := os.Setenv("GITHUB_AUTH_TOKEN", analyzeArgs.Token); err != nil {
		return err
	}

	err := validateAnalyzeArgs()
	if err != nil {
		return err
	}

	errFile, err := setErrorFile(analyzeArgs.ErrorFile)
	if err != nil {
		return err
	}
	defer func() {
		if errlog.HadErrors() {
			screen.Printf("Some errors raised during the execution. Check %s for more details", errFile.Name())
		}
	}()

	err = setOutputFile(analyzeArgs.OutputFile)
	if err != nil {
		return err
	}

	err = InitColorPackage(analyzeArgs.ColorWhen)
	if err != nil {
		return err
	}

	executor, err := setupExecutor(&analyzeArgs)
	if err != nil {
		return err
	}

	screen.Printf("Note: to get the OpenSSF scorecard results for the organization repositories use the --scorecard option\n\n")

	return executor.Run()
}
