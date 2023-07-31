package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"log"

	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type serverCollector struct {
	collectors.BaseCollector
	Client   *gitlab.Client
	Context  context.Context
	isServer bool
}

func NewServerCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &serverCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.Enterprise),
		Client:        client,
		Context:       ctx,
		isServer:      client.IsServer(),
	}
	return c
}

func (c *serverCollector) CollectTotalEntities() int {
	if c.isServer {
		return 1
	}

	return 0
}

func (c *serverCollector) Collect() collectors.SubCollectorChannels {
	return c.WrappedCollection(func() {
		if !c.isServer {
			return
		}

		settings, _, err := c.Client.Client().Settings.GetSettings()

		if err != nil {
			log.Printf("failed to collect server settings %s", err)
			return
		}

		entity := gitlab_collected.NewServer(c.Client.ServerUrl(), settings)
		c.CollectDataWithContext(entity, entity.CanonicalLink(),
			newServerCollectionContext())
		c.CollectionChangeByOne()
	})
}
