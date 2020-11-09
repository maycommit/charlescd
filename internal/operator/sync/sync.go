package sync

import (
	"charlescd/internal/manager/project"
	"charlescd/internal/operator/repository"
	"charlescd/internal/utils/git"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

func Start(client dynamic.Interface, disco discovery.DiscoveryInterface, res *unstructured.Unstructured) error {
	projects, err := project.GetProjectsByResource(*res)
	if err != nil {
		return err
	}

	for _, project := range projects {
		_, err := git.CloneAndOpenRepository(project)
		if err != nil {
			return err
		}

		manifests, err := repository.ParseManifests(project)
		if err != nil {
			return err
		}

		for _, manifest := range manifests {
			err := unstructured.SetNestedField(manifest.Object, res.GetName(), "metadata", "labels", "circle")
			if err != nil {
				return err
			}

			gv := manifest.GroupVersionKind().GroupVersion()
			resources, err := disco.ServerResourcesForGroupVersion(gv.String())
			if err != nil {
				return err
			}

			var apiResource v1.APIResource
			for _, r := range resources.APIResources {
				if r.Kind == manifest.GetKind() {
					apiResource = r
					break
				}
			}

			gvr := gv.WithResource(apiResource.Name)
			_, err = client.Resource(gvr).Namespace("default").Create(context.TODO(), manifest, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}