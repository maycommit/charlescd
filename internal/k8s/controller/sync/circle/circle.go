package circle

import (
	"log"

	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"

	cache2 "k8s.io/client-go/tools/cache"

	custominformer "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"
)

func Run(stopCh <-chan struct{}, factory custominformer.SharedInformerFactory, appCache *cache.Cache) {
	circleInformer := factory.Circle().V1alpha1().Circles().Informer()
	circleInformer.AddEventHandler(cache2.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			circle := obj.(*v1alpha1.Circle)
			circle.Kind = "Circle"
			circle.APIVersion = "circlerr.io/v1alpha1"
			appCache.Circles.Put(circle.GetName(), *circle)
		},
		UpdateFunc: func(old, new interface{}) {
			_ = old.(*v1alpha1.Circle)
			newCircle := new.(*v1alpha1.Circle)
			newCircle.Kind = "Circle"
			newCircle.APIVersion = "circlerr.io/v1alpha1"
			appCache.Circles.Put(newCircle.GetName(), *newCircle)

			if newCircle.Status.Projects != nil {
				for _, p := range newCircle.Status.Projects {
					if newCircle.Status.Status == health.HealthStatusHealthy {
						appCache.Projects.Get(p.Name).SetRoute(newCircle.GetName(), project.ProjectRoute{
							CircleID:    string(newCircle.UID),
							ReleaseName: newCircle.Spec.Release.Name,
						})
					} else {
						appCache.Projects.Get(p.Name).DeleteRoute(newCircle.GetName())
					}
				}
			}

		},
		DeleteFunc: func(obj interface{}) {
			circle := obj.(*v1alpha1.Circle)
			appCache.Circles.LogicDeletion(circle.GetName())

			if circle.Status.Projects != nil {
				for _, p := range circle.Status.Projects {
					appCache.Projects.Get(p.Name).DeleteRoute(circle.GetName())
				}
			}

		},
	})

	_ = circleInformer.SetWatchErrorHandler(func(r *cache2.Reflector, err error) {
		log.Fatalln(err)
	})

	circleInformer.Run(stopCh)
}
