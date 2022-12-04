package github

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Legit-Labs/legitify/internal/clients/github/types"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	commontypes "github.com/Legit-Labs/legitify/internal/common/types"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"

	gh "github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const experimentalApiAcceptHeader = "application/vnd.github.hawkgirl-preview+json"
const scopeHttpHeader = "X-OAuth-Scopes"

type Client struct {
	client           *gh.Client
	orgs             []string
	graphQLClient    *githubv4.Client
	context          context.Context
	orgsCache        []githubcollected.ExtendedOrg
	cacheLock        sync.RWMutex
	scopes           permissions.TokenScopes
	graphQLRawClient *http.Client
	serverUrl        string
}

func isBadRequest(err error) bool {
	return err.Error() == "Bad credentials"
}

func newHttpClients(ctx context.Context, token string) (client *http.Client, graphQL *http.Client) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	acceptHeader := experimentalApiAcceptHeader
	clientWithAcceptHeader := NewClientWithAcceptHeader(tc.Transport, &acceptHeader)

	return tc, clientWithAcceptHeader
}

func NewClient(ctx context.Context, token string, githubEndpoint string, org []string, fillCache bool) (*Client, error) {
	client := &Client{
		orgs:      org,
		context:   ctx,
		serverUrl: strings.TrimRight(githubEndpoint, "/"),
	}

	if err := client.initClients(ctx, token); err != nil {
		return nil, err
	}

	scopes, err := client.collectTokenScopes()
	if err != nil {
		return nil, err
	}
	client.scopes = scopes

	if fillCache {
		if err := client.fillCache(); err != nil {
			return nil, err
		}
	}

	if client.IsGithubCloud() {
		log.Printf("Using Github Cloud")
	} else {
		log.Printf("Using Github Enterprise Endpoint: %s\n\n", client.serverUrl)
	}

	return client, nil
}

func (c *Client) Client() *gh.Client {
	return c.client
}

func (c *Client) GraphQLClient() *githubv4.Client {
	return c.graphQLClient
}

func (c *Client) IsGithubCloud() bool {
	return c.serverUrl == ""
}

func (c *Client) initClients(ctx context.Context, token string) error {
	if err := c.validateToken(token); err != nil {
		return err
	}

	var ghClient *gh.Client
	var graphQLClient *githubv4.Client

	rawClient, graphQLRawClient := newHttpClients(ctx, token)
	if c.IsGithubCloud() {
		ghClient = gh.NewClient(rawClient)
		graphQLClient = githubv4.NewClient(graphQLRawClient)
	} else {
		var err error
		ghClient, err = gh.NewEnterpriseClient(c.serverUrl, c.serverUrl, rawClient)
		if err != nil {
			return err
		}
		graphQLClient = githubv4.NewEnterpriseClient(c.getGitHubGraphURL(), graphQLRawClient)

	}

	c.graphQLRawClient = graphQLRawClient
	c.client = ghClient
	c.graphQLClient = graphQLClient
	return nil
}

// Note: tokens before April 2021 did not have the ghp_ prefix.
var githubTokenPattern = regexp.MustCompile("(ghp_)?[A-Za-z0-9_]{36}")

func (c *Client) validateToken(token string) error {
	if token == "" {
		return fmt.Errorf("missing token")
	} else if strings.HasPrefix(token, "github_pat_") {
		return fmt.Errorf("GitHub fine-grained tokens are not supported at this moment, please use classic PAT")
	} else if !githubTokenPattern.MatchString(token) {
		return fmt.Errorf("GitHub token seems invalid (expected pattern: '%v')", githubTokenPattern)
	}

	return nil
}

func (c *Client) getGitHubGraphURL() string {
	if c.IsGithubCloud() {
		return "https://api.github.com/graphql"
	}

	return c.serverUrl + "/api/graphql"
}

func (c *Client) fillCache() error {
	_, err := c.CollectOrganizations()
	if err != nil && isBadRequest(err) {
		return fmt.Errorf("invalid token (make sure it's not expired or revoked)")
	}

	if len(c.orgsCache) == 0 {
		if len(c.orgs) != 0 {
			return fmt.Errorf("token doesn't have access to the requested organizations")
		} else {
			return fmt.Errorf("token doesn't have access to any organization")
		}
	}

	return nil
}

