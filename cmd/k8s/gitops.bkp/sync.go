package main

import (
	"fmt"
	"log"

	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/annotation"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/diff"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/argoproj/gitops-engine/pkg/utils/tracing"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/klog/klogr"
)

type SyncOpts struct {
	config       *rest.Config
	namespace    string
	clusterCache cache.ClusterCache
	kubectl      kube.Kubectl
}

type resourceInfo struct {
	manageMark string
}

func NewSyncOpts(
	config *rest.Config,
	namespace string,
) SyncOpts {
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

	err := clusterCache.EnsureSynced()
	if err != nil {
		log.Fatalln(err)
	}

	return SyncOpts{
		config:       config,
		namespace:    namespace,
		clusterCache: clusterCache,
		kubectl: &kube.KubectlCmd{
			Log:    klogr.New(),
			Tracer: tracing.NopTracer{},
		},
	}
}

func (o *SyncOpts) Sync(manifests []*unstructured.Unstructured, opts ...sync.SyncOpt) apperror.Error {
	managedResources, err := o.clusterCache.GetManagedLiveObjs(manifests, func(r *cache.Resource) bool {
		return r.Info.(*resourceInfo).manageMark != ""
	})
	if err != nil {
		return apperror.New("Sync failed", err.Error()).AddOperation("gitops.Sync")
	}

	result := sync.Reconcile(manifests, managedResources, o.namespace, o.clusterCache)
	diffRes, err := diff.DiffArray(result.Target, result.Live)
	if err != nil {
		return nil
	}

	opts = append(opts, sync.WithSkipHooks(!diffRes.Modified))
	syncCtx, err := sync.NewSyncContext("", result, o.config, o.config, o.kubectl, o.namespace, opts...)
	if err != nil {
		return nil
	}

	syncCtx.Sync()
	phase, message, _ := syncCtx.GetState()
	if phase.Completed() {
		if phase == common.OperationError {
			err = fmt.Errorf("sync operation failed: %s", message)
		}
		if err != nil {
			return apperror.New("Sync failed", err.Error()).AddOperation("gitops.OperationError")
		}
	}

	return nil
}
