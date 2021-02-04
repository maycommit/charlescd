package project

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/pkg/customerror"
	projectApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/project/v1alpha1"
	circlerrExternalversions "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	k8sClientCache "k8s.io/client-go/tools/cache"
)

type handler struct {
	appCache cache.Cache
}

func (h *handler) addFunc(obj interface{}) {
	// TODO: ADD TO CACHE
	newProject := obj.(*projectApi.Project)

	h.appCache.Projects().Put(newProject.GetName(), *newProject)
}

func (h *handler) updateFunc(oldObj interface{}, newObj interface{}) {
	// TODO: UPDATE CACHE
	newProject := newObj.(*projectApi.Project)
	h.appCache.Projects().Put(newProject.GetName(), *newProject)
}

func (h *handler) deleteFunc(obj interface{}) {
	// TODO: UPDATE CACHE
	project := obj.(*projectApi.Project)
	h.appCache.Projects().Delete(project.GetName())
}

func New(
	stopCh <-chan struct{},
	factory circlerrExternalversions.SharedInformerFactory,
	appCache cache.Cache) {

	h := &handler{
		appCache: appCache,
	}

	projectInformer := factory.Project().V1alpha1().Projects().Informer()
	projectInformer.AddEventHandler(k8sClientCache.ResourceEventHandlerFuncs{
		AddFunc:    h.addFunc,
		UpdateFunc: h.updateFunc,
		DeleteFunc: h.deleteFunc,
	})

	projectInformer.SetWatchErrorHandler(func(r *k8sClientCache.Reflector, err error) {
		logrus.Warningln(customerror.LogFields(customerror.New("Failed project handler!", err, nil, "handler.project.New")))
	})

	projectInformer.Run(stopCh)
}
