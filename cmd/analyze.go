package cmd

import (
	"github.com/Legit-Labs/legitify/cmd/common_options"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	"log"
	"os"
	"strings"

	"github.com/Legit-Labs/legitify/internal/opa"

	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/common/namespace"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/outputer"
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
	argPoliciesPath = "policies-path"
	argNamespace    = "namespace"
	argOutputFormat = "output-format"
	argOutputScheme = "output-scheme"
	argColor        = "color"
	argScorecard    = "scorecard"
	argFailedOnly   = "failed-only"
)

type args struct {
	Token         string
	Orgs          []string
	PoliciesPath  []string
	Namespaces    []string
	ColorWhen     string
	OutputFile    string
	ErrorFile     string
	OutputFormat  string
	OutputScheme  string
	ScorecardWhen string
	FailedOnly    bool
}

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
	flags.StringVarP(&analyzeArgs.Token, common_options.ArgToken, "t", "", "token to authenticate with github (required unless environment variable GITHUB_TOKEN is set)")
	flags.StringSliceVarP(&analyzeArgs.Orgs, argOrg, "", nil, "specific organizations to collect")
	flags.StringSliceVarP(&analyzeArgs.PoliciesPath, argPoliciesPath, "p", []string{}, "directory containing opa policies")
	flags.StringSliceVarP(&analyzeArgs.Namespaces, argNamespace, "n", namespace.All, "which namespace to run")
	flags.StringVarP(&analyzeArgs.OutputFile, common_options.ArgOutputFile, "", "", "output file, defaults to stdout")
	flags.StringVarP(&analyzeArgs.ErrorFile, common_options.ArgErrorFile, "", "error.log", "error log path")
	flags.StringVarP(&analyzeArgs.OutputFormat, argOutputFormat, "f", formatter.Human, "output format "+formats)
	flags.StringVarP(&analyzeArgs.OutputScheme, argOutputScheme, "", converter.DefaultScheme, "output scheme "+schemeTypes)
	flags.StringVarP(&analyzeArgs.ColorWhen, argColor, "", DefaultColorOption, "when to use coloring "+colorWhens)
	flags.StringVarP(&analyzeArgs.ScorecardWhen, argScorecard, "", DefaultScOption, "Whether to run additional scorecard checks "+scorecardWhens)
	flags.BoolVarP(&analyzeArgs.FailedOnly, argFailedOnly, "", false, "Only show violated policied (do not show succeeded/skipped)")

	return analyzeCmd
}

func validateAnalyzeArgs() error {
	if err := github.IsTokenValid(analyzeArgs.Token); err != nil {
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

	return nil
}

func executeAnalyzeCommand(cmd *cobra.Command, _args []string) error {
	if analyzeArgs.Token == "" {
		analyzeArgs.Token = viper.GetString(common_options.EnvToken)
	}

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

	ctx := context_utils.NewContextWithOrg(analyzeArgs.Orgs)
	ctx = context_utils.NewContextWithScorecard(ctx,
		IsScorecardEnabled(analyzeArgs.ScorecardWhen),
		IsScorecardVerbose(analyzeArgs.ScorecardWhen))

	if !IsScorecardEnabled(analyzeArgs.ScorecardWhen) {
		stdErrLog.Printf("Note: to get the OpenSSF scorecard results for the organization repositories use the --scorecard option\n\n")
	}

	githubClient, err := github.NewClient(ctx, analyzeArgs.Token, analyzeArgs.Orgs)
	if err != nil {
		return err
	}

	ctx = context_utils.NewContextWithTokenScopes(ctx, githubClient.Scopes())

	opaEngine, err := opa.Load(analyzeArgs.PoliciesPath)
	if err != nil {
		return err
	}

	manager := collectors.NewCollectorsManager(ctx, analyzeArgs.Namespaces, githubClient)
	analyzer := analyzers.NewAnalyzer(ctx, opaEngine, skippers.NewSkipper(ctx))
	enricherManager := enricher.NewEnricherManager(ctx)
	out := outputer.NewOutputer(ctx, analyzeArgs.OutputFormat, analyzeArgs.OutputScheme, analyzeArgs.FailedOnly)

	stdErrLog.Printf("Gathering collection metadata...")
	collectionMetadata := manager.CollectMetadata()
	progressBar := progressbar.NewProgressBar(collectionMetadata)

	// TODO progressBar should run before collection starts and wait for channels to read from
	collectionChannels := manager.Collect()
	pWaiter := progressBar.Run(collectionChannels.Progress)
	analyzedDataChan := analyzer.Analyze(collectionChannels.Collected)
	enrichedDataChan := enricherManager.Enrich(analyzedDataChan)
	outputWaiter := out.Digest(enrichedDataChan)

	// Wait for progress bars to finish before outputing
	pWaiter.Wait()

	// Wait for output to be digested
	outputWaiter.Wait()

	err = out.Output(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
