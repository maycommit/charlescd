package cache

import (
	gitopsCache "github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/cluster"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"

	"k8s.io/client-go/rest"
)

type cache struct {
	projects *project.ProjectsCache
	circles  *circle.CirclesCache
	cluster gitopsCache.ClusterCache
}

type Cache interface {
	Projects() *project.ProjectsCache
	Circles() *circle.CirclesCache
	Cluster() gitopsCache.ClusterCache
}

func (c *cache) Projects() *project.ProjectsCache {
	return c.projects
}

func (c *cache) Circles() *circle.CirclesCache {
	return c.circles
}

func (c *cache) Cluster() gitopsCache.ClusterCache {
	return c.cluster
}

func New(config *rest.Config) Cache {
	return &cache{
		projects: project.NewProjectCache(),
		circles:  circle.NewCirclesCache(),
		cluster: cluster.NewClusterCache(config, []string{}),
	}
}
