package githubcollected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type Enterprise struct {
	MembersCanChangeRepositoryVisibilitySetting string `json:"members_can_change_repository_visibility"`
	EnterpriseName                              string `json:"name"`
	Url                                         string `json:"url"`
	Id                                          int64  `json:"id"`
	UserRole                                    string
}

func NewEnterprise(MembersCanChangeRepositoryVisibilitySetting string, Name string, Url string, Id int64, isAdmin bool) Enterprise {
	UserRole := permissions.EnterpriseNonAdminRole
	if isAdmin {
		UserRole = permissions.EnterpriseAdminRole
	}
	return Enterprise{MembersCanChangeRepositoryVisibilitySetting: MembersCanChangeRepositoryVisibilitySetting,
		EnterpriseName: Name,
		Url:            Url,
		Id:             Id,
		UserRole:       UserRole,
	}
}

func (o Enterprise) ViolationEntityType() string {
	return namespace.Enterprise
}

func (o Enterprise) CanonicalLink() string {
	return o.Url
}

func (o Enterprise) Name() string {
	return o.EnterpriseName
}

func (o Enterprise) ID() int64 {
	return o.Id
}
