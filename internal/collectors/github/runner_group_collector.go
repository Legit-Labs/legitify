package github

import (
	"log"
	"sync"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/clients/github/pagination"
	ghcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/google/go-github/v49/github"
	"golang.org/x/net/context"
)

type runnersCollector struct {
	collectors.BaseCollector
	client  *ghclient.Client
	context context.Context
}

func NewRunnersCollector(ctx context.Context, client *ghclient.Client) collectors.Collector {
	c := &runnersCollector{
		client:  client,
		context: ctx,
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (c *runnersCollector) Namespace() namespace.Namespace {
	return namespace.RunnerGroup
}

func (c *runnersCollector) collectForOrg(orgName string) ([]*github.RunnerGroup, error) {
	mapper := func(rg *github.RunnerGroups) []*github.RunnerGroup {
		if rg == nil {
			return []*github.RunnerGroup{}
		}
		return rg.RunnerGroups
	}
	result := pagination.NewMapper(c.client.Client().Actions.ListOrganizationRunnerGroups, nil, mapper).Sync(c.context, orgName)
	if result.Err != nil {
		// TODO check for permission error and leave logging for the caller
		log.Printf("Error collecting runner groups for %s - %v", orgName, result.Err)
	}
	return result.Collected, result.Err
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
		org := org
		gw.Do(func() {
			result, err := c.collectForOrg(org.Name())
			if err != nil {
				// TODO check for permission error and consider failing execution for other errors
				log.Printf("Error collecting runner groups for %s - %v", org.Name(), err)
			} else {
				mutex.Lock()
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
			groups, err := c.collectForOrg(org.Name())
			if err != nil {
				return
			}

			for _, rg := range groups {
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
