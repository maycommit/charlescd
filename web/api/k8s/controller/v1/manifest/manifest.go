package manifest

import (
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

func Get(kubeClient dynamic.Interface, discoveryClient *discovery.DiscoveryClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		group := r.URL.Query().Get("group")
		version := r.URL.Query().Get("version")
		kind := r.URL.Query().Get("kind")
		name := mux.Vars(r)["manifestName"]

		resource := ""
		gv := schema.GroupVersion{Group: group, Version: version}
		res, err := discoveryClient.ServerResourcesForGroupVersion(gv.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		for _, i := range res.APIResources {
			if i.Kind == kind {
				resource = i.Name
				break
			}
		}

		s := schema.GroupVersionResource{group, version, resource}
		manifest, err := kubeClient.Resource(s).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}


		//b, err := manifest.MarshalJSON()
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	json.NewEncoder(w).Encode(err)
		//	return
		//}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(manifest)
	}
}
