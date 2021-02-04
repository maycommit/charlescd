package override

import (
	"fmt"
	"log"
	"strings"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func overrideName(manifest *unstructured.Unstructured, projectName, releaseName string) {
	manifest.SetName(fmt.Sprintf("%s-%s", projectName, releaseName))
}

func overrideAnnotations(manifest *unstructured.Unstructured, projectName, circleName, releaseName string) {
	annotations := manifest.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations[annotation.CircleAnnotation] = circleName
	annotations[annotation.ProjectAnnotation] = projectName
	annotations[annotation.ReleaseAnnotation] = releaseName

	manifest.SetAnnotations(annotations)

	_, ok, err := unstructured.NestedMap(manifest.Object, "spec", "template")
	if err != nil {
		log.Fatalln(err)
	}

	if ok {
		templateMetadataSpec, ok, err := unstructured.NestedStringMap(manifest.Object, "spec", "template", "metadata", "annotations")
		if err != nil {
			log.Fatalln(err)
		}

		templateAnnotations := map[string]string{}
		if ok {
			templateAnnotations = templateMetadataSpec
		}

		templateAnnotations[annotation.CircleAnnotation] = circleName
		templateAnnotations[annotation.ProjectAnnotation] = projectName
		templateAnnotations[annotation.ReleaseAnnotation] = releaseName

		err = unstructured.SetNestedStringMap(manifest.Object, templateAnnotations, "spec", "template", "metadata", "annotations")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func overrideImage(manifest *unstructured.Unstructured, image string) error {
	b, err := manifest.MarshalJSON()
	if err != nil {
		return err
	}

	newManifest := strings.ReplaceAll(string(b), "{{ circlerr.image }}", image)

	err = manifest.UnmarshalJSON([]byte(newManifest))
	return err
}

func Do(manifests []*unstructured.Unstructured, projectName string, circle *circle.CircleCache) ([]*unstructured.Unstructured, error) {
	for _, m := range manifests {
		overrideName(m, projectName, circle.Circle().Spec.Release.Name)
		overrideAnnotations(m, projectName, circle.Circle().Name, circle.Circle().Spec.Release.Name)
		//err := overrideImage(m, circle.GetProject(projectName).Image)
		//if err != nil {
		//	return nil, err
		//}
	}

	return manifests, nil
}
