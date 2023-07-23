package errlog

import (
	"fmt"
	"sync"

	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/iancoleman/orderedmap"
)

type PermLog struct {
	permissions *orderedmap.OrderedMap
	lock        sync.Mutex
}

func NewPermLog() *PermLog {
	return &PermLog{
		permissions: orderedmap.New(),
	}
}

func (p *PermLog) MarshalJSON() ([]byte, error) {
	p.sortEntities()
	return p.permissions.MarshalJSON()
}

func (p *PermLog) Add(issue PermIssue) {
	p.lock.Lock()
	defer p.lock.Unlock()

	permissions, ok := p.permissions.Get(issue.Permission)
	if !ok {
		permissions = orderedmap.New()
	}
	typedEntities := permissions.(*orderedmap.OrderedMap)

	effects, ok := typedEntities.Get(issue.Entity)
	if !ok {
		effects = make(effectSet)
	}
	typedEffects := effects.(effectSet)
	typedEffects[issue.Effect] = true

	typedEntities.Set(issue.Entity, typedEffects)
	p.permissions.Set(issue.Permission, typedEntities)
}

func (p *PermLog) Empty() bool {
	return len(p.permissions.Keys()) == 0
}

func (p *PermLog) sortEntities() {
	for _, permission := range p.permissions.Keys() {
		entity := map_utils.UnsafeGet[*orderedmap.OrderedMap](p.permissions, permission)
		entity.Sort(lessByEntityName)
		p.permissions.Set(permission, entity)
	}
}

func lessByEntityName(a *orderedmap.Pair, b *orderedmap.Pair) bool {
	aName := a.Key()
	bName := b.Key()
	return aName < bName
}

type PermIssue struct {
	Permission string
	Entity     string
	Effect     string
}

func NewPermIssue(permission string, entity string, namespace string, effect string) PermIssue {
	return PermIssue{
		Permission: permission,
		Entity:     fmt.Sprintf("%s:%s", namespace, entity),
		Effect:     effect,
	}
}
