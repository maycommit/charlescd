package main

import (
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/maycommit/circlerr/internal/k8s/controller/utils/kube"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var repoUrl = "https://github.com/maycommit/circlerr"

func init() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	masterUrl := flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	gitDir := flag.String("gitdir", "./tmp/git", "")
	flag.Parse()

	os.Setenv("GIT_DIR", *gitDir)
	os.Setenv("KUBECONFIG", *kubeconfig)
	os.Setenv("MASTER_URL", *masterUrl)
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
				logrus.Fatalln(err.AddOperation("gitops.main.RemoteSync").LogFields())
				return
			}

			if currentRevision != revision {
				manifests, err = Parse(repoUrl, "examples/manage")
				if err != nil {
					logrus.Fatalln(err.AddOperation("gitops.main.Parse").LogFields())
					return
				}

				revision = currentRevision
			}

			err = syncOpts.Sync(manifests)
			if err != nil {
				logrus.Fatalln(err.AddOperation("gitops.main.Sync").LogFields())
				return
			}
		}
	}
}
