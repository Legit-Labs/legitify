package collectors

import (
	"log"

	"github.com/google/go-github/v44/github"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/shurcooL/githubv4"
	"golang.org/x/net/context"
)

type organizationCollector struct {
	baseCollector
	Client  ghclient.Client
	Context context.Context
}

var orgSamlQuery struct {
	Organization struct {
		SamlIdentityProvider struct {
			ExternalIdentities struct {
				TotalCount int
			} `graphql:"externalIdentities(first: 1)"`
		}
	} `graphql:"organization(login: $login)"`
}

func newOrganizationCollector(ctx context.Context, client ghclient.Client) collector {
	c := &organizationCollector{
		Client:  client,
		Context: ctx,
	}
	initBaseCollector(&c.baseCollector, c)
	return c
}

func (c *organizationCollector) Namespace() namespace.Namespace {
	return namespace.Organization
}

func (c *organizationCollector) CollectMetadata() Metadata {
	orgs, err := c.Client.CollectOrganizations()
	res := Metadata{}

	if err != nil {
		log.Printf("failed to collect organizations %s", err)
	} else {
		res.TotalEntities = len(orgs)
	}

	return res
}

func (c *organizationCollector) Collect() subCollectorChannels {
	return c.wrappedCollection(func() {
		orgs, err := c.Client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		c.totalCollectionChange(len(orgs))
		gw := group_waiter.New()
		for _, org := range orgs {
			org := org
			gw.Do(func() {
				extend := c.collectExtraData(&org)
				c.collectData(org, extend, *extend.Organization.HTMLURL, []permissions.Role{org.Role})
				c.collectionChangeByOne()
			})
		}
		gw.Wait()
	})
}

func (c *organizationCollector) collectExtraData(org *ghcollected.ExtendedOrg) ghcollected.Organization {
	samlEnabled, err := c.collectOrgSamlData(org.Name())

	if err != nil {
		samlEnabled = nil
		log.Printf("failed to collect saml data for %s, %s", org.Name(), err)
	}

	hooks, err := c.collectOrgWebhooks(org.Name())
	if err != nil {
		hooks = nil
		log.Printf("failed to collect webhooks data for %s, %s", org.Name(), err)
	}

	return ghcollected.Organization{
		Organization: org,
		SamlEnabled:  samlEnabled,
		Hooks:        hooks,
	}
}

func (c *organizationCollector) collectOrgWebhooks(org string) ([]*github.Hook, error) {
	var result []*github.Hook

	err := ghclient.PaginateResults(func(opts *github.ListOptions) (*github.Response, error) {
		hooks, resp, err := c.Client.Client().Organizations.ListHooks(c.Context, org, opts)
		if err != nil {
			if resp.Response.StatusCode == 404 {
				perm := newMissingPermission(permissions.OrgHookAdmin, org,
					"Cannot read organization webhooks", namespace.Organization)
				c.issueMissingPermissions(perm)
			}
			return nil, err
		}
		result = append(result, hooks...)
		return resp, nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *organizationCollector) collectOrgSamlData(org string) (*bool, error) {
	variables := map[string]interface{}{
		"login": githubv4.String(org),
	}

	err := c.Client.GraphQLClient().Query(c.Context, &orgSamlQuery, variables)

	if err != nil {
		return nil, err
	}

	samlEnabled := orgSamlQuery.Organization.SamlIdentityProvider.ExternalIdentities.TotalCount > 0

	return &samlEnabled, nil

}
