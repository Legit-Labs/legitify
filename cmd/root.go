package cmd

import (
	"log"

	"github.com/Legit-Labs/legitify/internal/screen"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "legitify",
	Short: "Strengthen the security posture of your GitHub organization!",
	Long:  `Detect and remediate misconfigurations, security and compliance issues across all your GitHub assets with ease.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

const logo = ` ___      _______  _______  ___   _______  ___   _______  __   __
|   |    |       ||       ||   | |       ||   | |       ||  | |  |
|   |    |    ___||    ___||   | |_     _||   | |    ___||  |_|  |
|   |    |   |___ |   | __ |   |   |   |  |   | |   |___ |       |
|   |___ |    ___||   ||  ||   |   |   |  |   | |    ___||_     _|
|       ||   |___ |   |_| ||   |   |   |  |   | |   |      |   |
|_______||_______||_______||___|   |___|  |___| |___|      |___|`
const brand = `Legit Security`

func Execute() {
	if screen.IsTty() {
		logoColored := color.New(color.FgMagenta, color.Bold).Sprintf("%s", logo)
		brandColored := color.New(color.Bold).Sprintf("%s", brand)
		screen.Printf("%s\nBy %s\n\n", logoColored, brandColored)
	}
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error executing command: %s", err)
	}
}
