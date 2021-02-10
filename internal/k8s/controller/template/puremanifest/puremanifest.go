package puremanifest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache/project"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/template/override"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/git"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PureManifest struct {
	projectName string
	project     *project.ProjectCache
}

func NewPureManifest(projectName string, project *project.ProjectCache) PureManifest {
	return PureManifest{projectName: projectName, project: project}
}

func (p PureManifest) ParseManifests(circleName string, circle *circle.CircleCache) ([]*unstructured.Unstructured, apperror.Error) {
	var res []*unstructured.Unstructured
	gitDirOut := git.GetOutDir(p.project.RepoURL)
	if err := filepath.Walk(filepath.Join(gitDirOut, p.project.Path), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(info.Name()); ext != ".json" && ext != ".yml" && ext != ".yaml" {
			return nil
		}
		data, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}
		items, err := kube.SplitYAML(data)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %v", path, err)
		}

		res = append(res, items...)
		return nil
	}); err != nil {
		return nil, apperror.New("Parse manifests failed", err.Error()).AddOperation("puremanifest.ParseManifests.Walk")
	}

	manifests, err := override.Do(res, p.projectName, circle)
	if err != nil {
		return nil, apperror.New("Parse manifests failed", err.Error()).AddOperation("puremanifest.ParseManifests.Do")
	}

	return manifests, nil
}
