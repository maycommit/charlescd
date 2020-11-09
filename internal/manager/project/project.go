package project

import (
	"context"
	"encoding/json"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Project struct {
	Name       string   `json:"name"`
	Tag        string   `json:"tag"`
	Repository string   `json:"repository"`
	Paths      []string `json:"paths"`
	Token      string   `json:"token"`
}

var Resource = schema.GroupVersionResource{
	Group:    "charlecd.io",
	Version:  "v1",
	Resource: "projects",
}

func GetProjectsByResource(res unstructured.Unstructured) ([]Project, error) {
	specProjects, ok, err := unstructured.NestedSlice(res.Object, "spec", "projects")
	if err != nil {
		return nil, err
	}

	var projects []Project
	if ok {
		for i, project := range specProjects {
			projects = append(projects, Project{
				Name:       project.(map[string]interface{})["name"].(string),
				Tag:        project.(map[string]interface{})["tag"].(string),
				Token:      project.(map[string]interface{})["token"].(string),
				Repository: project.(map[string]interface{})["repository"].(string),
			})

			for _, path := range project.(map[string]interface{})["paths"].([]interface{}) {
				projects[i].Paths = append(projects[i].Paths, path.(string))
			}
		}
	}

	return projects, nil
}

func CreateProject(client dynamic.Interface, data io.ReadCloser) error {
	newProject := Project{}
	err := json.NewDecoder(data).Decode(&newProject)
	if err != nil {
		return err
	}

	projectObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "charlecd.io/v1",
			"kind":       "Project",
			"metadata": map[string]interface{}{
				"name": newProject.Name,
			},
			"spec": map[string]interface{}{
				"repository": newProject.Repository,
				"paths": newProject.Paths,
				"token": newProject.Token,
			},
		},
	}

	_, err = client.Resource(Resource).Namespace("default").Create(context.TODO(), projectObject, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}