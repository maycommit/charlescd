package main

import (
	v1 "charlescd/cmd/manager/api/v1"
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	circleclientset "charlescd/pkg/client/clientset/versioned"
	circlepb "charlescd/pkg/grpc/circle"
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

	client := circleclientset.NewForConfigOrDie(config)

	conn, err := grpc.Dial(":9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	grpcClient := circlepb.NewCircleServiceClient(conn)

	r := mux.NewRouter()
	{
		r.HandleFunc("/circles", v1.CircleCreate(client)).Methods("POST")
		r.HandleFunc("/circles/{name}/deploy", v1.CircleDeploy(client)).Methods("POST")
		r.HandleFunc("/circles", v1.CircleFindAll(client)).Methods("GET")
		r.HandleFunc("/circles/{name}", v1.CircleShow(client, grpcClient)).Methods("GET")
	}
	log.Println("Start manager on port 8080...")
	log.Println(http.ListenAndServe(":8080", r))
}
