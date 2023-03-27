package cmd

import (
	"github.com/Legit-Labs/legitify/internal/common/scm_type"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/Legit-Labs/legitify/internal/screen"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type args struct {
	Token                      string
	OpenAIToken                string
	Endpoint                   string
	ScmType                    scm_type.ScmType
	Organizations              []string
	Repositories               []string
	PoliciesPath               []string
	Namespaces                 []string
	ColorWhen                  string
	OutputFile                 string
	ErrorFile                  string
	OutputFormat               string
	OutputScheme               string
	ScorecardWhen              string
	InputFile                  string
	FailedOnly                 bool
	SimulateSecondaryRateLimit bool
}

const (
	ArgErrorFile  = "error-file"
	ArgOutputFile = "output-file"
	ArgToken      = "token"
	ArgServerUrl  = "server-url"
	ScmType       = "scm"
)

const (
	EnvToken     = "github_token"
	NewEnvToken  = "legitify_token"
	EnvServerUrl = "server_url"
)

func (a *args) addOutputOptions(flags *pflag.FlagSet) {
	colorWhens := toOptionsString(ColorOptions())
	flags.StringVarP(&a.OutputFile, ArgOutputFile, "o", "", "output file, defaults to stdout")
	flags.StringVarP(&a.ErrorFile, ArgErrorFile, "e", "error.log", "error log path")
	flags.StringVarP(&a.ColorWhen, argColor, "", DefaultColorOption, "when to use coloring "+colorWhens)
}

func (a *args) applyOutputOptions() (preExitHook func(), err error) {
	if err := setOutputFile(a.OutputFile); err != nil {
		return nil, err
	}

	if err := InitColorPackage(a.ColorWhen); err != nil {
		return nil, err
	}

	errFile, err := setErrorFile(a.ErrorFile)
	if err != nil {
		return nil, err
	}

	return func() {
		if errlog.HadErrors() {
			screen.Printf("\n\nSome errors raised during the execution. Check %s for more details", errFile.Name())
		}
	}, nil
}

func (a *args) addCommonCollectionOptions(flags *pflag.FlagSet) {
	flags.StringVarP(&a.Token, ArgToken, "t", "", "token to authenticate with github/gitlab (required unless environment variable LEGITIFY_AUTH_TOKEN is set)")
	flags.StringVarP(&a.Endpoint, ArgServerUrl, "", "", "github/gitlab endpoint to use instead of the Cloud API (can be set via the environment variable SERVER_URL)")
	flags.StringVarP(&a.ScmType, ScmType, "", scm_type.GitHub, "server type (GitHub, GitLab), defaults to GitHub")
}

func (a *args) applyCommonCollectionOptions() error {
	if err := a.validateCommonCollectionOptions(); err != nil {
		return err
	}

	if a.Token == "" {
		// backwards compatibility: support both LEGITIFY_TOKEN and GITHUB_TOKEN environment variables.
		// In the future we'll remove the GITHUB_TOKEN option
		a.Token = viper.GetString(NewEnvToken)
		if a.Token == "" {
			a.Token = viper.GetString(EnvToken)
		}
	}

	if a.Endpoint == "" {
		a.Endpoint = viper.GetString(EnvServerUrl)
	}

	return nil
}

func (a *args) validateCommonCollectionOptions() error {
	if err := scm_type.Validate(a.ScmType); err != nil {
		return err
	}

	return nil
}

func (a *args) addSchemeOutputOptions(flags *pflag.FlagSet) {
	a.addOutputOptions(flags)

	formats := toOptionsString(formatter.OutputFormats())
	schemeTypes := toOptionsString(scheme.SchemeTypes())

	flags.StringVarP(&a.OutputFormat, argOutputFormat, "f", formatter.Human, "output format "+formats)
	flags.StringVarP(&a.OutputScheme, argOutputScheme, "", scheme.DefaultScheme, "output scheme "+schemeTypes)
	flags.BoolVarP(&a.FailedOnly, argFailedOnly, "", false, "Only show violated policies (do not show succeeded/skipped)")
}

func (a *args) applySchemeOutputOptions() (preExitHook func(), err error) {
	if err := a.validateSchemeOutputOptions(); err != nil {
		return nil, err
	}

	if preExitHook, err := a.applyOutputOptions(); err != nil {
		return nil, err
	} else {
		return preExitHook, nil
	}
}

func (a *args) validateSchemeOutputOptions() error {
	if err := converter.ValidateOutputScheme(a.OutputScheme); err != nil {
		return err
	}

	if err := formatter.ValidateOutputFormat(a.OutputFormat, a.OutputScheme); err != nil {
		return err
	}

	return nil
}
