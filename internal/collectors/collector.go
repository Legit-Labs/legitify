package collectors

import (
	"github.com/Legit-Labs/legitify/internal/collected"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type Metadata struct {
	TotalEntities int
}

type CollectionMetric struct {
	CollectionChange int
	Finished         bool
	Namespace        string
}

type CollectedDataContext interface {
	IsEnterprise() bool
	Roles() []permissions.Role
}

type CollectedData struct {
	Context       CollectedDataContext
	Entity        collected.Entity
	Namespace     namespace.Namespace
	CanonicalLink string
}

type SubCollectorChannels struct {
	Collected         <-chan CollectedData
	Progress          <-chan CollectionMetric
	MissingPermission <-chan MissingPermission
}

type Collector interface {
	Collect() SubCollectorChannels
	Namespace() namespace.Namespace
	CollectMetadata() Metadata
}
