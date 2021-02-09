package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/go-git/go-git/v5"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	circleCache "github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	"github.com/maycommit/circlerr/internal/k8s/controller/template"
	gitutils "github.com/maycommit/circlerr/internal/k8s/controller/utils/git"
	circleApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"
	circlerrVersioned "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"

	gitopsEngineCache "github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/diff"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/argoproj/gitops-engine/pkg/utils/tracing"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	cacheUtils "github.com/maycommit/circlerr/internal/k8s/controller/utils/cache"
	"github.com/maycommit/circlerr/pkg/customerror"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/klog/klogr"
)

type Engine struct {
	appCache  cache.Cache
	kubectl   kube.Kubectl
	config    *rest.Config
	clientset *circlerrVersioned.Clientset
}

func New(config *rest.Config, appCache cache.Cache, clientset *circlerrVersioned.Clientset) *Engine {
	e := &Engine{
		config:    config,
		appCache:  appCache,
		clientset: clientset,
		kubectl: &kube.KubectlCmd{
			Log:    klogr.New(),
			Tracer: tracing.NopTracer{},
		},
	}
	return e
}

func (e *Engine) getManifests(circleName string, circle *circleCache.CircleCache) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}

	for _, cp := range circle.Circle().Spec.Release.Projects {
		projectCache := e.appCache.Projects().Get(cp.Name)
		gitOptions := git.CloneOptions{
			URL: projectCache.RepoURL,
		}
		remoteRevision, err := gitutils.SyncRepository(gitOptions, projectCache.GetRevision())
		if err != nil {
			return nil, err
		}

		if len(circle.Manifests()) <= 0 || remoteRevision != projectCache.GetRevision() {
			t := template.NewTemplate(cp.Name, projectCache)
			newManifests, err := t.ParseManifests(circleName, circle)
			if err != nil {
				return nil, err
			}

			projectCache.SetRevision(remoteRevision)
			e.appCache.Circles().Get(circleName).SetManifests(newManifests)
			manifests = append(manifests, newManifests...)
		}
	}

	return manifests, nil
}

func (e *Engine) syncCluster(
	manifests []*unstructured.Unstructured,
	namespace string) ([]common.ResourceSyncResult, error) {

	managedResources, err := e.appCache.Cluster().GetManagedLiveObjs(manifests, func(r *gitopsEngineCache.Resource) bool {
		return cacheUtils.IsManagedResource(r)
	})
	if err != nil {
		return nil, customerror.New("Sync cluster failed!", err, nil, "engine.SyncCluster.GetManagedLiveObjs")
	}

	result := sync.Reconcile(manifests, managedResources, namespace, e.appCache.Cluster())
	diffRes, err := diff.DiffArray(result.Target, result.Live)
	if err != nil {
		return nil, customerror.New("Sync cluster failed!", err, nil, "engine.SyncCluster.GetManagedLiveObjs")
	}

	opts := []sync.SyncOpt{sync.WithSkipHooks(!diffRes.Modified), sync.WithOperationSettings(false, true, true, false)}
	syncCtx, err := sync.NewSyncContext("", result, e.config, e.config, e.kubectl, namespace, opts...)
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

func (e *Engine) updateCircleStatus(circleName string, circle *circleCache.CircleCache, results []common.ResourceSyncResult) error {
	resourcesStatus := []circleApi.ResourceStatus{}
	updateCircle := circle.Circle()

	namespace := circle.Circle().Spec.Destination.Namespace
	topLevelResources := e.appCache.Cluster().GetNamespaceTopLevelResources(namespace)

	for _, res := range results {
		t := topLevelResources[res.ResourceKey]

		if t.Resource != nil {
			healthStatus, err := health.GetResourceHealth(t.Resource, nil)
			if err != nil {
				return err
			}

			r := circleApi.ResourceStatus{
				Group:             res.ResourceKey.Group,
				Kind:              res.ResourceKey.Kind,
				Name:              res.ResourceKey.Name,
				Namespace:         res.ResourceKey.Namespace,
				Health:            healthStatus,
				Status:            string(res.Status),
				CreationTimestamp: *t.CreationTimestamp,
			}

			resourcesStatus = append(resourcesStatus, r)
		}

	}

	updateCircle.Status = circleApi.CircleStatus{
		Projects: []circleApi.ProjectStatus{
			{Resources: resourcesStatus},
		},
	}
	_, err := e.clientset.CircleV1alpha1().Circles("default").Update(
		context.TODO(),
		&updateCircle,
		metav1.UpdateOptions{},
	)
	return err
}

func (e *Engine) wave(circleName string, circle *circleCache.CircleCache) error {
	manifests, err := e.getManifests(circleName, circle)
	if err != nil {
		return err
	}

	namespace := circle.Circle().Spec.Destination.Namespace
	results, err := e.syncCluster(manifests, namespace)
	if err != nil {
		return err
	}

	err = e.updateCircleStatus(circleName, circle, results)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) Start() error {
	err := e.appCache.Cluster().EnsureSynced()
	if err != nil {
		logrus.Fatalln(customerror.LogFields(customerror.New("Failed ensure cache", err, nil, "engine.Start")))
	}

	for {
		e.appCache.Circles().IterateAllCircles(func(circleName string, circle *circleCache.CircleCache) {
			if circle.Circle().Spec.Release != nil && !circle.IsDeletion() {
				err := e.wave(circleName, circle)
				if err != nil {
					logrus.Warnln(customerror.LogFields(err))
					return
				}
			}
		})

		time.Sleep(3 * time.Second)
	}
}
