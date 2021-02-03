package circle

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	appCache "github.com/maycommit/circlerr/internal/k8s/controller/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"
	"github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"
	clientset "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"

	cache "github.com/argoproj/gitops-engine/pkg/cache"

	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

type Circle struct {
	Name    string `json:"name"`
	Managed bool   `json:"managed"`
	v1alpha1.CircleSpec
}

type CircleListItem struct {
	Circle
	Status v1alpha1.CircleStatus `json:"status"`
}

type ResourceParent struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Controller bool   `json:"controller"`
}

type Resource struct {
	Ref     v1alpha1.ResourceStatus `json:"ref"`
	Parents []ResourceParent        `json:"parents"`
}

type Node struct {
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
}

type CircleTree struct {
	Nodes []Node `json:"projects"`
}

func Create(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var newCircle Circle
		err := json.NewDecoder(r.Body).Decode(&newCircle)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		circleRes := &v1alpha1.Circle{
			ObjectMeta: metav1.ObjectMeta{
				Name: newCircle.Name,
			},
			Spec: newCircle.CircleSpec,
		}

		_, err = client.CircleV1alpha1().Circles(newCircle.Destination.Namespace).Create(context.TODO(), circleRes, metav1.CreateOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCircle)
	}
}

func List(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		circles := []CircleListItem{}
		res, err := client.CircleV1alpha1().Circles(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		for _, item := range res.Items {

			newCircle := CircleListItem{
				Circle: Circle{
					Name:       item.GetName(),
					Managed:    false,
					CircleSpec: item.Spec,
				},
				Status: item.Status,
			}

			annotations := item.GetAnnotations()

			if annotations != nil && annotations[annotation.ManageAnnotation] != "" {
				newCircle.Managed = true
			}

			circles = append(circles, newCircle)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(circles)
	}
}

func Get(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]

		res, err := client.CircleV1alpha1().Circles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		newCircle := CircleListItem{}
		newCircle.Name = res.GetName()
		newCircle.CircleSpec = res.Spec
		newCircle.Status = res.Status

		annotations := res.GetAnnotations()

		if annotations != nil && annotations[annotation.ManageAnnotation] != "" {
			newCircle.Managed = true
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newCircle)
	}
}

func Update(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func Delete(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]

		err := client.CircleV1alpha1().Circles(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
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

func RemoveProject(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]
		projectName := mux.Vars(r)["projectName"]

		err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			res, err := client.CircleV1alpha1().Circles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			currentProject := []v1alpha1.CircleProject{}
			if res.Spec.Release != nil {
				for _, project := range res.Spec.Release.Projects {
					if project.Name != projectName {
						currentProject = append(currentProject, project)
					}
				}

				res.Spec.Release.Projects = currentProject
			}

			_, err = client.CircleV1alpha1().Circles(namespace).Update(context.TODO(), res, metav1.UpdateOptions{})
			return err
		})
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

func AddProject(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]

		var projects []v1alpha1.CircleProject
		err := json.NewDecoder(r.Body).Decode(&projects)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			res, err := client.CircleV1alpha1().Circles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			res.Spec.Release.Projects = append(res.Spec.Release.Projects, projects...)

			_, err = client.CircleV1alpha1().Circles(namespace).Update(context.TODO(), res, metav1.UpdateOptions{})
			return err
		})
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

func RemoveRelease(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]

		err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			res, err := client.CircleV1alpha1().Circles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			res.Spec.Release = nil

			_, err = client.CircleV1alpha1().Circles(namespace).Update(context.TODO(), res, metav1.UpdateOptions{})
			return err
		})
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

func AddRelease(client *clientset.Clientset) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]

		var newRelease *v1alpha1.CircleRelease
		err := json.NewDecoder(r.Body).Decode(&newRelease)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			res, err := client.CircleV1alpha1().Circles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			res.Spec.Release = newRelease

			_, err = client.CircleV1alpha1().Circles(namespace).Update(context.TODO(), res, metav1.UpdateOptions{})
			return err
		})
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(nil)
	}
}

func Tree(appCache *appCache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		name := mux.Vars(r)["circleID"]
		currentCircle := appCache.Circles.Get(name)
		if currentCircle == nil {
			json.NewEncoder(w).Encode(nil)
			return
		}
		projects := currentCircle.Status.Projects
		cl := appCache.Cluster.Get()

		circleTree := CircleTree{
			Nodes: []Node{},
		}

		for _, proj := range projects {
			resources := []Resource{}
			for _, res := range proj.Resources {
				resKey := kube.NewResourceKey(
					res.Group,
					res.Kind,
					namespace,
					res.Name,
				)

				cl.IterateHierarchy(resKey, func(resource *cache.Resource, namespaceResources map[kube.ResourceKey]*cache.Resource) {
					node := Resource{}

					node.Ref = v1alpha1.ResourceStatus{
						Kind:              resource.Ref.Kind,
						Name:              resource.Ref.Name,
						Group:             resource.ResourceKey().Group,
						CreationTimestamp: *resource.CreationTimestamp.DeepCopy(),
					}

					if resource.Resource != nil {

						node.Ref.Version = resource.Resource.GroupVersionKind().Version

						status, err := health.GetResourceHealth(resource.Resource, nil)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							json.NewEncoder(w).Encode(err)
							return
						}

						if status != nil {
							node.Ref.Health = &v1alpha1.ResourceHealth{
								Status:  status.Status,
								Message: status.Message,
							}
						}

					}

					for _, parent := range resource.OwnerRefs {

						newParent := ResourceParent{
							Kind: parent.Kind,
							Name: parent.Name,
						}

						if parent.Controller != nil {
							newParent.Controller = *parent.Controller
						}

						node.Parents = append(node.Parents, newParent)
					}

					resources = append(resources, node)
				})
			}

			circleTree.Nodes = append(circleTree.Nodes, Node{
				Name:      proj.Name,
				Resources: resources,
			})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(circleTree)
	}
}
