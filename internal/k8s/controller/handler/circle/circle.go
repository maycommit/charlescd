package circle

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/pkg/customerror"
	circlerrExternalversions "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	k8sClientCache "k8s.io/client-go/tools/cache"
)

type handler struct {
	appCache *cache.Cache
}

func (h *handler) addFunc(obj interface{}) {
	// TODO: ADD TO CACHE
}

func (h *handler) updateFunc(oldObj interface{}, newObj interface{}) {
	// TODO: UPDATE CACHE
	// TODO: VERIFY CIRCLE STATUS FOR UPDATE ROUTES
}

func (h *handler) deleteFunc(obj interface{}) {
	// TODO: SET CIRCLE DELETION
	// TODO: UPDATE ROUTES
}

func New(
	stopCh <-chan struct{},
	factory circlerrExternalversions.SharedInformerFactory,
	appCache *cache.Cache) {

	h := &handler{
		appCache: appCache,
	}

	circleInformer := factory.Circle().V1alpha1().Circles().Informer()
	circleInformer.AddEventHandler(k8sClientCache.ResourceEventHandlerFuncs{
		AddFunc:    h.addFunc,
		UpdateFunc: h.updateFunc,
		DeleteFunc: h.deleteFunc,
	})

	circleInformer.SetWatchErrorHandler(func(r *k8sClientCache.Reflector, err error) {
		logrus.Warningln(customerror.LogFields(customerror.New("Failed circle handler!", err, nil, "handler.circle.New")))
	})

	circleInformer.Run(stopCh)
}
