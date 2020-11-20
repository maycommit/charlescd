package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	appclientset "github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
	appinformers "github.com/argoproj/argo-cd/pkg/client/informers/externalversions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	appClient := appclientset.NewForConfigOrDie(config)
	resyncDuration := time.Duration(180) * time.Second
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInformerFactory := appinformers.NewFilteredSharedInformerFactory(
		appClient,
		resyncDuration,
		"default",
		func(options *metav1.ListOptions) {},
	)

	informer := appInformerFactory.Argoproj().V1alpha1().Applications().Informer()
	// lister := appInformerFactory.Argoproj().V1alpha1().Applications().Lister()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Add obj")

			fmt.Println(obj)
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("Update obj")
			fmt.Println(old)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete obj")
			fmt.Println(obj)
		},
	})

	informer.Run(ctx.Done())
}
