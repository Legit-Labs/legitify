package github

import "github.com/google/go-github/v44/github"

func PaginateResults(api func(opts *github.ListOptions) (*github.Response, error)) error {
	var opts github.ListOptions

	for {
		resp, err := api(&opts)

		if err != nil {
			return err
		}

		if resp.NextPage == 0 {
			return nil
		}

		opts.Page = resp.NextPage
	}
}
