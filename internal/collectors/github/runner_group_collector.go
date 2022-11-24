package github

import (
	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/google/go-github/v44/github"
	"golang.org/x/net/context"
	"log"
	"sync"
)

type runnersCollector struct {
	collectors.BaseCollector
	client  ghclient.Client
	context context.Context
	cache   map[string][]*github.RunnerGroup
}

func NewRunnersCollector(ctx context.Context, client ghclient.Client) collectors.Collector {
	c := &runnersCollector{
		client:  client,
		context: ctx,
		cache:   make(map[string][]*github.RunnerGroup),
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (c *runnersCollector) Namespace() namespace.Namespace {
	return namespace.RunnerGroup
}

func (c *runnersCollector) CollectMetadata() collectors.Metadata {
	gw := group_waiter.New()
	orgs, err := c.client.CollectOrganizations()
	if err != nil {
		log.Printf("failed to collection organizations %s", err)
		return collectors.Metadata{}
	}

	totalCount := 0
	var mutex = &sync.RWMutex{}
	for _, org := range orgs {
		gw.Do(func() {
			org := org
			result := make([]*github.RunnerGroup, 0)
			err := ghclient.PaginateResults(func(opts *github.ListOptions) (*github.Response, error) {
				runners, resp, err := c.client.Client().Actions.ListOrganizationRunnerGroups(c.context, org.Name(), opts)

				if err != nil {
					return nil, err
				}

				result = append(result, runners.RunnerGroups...)
				return resp, nil
			})

			if err != nil {
				log.Printf("Error collecting runner groups for %s - %v", org.Name(), err)
			} else {
				mutex.Lock()
				c.cache[org.Name()] = result
				totalCount = totalCount + len(result)
				mutex.Unlock()
			}
		})
	}

	gw.Wait()
	return collectors.Metadata{
		TotalEntities: totalCount,
	}
}

func (c *runnersCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		orgs, err := c.client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		for _, org := range orgs {
			cached := c.cache[org.Name()]

			for _, rg := range cached {
				c.CollectionChangeByOne()

				c.CollectData(org,
					ghcollected.RunnerGroup{
						Organization: org,
						RunnerGroup:  rg,
					},
					org.CanonicalLink(),
					[]permissions.Role{org.Role})
			}
		}
	})
}
