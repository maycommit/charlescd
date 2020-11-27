package main

import (
	"charlescd/internal/controller"
	"context"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"

	circleclientset "charlescd/pkg/client/clientset/versioned"
	circleinformers "charlescd/pkg/client/informers/externalversions"
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

	kubeClient := kubernetes.NewForConfigOrDie(config)
	circleClient := circleclientset.NewForConfigOrDie(config)
	circleInformerFactory := circleinformers.NewSharedInformerFactory(circleClient, 0)


	ctrl := controller.NewController(
		"default",
		kubeClient,
		circleClient,
		circleInformerFactory.Circle().V1alpha1().Circles(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl.Run(ctx)
}
