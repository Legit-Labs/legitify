package github

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/Legit-Labs/legitify/internal/clients/github/pagination"
	"github.com/Legit-Labs/legitify/internal/clients/github/transport"
	"github.com/Legit-Labs/legitify/internal/clients/github/types"
	commontransport "github.com/Legit-Labs/legitify/internal/clients/transport"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/slice_utils"
	commontypes "github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/Legit-Labs/legitify/internal/screen"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"

	gh "github.com/google/go-github/v53/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const scopeHttpHeader = "X-OAuth-Scopes"

type Client struct {
	client           *gh.Client
	orgs             []string
	graphQLClient    *githubv4.Client
	context          context.Context
	scopes           permissions.TokenScopes
	graphQLRawClient *http.Client
	serverUrl        string
	once             sync.Once
	enterprises      []string
}

func NewClient(ctx context.Context, token string, githubEndpoint string, org []string, enterprises []string) (*Client, error) {
	client := &Client{
		orgs:        org,
		context:     ctx,
		serverUrl:   strings.TrimRight(githubEndpoint, "/"),
		enterprises: enterprises,
	}

	if err := client.initClients(ctx, token); err != nil {
		return nil, err
	}

	scopes, err := client.collectTokenScopes()
	if err != nil {
		return nil, err
	}
	client.scopes = scopes

	client.printInstanceTypeMessage()

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

	rawClient, graphQLRawClient, err := newHttpClients(ctx, token)
	if err != nil {
		return err
	}

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

func (c *Client) Scopes() permissions.TokenScopes {
	return c.scopes
}

func (c *Client) Orgs() []string {
	return c.orgs
}

func (c *Client) inRealOrgs(org string, realOrgs []string) bool {
	for _, realOrg := range realOrgs {
		if strings.EqualFold(org, realOrg) {
			return true
		}
	}

	return false
}

func (c *Client) setOrgsList(realOrgs []string) error {
	if len(c.orgs) == 0 {
		c.orgs = realOrgs
		return nil
	}

	for _, userRequestedOrg := range c.orgs {
		if !c.inRealOrgs(userRequestedOrg, realOrgs) {
			return fmt.Errorf("user has no access to the requested organization: %s", userRequestedOrg)
		}
	}

	return nil
}

func (c *Client) CollectOrganizations() ([]githubcollected.ExtendedOrg, error) {
	realOrgs, err := c.collectOrgsList()
	if err != nil {
		return nil, err
	}
	c.once.Do(func() {
		err = c.setOrgsList(realOrgs)
	})
	if err != nil {
		return nil, err
	}

	orgs, err := c.collectSpecificOrganizations()
	if err != nil {
		return nil, err
	}

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
			ID:   int(*o.ID),
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
	mapper := func(orgs []*gh.Organization) []string {
		if orgs == nil {
			return []string{}
		}
		return slice_utils.Map(orgs, func(o *gh.Organization) string {
			return *o.Login
		})
	}
	res, err := pagination.NewMapper(c.Client().Organizations.List, nil, mapper).Sync(c.context, "")
	if err != nil {
		return nil, err
	}

	return res.Collected, nil
}

func (c *Client) collectSpecificOrganizations() ([]githubcollected.ExtendedOrg, error) {
	res := make([]githubcollected.ExtendedOrg, 0, len(c.orgs))

	for _, o := range c.orgs {
		org, err := c.Organization(o)

		if err != nil {
			log.Printf("failed to list org %v: %v", o, err)
			continue
		}

		res = append(res, *org)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("could not list any organization")
	}

	return res, nil
}

func (c *Client) Organization(login string) (*githubcollected.ExtendedOrg, error) {
	org, _, err := c.Client().Organizations.Get(c.context, login)

	if err != nil {
		return nil, err
	}

	role, err := c.getRole(*org.Login)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	result := githubcollected.NewExtendedOrg(org, role)

	return &result, nil
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

func (c *Client) printInstanceTypeMessage() {
	var instanceTypeMsg string
	if c.IsGithubCloud() {
		instanceTypeMsg = "Using Github Cloud"
	} else {
		instanceTypeMsg = fmt.Sprintf("Using Github Enterprise Endpoint: %s", c.serverUrl)
	}
	screen.Printf("%s\n", instanceTypeMsg)
}

type samlError struct {
	organization string
}

func (se *samlError) Error() string {
	return fmt.Sprintf("Token is not SAML authorized for organization: %s.\nPlease go to https://github.com/settings/tokens and authorize.", se.organization)
}

func newHttpClients(ctx context.Context, token string) (client *http.Client, graphQL *http.Client, err error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := &oauth2.Transport{
		Base:   commontransport.NewCacheTransport(),
		Source: ts,
	}

	rateLimitWaiter, err := transport.NewRateLimitWaiter(ctx, tc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create rate limiter: %v", err)
	}

	clientWithSecondaryRateLimit := commontransport.NewCacheTracker(rateLimitWaiter)
	clientWithAcceptHeader := transport.NewGraphQL(rateLimitWaiter.Transport)

	return clientWithSecondaryRateLimit, clientWithAcceptHeader, nil
}

var enterpriseQuery struct {
	Enterprise struct {
		OwnerInfo struct {
			MembersCanChangeRepositoryVisibilitySetting   string
			AllowPrivateRepositoryForkingSetting          string
			MembersCanInviteCollaboratorsSetting          string
			TwoFactorRequiredSetting                      string
			MembersCanCreatePublicRepositoriesSetting     bool
			DefaultRepositoryPermissionSetting            string
			MembersCanDeleteRepositoriesSetting           string
			NotificationDeliveryRestrictionEnabledSetting string
			SamlIdentityProvider                          struct {
				ExternalIdentities struct {
					TotalCount int
				} `graphql:"externalIdentities(first: 1)"`
			}
		}
		Name          string
		Url           string
		Id            string
		DatabaseId    int64
		ViewerIsAdmin bool
	} `graphql:"enterprise(slug: $slug)"`
}

func (c *Client) CollectEnterprises() ([]githubcollected.Enterprise, error) {
	if len(c.enterprises) == 0 {
		return nil, nil
	}
	enterprises, err := c.collectSpecificEnterprises()
	if err != nil {
		return nil, err
	}

	return enterprises, nil
}

func (c *Client) collectSpecificEnterprises() ([]githubcollected.Enterprise, error) {
	res := make([]githubcollected.Enterprise, 0, len(c.enterprises))

	for _, enterprise := range c.enterprises {

		variables := map[string]interface{}{
			"slug": githubv4.String(enterprise),
		}

		err := c.GraphQLClient().Query(c.context, &enterpriseQuery, variables)
		if err != nil {
			log.Printf("failed to get enterprise %v: %v", enterprise, err)
			return nil, err
		}
		if enterpriseQuery.Enterprise.DatabaseId == 0 {
			log.Printf("Failed to get enterprise %v . User is not a member of this enterprise", enterprise)
			return nil, err
		}
		samlEnabled := enterpriseQuery.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.TotalCount > 0
		codeAndSecurityPolicySettings, err := c.GetSecurityAndAnalysisForEnterprise(enterprise)
		if err != nil {
			log.Printf("failed to get code security settings for enterprise %v: %v", enterprise, err)
		}
		newEnter := githubcollected.NewEnterprise(
			enterpriseQuery.Enterprise.OwnerInfo.MembersCanChangeRepositoryVisibilitySetting,
			enterpriseQuery.Enterprise.Name,
			enterpriseQuery.Enterprise.Url,
			enterpriseQuery.Enterprise.DatabaseId,
			enterpriseQuery.Enterprise.ViewerIsAdmin,
			enterpriseQuery.Enterprise.OwnerInfo.AllowPrivateRepositoryForkingSetting,
			enterpriseQuery.Enterprise.OwnerInfo.MembersCanInviteCollaboratorsSetting,
			enterpriseQuery.Enterprise.OwnerInfo.MembersCanCreatePublicRepositoriesSetting,
			enterpriseQuery.Enterprise.OwnerInfo.TwoFactorRequiredSetting,
			enterpriseQuery.Enterprise.OwnerInfo.DefaultRepositoryPermissionSetting,
			enterpriseQuery.Enterprise.OwnerInfo.MembersCanDeleteRepositoriesSetting,
			enterpriseQuery.Enterprise.OwnerInfo.NotificationDeliveryRestrictionEnabledSetting,
			samlEnabled,
			codeAndSecurityPolicySettings)
		res = append(res, newEnter)

	}

	return res, nil
}

func (c *Client) GetRulesForBranch(organization, repository, branch string) ([]*types.RepositoryRule, error) {
	url := fmt.Sprintf("repos/%v/%v/rules/branches/%v", organization, repository, branch)
	req, err := c.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var rules []*types.RepositoryRule
	_, err = c.client.Do(c.context, req, &rules)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		specific, _, err := c.Client().Repositories.GetRuleset(c.context, organization, repository, rule.Id, true)
		if err != nil {
			continue
		}

		rule.Ruleset = specific
	}

	return rules, nil

}

func (c *Client) GetSecurityAndAnalysisForEnterprise(enterprise string) (*types.AnalysisAndSecurityPolicies, error) {
	url := fmt.Sprintf("/api/v3/enterprises/%v/code_security_and_analysis", enterprise)
	req, err := c.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var p types.AnalysisAndSecurityPolicies
	_, err = c.client.Do(c.context, req, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
