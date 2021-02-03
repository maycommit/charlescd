package cache

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/cluster"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"

	"k8s.io/client-go/rest"
)

type Cache struct {
	Projects *project.ProjectsCache
	Circles  *circle.CirclesCache
	Cluster  cluster.ClusterCache
}

func NewCache(config *rest.Config, namespace string) *Cache {
	return &Cache{
		Projects: project.NewProjectCache(),
		Circles:  circle.NewCircleCache(),
		Cluster:  cluster.NewClusterCache(config, []string{}),
	}
}
