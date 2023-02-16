package cmd

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/screen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newAnalyzeGptCommand())
}

var analyzeGptArgs args

func newAnalyzeGptCommand() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:          "gpt-analysis",
		Short:        `Analyze your GitHub/GitLab assets for security issues with GPT`,
		RunE:         executeAnalyzeGPTCommand,
		SilenceUsage: true,
	}

	viper.AutomaticEnv()
	flags := analyzeCmd.Flags()
	analyzeGptArgs.addCommonOptions(flags)

	flags.StringSliceVarP(&analyzeGptArgs.Organizations, argOrg, "", nil, "specific organizations to collect")
	flags.StringSliceVarP(&analyzeGptArgs.Repositories, argRepository, "", nil, "specific repositories to collect (--repo owner/repo_name (e.g. ossf/scorecard)")
	flags.StringVarP(&analyzeGptArgs.OpenAIToken, "openai-token", "", "", "token to authenticate with openai API")

	return analyzeCmd
}

func validateAnalyzeGPTArgs() error {
	if err := analyzeGptArgs.validateCommonOptions(); err != nil {
		return err
	}

	if len(analyzeGptArgs.Organizations) == 0 && len(analyzeGptArgs.Repositories) == 0 {
		return fmt.Errorf("must specificy at least one organization or repository")
	}

	if len(analyzeGptArgs.Organizations) != 0 && len(analyzeGptArgs.Repositories) != 0 {
		return fmt.Errorf("cannot use --org & --repo options together")
	}

	if analyzeGptArgs.OpenAIToken == "" {
		return fmt.Errorf("must provide openai API token")
	}

	return nil
}

func setup() (*analyzeGPTExecutor, error) {
	if len(analyzeGptArgs.Repositories) != 0 {
		analyzeGptArgs.Namespaces = []namespace.Namespace{namespace.Repository}
	} else if len(analyzeGptArgs.Organizations) != 0 {
		analyzeGptArgs.Namespaces = []namespace.Namespace{namespace.Organization}
	}

	switch analyzeArgs.ScmType {
	case scm_type.GitHub:
		return setupGitHubGPTExecutor(&analyzeGptArgs)
	case scm_type.GitLab:
		return setupGitLabGPTExecutor(&analyzeGptArgs)
	default:
		// shouldn't happen since scm type is validated before
		return nil, fmt.Errorf("invalid scm type %s", analyzeArgs.ScmType)
	}
}

func executeAnalyzeGPTCommand(cmd *cobra.Command, _args []string) error {
	analyzeGptArgs.ApplyEnvVars()

	err := validateAnalyzeGPTArgs()
	if err != nil {
		return err
	}

	errFile, err := setErrorFile(analyzeGptArgs.ErrorFile)
	if err != nil {
		return err
	}
	defer func() {
		if errlog.HadErrors() {
			screen.Printf("Some errors raised during the execution. Check %s for more details", errFile.Name())
		}
	}()

	err = setOutputFile(analyzeGptArgs.OutputFile)
	if err != nil {
		return err
	}

	executor, err := setup()
	if err != nil {
		return err
	}

	return executor.Run()
}
