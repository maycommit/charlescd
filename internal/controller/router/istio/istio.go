package istio

import (
	"charlescd/pkg/apis/circle/v1alpha1"
	"encoding/json"

	"charlescd/internal/utils/kube"

	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	istioclient "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1"
)

type IstioUseCases interface {
	GetVirtualServiceManifests(circle v1alpha1.Circle) ([]*unstructured.Unstructured, error)
}

type IstioRouter struct {
	istioClient istioclient.NetworkingV1beta1Interface
}

func NewIstioRouter(istioClient istioclient.NetworkingV1beta1Interface) IstioUseCases {
	return &IstioRouter{
		istioClient: istioClient,
	}
}

func (i *IstioRouter) getRoutes(project v1alpha1.ProjectStatus) []*networkingv1beta1.HTTPRoute {
	return []*networkingv1beta1.HTTPRoute{
		{
			Route: []*networkingv1beta1.HTTPRouteDestination{
				{
					Destination: &networkingv1beta1.Destination{
						Host: "guestbook-ui",
					},
				},
			},
		},
	}
}

func (i IstioRouter) getNewVirtualService() {}

func (i *IstioRouter) GetVirtualServiceManifests(circle v1alpha1.Circle) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}

	for _, project := range circle.Status.Projects {
		vs := &v1beta1.VirtualService{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "networking.istio.io/v1beta1",
				Kind:       "VirtualService",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: project.Name,
				Annotations: map[string]string{
					"charlescd.io/router": project.Name,
				},
			},
			Spec: networkingv1beta1.VirtualService{
				Hosts: []string{},
				Http:  i.getRoutes(project),
			},
		}

		b, err := json.Marshal(vs)
		if err != nil {
			return nil, err
		}

		ms, err := kube.SplitJSON(b)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, ms...)
	}

	return manifests, nil
}