func (c *Client) Scopes() permissions.TokenScopes {
	return c.scopes
}

func (c *Client) Orgs() []string {
	return c.orgs
}

func (c *Client) setOrgsList(realOrgs []string) error {
	if len(c.orgs) == 0 {
		c.orgs = realOrgs
	} else {
		for _, userRequestedOrg := range c.orgs {
			inRealOrgs := false
			for _, realOrg := range realOrgs {
				if strings.EqualFold(userRequestedOrg, realOrg) {
					inRealOrgs = true
					break
				}
			}
			if !inRealOrgs {
				return fmt.Errorf("User has no access to the requested organization: %s\n", userRequestedOrg)
			}
		}
	}

	return nil
}

func (c *Client) CollectOrganizations() ([]githubcollected.ExtendedOrg, error) {
	c.cacheLock.RLock()
	if c.orgsCache != nil {
		return c.orgsCache, nil
	}
	c.cacheLock.RUnlock()

	realOrgs, err := c.collectOrgsList()
	if err != nil {
		return nil, err
	}
	if err := c.setOrgsList(realOrgs); err != nil {
		return nil, err
	}

	orgs, err := c.collectSpecificOrganizations()
	if err != nil {
		return nil, err
	}

	c.cacheLock.Lock()
	c.orgsCache = orgs
	c.cacheLock.Unlock()

	return orgs, nil
}

func (c *Client) Organizations() ([]commontypes.Organization, error) {
	raw, err := c.CollectOrganizations()
	if err != nil {
		return nil, err
	}

	var result []commontypes.Organization
	for _, o := range raw {
		result = append(result, commontypes.Organization{
			Name: o.Name(),
			Role: o.Role,
		})
	}

	return result, nil
}

type orgPermissionQuery struct {
	Organization struct {
		ViewerCanAdminister *bool `graphql:"viewerCanAdminister"`
		// we need to fetch the repositories as well to test whether the token is SAML authorized
		Repositories struct {
			Nodes []struct {
				Id *githubv4.String
			}
		} `graphql:"repositories(first: 1)"`
	} `graphql:"organization(login: $login)"`
}

func isMissingScopeError(err error) bool {
	const msg = "Your token has not been granted the required scopes to execute this query"
	return strings.HasPrefix(err.Error(), msg)
}

const samlErrorMsg = "Resource protected by organization SAML enforcement. " +
	"You must grant your Personal Access token access to this organization."

func isMissingSamlAuthenticationError(err error) bool {
	return err != nil && err.Error() == samlErrorMsg
}

func (c *Client) getRole(orgName string) (permissions.OrganizationRole, error) {
	variables := map[string]interface{}{
		"login": githubv4.String(orgName),
	}
	query := orgPermissionQuery{}

	if err := c.GraphQLClient().Query(c.context, &query, variables); err != nil {
		if isMissingSamlAuthenticationError(err) {
			return permissions.OrgRoleNone, &samlError{organization: orgName}
		}

		if isMissingScopeError(err) {
			// In case the token is missing org:read, default to member.
			// We only list organizations of which the user is a member.
			return permissions.OrgRoleMember, nil
		}

		return permissions.OrgRoleNone, err
	}

	return permissions.GetOrgRole(query.Organization.ViewerCanAdminister), nil
}

func (c *Client) collectTokenScopes() (permissions.TokenScopes, error) {
	var buf bytes.Buffer
	resp, err := c.graphQLRawClient.Post(c.getGitHubGraphURL(), "application/json", &buf)
	if err != nil {
		return nil, err
	}

	scopesList := resp.Header.Get(scopeHttpHeader)
	parsed := strings.Split(scopesList, ", ")
	scopes := permissions.ParseTokenScopes(parsed)

	return scopes, nil
}

