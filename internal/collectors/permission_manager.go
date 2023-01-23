package collectors

import (
	"fmt"
	"strings"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/errlog"
	"github.com/iancoleman/orderedmap"
)

type MissingPermission struct {
	Permission string
	Entity     string
	Effect     string
	Namespace  namespace.Namespace
}

func NewMissingPermission(permission, entity, effect string, namespace namespace.Namespace) MissingPermission {
	return MissingPermission{
		Permission: permission,
		Entity:     entity,
		Effect:     effect,
		Namespace:  namespace,
	}
}

type effectSet = map[string]bool

func CollectMissingPermissions(missingPermissionChan chan MissingPermission) {
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
		entity := utils.UnsafeGet[*orderedmap.OrderedMap](permMap, permission)
		entity.Sort(lessByEntityName)
		permMap.Set(permission, entity)
	}

	// Build missing permissions string
	for _, permission := range permMap.Keys() {
		errlog.PermIssueF("missing permission: \"%s\" on:\n", permission)
		entity := utils.UnsafeGet[*orderedmap.OrderedMap](permMap, permission)
		for _, entityName := range entity.Keys() {
			effects := utils.UnsafeGet[effectSet](entity, entityName)
			filteredEffects := []string{}
			for effect := range effects {
				if effect == "" {
					continue
				}
				filteredEffects = append(filteredEffects, effect)
			}
			effectString := strings.Join(filteredEffects, ", ")
			errlog.PermIssueF("    - %s [%s]\n", entityName, effectString)
		}
	}
}
