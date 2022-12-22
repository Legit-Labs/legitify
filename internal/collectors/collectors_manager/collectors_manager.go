package collectors_manager

import (
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"

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
}

func NewCollectorsManager(initiatedCollectors []collectors.Collector) CollectorManager {
	return &manager{
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
			c := c
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
