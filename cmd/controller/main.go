package main

import (
	"charlescd/internal/controller/sync"
	"context"
	"flag"
	"log"
	"path/filepath"
	"time"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2/klogr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	circleclientset "charlescd/pkg/client/clientset/versioned"
)

const (
	namespace = "default"
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

	kubeClient := dynamic.NewForConfigOrDie(config)
	circleClient := circleclientset.NewForConfigOrDie(config)
	disco, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	clusterCache := sync.ClusterCache(config, []string{}, klogr)
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ticker.C:
			list, err := circleClient.CircleV1alpha1().Circles(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Fatalln(err)
			}

			for _, obj := range list.Items {
				syncConfig := sync.SyncConfig{
					ClusterCache: clusterCache,
					KubeClient:   kubeClient,
					Config:       config,
					Client:       circleClient.CircleV1alpha1().Circles(namespace),
					Disco:        disco,
					Circle:       &obj,
					Namespace:    namespace,
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
