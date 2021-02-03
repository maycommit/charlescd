package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache/circle"
	v1alpha12 "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"
	"github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appCache "github.com/maycommit/circlerr/internal/k8s/controller/cache"
	apperror "github.com/maycommit/circlerr/internal/k8s/controller/error"
	"github.com/maycommit/circlerr/internal/k8s/controller/router"
	"github.com/maycommit/circlerr/internal/k8s/controller/template"
	"github.com/maycommit/circlerr/internal/k8s/controller/utils/git"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/diff"
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/sync"
	"github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	kubeutil "github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

func (o *SyncOpts) getManifests(circleName string, circle *circle.CircleCache) ([]*unstructured.Unstructured, apperror.Error) {
	manifests := []*unstructured.Unstructured{}

	for _, circleProject := range circle.GetRelease().Projects {
		revision, err := git.SyncRepository(circleProject.Name, o.appCache.Projects.Get(circleProject.Name))
		if err != nil {
			return nil, err.AddOperation("cluster.getManifests.SyncRepository")
		}

		project := o.appCache.Projects.Get(circleProject.Name)
		if len(circle.GetManifests()) <= 0 || revision != project.GetRevision() {
			t := template.NewTemplate(circleProject.Name, project)
			newManifests, err := t.ParseManifests(circleName, circle)
			if err != nil {
				return nil, err
			}

			o.appCache.Circles.Get(circleName).SetManifests(newManifests)
			manifests = append(manifests, newManifests...)
		}

	}

	return manifests, nil
}

func (o *SyncOpts) Sync(manifests []*unstructured.Unstructured, opts ...sync.SyncOpt) ([]common.ResourceSyncResult, apperror.Error) {
	managedResources, err := o.appCache.Cluster.Get().GetManagedLiveObjs(manifests, func(r *cache.Resource) bool {
		return o.appCache.Cluster.IsManagedResource(r)
	})
	if err != nil {
		return nil, apperror.New("sync failed", err.Error()).AddOperation("cluster.Sync.GetManagedLiveObjs")
	}

	result := sync.Reconcile(manifests, managedResources, o.namespace, o.appCache.Cluster.Get())
	diffRes, err := diff.DiffArray(result.Target, result.Live)
	if err != nil {
		return nil, apperror.New("sync failed", err.Error()).AddOperation("cluster.Sync.Reconcile")
	}

	opts = append(opts, sync.WithSkipHooks(!diffRes.Modified))
	syncCtx, err := sync.NewSyncContext("", result, o.config, o.config, o.kubectl, o.namespace, opts...)
	if err != nil {
		return nil, apperror.New("sync failed", err.Error()).AddOperation("cluster.Sync.NewSyncContext")
	}

	syncCtx.Sync()
	phase, message, resources := syncCtx.GetState()
	if phase.Completed() {
		if phase == common.OperationError {
			err = fmt.Errorf("sync operation failed: %s", message)
		}

		if err != nil {
			return resources, apperror.New("sync failed", err.Error()).AddOperation("cluster.Sync.OperationError")
		}
	}

	return resources, nil
}

func (o *SyncOpts) updateProjectRoutes() {
	for projectName, project := range o.appCache.Projects.List() {
		for circleName := range project.GetRoutes() {
			currentCircle := o.appCache.Circles.Get(circleName)
			if currentCircle != nil {
				if currentCircle.Status.Status != health.HealthStatusHealthy {
					o.appCache.Projects.Get(projectName).DeleteRoute(circleName)
				}
			} else {
				o.appCache.Projects.Get(projectName).DeleteRoute(circleName)
			}
		}
	}

}

type SyncOpts struct {
	config         *rest.Config
	kubeClient     dynamic.Interface
	namespace      string
	appCache       *appCache.Cache
	router         router.UseCases
	kubectl        kube.Kubectl
	circlerrClient *versioned.Clientset
}

func New(
	config *rest.Config,
	kubeClient dynamic.Interface,
	namespace string,
	appCache *appCache.Cache,
	router router.UseCases,
	kubectl kubeutil.Kubectl,
	circlerrClient *versioned.Clientset,
) SyncOpts {
	return SyncOpts{
		config:         config,
		kubeClient:     kubeClient,
		namespace:      namespace,
		appCache:       appCache,
		router:         router,
		kubectl:        kubectl,
		circlerrClient: circlerrClient,
	}

}

func (o *SyncOpts) addErrorToCircle(circleName string, errStr []string) error {
	circle := o.appCache.Circles.Get(circleName)
	updateCircle := circle.Circle

	if len(errStr) > 0 && len(updateCircle.Status.Errors) <= 0 {
		updateCircle.Status.Status = health.HealthStatusDegraded
		updateCircle.Status.Projects = []v1alpha12.ProjectStatus{}
		updateCircle.Status.Errors = append(updateCircle.Status.Errors, errStr...)
	}

	if len(errStr) <= 0 && len(updateCircle.Status.Errors) > 0 {
		updateCircle.Status.Status = health.HealthStatusProgressing
		updateCircle.Status.Errors = []string{}
	}

	_, err := o.circlerrClient.CircleV1alpha1().Circles(o.namespace).Update(context.TODO(), &updateCircle, metav1.UpdateOptions{})
	if err != nil && errors.IsNotFound(err) {
		return nil
	}

	return err
}

func (o *SyncOpts) Run(stopCh chan struct{}) {
	goError := o.appCache.Cluster.Get().EnsureSynced()
	if goError != nil {
		errorLogFields := apperror.New("Run cluster sync error", goError.Error()).
			AddOperation("cluster.Run.EnsureSynced").
			LogFields()

		logrus.WithFields(errorLogFields).Warn()
		return
	}

	ticker := time.NewTicker(time.Second * 4)
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:

			manifests := []*unstructured.Unstructured{}
			for circleName, circle := range o.appCache.Circles.List() {
				if circle.Spec.Release == nil || circle.GetDeletion() {
					continue
				}
				newManifests, err := o.getManifests(circleName, circle)
				if err != nil {
					updateErr := o.addErrorToCircle(circleName, []string{err.Error()})
					if updateErr != nil {
						panic(updateErr)
					}

					logrus.WithFields(err.AddOperation("cluster.Run.getManifests").LogFields()).Error()
					continue
				}

				o.addErrorToCircle(circleName, []string{})
				manifests = append(manifests, newManifests...)
			}

			_, err := o.Sync(manifests, sync.WithOperationSettings(false, true, false, false))
			if err != nil {
				logrus.WithFields(err.AddOperation("cluster.Run.Sync").LogFields()).Error()
				return
			}

			err = o.refreshCircles()
			if err != nil {
				logrus.WithFields(err.AddOperation("cluster.Run.refreshCircles").LogFields()).Error()
				continue
			}

			if o.router != nil {
				err = o.router.Sync()
				if err != nil {
					logrus.WithFields(err.AddOperation("cluster.Run.RouterSync").LogFields()).Error()
					continue
				}
			}
		}
	}
}
