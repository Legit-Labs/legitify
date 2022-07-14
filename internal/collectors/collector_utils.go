package collectors

import (
	"fmt"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

func fullRepoName(org string, repo string) string {
	return fmt.Sprintf("%s/%s", org, repo)
}

func wrapCollectorChans(collected chan CollectedData, progress chan CollectionMetric) CollectorChannels {
	return CollectorChannels{
		Collected: collected,
		Progress:  progress,
	}
}

type collectedDataContext struct {
	organization githubcollected.ExtendedOrg
	roles        []permissions.Role
}

func (c *collectedDataContext) IsEnterprise() bool {
	return c.organization.IsEnterprise()
}

func (c *collectedDataContext) Roles() []permissions.Role {
	return c.roles
}

type baseCollector struct {
	collector
	collectedChan   chan CollectedData
	progressChan    chan CollectionMetric
	missingPermChan chan missingPermission
}

func initBaseCollector(b *baseCollector, c collector) {
	b.collector = c
}

func (c baseCollector) collectData(org githubcollected.ExtendedOrg, entity githubcollected.CollectedEntity, canonicalLink string, viewerRoles []permissions.Role) {

	c.collectedChan <- CollectedData{
		Entity:        entity,
		Namespace:     c.Namespace(),
		CanonicalLink: canonicalLink,
		Context: &collectedDataContext{
			roles:        viewerRoles,
			organization: org,
		},
	}
}

func (c baseCollector) totalCollectionChange(total int) {
	c.progressChan <- CollectionMetric{
		Namespace:             c.Namespace(),
		TotalCollectionChange: total,
	}
}

func (c baseCollector) collectionChange(change int) {
	c.progressChan <- CollectionMetric{
		Namespace:        c.Namespace(),
		CollectionChange: change,
	}
}

func (c baseCollector) collectionChangeByOne() {
	c.collectionChange(1)
}

func (c baseCollector) issueMissingPermissions(missingPermissions ...missingPermission) {
	for _, p := range missingPermissions {
		c.missingPermChan <- p
	}
}

func (b *baseCollector) makeChannels() {
	b.collectedChan = make(chan CollectedData)
	b.progressChan = make(chan CollectionMetric)
	b.missingPermChan = make(chan missingPermission)
}

func (b *baseCollector) closeChannels() {
	close(b.collectedChan)
	close(b.progressChan)
	close(b.missingPermChan)
}

func (b *baseCollector) getChannels() subCollectorChannels {
	return subCollectorChannels{
		CollectorChannels: wrapCollectorChans(b.collectedChan, b.progressChan),
		MissingPermission: b.missingPermChan,
	}
}

func (b *baseCollector) wrappedCollection(collection func()) subCollectorChannels {
	b.makeChannels()
	go func() {
		defer b.closeChannels()
		collection()
	}()
	return b.getChannels()
}
