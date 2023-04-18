package githubcollected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
)

type Enterprise struct {
	MembersCanChangeRepositoryVisibilitySetting string
	EnterpriseName                              string
	Url                                         string
	Id                                          int64
	UserRole                                    permissions.OrganizationRole
}

func NewEnterprise(MembersCanChangeRepositoryVisibilitySetting string, Name string, Url string, Id int64) Enterprise {
	return Enterprise{MembersCanChangeRepositoryVisibilitySetting: MembersCanChangeRepositoryVisibilitySetting,
		EnterpriseName: Name,
		Url:            Url,
		Id:             Id,
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
