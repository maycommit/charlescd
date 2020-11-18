package sync

import (
	"charlescd/internal/manager/circle"
	"charlescd/internal/operator/repository"
	"charlescd/internal/utils"
	"context"
	"fmt"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/go-logr/logr"

	//"github.com/libopenstorage/openstorage/api/client"
	"k8s.io/client-go/rest"
	"os"
	"text/tabwriter"

	"github.com/argoproj/gitops-engine/pkg/cache"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
)

const (
	annotationGCMark = "charlescd.io/circle"
)

type resourceInfo struct {
	gcMark string
}

type SyncResource struct {

}

type SyncConfig struct {
	ClusterCache cache.ClusterCache
	Config *rest.Config
	Disco discovery.DiscoveryInterface
	CircleRes *unstructured.Unstructured
	Namespace string
	Prune bool
	Log logr.Logger
}

func ClusterCache(config *rest.Config, namespaces []string, log logr.Logger) cache.ClusterCache {
	clusterCache := cache.NewClusterCache(config,
		cache.SetNamespaces(namespaces),
		cache.SetLogr(log),
		cache.SetPopulateResourceInfoHandler(func(un *unstructured.Unstructured, isRoot bool) (info interface{}, cacheManifest bool) {
			// store gc mark of every resource
			gcMark := un.GetAnnotations()[annotationGCMark]
			info = &resourceInfo{gcMark: un.GetAnnotations()[annotationGCMark]}
			// cache resources that has that mark to improve performance
			cacheManifest = gcMark != ""
			return
		}),
	)

	return clusterCache
}

func Start(syncConfig SyncConfig) error {


	circleResources, err := circle.GetResourcesByResource(*syncConfig.CircleRes)
	if err != nil {
		return err
	}

	gitOpsEngine := engine.NewEngine(syncConfig.Config, syncConfig.ClusterCache, engine.WithLogr(syncConfig.Log))
	_, err = gitOpsEngine.Run()
	if err != nil {
		return err
	}

	manifests := []*unstructured.Unstructured{}
	for _, resource := range circleResources {
		manifests, err = repository.ParseManifests(resource)
		if err != nil {
		return err
		}
	}

	for _, manifest := range manifests {
		annotations := manifest.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations[annotationGCMark] = syncConfig.CircleRes.GetName()
		manifest.SetAnnotations(annotations)
	}


	result, err := gitOpsEngine.Sync(context.Background(), manifests, func(r *cache.Resource) bool {
		return r.Info.(*resourceInfo).gcMark == syncConfig.CircleRes.GetName()
	}, utils.NewSHA1Hash(), syncConfig.Namespace, sync.WithPrune(syncConfig.Prune), sync.WithLogr(syncConfig.Log))
	if err != nil {
		syncConfig.Log.Error(err, "Failed to synchronize cluster state")
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "RESOURCE\tRESULT\n")
	for _, res := range result {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", res.ResourceKey.String(), res.Message)
	}
	_ = w.Flush()



	return nil
}
