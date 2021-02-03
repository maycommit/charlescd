package istio

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"

	"k8s.io/client-go/rest"

	istioclient "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1"
)

type IstioRouter struct {
	appcache         *cache.Cache
	virtualService   istioclient.VirtualServiceInterface
	destinationRules istioclient.DestinationRuleInterface
}

func NewIstioRouter(appcache *cache.Cache, config *rest.Config, namespace string) IstioRouter {
	return IstioRouter{
		appcache:         appcache,
		virtualService:   istioclient.NewForConfigOrDie(config).VirtualServices(namespace),
		destinationRules: istioclient.NewForConfigOrDie(config).DestinationRules(namespace),
	}
}

func (r IstioRouter) Sync() apperror.Error {
	for projectName, project := range r.appcache.Projects.List() {
		err := r.manageVirtualServices(projectName, project)
		if err != nil {
			return err
		}
	}

	return nil
}
