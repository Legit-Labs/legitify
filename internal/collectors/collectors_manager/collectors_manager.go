package collectors_manager

import (
	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
)

type CollectorManager interface {
	Collect() <-chan collectors.CollectedData
}

type manager struct {
	collectors []collectors.Collector
}

func NewCollectorsManager(initiatedCollectors []collectors.Collector) CollectorManager {
	return &manager{
		collectors: initiatedCollectors,
	}
}

func (m *manager) Collect() <-chan collectors.CollectedData {
	collectedChan := make(chan collectors.CollectedData)

	// require all collection bars and the metadata bar
	progressbar.Report(progressbar.NewMinimalBars(len(m.collectors) + 1))

	// init the metadata bar
	const metadataBarName = "metadata"
	progressbar.Report(progressbar.NewRequiredBar(metadataBarName, len(m.collectors)))

	go func() {
		defer close(collectedChan)

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
				totalEntities := c.CollectTotalEntities()
				progressbar.Report(progressbar.NewRequiredBar(c.Namespace(), totalEntities))
				progressbar.Report(progressbar.NewUpdate(metadataBarName, 1))

				if totalEntities == 0 {
					return
				}

				pb := collectionChannels.Progress
				collected := collectionChannels.Collected
				perm := collectionChannels.MissingPermission

				for {
					select {
					case x, ok := <-pb:
						if !ok {
							pb = nil
						} else {
							progressbar.Report(x)
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
			})
		}
		gw.Wait()
		close(missingPermissionsChannel)
		permWait.Wait()
	}()

	return collectedChan
}
