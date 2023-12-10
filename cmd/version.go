package cmd

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

const (
	versionCmdText = "version"
)

var versionCmd = &cobra.Command{
	Use:   versionCmdText,
	Short: "Print the version number",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetVersion())
	},
}

func GetVersion() string {
	return version.ReadableVersion
}

func GetVersionLean() string {
	return version.ReadableVersionLean
}
