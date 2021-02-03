package main

import (
	"flag"
	"os"
	"time"

	"github.com/maycommit/circlerr/internal/k8s/controller/utils/kube"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var repoUrl = "https://github.com/circlerr/manage-controller-example"

func init() {
	fGitDir := flag.String("gitdir", "./tmp/git", "")
	fKubeconfigPath := flag.String("kubepath", os.Getenv("KUBECONFIG_PATH"), "")
	fK8sConnType := flag.String("k8sconntype", os.Getenv("K8S_CONN_TYPE"), "")
	flag.Parse()

	os.Setenv("GIT_DIR", *fGitDir)
	os.Setenv("KUBECONFIG_PATH", *fKubeconfigPath)
	os.Setenv("K8S_CONN_TYPE", *fK8sConnType)
}

func main() {
	revision := ""
	manifests := []*unstructured.Unstructured{}
	config, err := kube.GetClusterConfig()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(5 * time.Second)
	syncOpts := NewSyncOpts(config, "default")

	for {
		select {
		case <-ticker.C:
			currentRevision, err := RemoteSync(repoUrl)
			if err != nil {
				err.AddOperation("gitops.main.RemoteSync").LogFields()
				return
			}

			if currentRevision != revision {
				manifests, err = Parse(repoUrl, "")
				if err != nil {
					err.AddOperation("gitops.main.Parse").LogFields()
					return
				}

				revision = currentRevision
			}

			err = syncOpts.Sync(manifests)
			if err != nil {
				err.AddOperation("gitops.main.Sync").LogFields()
				return
			}
		}
	}
}
