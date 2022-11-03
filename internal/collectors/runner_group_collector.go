package collectors

import (
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/google/go-github/v44/github"
	"log"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type runnersCollector struct {
	baseCollector
	client  ghclient.Client
	context context.Context
	cache   map[string][]*github.RunnerGroup
}

func newRunnersCollector(ctx context.Context, client ghclient.Client) collector {
	c := &runnersCollector{
		client:  client,
		context: ctx,
		cache:   make(map[string][]*github.RunnerGroup),
	}
	initBaseCollector(&c.baseCollector, c)
	return c
}

func (c *runnersCollector) Namespace() namespace.Namespace {
	return namespace.RunnerGroup
}

func (c *runnersCollector) CollectMetadata() Metadata {
	gw := group_waiter.New()
	orgs, err := c.client.CollectOrganizations()
	if err != nil {
		log.Printf("failed to collection organizations %s", err)
		return Metadata{}
	}

	totalCount := 0
	for _, org := range orgs {
		gw.Do(func() {
			org := org
			result := make([]*github.RunnerGroup, 0)
			err := ghclient.PaginateResults(func(opts *github.ListOptions) (*github.Response, error) {
				runners, resp, err := c.client.Client().Actions.ListOrganizationRunnerGroups(c.context, *org.Login, opts)

				if err != nil {
					log.Printf("err collection runner groups for %s - %v", *org.Login, err)
					return nil, err
				}

				result = append(result, runners.RunnerGroups...)
				return resp, nil
			})

			if err != nil {
				log.Printf("Failed to collected runner group metadata for organization %s", *org.Login)
			} else {
				c.cache[*org.Login] = result
				totalCount = totalCount + len(result)
			}
		})
	}

	gw.Wait()
	return Metadata{
		totalCount,
	}
}

func (c *runnersCollector) Collect() subCollectorChannels {
	return c.wrappedCollection(func() {
		orgs, err := c.client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		for _, org := range orgs {
			cached := c.cache[*org.Login]

			for _, rg := range cached {
				c.collectionChangeByOne()

				c.collectData(org,
					ghcollected.RunnerGroup{
						Organization: org,
						RunnerGroup:  rg,
					},
					*org.HTMLURL,
					[]permissions.Role{org.Role})
			}
		}
	})
}
