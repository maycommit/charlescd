package istio

import (
	"context"
	"fmt"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"

	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (i *IstioRouter) newRoute(circleID, hostName string) *networkingv1beta1.HTTPRoute {
	return &networkingv1beta1.HTTPRoute{
		Name: fmt.Sprintf("%s-router", hostName),
		Match: []*networkingv1beta1.HTTPMatchRequest{
			{
				Headers: map[string]*networkingv1beta1.StringMatch{
					"cookie": {
						MatchType: &networkingv1beta1.StringMatch_Regex{
							Regex: fmt.Sprintf(".*x-circle-id=%s.*", circleID),
						},
					},
				},
			},
			{
				Headers: map[string]*networkingv1beta1.StringMatch{
					"x-circle-id": {
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

func (i *IstioRouter) getRoutes(routeMap map[string]project.ProjectRoute) []*networkingv1beta1.HTTPRoute {
	routes := []*networkingv1beta1.HTTPRoute{}

	for _, projectRoute := range routeMap {
		routes = append(routes, i.newRoute(projectRoute.CircleID, projectRoute.ReleaseName))
	}

	return routes
}

func (i IstioRouter) getNewVirtualService(projectName string, routeMap map[string]project.ProjectRoute) *v1beta1.VirtualService {
	return &v1beta1.VirtualService{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "VirtualService",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: projectName,
			Annotations: map[string]string{
				annotation.RouterAnnotation: projectName,
			},
		},
		Spec: networkingv1beta1.VirtualService{
			Hosts: []string{projectName},
			Http:  i.getRoutes(routeMap),
		},
	}
}

func (i IstioRouter) manageVirtualServices(projectName string, projectCache *project.ProjectCache) apperror.Error {

	newVirtualService := i.getNewVirtualService(projectName, projectCache.GetRoutes())

	if projectCache.GetRoutes() == nil {
		err := i.virtualService.Delete(context.TODO(), projectName, metav1.DeleteOptions{})
		if err != nil && errors.IsNotFound(err) {
			return nil
		}

		if err != nil {
			return apperror.New("manageVirtualServices failed", err.Error()).AddOperation("istio.manageVirtualServices.Delete")
		}

		return nil
	}

	_, err := i.virtualService.Create(context.TODO(), newVirtualService, metav1.CreateOptions{})
	if err != nil && errors.IsAlreadyExists(err) {
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			currentVS, err := i.virtualService.Get(context.TODO(), projectName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			currentVS.Spec = newVirtualService.Spec
			_, err = i.virtualService.Update(context.TODO(), currentVS, metav1.UpdateOptions{})
			return err
		})

		if err != nil {
			return apperror.New("manageVirtualServices failed", err.Error()).AddOperation("istio.manageVirtualServices.RetryOnConflict")
		}
	}

	if err != nil {
		return apperror.New("manageVirtualServices failed", err.Error()).AddOperation("istio.manageVirtualServices.Create")
	}

	return nil
}
