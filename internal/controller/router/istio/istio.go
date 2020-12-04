package istio

import (
	"charlescd/pkg/apis/circle/v1alpha1"
	"encoding/json"
	"fmt"

	"charlescd/internal/utils/kube"

	"github.com/argoproj/gitops-engine/pkg/cache"
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
	istioClient  istioclient.NetworkingV1beta1Interface
	clusterCache cache.ClusterCache
}

func NewIstioRouter(
	istioClient istioclient.NetworkingV1beta1Interface,
	clusterCache cache.ClusterCache,
) IstioUseCases {
	return &IstioRouter{
		istioClient:  istioClient,
		clusterCache: clusterCache,
	}
}

func (i *IstioRouter) getDestination() {

}

func (i *IstioRouter) newRoute(routerName, circleID, hostName string) *networkingv1beta1.HTTPRoute {
	return &networkingv1beta1.HTTPRoute{
		Name: routerName,
		Match: []*networkingv1beta1.HTTPMatchRequest{
			{
				Headers: map[string]*networkingv1beta1.StringMatch{
					"cookie": &networkingv1beta1.StringMatch{
						MatchType: &networkingv1beta1.StringMatch_Regex{
							Regex: fmt.Sprintf(".*x-circle-id=%s.*", circleID),
						},
					},
				},
			},
			{
				Headers: map[string]*networkingv1beta1.StringMatch{
					"x-circle-id": &networkingv1beta1.StringMatch{
						MatchType: &networkingv1beta1.StringMatch_Exact{
							Exact: circleID,
						},
					},
				},
			},
		},
		Route: []*networkingv1beta1.HTTPRouteDestination{
			{
				Destination: &networkingv1beta1.Destination{
					Host: hostName,
				},
			},
		},
	}
}

func (i *IstioRouter) getRoutes(routeMap map[string]*networkingv1beta1.HTTPRoute) []*networkingv1beta1.HTTPRoute {
	routes := []*networkingv1beta1.HTTPRoute{}

	for _, route := range routeMap {
		routes = append(routes, route)
	}

	return routes
}

func (i IstioRouter) getNewVirtualService(project v1alpha1.ProjectStatus, circleID, releaseName string) *v1beta1.VirtualService {
	routerName := fmt.Sprintf("%s-%s-router", project.Name, releaseName)
	routeMap := map[string]*networkingv1beta1.HTTPRoute{
		routerName: i.newRoute(routerName, circleID, fmt.Sprintf("%s-%s", project.Name, releaseName)),
	}

	fmt.Println("------------NEW ROUTE MAP -------------", routeMap)

	return &v1beta1.VirtualService{
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
			Hosts: []string{project.Name},
			Http:  i.getRoutes(routeMap),
		},
	}
}

func (i IstioRouter) getCurrentVirtualService(
	project v1alpha1.ProjectStatus,
	releaseName string,
	circleID string,
	vs *v1beta1.VirtualService,
) *v1beta1.VirtualService {

	routeMap := map[string]*networkingv1beta1.HTTPRoute{}
	for _, route := range vs.Spec.Http {
		routeMap[route.Name] = route
	}

	hostName := fmt.Sprintf("%s-%s", project.Name, releaseName)
	routerName := fmt.Sprintf("%s-router", hostName)
	routeMap[routerName] = i.newRoute(routerName, circleID, hostName)

	fmt.Println("------------CURRENT ROUTE MAP -------------", routeMap)

	vs.Spec.Http = i.getRoutes(routeMap)
	return vs
}

func (i *IstioRouter) GetVirtualServiceManifests(circle v1alpha1.Circle) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}

	for _, project := range circle.Status.Projects {

		var vsManifest *unstructured.Unstructured
		res := i.clusterCache.GetNamespaceTopLevelResources("default")
		for key, item := range res {
			if key.Kind == "VirtualService" && key.Name == project.Name {
				vsManifest = item.Resource
				break
			}

		}

		var currentVirtualService *v1beta1.VirtualService
		var releaseName = circle.Spec.Release.Name
		circleID := circle.ObjectMeta.UID
		if vsManifest != nil {
			b, err := vsManifest.MarshalJSON()
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(b, &currentVirtualService)
			if err != nil {
				return nil, err
			}

			currentVirtualService = i.getCurrentVirtualService(project, releaseName, string(circleID), currentVirtualService)
		} else {
			currentVirtualService = i.getNewVirtualService(project, string(circleID), releaseName)
		}

		b, err := json.Marshal(currentVirtualService)
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
