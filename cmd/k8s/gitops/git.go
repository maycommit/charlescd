package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"
	gitutils "github.com/maycommit/circlerr/internal/k8s/controller/utils/git"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
)

func RemoteSync(repoUrl string) (string, apperror.Error) {
	r, cloneAndOpenRepositoryError := gitutils.CloneAndOpenRepository(repoUrl)
	if cloneAndOpenRepositoryError != nil {
		return "", cloneAndOpenRepositoryError
	}

	w, err := r.Worktree()
	if err != nil {
		return "", apperror.New("Remote sync failed", err.Error()).AddOperation("gitops.RemoteSync.Worktree")
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return "", apperror.New("Remote sync failed", err.Error()).AddOperation("gitops.RemoteSync.Pull")
	}

	h, err := r.ResolveRevision(plumbing.Revision("HEAD"))
	if err != nil {
		return "", nil
	}

	return h.String(), nil
}

func Parse(repoUrl, repoPath string) ([]*unstructured.Unstructured, apperror.Error) {
	var res []*unstructured.Unstructured
	gitDirOut := fmt.Sprintf("%s/%s", os.Getenv("GIT_DIR"), repoUrl)
	if err := filepath.Walk(filepath.Join(gitDirOut, repoPath), func(path string, infoRoot os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if (infoRoot.IsDir() && infoRoot.Name() == "circles") || (infoRoot.IsDir() && infoRoot.Name() == "projects") {
			if err := filepath.Walk(filepath.Join(gitDirOut, repoPath, infoRoot.Name()), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				if ext := filepath.Ext(path); ext != ".json" && ext != ".yml" && ext != ".yaml" {
					return nil
				}
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				var unstruct *unstructured.Unstructured
				d := kubeyaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 4096)
				err = d.Decode(&unstruct)
				if err != nil {
					return fmt.Errorf("failed to parse %s: %v", path, err)
				}

				res = append(res, unstruct)
				return nil

			}); err != nil {
				return err
			}
			return nil
		}

		return nil
	}); err != nil {
		return nil, apperror.New("Parse failed", err.Error()).AddOperation("gitops.Parse.Walk")
	}

	for _, r := range res {
		annotations := r.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations[annotation.ManageAnnotation] = r.GetCreationTimestamp().String()

		r.SetAnnotations(annotations)
	}

	return res, nil
}
