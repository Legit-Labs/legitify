package gitlab

import (
	"log"

	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	gitlab2 "github.com/xanzy/go-gitlab"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type groupCollector struct {
	collectors.BaseCollector
	Client  *gitlab.Client
	Context context.Context
}

func NewGroupCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &groupCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.Organization),
		Client:        client,
		Context:       ctx,
	}
	return c
}

func (c *groupCollector) CollectTotalEntities() int {
	groups, err := c.Client.Groups()
	if err != nil {
		log.Printf("failed to collect groups %s", err)
		return 0
	}

	return len(groups)
}

func (c *groupCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		groups, err := c.Client.Groups()
		if err != nil {
			log.Printf("failed to collect groups %s", err)
			return
		}

		gw := group_waiter.New()

		for _, g := range groups {
			g := g
			gw.Do(func() {
				fullGroup, _, err := c.Client.Client().Groups.GetGroup(g.ID, &gitlab2.GetGroupOptions{})
				if err != nil {
					log.Printf("failed to query group: %d - %s", g.ID, g.Name)
					return
				}

				isPremium := c.Client.IsGroupPremium(g)

				hooks, err := c.Client.GroupHooks(fullGroup.ID)

				if err != nil {
					log.Printf("failed to query group hooks: %d - %s", g.ID, g.Name)
				}

				entity := gitlab_collected.Organization{
					Group: fullGroup,
					Hooks: hooks,
				}

				c.CollectDataWithContext(entity, g.WebURL,
					newCollectionContext(g, []permissions.RepositoryRole{permissions.RepoRoleAdmin}, isPremium))
				c.CollectionChangeByOne()
			})
		}

		gw.Wait()
	})
}
