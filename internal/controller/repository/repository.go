package repository

import (
	"charlescd/internal/env"
	"charlescd/pkg/apis/circle/v1alpha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ParseManifests(project v1alpha1.CircleProject) ([]*unstructured.Unstructured, error) {
	var res []*unstructured.Unstructured
	gitDirOut := fmt.Sprintf("%s/%s", env.Get("GIT_DIR"), project.Name)
	if err := filepath.Walk(filepath.Join(gitDirOut, project.Path), func(path string, info os.FileInfo, err error) error {
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
		return nil, err
	}

	return res, nil
}
