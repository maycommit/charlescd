package sync

import (
	"charlescd/internal/controller/repository"
	"charlescd/pkg/apis/circle/v1alpha1"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
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
	"k8s.io/apimachinery/pkg/runtime/schema"

	circleInterface "charlescd/pkg/client/clientset/versioned/typed/circle/v1alpha1"

	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	istioclient "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
)

const (
	circleAnnotation  = "charlescd.io/circle"
	projectAnnotation = "charlescd.io/project"
)

type resourceInfo struct {
	circleMark  string
	projectMark string
}

type SyncResource struct {
}

type SyncConfig struct {
	Ctx          context.Context
	GitopsEngine engine.GitOpsEngine
	ClusterCache cache.ClusterCache
	Config       *rest.Config
	KubeClient   dynamic.Interface
	IstioClient  istioclient.NetworkingV1beta1Interface
	Client       circleInterface.CircleInterface
	Disco        discovery.DiscoveryInterface
	Namespace    string
	Prune        bool
	Log          logr.Logger
	StopEngine   engine.StopFunc
}

func ClusterCache(config *rest.Config, namespaces []string, log logr.Logger) cache.ClusterCache {
	clusterCache := cache.NewClusterCache(config,
		cache.SetNamespaces(namespaces),
		// cache.SetLogr(log),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {
			// store gc mark of every resource
			circleMark := un.GetAnnotations()[circleAnnotation]
			projectMark := un.GetAnnotations()[projectAnnotation]
			info = &resourceInfo{
				projectMark: un.GetAnnotations()[projectAnnotation],
				circleMark:  un.GetAnnotations()[circleAnnotation],
			}
			// cache resources that has that mark to improve performance

			cacheManifest = circleMark != "" && projectMark != ""
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

func GetManifests(circle v1alpha1.Circle) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}

	if circle.Spec.Release == nil {
		return manifests, nil
	}

	projects := circle.Spec.Release.Projects
	for _, project := range projects {
		r, err := cloneAndOpenRepository(project)
		if err != nil {
			return nil, err
		}

		w, err := r.Worktree()
		if err != nil {
			return nil, err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}

		newManifests, err := repository.ParseManifests(project)
		if err != nil {
			return nil, err
		}

		for _, newManifest := range newManifests {

			annotations := newManifest.GetAnnotations()
			if annotations == nil {
				annotations = make(map[string]string)
			}
			annotations[circleAnnotation] = circle.GetName()
			annotations[projectAnnotation] = project.Name
			newManifest.SetAnnotations(annotations)
		}

		manifests = append(manifests, newManifests...)
	}

	return manifests, nil
}

func (syncConfig SyncConfig) getProjectMap(result []common.ResourceSyncResult) (map[string][]v1alpha1.ResourceStatus, error) {
	projectMap := map[string][]v1alpha1.ResourceStatus{}
	for _, res := range result {

		mapResourceKey := syncConfig.ClusterCache.GetNamespaceTopLevelResources(syncConfig.Namespace)
		resKey := kube.NewResourceKey(
			res.ResourceKey.Group,
			res.ResourceKey.Kind,
			res.ResourceKey.Namespace,
			res.ResourceKey.Name,
		)

		resourceStatus := v1alpha1.ResourceStatus{
			Kind:  res.ResourceKey.Kind,
			Group: res.ResourceKey.Group,
			Name:  res.ResourceKey.Name,
		}

		var status *health.HealthStatus
		var err error

		if mapResourceKey[resKey] != nil {
			status, err = health.GetResourceHealth(mapResourceKey[resKey].Resource, nil)
			if err != nil {
				return nil, err
			}
		}

		if status != nil {
			resourceStatus.Health = &v1alpha1.ResourceHealth{
				Status:  status.Status,
				Message: status.Message,
			}
		}

		if mapResourceKey[resKey] != nil {
			resource := mapResourceKey[resKey].Resource
			projectName := resource.GetAnnotations()[projectAnnotation]
			projectMap[projectName] = append(projectMap[projectName], resourceStatus)
		}

	}

	return projectMap, nil
}

func (syncConfig SyncConfig) getResourceByGroupVersionAndKind(gv schema.GroupVersion, kind string) string {
	resource := ""
	r, _ := syncConfig.Disco.ServerResourcesForGroupVersion(gv.String())
	for _, i := range r.APIResources {
		if i.Kind == kind {
			resource = i.Name
			break
		}
	}

	return resource
}

