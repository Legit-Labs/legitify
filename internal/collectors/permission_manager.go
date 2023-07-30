package collectors

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/errlog"
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

func CollectMissingPermissions(missingPermissionChan chan MissingPermission) {
	for permission := range missingPermissionChan {
		errlog.AddPermIssue(errlog.NewPermIssue(permission.Permission, permission.Entity, permission.Namespace, permission.Effect))
	}
}
