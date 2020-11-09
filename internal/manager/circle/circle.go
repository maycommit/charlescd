package circle

import (
	"charlescd/internal/manager/project"
	"context"
	"encoding/json"
	"io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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


type Circle struct {
	Name         string        `json:"name"`
	Segments     []Segment     `json:"segments"`
	Environments []Environment `json:"environments"`
	Projects     []project.Project     `json:"projects"`
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
				"projects":     newCircle.Projects,
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



func ListCircles(client dynamic.Interface) ([]Circle, error) {
	list, err := client.Resource(Resource).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var circles []Circle
	for _, res := range list.Items {
		segments, err := getSegmentsByResource(res)
		if err != nil {
			return nil, err
		}

		environments, err := getEnvironmentsByResource(res)
		if err != nil {
			return nil, err
		}

		projects, err := project.GetProjectsByResource(res)
		if err != nil {
			return nil, err
		}

		circles = append(circles, Circle{
			Name:         res.GetName(),
			Segments:     segments,
			Environments: environments,
			Projects: projects,
		})
	}

	return circles, nil
}