func (syncConfig SyncConfig) sync(circleName string, manifests []*unstructured.Unstructured) ([]common.ResourceSyncResult, error) {
	return syncConfig.GitopsEngine.Sync(context.Background(), manifests, func(r *cache.Resource) bool {
		return r.Info.(*resourceInfo).circleMark == circleName
	}, "", syncConfig.Namespace, sync.WithOperationSettings(false, syncConfig.Prune, true, false), sync.WithLogr(syncConfig.Log))
}

func (syncConfig SyncConfig) prune(result []common.ResourceSyncResult) error {
	// Delete resource manually with prune false
	for _, res := range result {
		if res.Status == "PruneSkipped" {
			gv := schema.GroupVersion{
				Group:   res.ResourceKey.Group,
				Version: res.Version,
			}
			resource := syncConfig.getResourceByGroupVersionAndKind(gv, res.ResourceKey.Kind)
			gvr := gv.WithResource(resource)
			err := syncConfig.KubeClient.Resource(gvr).Namespace(syncConfig.Namespace).Delete(context.Background(), res.ResourceKey.Name, metav1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (syncConfig SyncConfig) printResources(result []common.ResourceSyncResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "RESOURCE\tSTATUS\tSYNCPHASE\tRESULT\n")
	for _, res := range result {

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", res.ResourceKey.String(), res.Status, res.SyncPhase, res.Message)
	}
	_ = w.Flush()
}

func (syncConfig SyncConfig) Do(circleName string, manifests []*unstructured.Unstructured, isDeletionCircle bool) error {

	result, err := syncConfig.sync(circleName, manifests)
	if err != nil {
		return err
	}

	err = syncConfig.prune(result)
	if err != nil {
		return err
	}

	syncConfig.printResources(result)

	if !isDeletionCircle {
		projectMap, err := syncConfig.getProjectMap(result)
		if err != nil {
			return err
		}

		projectsHealth := true
		projectsStatus := []v1alpha1.ProjectStatus{}
		for projectName, resources := range projectMap {

			for _, res := range resources {
				if res.Health.Status != health.HealthStatusHealthy {
					projectsHealth = false
					break
				}
			}

			projectsStatus = append(projectsStatus, v1alpha1.ProjectStatus{
				Name:      projectName,
				Resources: resources,
			})
		}

		if projectsHealth {

			for projectName := range projectMap {
				vs := &v1beta1.VirtualService{
					TypeMeta: metav1.TypeMeta{
						Kind: "VirtualService",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: projectName,
					},
					Spec: networkingv1beta1.VirtualService{
						Hosts: []string{},
						Http: []*networkingv1beta1.HTTPRoute{
							{
								Route: []*networkingv1beta1.HTTPRouteDestination{
									{
										Destination: &networkingv1beta1.Destination{
											Host: "guestbook-ui",
										},
									},
								},
							},
						},
					},
				}

				b, err := json.Marshal(vs)
				if err != nil {
					return err
				}

				fmt.Println("------------VSB---------", string(b))

				var un *unstructured.Unstructured
				err = json.Unmarshal(b, &un)
				if err != nil {
					return err
				}
				fmt.Println("------------VSJ---------", un.Object)

				_, err = syncConfig.IstioClient.VirtualServices(syncConfig.Namespace).Create(context.TODO(), vs, metav1.CreateOptions{})
				if err != nil && k8sErrors.IsAlreadyExists(err) {
					return retry.RetryOnConflict(retry.DefaultBackoff, func() error {

						vs, err := syncConfig.IstioClient.VirtualServices(syncConfig.Namespace).Get(context.TODO(), projectName, metav1.GetOptions{})
						if err != nil {
							return err
						}

						_, err = syncConfig.IstioClient.VirtualServices(syncConfig.Namespace).Update(context.TODO(), vs, metav1.UpdateOptions{})
						return err
					})
				}

				if err != nil {
					return err
				}

			}

		}

		return syncConfig.updateCircle(circleName, projectsStatus)

	}

	return nil
}

func (syncConfig SyncConfig) updateCircle(circleName string, projectStatusList []v1alpha1.ProjectStatus) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := syncConfig.Client.Get(context.TODO(), circleName, metav1.GetOptions{})
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
