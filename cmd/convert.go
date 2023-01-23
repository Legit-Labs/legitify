package cmd

import (
	"fmt"
	"os"

	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/Legit-Labs/legitify/internal/screen"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newConvertCommand())
}

const (
	argInputFile = "input-file"
)

var convertArgs args

func newConvertCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "convert",
		Short:        `Convert analyze output to a different scheme/format (input must be a flattened json)`,
		RunE:         executeConvertCommand,
		SilenceUsage: true,
	}

	formats := toOptionsString(formatter.OutputFormats())
	schemeTypes := toOptionsString(scheme.SchemeTypes())
	colorWhens := toOptionsString(ColorOptions())

	viper.AutomaticEnv()
	flags := cmd.Flags()
	convertArgs.addOutputOptions(flags)

	flags.StringVar(&convertArgs.InputFile, argInputFile, "", "the input file")
	flags.StringVarP(&convertArgs.OutputFormat, argOutputFormat, "f", formatter.Human, "output format "+formats)
	flags.StringVarP(&convertArgs.OutputScheme, argOutputScheme, "", scheme.DefaultScheme, "output scheme "+schemeTypes)
	flags.StringVarP(&convertArgs.ColorWhen, argColor, "", DefaultColorOption, "when to use coloring "+colorWhens)
	flags.BoolVarP(&convertArgs.FailedOnly, argFailedOnly, "", false, "Only show violated policied (do not show succeeded/skipped)")

	return cmd
}

func validateConvertArgs() error {
	if convertArgs.InputFile == "" {
		return fmt.Errorf("please provide an input file")
	}

	if err := converter.ValidateOutputScheme(convertArgs.OutputScheme); err != nil {
		return err
	}

	if err := formatter.ValidateOutputFormat(convertArgs.OutputFormat, convertArgs.OutputScheme); err != nil {
		return err
	}

	return nil
}

func executeConvertCommand(cmd *cobra.Command, _args []string) error {
	err := validateConvertArgs()
	if err != nil {
		return err
	}

	errFile, err := setErrorFile(convertArgs.ErrorFile)
	if err != nil {
		return err
	}
	defer func() {
		if errlog.HadErrors() {
			screen.Printf("Some errors raised during the execution. Check %s for more details", errFile.Name())
		}
	}()

	err = setOutputFile(convertArgs.OutputFile)
	if err != nil {
		return err
	}

	err = InitColorPackage(convertArgs.ColorWhen)
	if err != nil {
		return err
	}

	inputData, err := os.ReadFile(convertArgs.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	flattened, err := scheme.Unmarshal(inputData)
	if err != nil {
		return err
	}

	output, err := formatter.Format(convertArgs.OutputFormat, formatter.DefaultOutputIndent, flattened, convertArgs.FailedOnly)
	if err != nil {
		return fmt.Errorf("failed to format: %v", err)
	}

	os.Stdout.Write(output)
	return nil
}
