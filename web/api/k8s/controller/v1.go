package api

import (
	"fmt"

	"github.com/maycommit/circlerr/web/api/k8s/controller/v1/circle"
	"github.com/maycommit/circlerr/web/api/k8s/controller/v1/manifest"
	"github.com/maycommit/circlerr/web/api/k8s/controller/v1/project"
)

func (api *Api) newV1Api() {
	s := api.router.PathPrefix("/v1").Subrouter()
	{
		path := "/circles"
		s.HandleFunc(path, circle.Create(api.client)).Methods("POST")
		s.HandleFunc(path, circle.List(api.client)).Methods("GET")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}", path), circle.Get(api.client)).Methods("GET")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}/tree", path), circle.Tree(api.cache)).Methods("GET")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}", path), circle.Update(api.client)).Methods("PUT")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}", path), circle.Delete(api.client)).Methods("DELETE")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}/releases", path), circle.RemoveRelease(api.client)).Methods("DELETE")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}/releases", path), circle.AddRelease(api.client)).Methods("POST")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}/projects/{projectName}", path), circle.RemoveProject(api.client)).Methods("DELETE")
		s.HandleFunc(fmt.Sprintf("%s/{circleID}/projects", path), circle.AddProject(api.client)).Methods("POST")
	}
	{
		path := "/projects"
		s.HandleFunc(path, project.Create(api.client)).Methods("POST")
		s.HandleFunc(path, project.List(api.client, api.cache)).Methods("GET")
		s.HandleFunc(fmt.Sprintf("%s/{projectID}", path), project.Get(api.client)).Methods("GET")
		s.HandleFunc(fmt.Sprintf("%s/{projectID}", path), project.Update(api.client)).Methods("PUT")
		s.HandleFunc(fmt.Sprintf("%s/{projectID}", path), project.Delete(api.client)).Methods("DELETE")
	}
	{
		path := "/manifests"
		s.HandleFunc(fmt.Sprintf("%s/{manifestName}", path), manifest.Get(api.kubeClient, api.discoveryClient)).Methods("GET")
	}
}
