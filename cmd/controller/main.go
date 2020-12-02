package main

import (
	"charlescd/internal/controller/sync"
	"context"
	"flag"
	"log"
	"net"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2/klogr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"charlescd/cmd/controller/server"
	circleclientset "charlescd/pkg/client/clientset/versioned"
	circlepb "charlescd/pkg/grpc/circle"
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
