package collectors

import (
	"fmt"
	"log"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/iancoleman/orderedmap"
)

type missingPermission struct {
	Permission string
	Entity     string
	Effect     string
	Namespace  namespace.Namespace
}

func newMissingPermission(permission, entity, effect string, namespace namespace.Namespace) missingPermission {
	return missingPermission{
		Permission: permission,
		Entity:     entity,
		Effect:     effect,
		Namespace:  namespace,
	}
}

type effectSet = map[string]bool

func collectMissingPermissions(missingPermissionChan chan missingPermission) {
	permMap := orderedmap.New()
	for permission := range missingPermissionChan {
		entities, ok := permMap.Get(permission.Permission)
		if !ok {
			entities = orderedmap.New()
		}
		typedEntities := entities.(*orderedmap.OrderedMap)

		entityName := fmt.Sprintf("%s:%s", permission.Namespace, permission.Entity)
		effects, ok := typedEntities.Get(entityName)
		if !ok {
			effects = make(effectSet)
		}
		typedEffects := effects.(effectSet)
		typedEffects[permission.Effect] = true

		typedEntities.Set(entityName, typedEffects)
		permMap.Set(permission.Permission, typedEntities)
	}

	// Sort entities by entity name (within each permission)
	lessByEntityName := func(a *orderedmap.Pair, b *orderedmap.Pair) bool {
		aName := a.Key()
		bName := b.Key()
		return aName < bName
	}
	for _, permission := range permMap.Keys() {
		entity := utils.UnsafeGet(permMap, permission).(*orderedmap.OrderedMap)
		entity.Sort(lessByEntityName)
		permMap.Set(permission, entity)
	}

	// Build missing permissions string
	var sb strings.Builder
	for _, permission := range permMap.Keys() {
		sb.WriteString(fmt.Sprintf("missing permission: \"%s\" on:\n", permission))
		entity := utils.UnsafeGet(permMap, permission).(*orderedmap.OrderedMap)
		for _, entityName := range entity.Keys() {
			effects := utils.UnsafeGet(entity, entityName).(effectSet)
			filteredEffects := []string{}
			for effect := range effects {
				if effect == "" {
					continue
				}
				filteredEffects = append(filteredEffects, effect)
			}
			effectString := strings.Join(filteredEffects, ", ")
			sb.WriteString(fmt.Sprintf("    - %s [%s]\n", entityName, effectString))
		}
	}
	log.Print(sb.String())
}
