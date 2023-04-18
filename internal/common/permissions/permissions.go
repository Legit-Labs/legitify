package permissions

type OrganizationRole = string

const (
	OrgRoleNone OrganizationRole = "NONE"

	OrgRoleOwner  = "OWNER"
	OrgRoleMember = "MEMBER"
)

type RepositoryRole = string

const (
	RepoRoleNone RepositoryRole = "NONE"

	RepoRoleAdmin      = "ADMIN"
	RepoRoleMaintainer = "MAINTAIN"
	RepoRoleWrite      = "WRITE"
	RepoRoleTriage     = "TRIAGE"
	RepoRoleRead       = "READ"
)

type Role = string

func IsOrgRole(role Role) bool {
	return role == OrgRoleOwner || role == OrgRoleMember
}

func IsRepositoryRole(role Role) bool {
	return role == RepoRoleAdmin ||
		role == RepoRoleMaintainer ||
		role == RepoRoleWrite ||
		role == RepoRoleTriage ||
		role == RepoRoleRead
}

func IsEnterpriseRole(role Role) bool {
	return role == EnterpriseRead ||
		role == EnterpriseManageBilling ||
		role == EnterpriseManageRunners ||
		role == EnterpriseAdmin
}

func HasScope(requiredScope string, availableScopes TokenScopes, roles []Role) bool {
	hasPermission := false

	for _, role := range roles {
		switch {
		case IsOrgRole(role):
			hasPermission = HasOrgScope(requiredScope, availableScopes, role)
		case IsRepositoryRole(role):
			hasPermission = HasRepoScope(requiredScope, availableScopes, role)
		case IsEnterpriseRole(role):
			hasPermission = HasEnterpriseScope(requiredScope, availableScopes, role)
		}
		if hasPermission {
			return true
		}
	}

	return false
}

type TokenScope = string
type TokenScopes = map[TokenScope]bool

const (
	ScopeNone TokenScope = "None"

	RepoAdmin          = "repo"
	RepoRepoStatus     = "repo:status"
	RepoRepoDeployment = "repo_deployment"
	RepoPublicRepo     = "public_repo"
	RepoRepoInvite     = "repo:invite"
	RepoSecurityEvents = "security_events"
	RepoDelete         = "delete_repo"

	Workflow = "workflow"

	PackagesWrite  = "write:packages"
	PackagesRead   = "read:packages"
	PackagesDelete = "delete:packages"

	OrgAdmin = "admin:org"
	OrgWrite = "write:org"
	OrgRead  = "read:org"

	PublicKeyAdmin = "admin:public_key"
	PublicKeyWrite = "write:public_key"
	PublicKeyRead  = "read:public_key"

	OrgHookAdmin  = "admin:org_hook"
	RepoHookAdmin = "admin:repo_hook"
	RepoHookWrite = "write:repo_hook"
	RepoHookRead  = "read:repo_hook"

	Gist = "gist"

	Notifications = "notifications"

	UserAll    = "user"
	UserRead   = "read:user"
	UserEmail  = "read:email"
	UserFollow = "user:follow"

	DiscussionWrite = "write:discussion"
	DiscussionRead  = "read:discussion"

	EnterpriseAdmin         = "admin:enterprise"
	EnterpriseManageRunners = "manage_runners:enterprise"
	EnterpriseManageBilling = "manage_billing:enterprise"
	EnterpriseRead          = "read:enterprise"

	ProjectAll  = "project"
	ProjectRead = "read:project"

	GpgKeyAdmin = "admin:gpg_key"
	GpgKeyWrite = "write:gpg_key"
	GpgKeyRead  = "read:gpg_key"
)

func GetOrgRole(canAdminister *bool) OrganizationRole {
	switch {
	case canAdminister != nil && *canAdminister:
		return OrgRoleOwner
	default:
		// We only list organizations of which the user is a member
		return OrgRoleMember
	}
}

func initialScopes() TokenScopes {
	scopes := make(TokenScopes)
	scopes[RepoAdmin] = false
	scopes[RepoRepoStatus] = false
	scopes[RepoRepoDeployment] = false
	scopes[RepoPublicRepo] = false
	scopes[RepoRepoInvite] = false
	scopes[RepoSecurityEvents] = false
	scopes[RepoDelete] = false
	scopes[Workflow] = false
	scopes[PackagesWrite] = false
	scopes[PackagesRead] = false
	scopes[PackagesDelete] = false
	scopes[OrgAdmin] = false
	scopes[OrgWrite] = false
	scopes[OrgRead] = false
	scopes[PublicKeyAdmin] = false
	scopes[PublicKeyWrite] = false
	scopes[PublicKeyRead] = false
	scopes[OrgHookAdmin] = false
	scopes[RepoHookAdmin] = false
	scopes[RepoHookWrite] = false
	scopes[RepoHookRead] = false
	scopes[Gist] = false
	scopes[Notifications] = false
	scopes[UserAll] = false
	scopes[UserRead] = false
	scopes[UserEmail] = false
	scopes[UserFollow] = false
	scopes[DiscussionWrite] = false
	scopes[DiscussionRead] = false
	scopes[EnterpriseAdmin] = false
	scopes[EnterpriseManageRunners] = false
	scopes[EnterpriseManageBilling] = false
	scopes[EnterpriseRead] = false
	scopes[ProjectAll] = false
	scopes[ProjectRead] = false
	scopes[GpgKeyAdmin] = false
	scopes[GpgKeyWrite] = false
	scopes[GpgKeyRead] = false
	return scopes
}

