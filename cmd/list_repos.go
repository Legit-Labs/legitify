package cmd

import (
	"fmt"
	"sort"

	"github.com/Legit-Labs/legitify/internal/common/types"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newListReposCommand())
}

var listReposArgs args

func newListReposCommand() *cobra.Command {
	listReposCmd := &cobra.Command{
		Use:          "list-repos",
		Short:        `List repositories associated with a PAT`,
		RunE:         executeListReposCommand,
		SilenceUsage: true,
	}

	viper.AutomaticEnv()
	flags := listReposCmd.Flags()
	listReposArgs.addOutputOptions(flags)
	listReposArgs.addCommonCollectionOptions(flags)

	return listReposCmd
}

func executeListReposCommand(cmd *cobra.Command, _args []string) error {
	if err := listReposArgs.applyCommonCollectionOptions(); err != nil {
		return err
	}

	if preExit, err := listReposArgs.applyOutputOptions(); err != nil {
		return err
	} else {
		defer preExit()
	}

	client, err := provideGenericClient(&listReposArgs)
	if err != nil {
		return err
	}

	repositories, err := client.Repositories()
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		fmt.Printf("No repositories are associated with this PAT.\n")
	} else {
		fmt.Printf("Repositories:\n")
		fmt.Printf("-------------:\n")
		analyzable, notAnalyzable := groupByAnalyzable(repositories)

		if len(analyzable) > 0 {
			fmt.Println("Full analysis available for the following repositories:")
			for _, repo := range analyzable {
				fmt.Printf("  - %s (%s)\n", repo.String(), repo.Role)
			}
		}

		if len(notAnalyzable) > 0 {
			fmt.Println("Your permissions are NOT sufficient to analyze the following repositories:")
			for _, repo := range notAnalyzable {
				fmt.Printf("  - %s (%s)\n", repo.String(), repo.Role)
			}
		}
	}

	return nil
}

func groupByAnalyzable(repositories []types.RepositoryWithOwner) (analyzable []types.RepositoryWithOwner, notAnalyzable []types.RepositoryWithOwner) {
	for _, r := range repositories {
		if r.Role == "ADMIN" {
			analyzable = append(analyzable, r)
		} else {
			notAnalyzable = append(notAnalyzable, r)
		}
	}

	sort.Slice(analyzable, func(i, j int) bool {
		return analyzable[i].String() < analyzable[j].String()
	})

	sort.Slice(notAnalyzable, func(i, j int) bool {
		return notAnalyzable[i].String() < notAnalyzable[j].String()
	})

	return
}
