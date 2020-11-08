package main

import (
	v1 "charlescd/cmd/manager/api/v1"
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"k8s.io/client-go/dynamic"
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

	r := mux.NewRouter()
	{
		r.HandleFunc("/circles", v1.CircleCreate(client)).Methods("POST")
		r.HandleFunc("/circles", v1.CircleFindAll(client)).Methods("GET")
	}

	log.Println(http.ListenAndServe(":8080", r))
}
