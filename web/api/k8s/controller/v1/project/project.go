package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"
	"github.com/maycommit/circlerr/pkg/k8s/controller/apis/project/v1alpha1"
	clientset "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ProjectRoute struct {
	CircleName  string `json:"circleName"`
	CircleID    string `json:"circleId"`
	ReleaseName string `json:"releaseName"`
}

type Project struct {
	Name    string         `json:"name"`
	Managed bool           `json:"managed"`
	Routes  []ProjectRoute `json:"routes"`
	v1alpha1.ProjectSpec
}

func Create(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newProject Project
		err := json.NewDecoder(r.Body).Decode(&newProject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		projectRes := &v1alpha1.Project{
			ObjectMeta: metav1.ObjectMeta{
				Name: newProject.Name,
			},
			Spec: newProject.ProjectSpec,
		}

		_, err = client.ProjectV1alpha1().Projects("default").Create(context.TODO(), projectRes, metav1.CreateOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, newProject)
	}
}

func List(client *clientset.Clientset, appcache *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		projects := []Project{}
		res, err := client.ProjectV1alpha1().Projects(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		for _, item := range res.Items {

			newProject := Project{
				Name:        item.GetName(),
				Managed:     false,
				ProjectSpec: item.Spec,
				Routes:      []ProjectRoute{},
			}

			for circleName, route := range appcache.Projects.Get(item.GetName()).GetRoutes() {
				newProject.Routes = append(newProject.Routes, ProjectRoute{
					CircleName:  circleName,
					CircleID:    route.CircleID,
					ReleaseName: route.ReleaseName,
				})
			}

			annotations := item.GetAnnotations()

			if annotations != nil && annotations[annotation.ManageAnnotation] != "" {
				newProject.Managed = true
			}

			projects = append(projects, newProject)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(projects)
	}
}

func Get(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["projectID"]

		res, err := client.ProjectV1alpha1().Projects(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		newProject := Project{}
		newProject.Name = res.GetName()
		newProject.ProjectSpec = res.Spec

		annotations := res.GetAnnotations()

		if annotations != nil && annotations[annotation.ManageAnnotation] != "" {
			newProject.Managed = true
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newProject)
	}
}

func Update(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["projectID"]
		var newProject Project
		err := json.NewDecoder(r.Body).Decode(&newProject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		currentProject, err := client.ProjectV1alpha1().Projects("default").Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		currentProject.Spec.RepoURL = newProject.RepoURL
		currentProject.Spec.Path = newProject.Path
		currentProject.Spec.Token = newProject.Token
		currentProject.Spec.Template = newProject.Template

		_, err = client.ProjectV1alpha1().Projects("default").Update(context.TODO(), currentProject, metav1.UpdateOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, newProject)
	}
}

func Delete(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["projectID"]

		err := client.ProjectV1alpha1().Projects(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(nil)
	}
}
