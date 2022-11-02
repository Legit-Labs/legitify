package collectors

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"

	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
)

type CollectorChannels struct {
	Collected <-chan CollectedData
	Progress  <-chan CollectionMetric
}

type CollectorManager interface {
	Collect() CollectorChannels
	CollectMetadata() map[namespace.Namespace]Metadata
}

type manager struct {
	collectorCreators []newCollectorFunc
	collectors        []collector
	ctx               context.Context
	client            github.Client
}

type newCollectorFunc func(ctx context.Context, client github.Client) collector

var collectorsMapping = map[namespace.Namespace]newCollectorFunc{
	namespace.Repository:   newRepositoryCollector,
	namespace.Organization: newOrganizationCollector,
	namespace.Member:       newMemberCollector,
	namespace.Actions:      newActionCollector,
	namespace.Runners:      newRunnersCollector,
}

func NewCollectorsManager(ctx context.Context, ns []namespace.Namespace, client github.Client) CollectorManager {
	var selected []newCollectorFunc
	var collectors []collector
	for _, n := range ns {
		selected = append(selected, collectorsMapping[n])
		collectors = append(collectors, collectorsMapping[n](ctx, client))
	}

	return &manager{
		collectorCreators: selected,
		ctx:               ctx,
		client:            client,
		collectors:        collectors,
	}
}

func (m *manager) CollectMetadata() map[namespace.Namespace]Metadata {
	type metaDataPair struct {
		Namespace namespace.Namespace
		Metadata  Metadata
	}

	gw := group_waiter.New()
	ch := make(chan metaDataPair, len(m.collectorCreators))
	for _, c := range m.collectors {
		gw.Do(func() {
			ch <- metaDataPair{Namespace: c.Namespace(), Metadata: c.CollectMetadata()}
		})
	}
	gw.Wait()
	close(ch)

	res := make(map[namespace.Namespace]Metadata)
	for m := range ch {
		res[m.Namespace] = m.Metadata
	}

	return res
}

func (m *manager) createCollector(creator newCollectorFunc) collector {
	return creator(m.ctx, m.client)
}

func (m *manager) Collect() CollectorChannels {
	collectedChan := make(chan CollectedData)
	progressChan := make(chan CollectionMetric)

	go func() {
		defer close(collectedChan)
		defer close(progressChan)

		missingPermissionsChannel := make(chan missingPermission)
		permWait := group_waiter.New()
		permWait.Do(func() {
			collectMissingPermissions(missingPermissionsChannel)
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

				progressChan <- CollectionMetric{
					Finished:  true,
					Namespace: c.Namespace(),
				}
			})
		}
		gw.Wait()
		close(missingPermissionsChannel)
		permWait.Wait()
	}()

	return wrapCollectorChans(collectedChan, progressChan)
}
