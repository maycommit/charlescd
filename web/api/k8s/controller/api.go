package api

import (
	"log"
	"net/http"
	"time"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	"github.com/gorilla/mux"

	clientset "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"
)

type Api struct {
	// Dependencies
	cache           *cache.Cache
	client          *clientset.Clientset
	kubeClient      dynamic.Interface
	discoveryClient *discovery.DiscoveryClient

	//Server
	router *mux.Router
	server *http.Server
}

func NewApi(
	cache *cache.Cache,
	client *clientset.Clientset,
	kubeClient dynamic.Interface,
	discoveryClient *discovery.DiscoveryClient,
) Api {

	api := Api{
		cache:           cache,
		client:          client,
		kubeClient:      kubeClient,
		discoveryClient: discoveryClient,
		router:          mux.NewRouter().PathPrefix("/api").Subrouter(),
	}
	api.server = &http.Server{
		Handler: api.router,
		Addr:    ":8081",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	api.router.Use(LoggingMiddleware)
	api.router.Use(ValidatorMiddleware)
	api.health()
	api.newV1Api()

	return api
}

func (api *Api) health() {
	api.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(":)"))
		return
	})
}

func (api Api) Start() {
	log.Fatal(api.server.ListenAndServe())
}
