package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

func PaginateResults(api func(opts *gitlab.ListOptions) (*gitlab.Response, error),
	opts *gitlab.ListOptions) error {
	for {
		resp, err := api(opts)

		if err != nil {
			return err
		}

		if resp.NextPage == resp.TotalPages {
			return nil
		}

		opts.Page = resp.NextPage
	}
}
