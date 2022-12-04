package cmd

import (
	"fmt"
	"github.com/Legit-Labs/legitify/internal/common/types"

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
		Short:        `List organizations associated with a PAT`,
		RunE:         executeListOrgsCommand,
		SilenceUsage: true,
	}

	viper.AutomaticEnv()
	flags := listOrgsCmd.Flags()
	listOrgsArgs.addCommonOptions(flags)

	return listOrgsCmd
}

func validateListOrgsArgs() error {
	return listOrgsArgs.validateCommonOptions()
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

	client, err := provideGenericClient(&listOrgsArgs)
	if err != nil {
		return err
	}

	orgs, err := client.Organizations()
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
				fmt.Printf("  - %s (%s)\n", org.Name, org.Role)
			}
		}

		if len(member) > 0 {
			fmt.Println("Partial results available for the following organizations:")
			for _, org := range member {
				fmt.Printf("  - %s (%s)\n", org.Name, org.Role)
			}
		}
	}

	return nil
}

func groupByMembership(orgs []types.Organization) (owner []types.Organization, member []types.Organization) {
	for _, o := range orgs {
		if o.Role == permissions.OrgRoleOwner {
			owner = append(owner, o)
		} else {
			member = append(member, o)
		}
	}

	return
}
