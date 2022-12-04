package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	gitlab2 "github.com/xanzy/go-gitlab"
	"log"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type organizationCollector struct {
	collectors.BaseCollector
	Client  *gitlab.Client
	Context context.Context
}

func NewGroupCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &organizationCollector{
		Client:  client,
		Context: ctx,
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (c *organizationCollector) Namespace() namespace.Namespace {
	return namespace.Organization
}

func (c *organizationCollector) CollectMetadata() collectors.Metadata {
	_, response, err := c.Client.Client().Groups.ListGroups(&gitlab2.ListGroupsOptions{})
	res := collectors.Metadata{}

	if err != nil {
		log.Printf("failed to collect groups %s", err)
	} else {
		res.TotalEntities = response.TotalItems
	}

	return res
}

func (c *organizationCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		groups, err := c.Client.Groups()
		if err != nil {
			log.Printf("failed to collect groups %s", err)
			return
		}

		for _, g := range groups {
			entity := gitlab_collected.Organization{Group: *g}
			c.CollectDataWithContext(&entity, g.WebURL, newCollectionContext(g, []permissions.OrganizationRole{permissions.RepoRoleAdmin}))
			c.CollectionChangeByOne()
		}
	})
}
