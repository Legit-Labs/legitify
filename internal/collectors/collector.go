package collectors

import (
	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type Metadata struct {
	TotalEntities int
}

type CollectionMetric struct {
	TotalCollectionChange int
	CollectionChange      int
	Finished              bool
	Namespace             string
}

type CollectedDataContext interface {
	IsEnterprise() bool
	Roles() []permissions.Role
}

type CollectedData struct {
	Context       CollectedDataContext
	Entity        githubcollected.CollectedEntity
	Namespace     namespace.Namespace
	CanonicalLink string
}

type subCollectorChannels struct {
	CollectorChannels
	MissingPermission <-chan missingPermission
}

type collector interface {
	Collect() subCollectorChannels
	Namespace() namespace.Namespace
	CollectMetadata() Metadata
}
