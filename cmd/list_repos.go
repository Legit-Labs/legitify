package cmd

import (
	"context"
	"fmt"
	"github.com/Legit-Labs/legitify/cmd/common_options"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/shurcooL/githubv4"
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
		Short:        `List GitHub repositories associated with a PAT`,
		RunE:         executeListReposCommand,
		SilenceUsage: true,
	}

	viper.AutomaticEnv()
	flags := listReposCmd.Flags()
	flags.StringVarP(&listReposArgs.Token, common_options.ArgToken, "t", "", "token to authenticate with github (required unless environment variable GITHUB_TOKEN is set)")
	flags.StringVarP(&listReposArgs.OutputFile, common_options.ArgOutputFile, "o", "", "output file, defaults to stdout")
	flags.StringVarP(&listReposArgs.ErrorFile, common_options.ArgErrorFile, "e", "error.log", "error log path")

	return listReposCmd
}

func validateListReposArgs() error {
	if err := github.IsTokenValid(listReposArgs.Token); err != nil {
		return err
	}

	return nil
}

func executeListReposCommand(cmd *cobra.Command, _args []string) error {
	if listReposArgs.Token == "" {
		listReposArgs.Token = viper.GetString(common_options.EnvToken)
	}

	err := validateListReposArgs()
	if err != nil {
		return err
	}

	if err = setErrorFile(listReposArgs.ErrorFile); err != nil {
		return err
	}

	err = setOutputFile(listReposArgs.OutputFile)
	if err != nil {
		return err
	}

	ctx := context.Background()

	githubClient, err := github.NewClient(ctx, listReposArgs.Token, []string{}, true)
	if err != nil {
		return err
	}

	type repository struct {
		repoWithOwner string
		permission    string
	}
	var repositories []repository

	var query struct {
		Viewer struct {
			Repositories struct {
				PageInfo githubcollected.GitHubQLPageInfo
				Nodes    []struct {
					NameWithOwner    string
					ViewerPermission string
				}
			} `graphql:"repositories(first:50, after: $cursor)"`
		}
	}

	variables := map[string]interface{}{
		"cursor": (*githubv4.String)(nil),
	}

	for {
		err := githubClient.GraphQLClient().Query(ctx, &query, variables)
		if err != nil {
			return err
		}
		for _, r := range query.Viewer.Repositories.Nodes {
			repositories = append(repositories, repository{
				repoWithOwner: r.NameWithOwner,
				permission:    r.ViewerPermission,
			})
		}

		if !query.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}

		variables["$cursor"] = query.Viewer.Repositories.PageInfo.EndCursor
	}

	if len(repositories) == 0 {
		fmt.Printf("No repositories are associated with this PAT.\n")
	} else {
		fmt.Printf("Repositories:\n")
		for _, repo := range repositories {
			fmt.Printf("- %s (%s)\n", repo.repoWithOwner, repo.permission)
		}
	}

	return nil
}
