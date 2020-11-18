package main

import (
	"charlescd/internal/manager/circle"
	"charlescd/internal/operator/sync"
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/discovery"
	"k8s.io/klog/v2/klogr"
	"log"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
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

	disco, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	clusterCache := sync.ClusterCache(config, []string{}, klogr)

	resync := make(chan bool)

	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, v1.NamespaceAll, nil)
	i := f.ForResource(circle.Resource).Informer()
	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			logrus.Info("received add event!")
			resync <- true
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			logrus.Info("received update event!")
			resync <- true
		},
		DeleteFunc: func(obj interface{}) {
			logrus.Info("received delete event!")
			resync <- true
		},
	})

	stopChan := make(chan struct{})

	go func(i cache.SharedInformer) {
		log.Println("Start sync operator on port 8080...")
		i.Run(stopChan)
	}(i)

	ticker := time.NewTicker(3 * time.Second)

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
		case <-resync:
			fmt.Println("CIRCLE RESYNC")
		}
	}
}
