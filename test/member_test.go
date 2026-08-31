package test

import (
	"testing"
	"time"

	"github.com/Legit-Labs/legitify/internal/common/scm_type"

	githubcollected "github.com/Legit-Labs/legitify/internal/collected/github"
	gitlabcollected "github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	gitlab2 "github.com/xanzy/go-gitlab"
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
						LastActive: int(time.Now().AddDate(-1, -2, 0).UnixNano()),
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
		{
			name:             "should find too many admins",
			policyName:       "organization_has_too_many_admins",
			shouldBeViolated: true,
			args: memberMockConfiguration{
				hasLastActive: true,
				members: []githubcollected.OrganizationMember{
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    true,
					},
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    true,
					},
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    true,
					},
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    true,
					},
					{
						LastActive: int(time.Now().AddDate(0, -1, 0).UnixNano()),
						IsAdmin:    true,
					},
				},
			},
		},
	}

	for _, test := range tests {
		PolicyTestTemplate(t, test.name, newMemberMock(test.args),
			namespace.Member, test.policyName, test.shouldBeViolated, scm_type.GitHub)
	}
}

func gitlabMemberTestTemplate(t *testing.T, name string, mockData interface{}, testedPolicyName string, expectFailure bool) {
	PolicyTestTemplate(t, name, mockData, namespace.Member, testedPolicyName, expectFailure, scm_type.GitLab)
}

func TestGitlabCollaboratorTwoFactor(t *testing.T) {
	name := "Collaborator Should Have Two-Factor Authentication Enabled"
	testedPolicyName := "two_factor_authentication_is_disabled_for_a_collaborator"

	makeMockData := func(twoFactorEnabled bool) gitlabcollected.Member {
		return gitlabcollected.Member{User: &gitlab2.User{TwoFactorEnabled: twoFactorEnabled}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		gitlabMemberTestTemplate(t, name, makeMockData(options[expectFailure]), testedPolicyName, expectFailure)
	}
}

func TestGitlabExternalCollaboratorTwoFactor(t *testing.T) {
	name := "External Collaborator Should Have Two-Factor Authentication Enabled"
	testedPolicyName := "two_factor_authentication_is_disabled_for_an_external_collaborator"

	makeMockData := func(twoFactorEnabled bool) gitlabcollected.Member {
		return gitlabcollected.Member{User: &gitlab2.User{External: true, TwoFactorEnabled: twoFactorEnabled}}
	}

	options := map[bool]bool{
		false: true,
		true:  false,
	}

	for _, expectFailure := range bools {
		gitlabMemberTestTemplate(t, name, makeMockData(options[expectFailure]), testedPolicyName, expectFailure)
	}
}