func denormalizeScopes(scopes TokenScopes) TokenScopes {
	if scopes[RepoAdmin] {
		scopes[RepoRepoStatus] = true
		scopes[RepoRepoDeployment] = true
		scopes[RepoPublicRepo] = true
		scopes[RepoRepoInvite] = true
		scopes[RepoSecurityEvents] = true
		scopes[RepoDelete] = true

		// implicitly implied, although not shown in GH GUI.
		scopes[RepoHookAdmin] = true
		scopes[Workflow] = true
	}

	if scopes[RepoHookAdmin] {
		scopes[RepoHookWrite] = true
	}
	if scopes[RepoHookWrite] {
		scopes[RepoHookRead] = true
	}

	if scopes[OrgAdmin] {
		scopes[OrgWrite] = true
		scopes[OrgRead] = true

		// implicitly implied, although not shown in GH GUI.
		scopes[OrgHookAdmin] = true
		scopes[ProjectAll] = true
		scopes[PackagesWrite] = true
		scopes[PackagesDelete] = true
		scopes[DiscussionWrite] = true
	}

	if scopes[UserAll] {
		scopes[UserEmail] = true
		scopes[UserFollow] = true
		scopes[UserRead] = true

		// implicitly implied, although not shown in GH GUI.
		scopes[PublicKeyAdmin] = true
		scopes[GpgKeyAdmin] = true
		scopes[Notifications] = true
		scopes[Gist] = true
	}

	if scopes[PackagesWrite] {
		scopes[PackagesRead] = true
	}

	if scopes[PublicKeyAdmin] {
		scopes[PublicKeyWrite] = true
	}
	if scopes[PublicKeyWrite] {
		scopes[PublicKeyRead] = true
	}

	if scopes[GpgKeyAdmin] {
		scopes[GpgKeyWrite] = true
	}
	if scopes[GpgKeyWrite] {
		scopes[GpgKeyRead] = true
	}

	if scopes[DiscussionWrite] {
		scopes[DiscussionRead] = true
	}

	if scopes[EnterpriseAdmin] {
		scopes[EnterpriseManageBilling] = true
		scopes[EnterpriseManageRunners] = true
		scopes[EnterpriseRead] = true
	}

	if scopes[ProjectAll] {
		scopes[ProjectRead] = true
	}

	if scopes[GpgKeyAdmin] {
		scopes[GpgKeyWrite] = true
		scopes[GpgKeyRead] = true
	}

	return scopes
}

func ParseTokenScopes(scopesList []string) TokenScopes {
	scopes := initialScopes()

	for _, scope := range scopesList {
		scopes[scope] = true
	}
	scopes = denormalizeScopes(scopes)

	return scopes
}

var orgMemberValidScopes = map[TokenScope]bool{
	RepoAdmin:          false,
	RepoRepoStatus:     false,
	RepoRepoDeployment: false,
	RepoPublicRepo:     false,
	RepoRepoInvite:     false,
	RepoSecurityEvents: false,
	RepoDelete:         false,

	Workflow: false,

	PackagesWrite:  false,
	PackagesRead:   true,
	PackagesDelete: false,

	OrgAdmin: false,
	OrgWrite: false,
	OrgRead:  true,

	PublicKeyAdmin: true,
	PublicKeyWrite: true,
	PublicKeyRead:  true,

	OrgHookAdmin:  false,
	RepoHookAdmin: false,
	RepoHookWrite: false,
	RepoHookRead:  false,

	Gist: true,

	Notifications: true,

	UserAll:    true,
	UserRead:   true,
	UserEmail:  true,
	UserFollow: true,

	DiscussionWrite: true,
	DiscussionRead:  true,

	EnterpriseAdmin:         false,
	EnterpriseManageRunners: false,
	EnterpriseManageBilling: false,
	EnterpriseRead:          false,

	ProjectAll:  true,
	ProjectRead: true,

	GpgKeyAdmin: true,
	GpgKeyWrite: true,
	GpgKeyRead:  true,
}

func HasOrgScope(toCheck TokenScope, scopes TokenScopes, orgRole OrganizationRole) bool {
	switch orgRole {
	case OrgRoleOwner:
		return scopes[toCheck]
	case OrgRoleMember:
		if allowed, ok := orgMemberValidScopes[toCheck]; ok && allowed {
			return scopes[toCheck]
		}
	}

	return false
}