func (c *Client) collectOrgsList() ([]string, error) {
	var orgNames []string
	err := PaginateResults(func(opts *gh.ListOptions) (*gh.Response, error) {
		orgs, resp, err := c.Client().Organizations.List(c.context, "", opts)

		if err != nil {
			return nil, err
		}

		for _, o := range orgs {
			// The list-organizations API does not return all information,
			// so we only use it to pull the names
			orgNames = append(orgNames, *o.Login)
		}

		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	return orgNames, nil
}

func (c *Client) collectSpecificOrganizations() ([]githubcollected.ExtendedOrg, error) {
	res := make([]githubcollected.ExtendedOrg, 0)

	for _, o := range c.orgs {
		org, _, err := c.Client().Organizations.Get(c.context, o)

		if err != nil {
			return nil, err
		}

		role, err := c.getRole(*org.Login)
		if err != nil {
			log.Println(err.Error())
		} else {
			res = append(res, githubcollected.NewExtendedOrg(org, role))
		}
	}

	return res, nil
}

func (c *Client) GetActionsTokenPermissionsForOrganization(organization string) (*types.TokenPermissions, error) {
	u := fmt.Sprintf("orgs/%s/actions/permissions/workflow", organization)
	return c.GetActionsTokenPermissions(u)
}

func (c *Client) GetActionsTokenPermissionsForRepository(organization string, repository string) (*types.TokenPermissions, error) {
	u := fmt.Sprintf("repos/%s/%s/actions/permissions/workflow", organization, repository)
	return c.GetActionsTokenPermissions(u)
}

func (c *Client) GetActionsTokenPermissions(url string) (*types.TokenPermissions, error) {
	req, err := c.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	p := types.TokenPermissions{}
	_, err = c.client.Do(c.context, req, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) IsAnalyzable(repository commontypes.RepositoryWithOwner) (bool, error) {
	var repo struct {
		Repository struct {
			ViewerPermission githubv4.String
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"name":  githubv4.String(repository.Name),
		"owner": githubv4.String(repository.Owner),
	}

	err := c.GraphQLClient().Query(c.context, &repo, variables)
	if err != nil {
		return false, err
	}

	return repo.Repository.ViewerPermission == permissions.RepoRoleAdmin, nil
}

func uniqueRepositories(slice []commontypes.RepositoryWithOwner) []commontypes.RepositoryWithOwner {
	keys := make(map[string]bool)
	var list []commontypes.RepositoryWithOwner
	for _, entry := range slice {
		key := entry.String()
		if _, found := keys[key]; !found {
			keys[key] = true
			list = append(list, entry)
		}
	}
	return list
}

func (c *Client) Repositories() ([]commontypes.RepositoryWithOwner, error) {
	r1, err := c.getViewerRepositories()
	if err != nil {
		return nil, err
	}

	r2, err := c.getOrganizationsRepositories()
	if err != nil {
		return nil, err
	}

	return uniqueRepositories(append(r1, r2...)), nil
}

func (c *Client) getViewerRepositories() ([]commontypes.RepositoryWithOwner, error) {
	var repositories []commontypes.RepositoryWithOwner
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
		err := c.GraphQLClient().Query(c.context, &query, variables)
		if err != nil {
			return nil, err
		}
		for _, r := range query.Viewer.Repositories.Nodes {
			repositories = append(repositories, commontypes.NewRepositoryWithOwner(r.NameWithOwner, r.ViewerPermission))
		}

		if !query.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = query.Viewer.Repositories.PageInfo.EndCursor
	}

	return repositories, nil
}

func (c *Client) getOrganizationsRepositories() ([]commontypes.RepositoryWithOwner, error) {
	var repositories []commontypes.RepositoryWithOwner
	orgs, err := c.CollectOrganizations()
	if err != nil {
		return nil, err
	}

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
				"login":  githubv4.String(o.Name()),
			}

			for {
				err := c.GraphQLClient().Query(c.context, &query, variables)
				if err != nil {
					return
				}

				for _, r := range query.Organization.Repositories.Nodes {
					repositories = append(repositories, commontypes.NewRepositoryWithOwner(r.NameWithOwner, r.ViewerPermission))
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

type samlError struct {
	organization string
}

func (se *samlError) Error() string {
	return fmt.Sprintf("Token is not SAML authorized for organization: %s.\nPlease go to https://github.com/settings/tokens and authorize.", se.organization)
}
