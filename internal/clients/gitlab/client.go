package gitlab

import (
	"context"

	"github.com/Legit-Labs/legitify/internal/clients/gitlab/pagination"
	"github.com/Legit-Labs/legitify/internal/clients/gitlab/transport"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/slice_utils"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/patrickmn/go-cache"
	"github.com/xanzy/go-gitlab"
)

const (
	orgsCacheKeys   = "orgs"
	allGroupsFilter = ""
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

func NewClient(ctx context.Context, token string, endpoint string, orgs []string) (*Client, error) {
	var config []gitlab.ClientOptionFunc
	if endpoint != "" {
		config = []gitlab.ClientOptionFunc{
			gitlab.WithBaseURL(endpoint),
			gitlab.WithHTTPClient(transport.NewHttpClient()),
		}
	}

	git, err := gitlab.NewClient(token, config...)
	if err != nil {
		return nil, err
	}

	if len(orgs) == 0 {
		orgs = []string{allGroupsFilter}
	}

	result := &Client{
		context: ctx,
		client:  git,
		cache:   cache.New(cache.NoExpiration, cache.NoExpiration),
		orgs:    orgs,
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
	maintainerPermissions := gitlab.MaintainerPermissions
	opts := gitlab.ListProjectsOptions{MinAccessLevel: &maintainerPermissions}
	mapper := func(projects []*gitlab.Project) []types.RepositoryWithOwner {
		if projects == nil {
			return []types.RepositoryWithOwner{}
		}
		return slice_utils.Map(projects, func(p *gitlab.Project) types.RepositoryWithOwner {
			return types.NewRepositoryWithOwner(p.PathWithNamespace, permissions.RepoRoleAdmin)
		})
	}
	result, err := pagination.NewMapper(c.Client().Projects.ListProjects, opts, mapper).Sync()

	if err != nil {
		return nil, err
	}
	return result.Collected, nil
}

func (c *Client) GroupMembers(group *gitlab.Group) ([]*gitlab.GroupMember, error) {
	result, err := pagination.New[*gitlab.GroupMember](c.Client().Groups.ListGroupMembers, nil).Sync(group.ID)
	if err != nil {
		return nil, err
	}

	return result.Collected, nil
}

func (c *Client) Groups() ([]*gitlab.Group, error) {
	if groups, found := c.cache.Get(orgsCacheKeys); found {
		return groups.([]*gitlab.Group), nil
	}

	var result []*gitlab.Group

	ownedGroups := true
	for _, group := range c.orgs {
		opts := &gitlab.ListGroupsOptions{Owned: &ownedGroups, Search: &group}
		res, err := pagination.New[*gitlab.Group](c.Client().Groups.ListGroups, opts).Sync()
		if err != nil {
			return nil, err
		}
		result = append(result, res.Collected...)
	}

	return result, nil
}

func (c *Client) GroupHooks(gid int) ([]*gitlab.GroupHook, error) {
	result, err := pagination.New[*gitlab.GroupHook](c.Client().Groups.ListGroupHooks, nil).Sync(gid)
	if err != nil {
		return nil, err
	}

	return result.Collected, nil
}
