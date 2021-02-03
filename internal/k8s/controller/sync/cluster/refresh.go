package cluster

import (
	"context"
	"fmt"

	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"
	v1alpha12 "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (o *SyncOpts) getCircleTopLevelResources(circleName string) []map[kube.ResourceKey]*cache.Resource {
	topLevelResources := []map[kube.ResourceKey]*cache.Resource{}
	namespaceResources := o.appCache.Cluster.Get().GetNamespaceTopLevelResources(o.namespace)
	for resKey, resource := range namespaceResources {
		if resource.Resource != nil {
			circleAnnotation := resource.Resource.GetAnnotations()[annotation.CircleAnnotation]

			if circleAnnotation == circleName {
				topLevelResources = append(topLevelResources, map[kube.ResourceKey]*cache.Resource{
					resKey: resource,
				})
			}
		}
	}

	return topLevelResources
}

func (o *SyncOpts) getProjectsStatus(topLevelResources []map[kube.ResourceKey]*cache.Resource) ([]v1alpha12.ProjectStatus, apperror.Error) {
	projects := map[string][]v1alpha12.ResourceStatus{}
	projectStatus := []v1alpha12.ProjectStatus{}
	for _, resources := range topLevelResources {
		for _, resource := range resources {
			if resource.Resource != nil {
				projectName := resource.Resource.GetAnnotations()[annotation.ProjectAnnotation]
				resStatus := v1alpha12.ResourceStatus{
					Name:              resource.Resource.GetName(),
					Group:             resource.Resource.GroupVersionKind().Group,
					Version:           resource.Resource.GroupVersionKind().Version,
					Kind:              resource.Resource.GroupVersionKind().Kind,
					CreationTimestamp: resource.Resource.GetCreationTimestamp(),
				}

				status, err := health.GetResourceHealth(resource.Resource, nil)
				if err != nil {
					return nil, apperror.New("getProjectsStatus failed", err.Error()).AddOperation("cluster.getProjectsStatus.GetResourceHealth")
				}

				if status != nil {
					resStatus.Health = &v1alpha12.ResourceHealth{
						Status:  status.Status,
						Message: status.Message,
					}
				}

				projects[projectName] = append(projects[projectName], resStatus)
			}
		}
	}

	for projectName, resourcesStatus := range projects {
		status := health.HealthStatusHealthy
		for _, res := range resourcesStatus {
			if res.Health != nil && res.Health.Status == health.HealthStatusDegraded {
				status = res.Health.Status
				break
			}
		}

		projectStatus = append(projectStatus, v1alpha12.ProjectStatus{
			Name:      projectName,
			Status:    status,
			Resources: resourcesStatus,
		})
	}

	return projectStatus, nil
}

func (o *SyncOpts) updateCircle(resource *unstructured.Unstructured) apperror.Error {
	gvr := v1alpha12.SchemeGroupVersion.WithResource("circles")

	_, err := o.kubeClient.Resource(gvr).Namespace(o.namespace).Update(context.TODO(), resource, metav1.UpdateOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return apperror.New("updateCircle failed", err.Error()).AddOperation("cluster.updateCircle.Update")
	}

	return nil
}

func (o *SyncOpts) refreshCircles() apperror.Error {
	for circleName, circle := range o.appCache.Circles.List() {
		if circle.GetDeletion() {
			o.appCache.Circles.Delete(circleName)
			continue
		}

		newCircle := circle

		if circle.Status.Errors != nil && len(circle.Status.Errors) > 0 {
			return nil
		}

		if circle.GetRelease() != nil {
			topLevelResources := o.getCircleTopLevelResources(circleName)
			projectsStatus, err := o.getProjectsStatus(topLevelResources)
			if err != nil {
				return err
			}

			circleStatus := health.HealthStatusHealthy
			for _, proj := range projectsStatus {
				if proj.Status != health.HealthStatusHealthy {
					circleStatus = proj.Status
					break
				}
			}

			newCircle.Status = v1alpha12.CircleStatus{
				Status:   circleStatus,
				Projects: projectsStatus,
				Errors:   []string{},
			}
		} else {
			fmt.Println("----------------NEW CIRCLE---------------------")
			newCircle.Status = v1alpha12.CircleStatus{
				Errors: []string{},
			}
		}

		_, err := o.circlerrClient.CircleV1alpha1().Circles(o.namespace).Update(context.TODO(), &newCircle.Circle, metav1.UpdateOptions{})
		if err != nil && errors.IsNotFound(err) {
			return nil
		}

		if err != nil {
			return apperror.New("updateCircle failed", err.Error()).AddOperation("cluster.updateCircle.Update")
		}
	}

	return nil
}
