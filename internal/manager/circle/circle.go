package circle

import (
	"charlescd/pkg/apis/circle/v1alpha1"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"

	circleclientset "charlescd/pkg/client/clientset/versioned"
	circlepb "charlescd/pkg/grpc/circle"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

type Circle struct {
	Name         string                     `json:"name"`
	Release      v1alpha1.CircleRelease     `json:"release"`
	Destination  v1alpha1.CircleDestination `json:"destination"`
	Segments     []v1alpha1.Segment         `json:"segments"`
	Environments []v1alpha1.Environment     `json:"environments"`
}

var Resource = schema.GroupVersionResource{
	Group:    "charlecd.io",
	Version:  "v1",
	Resource: "circles",
}

func CreateCircle(client circleclientset.Interface, data io.ReadCloser) error {
	newCircle := Circle{}
	err := json.NewDecoder(data).Decode(&newCircle)
	if err != nil {
		return err
	}

	circleObject := &v1alpha1.Circle{
		ObjectMeta: metav1.ObjectMeta{
			Name: newCircle.Name,
		},
		Spec: v1alpha1.CircleSpec{
			Release:      newCircle.Release,
			Destination:  newCircle.Destination,
			Segments:     newCircle.Segments,
			Environments: newCircle.Environments,
		},
	}

	_, err = client.CircleV1alpha1().Circles("default").Create(context.TODO(), circleObject, metav1.CreateOptions{})
	if k8sErrors.IsAlreadyExists(err) {

		return retry.RetryOnConflict(retry.DefaultRetry, func() error {
			result, err := client.CircleV1alpha1().Circles("default").Get(context.TODO(), newCircle.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}

			result.Spec = circleObject.Spec
			_, err = client.CircleV1alpha1().Circles("default").Update(context.TODO(), result, metav1.UpdateOptions{})
			return err
		})
	}

	return nil
}

func ListCircles(client circleclientset.Interface) ([]Circle, error) {
	circles := []Circle{}
	list, err := client.CircleV1alpha1().Circles("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, circle := range list.Items {
		circles = append(circles, Circle{
			Name:         circle.GetName(),
			Release:      circle.Spec.Release,
			Destination:  circle.Spec.Destination,
			Segments:     circle.Spec.Segments,
			Environments: circle.Spec.Environments,
		})
	}

	return circles, nil
}

func GetCircle(client circleclientset.Interface, name string) (Circle, error) {
	circle, err := client.CircleV1alpha1().Circles("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return Circle{}, err
	}

	return Circle{
		Name:         circle.GetName(),
		Release:      circle.Spec.Release,
		Destination:  circle.Spec.Destination,
		Environments: circle.Spec.Environments,
		Segments:     circle.Spec.Segments,
	}, nil
}

func GetCircleTree(client circleclientset.Interface, grpcClient circlepb.CircleServiceClient, name string) (json.RawMessage, error) {
	tree, err := grpcClient.CircleTree(context.Background(), &circlepb.Circle{
		Name:      name,
		Namespace: "default",
	})
	if err != nil {
		return json.RawMessage{}, err
	}

	b, err := json.Marshal(tree)
	if err != nil {
		return json.RawMessage{}, err
	}

	return b, nil
}

func Deploy(client circleclientset.Interface, name string, data io.ReadCloser) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return nil
	})
}
