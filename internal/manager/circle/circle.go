package circle

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"io"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/retry"
)

const (
	StatusProcessing = "PROCESSING"
)

type Segment struct {
	Key       string `json:"key"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type Environment struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Project struct {
	Name       string   `json:"name"`
	Tag        string   `json:"tag"`
	Repository string   `json:"repository"`
	Paths      []string `json:"paths"`
	Token      string   `json:"token"`
}

func (p Project) GetGCMark(key kube.ResourceKey) string {
	h := sha256.New()
	_, _ = h.Write([]byte(fmt.Sprintf("%s/%s", p.Repository, strings.Join(p.Paths, ","))))
	_, _ = h.Write([]byte(strings.Join([]string{key.Group, key.Kind, key.Name}, "/")))
	return "sha256." + base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

type CircleResource struct {
	ReleaseName string `json:"releaseName"`
	Project Project `json:"project"`
	Tag         string `json:"tag"`
	Status      string `json:"status"`
	Error       string `json:"error"`
}

type Circle struct {
	Name         string           `json:"name"`
	Segments     []Segment        `json:"segments"`
	Environments []Environment    `json:"environments"`
	Resources    []CircleResource `json:"resources"`
}

var Resource = schema.GroupVersionResource{
	Group:    "charlecd.io",
	Version:  "v1",
	Resource: "circles",
}

func CreateCircle(client dynamic.Interface, data io.ReadCloser) error {
	newCircle := Circle{}
	err := json.NewDecoder(data).Decode(&newCircle)
	if err != nil {
		return err
	}

	circleObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "charlecd.io/v1",
			"kind":       "Circle",
			"metadata": map[string]interface{}{
				"name": newCircle.Name,
			},
			"spec": map[string]interface{}{
				"segments":     newCircle.Segments,
				"environments": newCircle.Environments,
				"resources":    newCircle.Resources,
			},
		},
	}

	_, err = client.Resource(Resource).Namespace("default").Create(context.TODO(), circleObject, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func getSegmentsByResource(res unstructured.Unstructured) ([]Segment, error) {
	specSegments, ok, err := unstructured.NestedSlice(res.Object, "spec", "segments")
	if err != nil {
		return nil, err
	}

	var segments []Segment
	if ok {
		for _, segment := range specSegments {
			segments = append(segments, Segment{
				Key:       segment.(map[string]interface{})["key"].(string),
				Condition: segment.(map[string]interface{})["condition"].(string),
				Value:     segment.(map[string]interface{})["value"].(string),
			})
		}
	}

	return segments, nil
}

func getEnvironmentsByResource(res unstructured.Unstructured) ([]Environment, error) {
	specEnvs, ok, err := unstructured.NestedSlice(res.Object, "spec", "environments")
	if err != nil {
		return nil, err
	}

	var environments []Environment
	if ok {
		for _, environment := range specEnvs {
			environments = append(environments, Environment{
				Key:   environment.(map[string]interface{})["key"].(string),
				Value: environment.(map[string]interface{})["value"].(string),
			})
		}
	}

	return environments, nil
}

func GetProjectByResource(res interface{}) (Project, error) {
	project := Project{
		Name:       res.(map[string]interface{})["name"].(string),
		Token:      res.(map[string]interface{})["token"].(string),
		Repository: res.(map[string]interface{})["repository"].(string),
	}

	for _, path := range res.(map[string]interface{})["paths"].([]interface{}) {
		project.Paths = append(project.Paths, path.(string))
	}

	return project, nil
}


func GetResourcesByResource(res unstructured.Unstructured) ([]CircleResource, error) {
	specEnvs, ok, err := unstructured.NestedSlice(res.Object, "spec", "resources")
	if err != nil {
		return nil, err
	}

	resources := []CircleResource{}
	if ok {
		for _, resource := range specEnvs {
			project, err := GetProjectByResource(resource.(map[string]interface{})["project"])
			if err != nil {
				return nil, err
			}

			resources = append(resources, CircleResource{
				ReleaseName: resource.(map[string]interface{})["releaseName"].(string),
				Project: project,
				Tag:         resource.(map[string]interface{})["tag"].(string),
			})
		}
	}

	return resources, nil
}

func ListCircles(client dynamic.Interface) ([]Circle, error) {
	list, err := client.Resource(Resource).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	circles := []Circle{}
	for _, res := range list.Items {
		segments, err := getSegmentsByResource(res)
		if err != nil {
			return nil, err
		}

		environments, err := getEnvironmentsByResource(res)
		if err != nil {
			return nil, err
		}

		resources, err := GetResourcesByResource(res)
		if err != nil {
			return nil, err
		}

		circles = append(circles, Circle{
			Name:         res.GetName(),
			Segments:     segments,
			Environments: environments,
			Resources:    resources,
		})
	}

	return circles, nil
}

func Deploy(client dynamic.Interface, name string, data io.ReadCloser) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {

		fmt.Println("NAME", name)
		result, getErr := client.Resource(Resource).Namespace("default").Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			return errors.New("failed get resource in cluster: " + getErr.Error())
		}

		resources, err := GetResourcesByResource(*result)
		if err != nil {
			return errors.New("failed getResourcesByResource: " + err.Error())
		}

		newResource := CircleResource{}
		err = json.NewDecoder(data).Decode(&newResource)
		if err != nil {
			return errors.New("fail to decode for json: " + err.Error())
		}

		newResource.Status = StatusProcessing

		resources = append(resources, newResource)

		nestedResources := []interface{}{}

		for _, r := range resources {
			nestedResources = append(nestedResources, r)
		}

		if err := unstructured.SetNestedSlice(result.Object, nestedResources, "spec", "resources"); err != nil {
			return errors.New("failed set nested field: " + err.Error())
		}

		_, updateErr := client.Resource(Resource).Namespace("default").Update(context.TODO(), result, metav1.UpdateOptions{})
		return errors.New("failed update: " + updateErr.Error())
	})
}
