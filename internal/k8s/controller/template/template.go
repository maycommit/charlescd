package template

import (
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/template/puremanifest"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type UseCases interface {
	ParseManifests(circleName string, circle *circle.CircleCache) ([]*unstructured.Unstructured, apperror.Error)
}

func NewTemplate(projectName string, project *project.ProjectCache) UseCases {
	return puremanifest.NewPureManifest(projectName, project)
}
