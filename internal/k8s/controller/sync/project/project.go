package project

import (
	"log"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/pkg/k8s/controller/apis/project/v1alpha1"
	custominformer "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"

	cache2 "k8s.io/client-go/tools/cache"
)

func Run(stopCh <-chan struct{}, factory custominformer.SharedInformerFactory, appCache *cache.Cache) {
	projectInformer := factory.Project().V1alpha1().Projects().Informer()
	projectInformer.AddEventHandler(cache2.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			project := obj.(*v1alpha1.Project)
			appCache.Projects.Put(project.GetName(), *project)
		},
		UpdateFunc: func(old, new interface{}) {
			_ = old.(*v1alpha1.Project)
			newProject := new.(*v1alpha1.Project)
			appCache.Projects.Put(newProject.GetName(), *newProject)
		},
		DeleteFunc: func(obj interface{}) {
			project := obj.(*v1alpha1.Project)
			appCache.Projects.Delete(project.GetName())
		},
	})

	projectInformer.SetWatchErrorHandler(func(r *cache2.Reflector, err error) {
		log.Fatalln(err)
	})

	projectInformer.Run(stopCh)
}
