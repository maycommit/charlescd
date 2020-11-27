package controller

import (
	circleclientset "charlescd/pkg/client/clientset/versioned"
	informers "charlescd/pkg/client/informers/externalversions/circle/v1alpha1"
	listers "charlescd/pkg/client/listers/circle/v1alpha1"
	"context"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	namespace       string
	kubeClientset   kubernetes.Interface
	circleClientset circleclientset.Interface
	circleLister    listers.CircleLister
	circleInformer  cache.SharedIndexInformer
	workqueue       workqueue.RateLimitingInterface
}

func NewController(
	namespace string,
	kubeClientset kubernetes.Interface,
	circleClientset circleclientset.Interface,
	circleInformer informers.CircleInformer,
) *Controller {
	ctrl := &Controller{
		namespace:       namespace,
		kubeClientset:   kubeClientset,
		circleClientset: circleClientset,
		circleLister:    circleInformer.Lister(),
		circleInformer:  circleInformer.Informer(),
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Circles"),
	}

	ctrl.circleInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var key string
			var err error
			if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
				log.Fatalln(err)
			}

			ctrl.workqueue.Add(key)
		},
	})

	return ctrl
}

func (ctrl *Controller) Run(ctx context.Context) {
	defer ctrl.workqueue.ShutDown()

	go ctrl.circleInformer.Run(ctx.Done())

	select {}
}
