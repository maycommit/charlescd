package cluster

import (
	cacheUtils "github.com/maycommit/circlerr/internal/k8s/controller/utils/cache"

	"k8s.io/client-go/rest"

	"github.com/argoproj/gitops-engine/pkg/cache"
)

func NewClusterCache(config *rest.Config, namespaces []string) cache.ClusterCache {
	return cache.NewClusterCache(config,
		cache.SetNamespaces(namespaces),
		cache.SetPopulateResourceInfoHandler(cacheUtils.ResourceInfoHandler),
	)
}
