package main

import (
	"flag"
	"os"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/router"
	"github.com/maycommit/circlerr/internal/k8s/controller/sync/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/sync/cluster"
	"github.com/maycommit/circlerr/internal/k8s/controller/sync/project"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/kube"
	api "github.com/maycommit/circlerr/web/api/k8s/controller"

	kubeutil "github.com/argoproj/gitops-engine/pkg/utils/kube"

	"github.com/argoproj/gitops-engine/pkg/utils/tracing"
	"k8s.io/client-go/discovery"
	"k8s.io/klog/klogr"

	"k8s.io/client-go/dynamic"

	clientset "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"
	custominformers "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"
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

	namespace := "default"
	client := clientset.NewForConfigOrDie(config)
	kubeClient := dynamic.NewForConfigOrDie(config)
	customInformerFactory := custominformers.NewSharedInformerFactory(client, 0)
	controllerCache := cache.NewCache(config, namespace)
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	var currentRouter router.UseCases
	if os.Getenv("ROUTER_TYPE") != "" {
		currentRouter = router.NewRouter(controllerCache, config, namespace)
	}

	stopCh := make(chan struct{})
	kubectl := &kubeutil.KubectlCmd{
		Log:    klogr.New(),
		Tracer: tracing.NopTracer{},
	}
	clusterSyncOpts := cluster.New(config, kubeClient, namespace, controllerCache, currentRouter, kubectl, client)

	go circle.Run(stopCh, customInformerFactory, controllerCache)
	go project.Run(stopCh, customInformerFactory, controllerCache)
	go clusterSyncOpts.Run(stopCh)

	api.NewApi(
		controllerCache,
		client,
		kubeClient,
		discoveryClient,
	).Start()
}
