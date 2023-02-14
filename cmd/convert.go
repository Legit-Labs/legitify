package cmd

import (
	"fmt"
	"os"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
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

	viper.AutomaticEnv()
	flags := cmd.Flags()
	convertArgs.addSchemeOutputOptions(flags)

	flags.StringVar(&convertArgs.InputFile, argInputFile, "", "the input file")

	return cmd
}

func validateConvertArgs() error {
	if convertArgs.InputFile == "" {
		return fmt.Errorf("please provide an input file")
	}

	return nil
}

func executeConvertCommand(cmd *cobra.Command, _args []string) error {
	err := validateConvertArgs()
	if err != nil {
		return err
	}

	if preExit, err := convertArgs.applySchemeOutputOptions(); err != nil {
		return err
	} else {
		defer preExit()
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
