package github

import (
	"log"
	"sync"
	"sync/atomic"

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
	client      *ghclient.Client
	context     context.Context
	orgLock     sync.Mutex
	groupsByOrg map[string][]*github.RunnerGroup
}

func NewRunnersCollector(ctx context.Context, client *ghclient.Client) collectors.Collector {
	c := &runnersCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.RunnerGroup),
		client:        client,
		context:       ctx,
		groupsByOrg:   make(map[string][]*github.RunnerGroup),
	}
	return c
}

func (c *runnersCollector) collectForOrg(orgName string) []*github.RunnerGroup {
	c.orgLock.Lock()
	defer c.orgLock.Unlock()

	if groups, ok := c.groupsByOrg[orgName]; ok {
		return groups
	}

	c.groupsByOrg[orgName] = nil
	mapper := func(rg *github.RunnerGroups) []*github.RunnerGroup {
		if rg == nil {
			return []*github.RunnerGroup{}
		}
		return rg.RunnerGroups
	}
	result, err := pagination.NewMapper(c.client.Client().Actions.ListOrganizationRunnerGroups, nil, mapper).Sync(c.context, orgName)
	if err != nil {
		perm := collectors.NewMissingPermission(permissions.OrgAdmin, orgName,
			"Cannot read organization runner groups", namespace.Organization)
		c.IssueMissingPermissions(perm)
	}
	c.groupsByOrg[orgName] = result.Collected

	return c.groupsByOrg[orgName]
}

func (c *runnersCollector) CollectTotalEntities() int {
	gw := group_waiter.New()
	orgs, err := c.client.CollectOrganizations()
	if err != nil {
		log.Printf("failed to collect organizations %s", err)
		return 0
	}

	var totalCount atomic.Int64
	for _, org := range orgs {
		org := org
		gw.Do(func() {
			groups := c.collectForOrg(org.Name())
			totalCount.Add(int64(len(groups)))
		})
	}
	gw.Wait()

	return int(totalCount.Load())
}

func (c *runnersCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		orgs, err := c.client.CollectOrganizations()

		if err != nil {
			log.Printf("failed to collect organizations %s", err)
			return
		}

		for _, org := range orgs {
			groups := c.collectForOrg(org.Name())
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
