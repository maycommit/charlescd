package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/diff"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	gitopsEngineKube "github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/argoproj/gitops-engine/pkg/utils/tracing"
	"github.com/go-git/go-git/v5"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"
	gitutils "github.com/maycommit/circlerr/internal/k8s/controller/utils/git"
	"github.com/maycommit/circlerr/pkg/customerror"

	"github.com/maycommit/circlerr/internal/k8s/controller/utils/kube"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"k8s.io/klog/klogr"
)

func init() {
	config := flag.String("config", "", "Path to config repository list.")
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	masterUrl := flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	gitDir := flag.String("gitdir", "./tmp/git", "")
	flag.Parse()

	os.Setenv("CONFIG", *config)
	os.Setenv("GIT_DIR", *gitDir)
	os.Setenv("KUBECONFIG", *kubeconfig)
	os.Setenv("MASTER_URL", *masterUrl)
}

type Repository struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

type RepositoryCache struct {
	Revision  string
	Manifests []*unstructured.Unstructured
}

type Repositories struct {
	Repositories []Repository `yaml:"repositories"`
}

var repositoryCache = map[Repository]*RepositoryCache{}

func loadRepositoriesInCache() error {
	conf := &Repositories{}
	configData, err := ioutil.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configData, conf)
	if err != nil {
		return err
	}

	for _, r := range conf.Repositories {
		if _, ok := repositoryCache[r]; !ok {
			repositoryCache[r] = &RepositoryCache{
				Revision:  "",
				Manifests: []*unstructured.Unstructured{},
			}
		}
	}

	return nil
}

func parseManifests(repoURL, repoPath, revision string) ([]*unstructured.Unstructured, error) {
	var manifests []*unstructured.Unstructured
	gitDirOut := gitutils.GetOutDir(repoURL)

	for _, path := range []string{"circles", "projects"} {
		fmt.Println(filepath.Join(gitDirOut, repoPath, path))

		if err := filepath.Walk(filepath.Join(gitDirOut, repoPath, path), func(path string, info os.FileInfo, err error) error {
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

			manifests = append(manifests, items...)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	//for _, m := range manifests {
	//	annotations := m.GetAnnotations()
	//	if annotations == nil {
	//		annotations = make(map[string]string)
	//	}
	//
	//	annotations[annotation.ManageAnnotation] = revision
	//	m.SetAnnotations(annotations)
	//}

	return manifests, nil
}

func syncRepositories() error {
	for r, c := range repositoryCache {
		gitOptions := git.CloneOptions{
			URL: r.Url,
		}

		remoteRevision, err := gitutils.SyncRepository(gitOptions, c.Revision)
		if err != nil {
			return err
		}

		if len(c.Manifests) <= 0 || c.Revision != remoteRevision {
			c.Manifests, err = parseManifests(r.Url, r.Path, remoteRevision)
			if err != nil {
				return err
			}
		}

		repositoryCache[r].Revision = remoteRevision
	}

	return nil
}

type resourceInfo struct {
	manageMark string
}

func syncCluster(
	kubectl gitopsEngineKube.Kubectl,
	config *rest.Config,
	clusterCache cache.ClusterCache,
	manifests []*unstructured.Unstructured,
	namespace string) ([]common.ResourceSyncResult, error) {

	managedResources, err := clusterCache.GetManagedLiveObjs(manifests, func(r *cache.Resource) bool {
		return r.Info.(*resourceInfo).manageMark != ""
	})
	if err != nil {
		return nil, customerror.New("Sync cluster failed!", err, nil, "engine.SyncCluster.GetManagedLiveObjs")
	}

	result := sync.Reconcile(manifests, managedResources, namespace, clusterCache)
	diffRes, err := diff.DiffArray(result.Target, result.Live)
	if err != nil {
		return nil, customerror.New("Sync cluster failed!", err, nil, "engine.SyncCluster.GetManagedLiveObjs")
	}

	opts := []sync.SyncOpt{sync.WithSkipHooks(!diffRes.Modified), sync.WithOperationSettings(false, true, true, false)}
	syncCtx, err := sync.NewSyncContext("", result, config, config, kubectl, namespace, opts...)
	if err != nil {
		return nil, customerror.New("Sync cluster failed!", err, nil, "engine.SyncCluster.GetManagedLiveObjs")
	}

	syncCtx.Sync()
	phase, message, resources := syncCtx.GetState()
	if phase.Completed() {
		if phase == common.OperationError {
			err = fmt.Errorf("sync operation failed: %s", message)
		}

		if err != nil {
			return resources, customerror.New("Sync cluster failed!", err, nil, "engine.SyncCluster.GetManagedLiveObjs")
		}
	}

	return resources, nil
}

func main() {

	err := loadRepositoriesInCache()
	if err != nil {
		klog.Fatalf("Load repositories by config file failed: %s\n", err.Error())
	}

	config, err := kube.GetClusterConfig()
	if err != nil {
		klog.Fatalf("Get cluster config failed: %s\n", err.Error())
	}

	clusterCache := cache.NewClusterCache(
		config,
		cache.SetNamespaces([]string{}),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {
			manageMark := un.GetAnnotations()[annotation.ManageAnnotation]
			info = &resourceInfo{
				manageMark: manageMark,
			}

			cacheManifest = manageMark != ""
			return
		}),
	)

	err = clusterCache.EnsureSynced()
	if err != nil {
		klog.Fatalf("Cache ensure synced failed: %s\n", err.Error())
	}

	kubectlCmd := gitopsEngineKube.KubectlCmd{
		Log:    klogr.New(),
		Tracer: tracing.NopTracer{},
	}

	err = syncRepositories()
	if err != nil {
		klog.Fatalf("Initial sync repository failed: %s\n", err.Error())
	}

	for {
		for _, c := range repositoryCache {
			_, err = syncCluster(&kubectlCmd, config, clusterCache, c.Manifests, "")
			if err != nil {
				klog.Fatalf("Reconcile repositories failed: %s\n", err.Error())
			}

			err = syncRepositories()
			if err != nil {
				klog.Fatalf("Sync repository failed on loop: %s\n", err.Error())
			}
		}

		time.Sleep(2 * time.Second)
	}
}
