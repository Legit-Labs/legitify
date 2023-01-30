package github

import (
	"fmt"
	"log"

	"github.com/Legit-Labs/legitify/internal/collectors"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"golang.org/x/net/context"
)

const (
	orgActionPermEffect = "Cannot read organization actions settings"
)

type actionCollector struct {
	collectors.BaseCollector
	client  *ghclient.Client
	context context.Context
}

func NewActionCollector(ctx context.Context, client *ghclient.Client) collectors.Collector {
	c := &actionCollector{
		client:  client,
		context: ctx,
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (c *actionCollector) Namespace() namespace.Namespace {
	return namespace.Actions
}

func (c *actionCollector) CollectTotalEntities() int {
	orgs, err := c.client.CollectOrganizations()
	if err != nil {
		log.Printf("failed to collect organizations %s", err)
		return 0
	}

	return len(orgs)
}

func (c *actionCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		orgs, err := c.client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		for _, org := range orgs {
			actionsPermissions, err1 := c.client.GetActionsTokenPermissionsForOrganization(org.Name())
			actionsData, _, err2 := c.client.Client().Organizations.GetActionsPermissions(c.context, org.Name())

			if err1 != nil || err2 != nil {
				entityName := fmt.Sprintf("%s/%s", namespace.Organization, org.Name())
				perm := collectors.NewMissingPermission(permissions.OrgAdmin, entityName, orgActionPermEffect, namespace.Organization)
				c.IssueMissingPermissions(perm)
			}

			c.CollectionChangeByOne()

			c.CollectData(org,
				ghcollected.OrganizationActions{
					Organization:       org,
					ActionsPermissions: actionsData,
					TokenPermissions:   actionsPermissions,
				},
				org.CanonicalLink(),
				[]permissions.Role{org.Role})
		}
	})
}
