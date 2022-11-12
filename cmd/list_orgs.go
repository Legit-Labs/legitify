package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Legit-Labs/legitify/internal/clients/github"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newListOrgsCommand())
}

var listOrgsArgs args

func newListOrgsCommand() *cobra.Command {
	listOrgsCmd := &cobra.Command{
		Use:          "list-orgs",
		Short:        `List GitHub organizations associated with a PAT`,
		RunE:         executeListOrgsCommand,
		SilenceUsage: true,
	}

	viper.AutomaticEnv()
	flags := listOrgsCmd.Flags()
	listOrgsArgs.AddCommonOptions(flags)

	return listOrgsCmd
}

func validateListOrgsArgs() error {
	return nil
}

func executeListOrgsCommand(cmd *cobra.Command, _args []string) error {
	listOrgsArgs.ApplyEnvVars()

	err := validateListOrgsArgs()
	if err != nil {
		return err
	}

	if err = setErrorFile(listOrgsArgs.ErrorFile); err != nil {
		return err
	}

	err = setOutputFile(listOrgsArgs.OutputFile)
	if err != nil {
		return err
	}

	stdErrLog := log.New(os.Stderr, "", 0)
	ctx := context.Background()
	githubClient, err := github.NewClient(ctx, listOrgsArgs.Token, listOrgsArgs.Endpoint, []string{}, true)
	if err != nil {
		return err
	}
	if !githubClient.IsGithubCloud() {
		stdErrLog.Printf("Using Github Enterprise Endpoint: %s\n\n", listOrgsArgs.Endpoint)
	}

	orgs, err := githubClient.CollectOrganizations()
	if err != nil {
		return err
	}

	if len(orgs) == 0 {
		fmt.Printf("No organizations are associated with this PAT.\n")
	} else {
		owner, member := groupByMembership(orgs)
		fmt.Printf("Organizations:\n")
		fmt.Printf("--------------:\n")

		if len(owner) > 0 {
			fmt.Println("Full analysis available for the following organizations:")
			for _, org := range owner {
				fmt.Printf("  - %s (%s)\n", org.Name(), org.Role)
			}
		}

		if len(member) > 0 {
			fmt.Println("Partial results available for the following organizations:")
			for _, org := range member {
				fmt.Printf("  - %s (%s)\n", org.Name(), org.Role)
			}
		}
	}

	return nil
}

func groupByMembership(orgs []githubcollected.ExtendedOrg) (owner []githubcollected.ExtendedOrg, member []githubcollected.ExtendedOrg) {
	for _, o := range orgs {
		if o.Role == permissions.OrgRoleOwner {
			owner = append(owner, o)
		} else {
			member = append(member, o)
		}
	}

	return
}
