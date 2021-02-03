package engine

import (
	"fmt"
	"time"

	gitopsEngineCache "github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/diff"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/argoproj/gitops-engine/pkg/utils/tracing"
	"github.com/maycommit/circlerr/internal/k8s/controller/cache"
	cacheUtils "github.com/maycommit/circlerr/internal/k8s/controller/utils/cache"
	"github.com/maycommit/circlerr/pkg/customerror"
	circleApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/klog/klogr"
)

type Engine struct {
	appCache cache.Cache
	kubectl  kube.Kubectl
	config   *rest.Config
}

func New(appCache cache.Cache) *Engine {
	e := &Engine{
		appCache: appCache,
		kubectl: &kube.KubectlCmd{
			Log:    klogr.New(),
			Tracer: tracing.NopTracer{},
		},
	}
	return e
}

func (e *Engine) GetManifests(circle circleApi.Circle) ([]*unstructured.Unstructured, error) {
	manifests := []*unstructured.Unstructured{}

	if circle.Spec.Release != nil {

	}
}

func (e *Engine) SyncCluster(
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

	opts := []sync.SyncOpt{sync.WithSkipHooks(!diffRes.Modified)}
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

func (e *Engine) IsIterableCircle(circle circleApi.Circle) bool {
	return circle.Spec.Release != nil
}

func (e *Engine) Wave(circle circleApi.Circle) {

}

func (e *Engine) Start() error {
	err := e.appCache.Cluster().EnsureSynced()
	if err != nil {
		logrus.Fatalln(customerror.LogFields(customerror.New("Failed ensyre cache", err, nil, "engine.Start")))
	}

	for {
		for _, c := range e.appCache.Circle().List() {
			e.Wave(c.Circle)
		}

		time.Sleep(3 * time.Second)
	}
}
