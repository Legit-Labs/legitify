package gitlab_collected

import (
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	gitlab2 "github.com/xanzy/go-gitlab"
)

type Repository struct {
	*gitlab2.Project
	Members                  []*gitlab2.ProjectMember       `json:"members,omitempty"`
	ProtectedBranches        []*gitlab2.ProtectedBranch     `json:"protected_branches"`
	Webhooks                 []*gitlab2.ProjectHook         `json:"webhooks"`
	PushRules                *gitlab2.ProjectPushRules      `json:"push_rules"`
	ApprovalConfiguration    *gitlab2.ProjectApprovals      `json:"approval_configuration"`
	ApprovalRules            []*gitlab2.ProjectApprovalRule `json:"approval_rules"`
	MinimumRequiredApprovals int                            `json:"minimum_required_approvals"`
}

func (r Repository) ViolationEntityType() string {
	return namespace.Repository
}

func (r Repository) CanonicalLink() string {
	return r.Project.WebURL
}

func (r Repository) Name() string {
	return r.Project.Name
}

func (r Repository) ID() int64 {
	return int64(r.Project.ID)
}
