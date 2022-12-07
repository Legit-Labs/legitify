package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	gitlab2 "github.com/xanzy/go-gitlab"
	"log"

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
		Client:  client,
		Context: ctx,
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (c *userCollector) Namespace() namespace.Namespace {
	return namespace.Member
}

func (c *userCollector) CollectMetadata() collectors.Metadata {
	opts := &gitlab2.ListUsersOptions{}

	_, resp, err := c.Client.Client().Users.ListUsers(opts)

	res := collectors.Metadata{}
	if err != nil {
		log.Printf("failed to collect users %s", err)
	} else {
		res.TotalEntities = resp.TotalItems
	}

	return res
}

func (c *userCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		userOptions := &gitlab2.ListUsersOptions{}

		gw := group_waiter.New()

		var result []*gitlab2.User
		err := gitlab.PaginateResults(func(opts *gitlab2.ListOptions) (*gitlab2.Response, error) {
			users, resp, err := c.Client.Client().Users.ListUsers(userOptions)

			if err != nil {
				return nil, err
			}

			cp := users
			gw.Do(func() {
				for _, user := range cp {
					entity := gitlab_collected.Member{
						User: user,
					}
					c.CollectDataWithContext(&entity, entity.CanonicalLink(),
						newCollectionContext(nil, []permissions.OrganizationRole{permissions.OrgRoleOwner}))
					c.CollectionChangeByOne()
				}
			})
			result = append(result, users...)
			return resp, nil

		}, &userOptions.ListOptions)
		gw.Wait()

		if err != nil {
			log.Printf("Failed to collect users")
			return
		}
	})
}
