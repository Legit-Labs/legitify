package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/Legit-Labs/legitify/cmd/common_options"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
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
	stdErrLog := log.New(os.Stderr, "", 0)
	githubEndpoint := viper.GetString(common_options.EnvGitHubEndpoint)

	githubClient, err := github.NewClient(ctx, listReposArgs.Token, githubEndpoint, []string{}, true)
	if err != nil {
		return err
	}
	if !githubClient.IsGithubCloud() {
		stdErrLog.Printf("Using Github Enterprise Endpoint: %s\n\n", githubEndpoint)
	}

	repositories, err := getRepositories(githubClient, ctx)
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
				fmt.Printf("  - %s (%s)\n", repo.repoWithOwner, repo.permission)
			}
		}

		if len(notAnalyzable) > 0 {
			fmt.Println("Your permissions are NOT sufficient to analyze the following repositories:")
			for _, repo := range notAnalyzable {
				fmt.Printf("  - %s (%s)\n", repo.repoWithOwner, repo.permission)
			}
		}
	}

	return nil
}

type repository struct {
	repoWithOwner string
	permission    string
}

func unique(slice []repository) []repository {
	keys := make(map[string]bool)
	list := []repository{}
	for _, entry := range slice {
		key := entry.repoWithOwner
		if _, value := keys[key]; !value {
			keys[key] = true
			list = append(list, entry)
		}
	}
	return list
}

func groupByAnalyzable(repositories []repository) (analyzable []repository, notAnalyzable []repository) {
	for _, r := range repositories {
		if r.permission == "ADMIN" {
			analyzable = append(analyzable, r)
		} else {
			notAnalyzable = append(notAnalyzable, r)
		}
	}

	sort.Slice(analyzable, func(i, j int) bool {
		return analyzable[i].repoWithOwner < analyzable[j].repoWithOwner
	})

	sort.Slice(notAnalyzable, func(i, j int) bool {
		return notAnalyzable[i].repoWithOwner < notAnalyzable[j].repoWithOwner
	})

	return
}

func getRepositories(githubClient github.Client, ctx context.Context) ([]repository, error) {
	r1, err := getViewerRepositories(githubClient, ctx)
	if err != nil {
		return nil, err
	}

	r2, err := getOrganizationRepositories(githubClient, ctx)
	if err != nil {
		return nil, err
	}

	return unique(append(r1, r2...)), nil
}

func getViewerRepositories(githubClient github.Client, ctx context.Context) ([]repository, error) {
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
			return nil, err
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

		variables["cursor"] = query.Viewer.Repositories.PageInfo.EndCursor
	}

	return repositories, nil
}

func getOrganizationRepositories(githubClient github.Client, ctx context.Context) ([]repository, error) {
	var repositories []repository
	orgs := githubClient.Orgs()
	gw := group_waiter.New()

	for _, o := range orgs {
		o := o
		gw.Do(func() {
			var query struct {
				Organization struct {
					Repositories struct {
						PageInfo githubcollected.GitHubQLPageInfo
						Nodes    []struct {
							NameWithOwner    string
							ViewerPermission string
						}
					} `graphql:"repositories(first: 50, after: $cursor)"`
				} `graphql:"organization(login: $login)"`
			}

			variables := map[string]interface{}{
				"cursor": (*githubv4.String)(nil),
				"login":  githubv4.String(o),
			}

			for {
				err := githubClient.GraphQLClient().Query(ctx, &query, variables)
				if err != nil {
					return
				}

				for _, r := range query.Organization.Repositories.Nodes {
					repositories = append(repositories, repository{
						repoWithOwner: r.NameWithOwner,
						permission:    r.ViewerPermission,
					})
				}

				if !query.Organization.Repositories.PageInfo.HasNextPage {
					break
				}

				variables["cursor"] = query.Organization.Repositories.PageInfo.EndCursor
			}
		})
	}

	gw.Wait()
	return repositories, nil
}
