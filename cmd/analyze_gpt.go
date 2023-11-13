package cmd

import (
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newAnalyzeGptCommand())
}

var analyzeGptArgs args

const argOpenAiToken = "openai_token"
const argOpenAiGptModel = "model"

func newAnalyzeGptCommand() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:          "gpt-analysis",
		Short:        `Analyze your GitHub/GitLab assets for security issues with GPT`,
		RunE:         executeAnalyzeGPTCommand,
		SilenceUsage: true,
	}

	flags := analyzeCmd.Flags()

	analyzeGptArgs.addCommonCollectionOptions(flags)
	analyzeGptArgs.addOutputOptions(flags)
	flags.StringSliceVarP(&analyzeGptArgs.Organizations, argOrg, "", nil, "specific organizations to collect")
	flags.StringSliceVarP(&analyzeGptArgs.Repositories, argRepository, "", nil, "specific repositories to collect (--repo owner/repo_name (e.g. ossf/scorecard)")
	flags.StringVarP(&analyzeGptArgs.OpenAIToken, argOpenAiToken, "", "", "token to authenticate with openai API")
	flags.StringVarP(&analyzeGptArgs.OpenAIGptModel, argOpenAiGptModel, "", gogpt.GPT3TextDavinci003, "gpt-model to use")
	viper.AutomaticEnv()

	return analyzeCmd
}

func applyAnalyzeGPTArgs() error {
	if err := analyzeGptArgs.validateCommonCollectionOptions(); err != nil {
		return err
	}

	if len(analyzeGptArgs.Organizations) == 0 && len(analyzeGptArgs.Repositories) == 0 {
		return fmt.Errorf("must specificy at least one organization or repository")
	}

	if len(analyzeGptArgs.Organizations) != 0 && len(analyzeGptArgs.Repositories) != 0 {
		return fmt.Errorf("cannot use --org & --repo options together")
	}

	if analyzeGptArgs.OpenAIToken == "" {
		analyzeGptArgs.OpenAIToken = viper.GetString(argOpenAiToken)
		if analyzeGptArgs.OpenAIToken == "" {
			return fmt.Errorf("must provide openai API token")
		}
	}

	return nil
}

func setup() (*analyzeGPTExecutor, error) {
	if len(analyzeGptArgs.Repositories) != 0 {
		analyzeGptArgs.Namespaces = append(analyzeGptArgs.Namespaces, namespace.Repository)
	} else if len(analyzeGptArgs.Organizations) != 0 {
		analyzeGptArgs.Namespaces = append(analyzeGptArgs.Namespaces, namespace.Organization)
	}

	switch analyzeGptArgs.ScmType {
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
	if err := analyzeGptArgs.applyCommonCollectionOptions(); err != nil {
		return err
	}

	err := applyAnalyzeGPTArgs()
	if err != nil {
		return err
	}

	preExitHook, err := analyzeGptArgs.applyOutputOptions()
	if err != nil {
		return err
	}
	defer preExitHook()

	executor, err := setup()
	if err != nil {
		return err
	}

	return executor.Run()
}
