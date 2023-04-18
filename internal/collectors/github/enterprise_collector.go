package github

import (
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"log"

	"github.com/Legit-Labs/legitify/internal/collectors"

	ghclient "github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type enterpriseCollector struct {
	collectors.BaseCollector
	Client  *ghclient.Client
	Context context.Context
}

func NewEnterpriseCollector(ctx context.Context, client *ghclient.Client) collectors.Collector {
	c := &enterpriseCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.Enterprise),
		Client:        client,
		Context:       ctx,
	}
	return c
}

func (c *enterpriseCollector) CollectTotalEntities() int {
	collectedEnterprises, err := c.Client.CollectEnterprises()
	if err != nil {
		log.Printf("failed to collect enterprise %s", err)
		return 0
	}

	return len(collectedEnterprises)
}

func (c *enterpriseCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		enterprises, err := c.Client.CollectEnterprises()

		if err != nil {
			log.Printf("failed to collect enterprise %s", err)
			return
		}

		gw := group_waiter.New()
		for _, enterprise := range enterprises {
			enterprise = enterprise
			gw.Do(func() {
				c.CollectionChangeByOne()
				c.CollectDataWithContext(enterprise, enterprise.Url, newEnterpriseContext([]permissions.Role{permissions.EnterpriseAdmin}))
			})
		}
		gw.Wait()
	})
}
