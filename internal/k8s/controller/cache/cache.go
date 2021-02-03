package cache

import (
	gitopsCache "github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/cluster"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"

	"k8s.io/client-go/rest"
)

type cache struct {
	project *project.ProjectsCache
	circle  *circle.CirclesCache
	cluster gitopsCache.ClusterCache
}

type Cache interface {
	Project() *project.ProjectsCache
	Circle() *circle.CirclesCache
	Cluster() gitopsCache.ClusterCache
}

func (c *cache) Project() *project.ProjectsCache {
	return c.project
}

func (c *cache) Circle() *circle.CirclesCache {
	return c.circle
}

func (c *cache) Cluster() gitopsCache.ClusterCache {
	return c.cluster
}

func New(config *rest.Config) Cache {
	return &cache{
		project: project.NewProjectCache(),
		circle:  circle.NewCircleCache(),
		cluster: cluster.NewClusterCache(config, []string{}),
	}
}
