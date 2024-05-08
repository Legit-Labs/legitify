package github

import (
	"context"
	"log"
	"net/http"

	"github.com/Legit-Labs/legitify/internal/collectors"

	"github.com/google/go-github/v53/github"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/clients/github/pagination"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/shurcooL/githubv4"
)

type organizationCollector struct {
	collectors.BaseCollector
	Client  *ghclient.Client
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

func NewOrganizationCollector(ctx context.Context, client *ghclient.Client) collectors.Collector {
	c := &organizationCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.Organization),
		Client:        client,
		Context:       ctx,
	}
	return c
}

func (c *organizationCollector) CollectTotalEntities() int {
	orgs, err := c.Client.CollectOrganizations()
	if err != nil {
		log.Printf("failed to collect organizations %s", err)
		return 0
	}

	return len(orgs)
}

func (c *organizationCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		orgs, err := c.Client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		gw := group_waiter.New()
		for _, org := range orgs {
			org := org
			gw.Do(func() {
				extend := c.collectExtraData(&org)
				c.CollectData(org, extend, *extend.Organization.HTMLURL, []permissions.Role{org.Role})
				c.CollectionChangeByOne()
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
		perm := collectors.NewMissingPermission(permissions.OrgHookAdmin, org.Name(),
			"Cannot collect webhook for organization", namespace.Organization)
		c.IssueMissingPermissions(perm)
	}

	secrets := c.collectOrgSecrets(org.Name())

	return ghcollected.Organization{
		Organization: org,
		SamlEnabled:  samlEnabled,
		Hooks:        hooks,
		OrgSecrets:   secrets,
	}
}

func (c *organizationCollector) collectOrgWebhooks(org string) ([]*github.Hook, error) {
	res, err := pagination.New[*github.Hook](c.Client.Client().Organizations.ListHooks, nil).Sync(c.Context, org)
	if err != nil {
		if res.Resp.Response.StatusCode == http.StatusNotFound {
			perm := collectors.NewMissingPermission(permissions.OrgHookAdmin, org,
				"Cannot read organization webhooks", namespace.Organization)
			c.IssueMissingPermissions(perm)
		}
		return nil, err
	}

	return res.Collected, nil
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

func (c *organizationCollector) collectOrgSecrets(org string) []*ghcollected.OrganizationSecret {
	secrets, err := c.Client.GetOrganizationSecrets(org)
	if err != nil {
		return nil
	}
	var orgSecrets []*ghcollected.OrganizationSecret
	for i := 0; i < len(secrets.Secrets); i++ {
		orgSecrets = append(orgSecrets, &ghcollected.OrganizationSecret{secrets.Secrets[i].Name, int(secrets.Secrets[i].UpdatedAt.Time.UnixNano())})
	}

	return orgSecrets
}
