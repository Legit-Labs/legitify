package gitlab

import (
	"log"
	"sync/atomic"

	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/clients/gitlab/pagination"
	"github.com/Legit-Labs/legitify/internal/collected/gitlab_collected"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
	"github.com/Legit-Labs/legitify/internal/common/permissions"
	"github.com/Legit-Labs/legitify/internal/common/types"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	gitlab2 "github.com/xanzy/go-gitlab"

	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"golang.org/x/net/context"
)

type repositoryCollector struct {
	collectors.BaseCollector
	Client  *gitlab.Client
	Context context.Context
}

func NewRepositoryCollector(ctx context.Context, client *gitlab.Client) collectors.Collector {
	c := &repositoryCollector{
		BaseCollector: collectors.NewBaseCollector(namespace.Repository),
		Client:        client,
		Context:       ctx,
	}
	return c
}

func (rc *repositoryCollector) CollectTotalEntities() int {
	repositories, exist := context_utils.GetRepositories(rc.Context)
	if exist {
		return len(repositories)
	}

	organizations, err := rc.Client.Organizations()
	if err != nil {
		log.Printf("failed to collect list of orgniazations to get repositories metadata %s", err)
		return 0
	}

	var total atomic.Int64
	gw := group_waiter.New()
	for _, org := range organizations {
		org := org
		gw.Do(func() {
			_, resp, err := rc.Client.Client().Groups.ListGroupProjects(org.ID, &gitlab2.ListGroupProjectsOptions{})
			if err != nil {
				log.Printf("failed to collect metadata for repositories of group %s (%d): %s", org.Name, org.ID, err)
			} else {
				total.Add(int64(resp.TotalItems))
			}
		})
	}
	gw.Wait()

	return int(total.Load())
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
			r := r
			gw.Do(func() {
				project, _, err := rc.Client.Client().Projects.GetProject(getRepositoryEncodedName(r), &gitlab2.GetProjectOptions{})
				if err != nil {
					log.Println(err.Error())
					return
				}

				rc.extendedCollection(project)
			})

		}
		gw.Wait()
	})
}

func getRepositoryEncodedName(repo types.RepositoryWithOwner) string {
	return repo.Owner + "/" + repo.Name
}

func (rc *repositoryCollector) extendProjectWithProtectedBranches(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	res, err := pagination.New[*gitlab2.ProtectedBranch](rc.Client.Client().ProtectedBranches.ListProtectedBranches, nil).Sync(int(project.ID()))
	if err != nil {
		log.Printf("failed to list projects %v", err)
		return project, err
	}

	extendedProject := project
	extendedProject.ProtectedBranches = res.Collected
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithMembers(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	res, err := pagination.New[*gitlab2.ProjectMember](rc.Client.Client().ProjectMembers.ListAllProjectMembers, nil).Sync(int(project.ID()))
	if err != nil {
		log.Printf("failed to list projects %v", err)
		return project, err
	}

	extendedProject := project
	extendedProject.Members = res.Collected
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithWebhooks(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	res, err := pagination.New[*gitlab2.ProjectHook](rc.Client.Client().Projects.ListProjectHooks, nil).Sync(int(project.ID()))
	if err != nil {
		log.Printf("failed to list project: %s webhook. error message: %s", project.Name(), err)
		return project, err
	}

	extendedProject := project
	extendedProject.Webhooks = res.Collected
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithPushRules(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	rules, _, err := rc.Client.Client().Projects.GetProjectPushRules(int(project.ID()))
	if err != nil {
		log.Printf("failed to get project push rule %s", err)
		return project, err
	}
	extendedProject := project
	extendedProject.PushRules = rules
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithMergeRequestApprovalRules(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	rules, _, err := rc.Client.Client().Projects.GetProjectApprovalRules(int(project.ID()))
	if err != nil {
		log.Printf("failed to get project merge request approval rules %s", err)
		return project, err
	}
	extendedProject := project
	extendedProject.ApprovalRules = rules
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithApprovalConfiguration(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	config, _, err := rc.Client.Client().Projects.GetApprovalConfiguration(int(project.ID()))
	if err != nil {
		log.Printf("failed to get project approval configuration %s", err)
		return project, err
	}
	extendedProject := project
	extendedProject.ApprovalConfiguration = config
	return extendedProject, nil
}

func (rc *repositoryCollector) extendProjectWithMinimumRequiredApprovals(project gitlab_collected.Repository) (gitlab_collected.Repository, error) {
	minRequiredApprovals := 0

	for _, rule := range project.ApprovalRules {
		if rule.ApprovalsRequired > minRequiredApprovals {
			minRequiredApprovals = rule.ApprovalsRequired
		}
	}

	project.MinimumRequiredApprovals = minRequiredApprovals
	return project, nil
}

func (rc *repositoryCollector) collectAll() collectors.SubCollectorChannels {
	return rc.WrappedCollection(func() {
		organizations, err := rc.Client.Organizations()
		if err != nil {
			log.Printf("failed to collect list of orgniazations to get repositories  %s", err)
			return
		}
		gw := group_waiter.New()
		for _, org := range organizations {
			org := org
			gw.Do(func() {
				ch := pagination.New[*gitlab2.Project](rc.Client.Client().Groups.ListGroupProjects, nil).Async(org.ID)
				for res := range ch {
					if res.Err != nil {
						log.Printf("failed to list projects for group %s (%d): %v", org.Name, org.ID, res.Err)
						return
					}
					for _, completedProject := range res.Collected {
						completedProject := completedProject
						gw.Do(func() {
							rc.extendedCollection(completedProject)
						})
					}
				}
			})
		}
		gw.Wait()
	})
}

func (rc *repositoryCollector) extendedCollection(completeProjectsList *gitlab2.Project) {
	proj := gitlab_collected.Repository{
		Project: completeProjectsList,
	}
	extensionFunctions := []func(gitlab_collected.Repository) (gitlab_collected.Repository, error){
		rc.extendProjectWithMembers,
		rc.extendProjectWithProtectedBranches,
		rc.extendProjectWithWebhooks,
		rc.extendProjectWithPushRules,
		rc.extendProjectWithMergeRequestApprovalRules,
		rc.extendProjectWithApprovalConfiguration,
		rc.extendProjectWithMinimumRequiredApprovals,
	}
	var err error
	for _, f := range extensionFunctions {
		proj, err = f(proj)
		if err != nil {
			log.Printf("project '%s' collection failed with error: %v", proj.Name(), err)
			break
		}
	}

	if err == nil {
		newContext := newCollectionContext(nil, []permissions.OrganizationRole{permissions.OrgRoleOwner})
		rc.CollectDataWithContext(proj, proj.Links.Self, &newContext)
	}
	rc.CollectionChangeByOne()
}
