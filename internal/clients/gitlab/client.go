package gitlab

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/patrickmn/go-cache"
	"github.com/xanzy/go-gitlab"
)

const (
	orgsCacheKeys = "orgs"
)

type Client struct {
	context context.Context
	client  *gitlab.Client
	cache   *cache.Cache
	orgs    []string
}

func (c *Client) Client() *gitlab.Client {
	return c.client
}

func NewClient(ctx context.Context, token string, endpoint string, orgs []string, fillCache bool) (*Client, error) {
	var config []gitlab.ClientOptionFunc
	if endpoint != "" {
		config = []gitlab.ClientOptionFunc{gitlab.WithBaseURL(endpoint)}
	}

	git, err := gitlab.NewClient(token, config...)
	if err != nil {
		return nil, err
	}

	// all groups
	if len(orgs) == 0 {
		orgs = []string{""}
	}

	result := &Client{
		context: ctx,
		client:  git,
		cache:   cache.New(cache.NoExpiration, cache.NoExpiration),
		orgs:    orgs,
	}

	if fillCache {
		if err := result.fillCache(); err != nil {
			return nil, err
		}
	}

	return result, nil
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

	groups, err := c.Groups()
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		result = append(result, types.Organization{
			Name: g.Name,
			Role: permissions.OrgRoleOwner,
		})
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) Repositories() ([]types.RepositoryWithOwner, error) {
	var result []types.RepositoryWithOwner

	dummy := gitlab.MaintainerPermissions
	options := gitlab.ListProjectsOptions{MinAccessLevel: &dummy}
	err := PaginateResults(func(opts *gitlab.ListOptions) (*gitlab.Response, error) {
		repos, resp, err := c.Client().Projects.ListProjects(&options)
		if err != nil {
			return nil, err
		}

		for _, r := range repos {
			result = append(result, types.NewRepositoryWithOwner(r.PathWithNamespace, permissions.RepoRoleAdmin))
		}

		return resp, nil
	}, &options.ListOptions)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) fillCache() error {
	if _, found := c.cache.Get(orgsCacheKeys); found {
		return nil
	}

	orgs, err := c.Groups()
	if err != nil {
		return err
	}

	c.cache.Set(orgsCacheKeys, orgs, cache.NoExpiration)

	return nil
}

func (c *Client) Groups() ([]*gitlab.Group, error) {
	if groups, found := c.cache.Get(orgsCacheKeys); found {
		return groups.([]*gitlab.Group), nil
	}

	var result []*gitlab.Group

	dummy := true
	for _, group := range c.orgs {
		options := gitlab.ListGroupsOptions{Owned: &dummy, Search: &group}

		err := PaginateResults(func(opts *gitlab.ListOptions) (*gitlab.Response, error) {
			groups, resp, err := c.Client().Groups.ListGroups(&options)
			if err != nil {
				return nil, err
			}

			result = append(result, groups...)

			return resp, nil
		}, &options.ListOptions)

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
