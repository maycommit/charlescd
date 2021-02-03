package router

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/router/istio"

	"k8s.io/client-go/rest"
)

type UseCases interface {
	Sync() apperror.Error
}

func NewRouter(appcache *cache.Cache, config *rest.Config, namespace string) UseCases {
	return istio.NewIstioRouter(appcache, config, namespace)
}
