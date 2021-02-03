package cache

import (
	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type resourceInfo struct {
	circleMark  string
	projectMark string
	releaseMark string
	routerMark  string
}

func IsManagedResource(r *cache.Resource) bool {
	return r.Info.(*resourceInfo).circleMark != ""
}

func ResourceInfoHandler(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {

	var circleMark string
	var projectMark string
	var releaseMark string

	circleMark = un.GetAnnotations()[annotation.CircleAnnotation]
	projectMark = un.GetAnnotations()[annotation.ProjectAnnotation]
	releaseMark = un.GetAnnotations()[annotation.ReleaseAnnotation]

	info = &resourceInfo{
		projectMark: projectMark,
		circleMark:  circleMark,
		releaseMark: releaseMark,
	}

	cacheManifest = circleMark != "" && projectMark != "" && releaseMark != ""
	return
}
