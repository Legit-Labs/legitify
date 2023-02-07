package collectors

import (
	"fmt"

	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/collected"
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

func FullRepoName(org string, repo string) string {
	return fmt.Sprintf("%s/%s", org, repo)
}

type collectedDataContext struct {
	organization githubcollected.ExtendedOrg
	roles        []permissions.Role
}

func (c *collectedDataContext) Premium() bool {
	return c.organization.IsEnterprise()
}

func (c *collectedDataContext) Roles() []permissions.Role {
	return c.roles
}

type BaseCollector struct {
	namespace       string
	collectedChan   chan CollectedData
	progressChan    chan progressbar.ChannelType
	missingPermChan chan MissingPermission
}

func NewBaseCollector(namespace string) BaseCollector {
	return BaseCollector{namespace: namespace}
}

func (b *BaseCollector) Namespace() string {
	return b.namespace
}

func (b *BaseCollector) CollectData(org githubcollected.ExtendedOrg, entity collected.Entity, canonicalLink string, viewerRoles []permissions.Role) {
	b.collectedChan <- CollectedData{
		Entity:        entity,
		Namespace:     b.namespace,
		CanonicalLink: canonicalLink,
		Context: &collectedDataContext{
			roles:        viewerRoles,
			organization: org,
		},
	}
}

func (b *BaseCollector) CollectDataWithContext(entity collected.Entity, canonicalLink string, ctx CollectedDataContext) {

	b.collectedChan <- CollectedData{
		Entity:        entity,
		Namespace:     b.namespace,
		CanonicalLink: canonicalLink,
		Context:       ctx,
	}
}

func (b *BaseCollector) CollectionChange(change int) {
	b.progressChan <- progressbar.NewUpdate(b.namespace, change)
}

func (b *BaseCollector) CollectionChangeByOne() {
	b.CollectionChange(1)
}

func (b *BaseCollector) IssueMissingPermissions(missingPermissions ...MissingPermission) {
	for _, p := range missingPermissions {
		b.missingPermChan <- p
	}
}

func (b *BaseCollector) makeChannels() {
	b.collectedChan = make(chan CollectedData)
	b.progressChan = make(chan progressbar.ChannelType)
	b.missingPermChan = make(chan MissingPermission)
}

func (b *BaseCollector) closeChannels() {
	b.progressChan <- progressbar.NewBarClose(b.namespace)
	close(b.collectedChan)
	close(b.progressChan)
	close(b.missingPermChan)
}

func (b *BaseCollector) getChannels() SubCollectorChannels {
	return SubCollectorChannels{
		Collected:         b.collectedChan,
		Progress:          b.progressChan,
		MissingPermission: b.missingPermChan,
	}
}

func (b *BaseCollector) WrappedCollection(collection func()) SubCollectorChannels {
	b.makeChannels()
	go func() {
		defer b.closeChannels()
		collection()
	}()
	return b.getChannels()
}
