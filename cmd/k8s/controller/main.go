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
	fGitDir := flag.String("gitdir", "./tmp/git", "")
	fRouterType := flag.String("router", os.Getenv("ROUTER_TYPE"), "")
	fKubeconfigPath := flag.String("kubepath", os.Getenv("KUBECONFIG_PATH"), "")
	fK8sConnType := flag.String("k8sconntype", os.Getenv("K8S_CONN_TYPE"), "")
	flag.Parse()

	os.Setenv("GIT_DIR", *fGitDir)
	os.Setenv("ROUTER_TYPE", *fRouterType)
	os.Setenv("KUBECONFIG_PATH", *fKubeconfigPath)
	os.Setenv("K8S_CONN_TYPE", *fK8sConnType)
}

func main() {
	config, err := kube.GetClusterConfig()
	if err != nil {
		panic(err)
	}

	circlerrClientset := circlerrVersioned.NewForConfigOrDie(config)
	// kubeClient := dynamic.NewForConfigOrDie(config)
	circlerrInformerFactory := circlerrExternalversions.NewSharedInformerFactory(circlerrClientset, 0)
	appCache := cache.New(config)
	e := engine.New(appCache)

	stopCh := make(chan struct{})
	go circleHandler.New(stopCh, circlerrInformerFactory, nil)
	go projectHandler.New(stopCh, circlerrInformerFactory, nil)
	go e.Start()

	// controllerCache := cache.NewCache(config, namespace)
	// discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	// var currentRouter router.UseCases
	// if os.Getenv("ROUTER_TYPE") != "" {
	// 	currentRouter = router.NewRouter(controllerCache, config, namespace)
	// }

	// kubectl := &kubeutil.KubectlCmd{
	// 	Log:    klogr.New(),
	// 	Tracer: tracing.NopTracer{},
	// }
	// clusterSyncOpts := cluster.New(config, kubeClient, namespace, controllerCache, currentRouter, kubectl, client)

	// go circle.Run(stopCh, customInformerFactory, controllerCache)
	// go project.Run(stopCh, customInformerFactory, controllerCache)
	// go clusterSyncOpts.Run(stopCh)

	// api.NewApi(
	// 	controllerCache,
	// 	client,
	// 	kubeClient,
	// 	discoveryClient,
	// ).Start()
}
