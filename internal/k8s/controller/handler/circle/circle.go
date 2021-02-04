package circle

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/pkg/customerror"
	circleApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"
	circlerrExternalversions "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	k8sClientCache "k8s.io/client-go/tools/cache"
)

type handler struct {
	appCache cache.Cache
}

func (h *handler) addFunc(obj interface{}) {
	newCircle := obj.(*circleApi.Circle)

	h.appCache.Circles().Add(newCircle.GetName(), *newCircle)
}

func (h *handler) updateFunc(oldObj interface{}, newObj interface{}) {
	// TODO: UPDATE CACHE
	// TODO: VERIFY CIRCLE STATUS FOR UPDATE ROUTES
	newCircle := newObj.(*circleApi.Circle)

	h.appCache.Circles().Set(newCircle.GetName(), *newCircle)
}

func (h *handler) deleteFunc(obj interface{}) {
	// TODO: UPDATE ROUTES
	newCircle := obj.(*circleApi.Circle)
	h.appCache.Circles().Get(newCircle.GetName()).SetDeletion()
}

func New(
	stopCh <-chan struct{},
	factory circlerrExternalversions.SharedInformerFactory,
	appCache cache.Cache) {

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
