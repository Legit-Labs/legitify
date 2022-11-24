package collectors_manager

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/collectors"
	github2 "github.com/Legit-Labs/legitify/internal/collectors/github"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"

	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
)

type CollectorChannels struct {
	Collected <-chan collectors.CollectedData
	Progress  <-chan collectors.CollectionMetric
}

type CollectorManager interface {
	Collect() CollectorChannels
	CollectMetadata() map[namespace.Namespace]collectors.Metadata
}

type manager struct {
	collectors []collectors.Collector
	ctx        context.Context
	client     github.Client
}

type newCollectorFunc func(ctx context.Context, client github.Client) collectors.Collector

var collectorsMapping = map[namespace.Namespace]newCollectorFunc{
	namespace.Repository:   github2.NewRepositoryCollector,
	namespace.Organization: github2.NewOrganizationCollector,
	namespace.Member:       github2.NewMemberCollector,
	namespace.Actions:      github2.NewActionCollector,
	namespace.RunnerGroup:  github2.NewRunnersCollector,
}

func NewCollectorsManager(ctx context.Context, ns []namespace.Namespace, client github.Client) CollectorManager {
	var initiatedCollectors []collectors.Collector
	for _, n := range ns {
		initiatedCollectors = append(initiatedCollectors, collectorsMapping[n](ctx, client))
	}

	return &manager{
		ctx:        ctx,
		client:     client,
		collectors: initiatedCollectors,
	}
}

func (m *manager) CollectMetadata() map[namespace.Namespace]collectors.Metadata {
	type metaDataPair struct {
		Namespace namespace.Namespace
		Metadata  collectors.Metadata
	}

	gw := group_waiter.New()
	ch := make(chan metaDataPair, len(m.collectors))
	for _, c := range m.collectors {
		c := c
		gw.Do(func() {
			ch <- metaDataPair{Namespace: c.Namespace(), Metadata: c.CollectMetadata()}
		})
	}
	gw.Wait()
	close(ch)

	res := make(map[namespace.Namespace]collectors.Metadata)
	for m := range ch {
		res[m.Namespace] = m.Metadata
	}

	return res
}

func (m *manager) Collect() CollectorChannels {
	collectedChan := make(chan collectors.CollectedData)
	progressChan := make(chan collectors.CollectionMetric)

	go func() {
		defer close(collectedChan)
		defer close(progressChan)

		missingPermissionsChannel := make(chan collectors.MissingPermission)
		permWait := group_waiter.New()
		permWait.Do(func() {
			collectors.CollectMissingPermissions(missingPermissionsChannel)
		})

		gw := group_waiter.New()
		for _, c := range m.collectors {
			collectionChannels := c.Collect()

			gw.Do(func() {
				pb := collectionChannels.Progress
				collected := collectionChannels.Collected
				perm := collectionChannels.MissingPermission

				for {
					select {
					case x, ok := <-pb:
						if !ok {
							pb = nil
						} else {
							progressChan <- x
						}
					case x, ok := <-collected:
						if !ok {
							collected = nil
						} else {
							collectedChan <- x
						}
					case x, ok := <-perm:
						if !ok {
							perm = nil
						} else {
							missingPermissionsChannel <- x
						}
					}

					if pb == nil && collected == nil && perm == nil {
						break
					}
				}

				progressChan <- collectors.CollectionMetric{
					Finished:  true,
					Namespace: c.Namespace(),
				}
			})
		}
		gw.Wait()
		close(missingPermissionsChannel)
		permWait.Wait()
	}()

	return CollectorChannels{
		Collected: collectedChan,
		Progress:  progressChan,
	}
}
