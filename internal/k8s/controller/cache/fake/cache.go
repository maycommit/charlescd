package fake

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"

	gitopsCache "github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/cache/mocks"

	"k8s.io/client-go/rest"
)

type resourceInfo struct {
	circleMark  string
	projectMark string
	releaseMark string
	routerMark  string
}

type cluster struct {
	cache *mocks.ClusterCache
}

func (c *cluster) IsManagedResource(r *gitopsCache.Resource) bool {
	return r.Info.(*resourceInfo).circleMark != ""
}

func (c *cluster) Get() gitopsCache.ClusterCache {

	return c.cache
}

func NewCache(config *rest.Config, namespace string) *cache.Cache {
	mockClusterCache := &mocks.ClusterCache{}

	return &cache.Cache{
		Projects: project.NewProjectCache(),
		Circles:  circle.NewCircleCache(),
		Cluster: &cluster{
			cache: mockClusterCache,
		},
	}
}
