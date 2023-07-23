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

type userCollector struct {
	collectors.BaseCollector
	Client  *gitlab.Client
	Context context.Context
}

func NewUserCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &userCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.Member),
		Client:        client,
		Context:       ctx,
	}
	return c
}

func (c *userCollector) CollectTotalEntities() int {
	groups, err := c.Client.Groups()
	if err != nil {
		log.Printf("failed to collect members: %v", err)
		return 0
	}

	totalGroupMembers := 0
	gw := group_waiter.New()
	for _, g := range groups {
		group := g
		gw.Do(func() {
			_, resp, err := c.Client.Client().Groups.ListGroupMembers(group.ID, &gitlab2.ListGroupMembersOptions{})
			if err != nil {
				log.Printf("Failed to get members for group %s", g.Name)
				return
			}
			totalGroupMembers += resp.TotalItems
		})

	}
	gw.Wait()

	return totalGroupMembers
}

func (c *userCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		allGroups, err := c.Client.Groups()
		if err != nil {
			return
		}

		gw := group_waiter.New()

		for _, group := range allGroups {
			g := group
			gw.Do(func() {
				c.collectGroupUsers(g)
			})
		}

		gw.Wait()
	})
}

func (c *userCollector) collectGroupUsers(group *gitlab2.Group) {
	gw := group_waiter.New()

	members, err := c.Client.GroupMembers(group)
	if err != nil {
		log.Printf("Failed to collect group members: %s - %s", group.Name, err)
		return
	}

	for _, m := range members {
		m := m
		gw.Do(func() {
			u, _, err := c.Client.Client().Users.GetUser(m.ID, gitlab2.GetUsersOptions{})
			if err != nil {
				log.Printf("failed to collect user %s - %s", m.Name, err)
				return
			}
			entity := gitlab_collected.Member{
				User: u,
			}
			c.CollectDataWithContext(&entity, entity.CanonicalLink(),
				newCollectionContext(nil, []permissions.OrganizationRole{permissions.OrgRoleOwner},
					c.Client.IsGroupPremium(group)))
			c.CollectionChangeByOne()
		})
	}

	gw.Wait()
	if err != nil {
		log.Printf("Failed to collect all users")
		return
	}
}
