package cluster

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"

	"k8s.io/client-go/rest"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type resourceInfo struct {
	circleMark  string
	projectMark string
	releaseMark string
	routerMark  string
}

type ClusterCache interface {
	IsManagedResource(r *cache.Resource) bool
	Get() cache.ClusterCache
}

type clusterCache struct {
	cache cache.ClusterCache
}

func (c *clusterCache) IsManagedResource(r *cache.Resource) bool {
	return r.Info.(*resourceInfo).circleMark != ""
}

func (c *clusterCache) Get() cache.ClusterCache {
	return c.cache
}

func NewClusterCache(config *rest.Config, namespaces []string) ClusterCache {
	c := &clusterCache{}

	c.cache = cache.NewClusterCache(config,
		cache.SetNamespaces(namespaces),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {

			var circleMark string
			var projectMark string
			var releaseMark string

			circleMark = un.GetAnnotations()[annotation.CircleAnnotation]
			projectMark = un.GetAnnotations()[annotation.ProjectAnnotation]
			releaseMark = un.GetAnnotations()[annotation.ReleaseAnnotation]

			info = &resourceInfo{
				projectMark: projectMark,
				circleMark:  circleMark,
				releaseMark: releaseMark,
			}

			cacheManifest = circleMark != "" && projectMark != "" && releaseMark != ""
			return
		}),
	)

	return c
}
