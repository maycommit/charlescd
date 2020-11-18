package repository

import (
	"bytes"
	"charlescd/internal/env"
	"charlescd/internal/manager/circle"
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

func SplitYAML(yamlData []byte) ([]*unstructured.Unstructured, error) {
	// Similar way to what kubectl does
	// https://github.com/kubernetes/cli-runtime/blob/master/pkg/resource/visitor.go#L573-L600
	// Ideally k8s.io/cli-runtime/pkg/resource.Builder should be used instead of this method.
	// E.g. Builder does list unpacking and flattening and this code does not.
	d := kubeyaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlData), 4096)
	var objs []*unstructured.Unstructured
	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			return objs, fmt.Errorf("failed to unmarshal manifest: %v", err)
		}
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		u := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(ext.Raw, u); err != nil {
			return objs, fmt.Errorf("failed to unmarshal manifest: %v", err)
		}
		objs = append(objs, u)
	}
	return objs, nil
}

func ParseManifests(circleResource circle.CircleResource) ([]*unstructured.Unstructured, error) {
	var res []*unstructured.Unstructured
	project := circleResource.Project
	gitDirOut := fmt.Sprintf("%s/%s", env.Get("GIT_DIR"), project.Name)
	for i := range project.Paths {
		if err := filepath.Walk(filepath.Join(gitDirOut, project.Paths[i]), func(path string, info os.FileInfo, err error) error {
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
			items, err := SplitYAML(data)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %v", path, err)
			}

			fmt.Println("READ: ", items)
			res = append(res, items...)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return res, nil
}
