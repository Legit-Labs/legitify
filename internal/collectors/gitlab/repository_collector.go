package gitlab

import (
	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/Legit-Labs/legitify/internal/context_utils"

	gitlab2 "github.com/xanzy/go-gitlab"
	"log"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"

	_ "github.com/xanzy/go-gitlab"
)

type repositoryCollector struct {
	collectors.BaseCollector
	Client  *gitlab.Client
	Context context.Context
}

func NewRepositoryCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &repositoryCollector{
		Client:  client,
		Context: ctx,
	}
	collectors.InitBaseCollector(&c.BaseCollector, c)
	return c
}

func (rc *repositoryCollector) Namespace() namespace.Namespace {
	return namespace.Repository
}

func (rc *repositoryCollector) CollectMetadata() collectors.Metadata {
	_, resp, err := rc.Client.Client().Projects.ListProjects(&gitlab2.ListProjectsOptions{Owned: gitlab2.Bool(true)})
	res := collectors.Metadata{}

	if err != nil {
		log.Printf("failed to collect metadata for repositories %s", err)
	} else {
		res.TotalEntities = resp.TotalItems
	}
	return res
}

func (rc *repositoryCollector) Collect() collectors.SubCollectorChannels {
	repositories, exist := context_utils.GetRepositories(rc.Context)

	if exist {
		return rc.collectSpecific(repositories)
	}

	return rc.collectAll()
}

func (rc *repositoryCollector) collectSpecific(repositories []types.RepositoryWithOwner) collectors.SubCollectorChannels {
	return rc.WrappedCollection(func() {
		gw := group_waiter.New()
		for _, r := range repositories {
			repo := r
			gw.Do(func() {
				project, _, err := rc.Client.Client().Projects.GetProject(repo.Owner+"%2F"+repo.Name, &gitlab2.GetProjectOptions{})
				if err != nil {
					log.Println(err.Error())
					return
				}
				proj := gitlab_collected.Repository{
					Project:                      project,
					VulnerabilityAlertsEnabled:   gitlab2.Bool(true),
					NoBranchProtectionPermission: true,
				}
				a := repositoryContext{testParam: true, isEnterprise: false, roles: nil}
				rc.CollectDataWithContext(proj, proj.Links.Self, &a)
			})

		}
		gw.Wait()
	})
}

func (rc *repositoryCollector) collectAll() collectors.SubCollectorChannels {
	return rc.WrappedCollection(func() {

		projects, _, err := rc.Client.Client().Projects.ListProjects(&gitlab2.ListProjectsOptions{Owned: gitlab2.Bool(true)})

		if err != nil {
			log.Printf("failed to collect metadata for repositories %s", err)
		}

		for _, project := range projects {
			if err != nil {
				log.Println(err.Error())
				return
			}
			members, _, err := rc.Client.Client().ProjectMembers.ListAllProjectMembers(project.ID, &gitlab2.ListProjectMembersOptions{})

			if err != nil {
				log.Printf("failed to collect project %s members: %s", project.Name, err)
				return
			}

			protectedBranches, _, err := rc.Client.Client().ProtectedBranches.ListProtectedBranches(project.ID, &gitlab2.ListProtectedBranchesOptions{})

			if err != nil {
				log.Printf("failed to collect project %s protected branches: %s", project.Name, err)
				return
			}

			proj := gitlab_collected.Repository{
				Project:                      project,
				Members:                      members,
				ProtectedBranches:            protectedBranches,
				VulnerabilityAlertsEnabled:   gitlab2.Bool(true),
				NoBranchProtectionPermission: true,
			}
			a := repositoryContext{testParam: true, isEnterprise: false, roles: nil}
			rc.CollectDataWithContext(proj, proj.Links.Self, &a)
		}

	})
}
