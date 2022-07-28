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

func (b baseCollector) collectData(org githubcollected.ExtendedOrg, entity githubcollected.CollectedEntity, canonicalLink string, viewerRoles []permissions.Role) {

	b.collectedChan <- CollectedData{
		Entity:        entity,
		Namespace:     b.Namespace(),
		CanonicalLink: canonicalLink,
		Context: &collectedDataContext{
			roles:        viewerRoles,
			organization: org,
		},
	}
}

func (b baseCollector) collectDataWithContext(entity githubcollected.CollectedEntity, canonicalLink string, ctx CollectedDataContext) {

	b.collectedChan <- CollectedData{
		Entity:        entity,
		Namespace:     b.Namespace(),
		CanonicalLink: canonicalLink,
		Context:       ctx,
	}
}

func (b baseCollector) totalCollectionChange(total int) {
	b.progressChan <- CollectionMetric{
		Namespace:             b.Namespace(),
		TotalCollectionChange: total,
	}
}

func (b baseCollector) collectionChange(change int) {
	b.progressChan <- CollectionMetric{
		Namespace:        b.Namespace(),
		CollectionChange: change,
	}
}

func (b baseCollector) collectionChangeByOne() {
	b.collectionChange(1)
}

func (b baseCollector) issueMissingPermissions(missingPermissions ...missingPermission) {
	for _, p := range missingPermissions {
		b.missingPermChan <- p
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
