package gitlab

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	context context.Context
	client  *gitlab.Client
}

func (c *Client) Client() *gitlab.Client {
	return c.client
}

func NewClient(ctx context.Context, token string, endpoint string, fillCache bool) (*Client, error) {
	var config []gitlab.ClientOptionFunc
	if endpoint != "" {
		config = []gitlab.ClientOptionFunc{gitlab.WithBaseURL(endpoint)}
	}

	git, err := gitlab.NewClient(token, config...)
	if err != nil {
		return nil, err
	}

	return &Client{
		context: ctx,
		client:  git,
	}, nil
}

func (c *Client) IsAnalyzable(repo types.RepositoryWithOwner) (bool, error) {
	_, _, err := c.Client().Projects.GetProject(repo.String(), nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) Scopes() permissions.TokenScopes {
	return permissions.TokenScopes{}
}

func (c *Client) Organizations() ([]types.Organization, error) {
	var result []types.Organization

	dummy := true
	options := gitlab.ListGroupsOptions{Owned: &dummy}

	err := PaginateResults(func(opts *gitlab.ListOptions) (*gitlab.Response, error) {
		groups, resp, err := c.Client().Groups.ListGroups(&options)
		if err != nil {
			return nil, err
		}

		for _, g := range groups {
			result = append(result, types.Organization{
				Name: g.Name,
				Role: permissions.OrgRoleOwner,
			})
		}

		return resp, nil
	}, &options.ListOptions)

	if err != nil {
		return nil, err
	}

	return result, nil
}
