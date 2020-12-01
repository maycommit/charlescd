package sync

import (
	"charlescd/internal/controller/repository"
	"charlescd/internal/utils"
	"charlescd/pkg/apis/circle/v1alpha1"
	"context"
	"fmt"
	"os"

	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/go-git/go-git/v5"

	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/go-logr/logr"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"

	"github.com/argoproj/gitops-engine/pkg/cache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	circleInterface "charlescd/pkg/client/clientset/versioned/typed/circle/v1alpha1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
)

const (
	circleAnnotation  = "charlescd.io/circle"
	projectAnnotation = "charlescd.io/circle"
)

type resourceInfo struct {
	gcMark string
}

type SyncResource struct {
}

type SyncConfig struct {
	ClusterCache cache.ClusterCache
	Config       *rest.Config
	KubeClient   dynamic.Interface
	Client       circleInterface.CircleInterface
	Disco        discovery.DiscoveryInterface
	Circle       *v1alpha1.Circle
	Namespace    string
	Prune        bool
	Log          logr.Logger
}

func ClusterCache(config *rest.Config, namespaces []string, log logr.Logger) cache.ClusterCache {
	clusterCache := cache.NewClusterCache(config,
		cache.SetNamespaces(namespaces),
		cache.SetLogr(log),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {
			// store gc mark of every resource
			gcMark := un.GetAnnotations()[circleAnnotation]
			info = &resourceInfo{gcMark: un.GetAnnotations()[circleAnnotation]}
			// cache resources that has that mark to improve performance
			cacheManifest = gcMark != ""
			return
		}),
	)

	return clusterCache
}

func cloneAndOpenRepository(project v1alpha1.CircleProject) (*git.Repository, error) {
	os.Setenv("GIT_DIR", "./tmp/git")

	gitDirOut := fmt.Sprintf("%s/%s", os.Getenv("GIT_DIR"), project.Name)

	r, err := git.PlainClone(gitDirOut, false, &git.CloneOptions{
		URL:      project.RepoURL,
		Progress: os.Stdout,
	})
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return nil, err
	}

	r, err = git.PlainOpen(gitDirOut)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Start(syncConfig SyncConfig) error {

	release := syncConfig.Circle.Spec.Release
	gitOpsEngine := engine.NewEngine(syncConfig.Config, syncConfig.ClusterCache, engine.WithLogr(syncConfig.Log))
	_, err := gitOpsEngine.Run()
	if err != nil {
		return err
	}

	manifests := []*unstructured.Unstructured{}
	for _, project := range release.Projects {
		r, err := cloneAndOpenRepository(project)
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}

		manifests, err = repository.ParseManifests(project)
		if err != nil {
			return err
		}

		for _, manifest := range manifests {

			annotations := manifest.GetAnnotations()
			if annotations == nil {
				annotations = make(map[string]string)
			}
			annotations[circleAnnotation] = syncConfig.Circle.GetName()
			annotations[projectAnnotation] = project.Name
			manifest.SetAnnotations(annotations)
		}
	}

	result, err := gitOpsEngine.Sync(context.Background(), manifests, func(r *cache.Resource) bool {
		return r.Info.(*resourceInfo).gcMark == syncConfig.Circle.GetName()
	}, utils.NewSHA1Hash(), syncConfig.Namespace, sync.WithPrune(syncConfig.Prune), sync.WithLogr(syncConfig.Log))
	if err != nil {
		syncConfig.Log.Error(err, "Failed to synchronize cluster state")
		return err
	}

	projectMap := map[string][]v1alpha1.ResourceStatus{}
	for _, res := range result {

		mapResourceKey := syncConfig.ClusterCache.GetNamespaceTopLevelResources(syncConfig.Namespace)
		resKey := kube.NewResourceKey(
			res.ResourceKey.Group,
			res.ResourceKey.Kind,
			res.ResourceKey.Namespace,
			res.ResourceKey.Name,
		)

		status, err := health.GetResourceHealth(mapResourceKey[resKey].Resource, nil)
		if err != nil {
			return err
		}

		resourceStatus := v1alpha1.ResourceStatus{
			Kind:  res.ResourceKey.Kind,
			Group: res.ResourceKey.Group,
			Name:  res.ResourceKey.Name,
			Health: &v1alpha1.ResourceHealth{
				Status:  status.Status,
				Message: status.Message,
			},
		}

		resource := mapResourceKey[resKey].Resource
		projectName := resource.GetAnnotations()[projectAnnotation]
		projectMap[projectName] = append(projectMap[projectName], resourceStatus)
	}

	projectsStatus := []v1alpha1.ProjectStatus{}
	for projectName, resources := range projectMap {
		projectsStatus = append(projectsStatus, v1alpha1.ProjectStatus{
			Name:      projectName,
			Resources: resources,
		})
	}

	if err := syncConfig.updateCircle(projectsStatus); err != nil {
		return err
	}

	return nil
}

func (syncConfig SyncConfig) updateCircle(projectStatusList []v1alpha1.ProjectStatus) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		name := syncConfig.Circle.GetName()
		result, err := syncConfig.Client.Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		result.Status = v1alpha1.CircleStatus{
			Projects: projectStatusList,
		}

		_, err = syncConfig.Client.Update(context.TODO(), result, metav1.UpdateOptions{})
		return err
	})
}
