package main

import (
	"flag"
	"os"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/engine"
	circleHandler "github.com/maycommit/circlerr/internal/k8s/controller/handler/circle"
	projectHandler "github.com/maycommit/circlerr/internal/k8s/controller/handler/project"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/kube"

	circlerrVersioned "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"
	circlerrExternalversions "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"
)

func init() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	masterUrl := flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	gitDir := flag.String("gitdir", "./tmp/git", "")
	flag.Parse()

	os.Setenv("GIT_DIR", *gitDir)
	os.Setenv("KUBECONFIG", *kubeconfig)
	os.Setenv("MASTER_URL", *masterUrl)
}

func main() {
	config, err := kube.GetClusterConfig()
	if err != nil {
		panic(err)
	}

	circlerrClientset := circlerrVersioned.NewForConfigOrDie(config)
	circlerrInformerFactory := circlerrExternalversions.NewSharedInformerFactory(circlerrClientset, 0)
	appCache := cache.New(config)
	e := engine.New(config, appCache, circlerrClientset)

	stopCh := make(chan struct{})
	go circleHandler.New(stopCh, circlerrInformerFactory, appCache)
	go projectHandler.New(stopCh, circlerrInformerFactory, appCache)
	e.Start()
}
