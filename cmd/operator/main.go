package main

import (
	"charlescd/internal/manager/circle"
	"charlescd/internal/operator/sync"
	"context"
	"flag"
	"log"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2/klogr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	klogr := klogr.New()
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	clusterCache := sync.ClusterCache(config, []string{}, klogr)

	ticker := time.NewTicker(3 * time.Second)

	factory := informers.NewSharedInformerFactory(clientset, 0)
	podInformer := factory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				newPod := obj.(*v1.Pod)
				newPod.GetAnnotations()
			},
			UpdateFunc: func(old, new interface{}) {
				oldPod := old.(*v1.Pod)
				newPod := new.(*v1.Pod)

				oldPod.GetAnnotations()
				newPod.GetAnnotations()
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				pod.GetAnnotations()
			},
		},
	)

	stop := make(chan struct{})
	go func(i cache.SharedIndexInformer) {
		i.Run(stop)
	}(podInformer.Informer())

	for {
		select {
		case <-ticker.C:
			list, err := client.Resource(circle.Resource).Namespace("default").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Fatalln(err)
			}

			for _, obj := range list.Items {
				syncConfig := sync.SyncConfig{
					ClusterCache: clusterCache,
					Config:       config,
					Disco:        disco,
					CircleRes:    &obj,
					Namespace:    "default",
					Prune:        true,
					Log:          klogr,
				}
				err := sync.Start(syncConfig)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}