var repoAdminValidScopes = map[TokenScope]bool{
	RepoAdmin:          true,
	RepoRepoStatus:     true,
	RepoRepoDeployment: true,
	RepoPublicRepo:     true,
	RepoRepoInvite:     true,
	RepoSecurityEvents: true,
	RepoDelete:         true,

	Workflow: true,

	PackagesWrite:  true,
	PackagesRead:   true,
	PackagesDelete: true,

	OrgAdmin: false,
	OrgWrite: false,
	OrgRead:  false,

	PublicKeyAdmin: true,
	PublicKeyWrite: true,
	PublicKeyRead:  true,

	OrgHookAdmin:  false,
	RepoHookAdmin: true,
	RepoHookWrite: true,
	RepoHookRead:  true,

	Gist: true,

	Notifications: true,

	UserAll:    true,
	UserRead:   true,
	UserEmail:  true,
	UserFollow: true,

	DiscussionWrite: true,
	DiscussionRead:  true,

	EnterpriseAdmin:         false,
	EnterpriseManageRunners: false,
	EnterpriseManageBilling: false,
	EnterpriseRead:          false,

	ProjectAll:  true,
	ProjectRead: true,

	GpgKeyAdmin: true,
	GpgKeyWrite: true,
	GpgKeyRead:  true,
}

var repoNonAdminValidScopes = map[TokenScope]bool{
	RepoAdmin:          false,
	RepoRepoStatus:     true,
	RepoRepoDeployment: true,
	RepoPublicRepo:     true,
	RepoRepoInvite:     false,
	RepoSecurityEvents: false,
	RepoDelete:         false,

	Workflow: true,

	PackagesWrite:  true,
	PackagesRead:   true,
	PackagesDelete: true,

	OrgAdmin: false,
	OrgWrite: false,
	OrgRead:  false,

	PublicKeyAdmin: true,
	PublicKeyWrite: true,
	PublicKeyRead:  true,

	OrgHookAdmin:  false,
	RepoHookAdmin: true,
	RepoHookWrite: true,
	RepoHookRead:  true,

	Gist: true,

	Notifications: true,

	UserAll:    true,
	UserRead:   true,
	UserEmail:  true,
	UserFollow: true,

	DiscussionWrite: true,
	DiscussionRead:  true,

	EnterpriseAdmin:         false,
	EnterpriseManageRunners: false,
	EnterpriseManageBilling: false,
	EnterpriseRead:          false,

	ProjectAll:  false,
	ProjectRead: true,

	GpgKeyAdmin: true,
	GpgKeyWrite: true,
	GpgKeyRead:  true,
}

var repoReadValidScopes = map[TokenScope]bool{
	RepoAdmin:          false,
	RepoRepoStatus:     true,
	RepoRepoDeployment: false,
	RepoPublicRepo:     true,
	RepoRepoInvite:     false,
	RepoSecurityEvents: false,
	RepoDelete:         false,

	Workflow: false,

	PackagesWrite:  false,
	PackagesRead:   true,
	PackagesDelete: false,

	OrgAdmin: false,
	OrgWrite: false,
	OrgRead:  true,

	PublicKeyAdmin: false,
	PublicKeyWrite: false,
	PublicKeyRead:  true,

	OrgHookAdmin:  false,
	RepoHookAdmin: true,
	RepoHookWrite: true,
	RepoHookRead:  true,

	Gist: true,

	Notifications: true,

	UserAll:    true,
	UserRead:   true,
	UserEmail:  true,
	UserFollow: true,

	DiscussionWrite: false,
	DiscussionRead:  true,

	EnterpriseAdmin:         false,
	EnterpriseManageRunners: false,
	EnterpriseManageBilling: false,
	EnterpriseRead:          false,

	ProjectAll:  false,
	ProjectRead: true,

	GpgKeyAdmin: true,
	GpgKeyWrite: true,
	GpgKeyRead:  true,
}

func HasRepoScope(toCheck TokenScope, scopes TokenScopes, repoRole RepositoryRole) bool {
	var mapping map[string]bool
	switch repoRole {
	case RepoRoleAdmin:
		mapping = repoAdminValidScopes
	case RepoRoleMaintainer:
		fallthrough
	case RepoRoleWrite:
		fallthrough
	case RepoRoleTriage:
		mapping = repoNonAdminValidScopes
	case RepoRoleRead:
		mapping = repoReadValidScopes
	}
	allowed, ok := mapping[toCheck]
	return ok && allowed && scopes[toCheck]
}

type EnterpriseRole = string

var enterpriseAdminValidScopes = map[TokenScope]bool{
	EnterpriseAdmin:         true,
	EnterpriseManageRunners: false,
	EnterpriseManageBilling: false,
	EnterpriseRead:          false,
}

func HasEnterpriseScope(toCheck TokenScope, scopes TokenScopes, enterpriseRole EnterpriseRole) bool {
	var mapping map[string]bool
	switch enterpriseRole {
	case EnterpriseAdmin:
		mapping = enterpriseAdminValidScopes
	}

	allowed, ok := mapping[toCheck]
	return ok && allowed && scopes[toCheck]
}
