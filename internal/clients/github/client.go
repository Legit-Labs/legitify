package github

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"

	"github.com/google/go-github/v44/github"
	gh "github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client interface {
	Client() *gh.Client
	GraphQLClient() *githubv4.Client
	CollectOrganizations() ([]githubcollected.ExtendedOrg, error)
	Scopes() permissions.TokenScopes
	Orgs() []string
}

const experimentalApiAcceptHeader = "application/vnd.github.hawkgirl-preview+json"
const scopeHttpHeader = "X-OAuth-Scopes"

type client struct {
	client        *gh.Client
	orgs          []string
	graphQLClient *githubv4.Client
	context       context.Context
	orgsCache     []githubcollected.ExtendedOrg
	cacheLock     sync.RWMutex
	scopes        permissions.TokenScopes
	rawClient     *http.Client
}

func IsTokenValid(token string) error {
	if token == "" {
		return fmt.Errorf("missing token")
	} else if len(token) != 40 {
		return fmt.Errorf("GitHub token seems invalid (should have 40 characters)")
	} else if !strings.HasPrefix(token, "ghp_") {
		return fmt.Errorf("GitHub token seems invalid (should start with \"ghp_\"")
	}

	return nil
}

func isBadRequest(err error) bool {
	return err.Error() == "Bad credentials"
}

func NewClient(ctx context.Context, token string, org []string) (Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token must be provided")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	ghClient := gh.NewClient(tc)

	acceptHeader := experimentalApiAcceptHeader
	clientWithAcceptHeader := NewClientWithAcceptHeader(tc.Transport, &acceptHeader)
	graphQLClient := githubv4.NewClient(&clientWithAcceptHeader)

	client := &client{
		client:        ghClient,
		orgs:          org,
		graphQLClient: graphQLClient,
		context:       ctx,
		rawClient:     &clientWithAcceptHeader,
	}

	// fill cache & token scopes
	_, err := client.CollectOrganizations()
	if err != nil && isBadRequest(err) {
		return nil, fmt.Errorf("invalid token (make sure it's not expired or revoked)")
	}

	if len(client.orgsCache) == 0 {
		if len(org) != 0 {
			return nil, fmt.Errorf("token doesn't have access to the requsted organizations")
		} else {
			return nil, fmt.Errorf("token doesn't have access to any organization")
		}
	}

	return client, nil
}

func (c *client) Client() *gh.Client {
	return c.client
}

func (c *client) GraphQLClient() *githubv4.Client {
	return c.graphQLClient

}

func (c *client) Scopes() permissions.TokenScopes {
	return c.scopes
}

func (c *client) Orgs() []string {
	return c.orgs
}

func (c *client) setOrgsList(realOrgs []string) error {
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

func (c *client) CollectOrganizations() ([]githubcollected.ExtendedOrg, error) {
	c.cacheLock.RLock()
	if c.orgsCache != nil {
		return c.orgsCache, nil
	}
	c.cacheLock.RUnlock()

	scopes, err := c.collectTokenScopes()
	if err != nil {
		return nil, err
	}
	c.scopes = scopes

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

func (c *client) getRole(orgName string) (permissions.OrganizationRole, error) {
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

func (c *client) collectTokenScopes() (permissions.TokenScopes, error) {
	graphQLUrl := "https://api.github.com/graphql"
	var buf bytes.Buffer
	resp, err := c.rawClient.Post(graphQLUrl, "application/json", &buf)
	if err != nil {
		return nil, err
	}

	scopesList := resp.Header.Get(scopeHttpHeader)
	parsed := strings.Split(scopesList, ", ")
	scopes := permissions.ParseTokenScopes(parsed)

	return scopes, nil
}

func (c *client) collectOrgsList() ([]string, error) {
	var orgNames []string
	err := PaginateResults(func(opts *github.ListOptions) (*github.Response, error) {
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

func (c *client) collectSpecificOrganizations() ([]githubcollected.ExtendedOrg, error) {
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

type samlError struct {
	organization string
}

func (se *samlError) Error() string {
	return fmt.Sprintf("Token is not SAML authorized for organization: %s.\nPlease go to https://github.com/settings/tokens and authorize.", se.organization)
}
