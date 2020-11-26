package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	circleclientset "charlescd/pkg/client/clientset/versioned"
	circleinformers "charlescd/pkg/client/informers/externalversions"
	circlelisters "charlescd/pkg/client/listers/circle/v1alpha1"
	informers "charlescd/pkg/client/informers/externalversions/circle/v1alpha1"
	listers "charlescd/pkg/client/listers/circle/v1alpha1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

type Controller struct {
	namespace string
	kubeClientset kubernetes.Interface
	circleClientset circleclientset.Interface
	circleLister listers.CircleLister
	workqueue workqueue.RateLimitingInterface
}

func NewController(
	namespace string,
	kubeClientset kubernetes.Interface,
	circleClientset circleclientset.Interface,
	circleInformer informers.CircleInformer,
) (*Controller, error) {
	controller := Controller{
		namespace: namespace,
		kubeClientset: kubeClientset,
		circleClientset: circleClientset,
		circleLister: circleInformer.Lister(),
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Circles"),
	}

	circleInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var key string
			var err error
			if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
				log.Fatalln(err)
			}

			controller.workqueue.Add(key)
		},
	})

}

func (ctrl Controller) newCircleInformerAndLister() {
	circleInformerFactory := circleinformers.NewSharedInformerFactory(
		ctrl.circleClientset,
		0,
	)
	informer := circleInformerFactory.
}
