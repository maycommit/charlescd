package main

import (
	"charlescd/internal/manager/circle"
	"charlescd/internal/operator/sync"
	"flag"
	"k8s.io/client-go/discovery"
	"log"
	"path/filepath"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
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

	f := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, v1.NamespaceAll, nil)
	i := f.ForResource(circle.Resource).Informer()
	i.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			logrus.Info("received add event!")
			err := sync.Start(client, disco, obj.(*unstructured.Unstructured))
			if err != nil {
				log.Fatalln(err)
			}
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			logrus.Info("received update event!")
			err := sync.Start(client, disco, obj.(*unstructured.Unstructured))
			if err != nil {
				log.Fatalln(err)
			}
		},
		DeleteFunc: func(obj interface{}) {
			logrus.Info("received delete event!")
			err := sync.Start(client, disco, obj.(*unstructured.Unstructured))
			if err != nil {
				log.Fatalln(err)
			}
		},
	})

	stopChan := make(chan struct{})
	log.Println("Start sync operator on port 8080...")
	i.Run(stopChan)
}
