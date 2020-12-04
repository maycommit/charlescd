package main

import (
	"charlescd/internal/controller/sync"
	"context"
	"flag"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/argoproj/gitops-engine/pkg/engine"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2/klogr"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"charlescd/cmd/controller/server"
	"charlescd/pkg/apis/circle/v1alpha1"
	circleclientset "charlescd/pkg/client/clientset/versioned"
	circleinformer "charlescd/pkg/client/informers/externalversions"
	circlepb "charlescd/pkg/grpc/circle"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
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
	istioClient := istioclient.NewForConfigOrDie(config)
	circleClient := circleclientset.NewForConfigOrDie(config)
	circleInformerFactory := circleinformer.NewSharedInformerFactory(circleClient, 0)
	disco, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	clusterCache := sync.ClusterCache(config, []string{}, klogr)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	circlepb.RegisterCircleServiceServer(s, &server.GRPCServer{
		ClusterCache:                     &clusterCache,
		CircleClientset:                  circleClient,
		UnimplementedCircleServiceServer: circlepb.UnimplementedCircleServiceServer{},
	})

	go func() {
		log.Println("GRPC server started on port 9090!")
		if err := s.Serve(lis); err != nil {
			log.Fatalln(err)
		}

	}()

	gitOpsEngine := engine.NewEngine(config, clusterCache)
	stopEngine, err := gitOpsEngine.Run()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	syncConfig := sync.SyncConfig{
		Ctx:          ctx,
		ClusterCache: clusterCache,
		KubeClient:   kubeClient,
		IstioClient:  istioClient.NetworkingV1beta1(),
		Config:       config,
		GitopsEngine: gitOpsEngine,
		Client:       circleClient.CircleV1alpha1().Circles(namespace),
		Disco:        disco,
		Namespace:    namespace,
		Prune:        false,
		Log:          klogr,
		StopEngine:   stopEngine,
	}

	circles, err := syncConfig.GetInitialCircleState()
	if err != nil {
		log.Fatalln(err)
	}

	circleInformer := circleInformerFactory.Circle().V1alpha1().Circles().Informer()
	circleInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			circle := obj.(*v1alpha1.Circle)
			circleName := circle.GetName()
			manifests, err := syncConfig.GetManifests(*circle)
			if err != nil {
				log.Fatalln(err)
			}
			circles[circleName] = sync.CircleState{
				Manifests: manifests,
				Synced:    false,
			}
		},
		UpdateFunc: func(old, new interface{}) {
			_ = old.(*v1alpha1.Circle)
			newCircle := new.(*v1alpha1.Circle)

			circleName := newCircle.GetName()
			manifests, err := syncConfig.GetManifests(*newCircle)
			if err != nil {
				log.Fatalln(err)
			}

			circles[circleName] = sync.CircleState{
				Manifests: manifests,
				Synced:    false,
			}
			// TODO: Implement diff circles for change sync status
		},
		DeleteFunc: func(obj interface{}) {
			circle := obj.(*v1alpha1.Circle)

			circles[circle.GetName()] = sync.CircleState{
				Manifests:   []*unstructured.Unstructured{},
				Synced:      false,
				ForDeletion: true,
			}
		},
	})

	stopCh := make(chan struct{})

	go circleInformer.Run(stopCh)

	// ticker := time.NewTicker(3 * time.Second)
	for {
		time.Sleep(2 * time.Second)
		for circleName, state := range circles {
			if state.Synced {
				continue
			}

			err := syncConfig.Do(circleName, state.Manifests, state.ForDeletion)
			if err != nil {
				log.Fatalln(err)
			}

			if state.ForDeletion {
				delete(circles, circleName)
			}
		}

	}
}
