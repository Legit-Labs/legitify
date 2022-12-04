package test

import (
	"testing"
	"time"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
)

type memberMockConfiguration struct {
	hasLastActive bool
	members       []githubcollected.OrganizationMember
}

func newMemberMock(config memberMockConfiguration) githubcollected.OrganizationMembers {
	return githubcollected.OrganizationMembers{
		Organization:  defaultOrg,
		HasLastActive: config.hasLastActive,
		Members:       config.members,
	}
}
func TestMember(t *testing.T) {
	tests := []struct {
		name             string
		policyName       string
		shouldBeViolated bool
		args             memberMockConfiguration
	}{
		{
			name:             "non admin should be stale",
			policyName:       "stale_member_found",
			shouldBeViolated: true,
			args: memberMockConfiguration{
				hasLastActive: true,
				members: []githubcollected.OrganizationMember{
					{
						LastActive: int(time.Now().AddDate(0, -9, 0).UnixNano()),
						IsAdmin:    false,
					},
				},
			},
		},
		{
			name:             "admin should be stale",
			policyName:       "stale_admin_found",
			shouldBeViolated: true,
			args: memberMockConfiguration{
				hasLastActive: true,
				members: []githubcollected.OrganizationMember{
					{
						LastActive: int(time.Now().AddDate(0, -9, 0).UnixNano()),
						IsAdmin:    true,
					},
				},
			},
		},
		{
			name:             "non admin should not be stale",
			policyName:       "stale_member_found",
			shouldBeViolated: false,
			args: memberMockConfiguration{
				hasLastActive: true,
				members: []githubcollected.OrganizationMember{
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    false,
					},
				},
			},
		},
		{
			name:             "admin should not be stale",
			policyName:       "stale_member_found",
			shouldBeViolated: false,
			args: memberMockConfiguration{
				hasLastActive: true,
				members: []githubcollected.OrganizationMember{
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		PolicyTestTemplateGitHub(t, test.name, newMemberMock(test.args),
			namespace.Member, test.policyName, test.shouldBeViolated)
	}
}
