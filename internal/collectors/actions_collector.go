package collectors

import (
	"fmt"
	"log"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"golang.org/x/net/context"
)

const (
	orgActionPermEffect = "Cannot read organization actions"
)

type actionCollector struct {
	baseCollector
	client  ghclient.Client
	context context.Context
}

func newActionCollector(ctx context.Context, client ghclient.Client) collector {
	c := &actionCollector{
		client:  client,
		context: ctx,
	}
	initBaseCollector(&c.baseCollector, c)
	return c
}

func (c *actionCollector) Namespace() namespace.Namespace {
	return namespace.Actions
}

func (c *actionCollector) CollectMetadata() Metadata {
	orgs, err := c.client.CollectOrganizations()
	res := Metadata{}

	if err != nil {
		log.Printf("failed to collect organizations %s", err)
	} else {
		res.TotalEntities = len(orgs)
	}

	return res
}

func (c *actionCollector) Collect() subCollectorChannels {
	return c.wrappedCollection(func() {
		orgs, err := c.client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		c.totalCollectionChange(len(orgs))

		for _, org := range orgs {
			actionsData, _, err := c.client.Client().Organizations.GetActionsPermissions(c.context, *org.Login)

			if err != nil {
				entityName := fmt.Sprintf("%s/%s", namespace.Organization, *org.Login)
				perm := newMissingPermission(permissions.OrgAdmin, entityName, orgActionPermEffect, namespace.Organization)
				c.issueMissingPermissions(perm)
			}

			c.collectionChangeByOne()

			c.collectData(org,
				ghcollected.OrganizationActions{
					Organization:       org,
					ActionsPermissions: actionsData,
				},
				*org.HTMLURL,
				[]permissions.Role{org.Role})
		}
	})
}
