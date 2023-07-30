package gitlab

import (
	"context"
	"fmt"

	"github.com/Legit-Labs/legitify/internal/clients/gitlab/pagination"
	"github.com/Legit-Labs/legitify/internal/clients/gitlab/transport"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/slice_utils"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/xanzy/go-gitlab"
)

const (
	allGroupsFilter = ""
)

type Client struct {
	context  context.Context
	client   *gitlab.Client
	orgs     []string
	isAdmin  bool
	endpoint string
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
		context:  ctx,
		client:   git,
		orgs:     orgs,
		isAdmin:  IsAdmin(git),
		endpoint: endpoint,
	}

	return result, nil
}

func (c *Client) ServerUrl() string {
	return c.endpoint
}

func (c *Client) IsServer() bool {
	return c.endpoint != ""
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
			ID:   g.ID,
		})
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) Repositories() ([]types.RepositoryWithOwner, error) {
	opts := &gitlab.ListProjectsOptions{}
	if !c.IsAdmin() {
		opts.MinAccessLevel = gitlab.AccessLevel(gitlab.MaintainerPermissions)
	}
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

func (c *Client) IsAdmin() bool {
	return c.isAdmin
}

func IsAdmin(client *gitlab.Client) bool {
	res, _, err := client.Users.CurrentUser()
	if err != nil {
		return false // assume false on error
	}
	return res.IsAdmin
}

func (c *Client) Group(name string) (*gitlab.Group, error) {
	ownedGroups := !c.IsAdmin() // list all groups as site admin
	opts := &gitlab.ListGroupsOptions{
		Owned:  &ownedGroups,
		Search: &name,
	}

	res, err := pagination.New[*gitlab.Group](c.Client().Groups.ListGroups, opts).Sync()
	if err != nil {
		return nil, err
	}

	for _, g := range res.Collected {
		if g.Path == name {
			return g, nil
		}
	}

	return nil, fmt.Errorf("couldn't find group %s", name)
}

func (c *Client) Groups() ([]*gitlab.Group, error) {
	var result []*gitlab.Group

	ownedGroups := !c.IsAdmin() // list all groups as site admin
	for _, group := range c.orgs {
		opts := &gitlab.ListGroupsOptions{
			Owned:  &ownedGroups,
			Search: &group,
		}
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

func (c *Client) GroupPlan(group *gitlab.Group) (string, error) {
	nss, _, err := c.Client().Namespaces.SearchNamespace(group.FullName)
	if err != nil {
		return "", err
	}

	for _, n := range nss {
		if n.Path == group.Path {
			return n.Plan, nil
		}
	}

	return "", fmt.Errorf("didn't find namespace for %s", group.Path)
}

func (c *Client) IsGroupPremium(group *gitlab.Group) bool {
	plan, err := c.GroupPlan(group)
	if err != nil {
		return false
	}

	return plan != "free"
}
