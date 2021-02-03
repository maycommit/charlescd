package cluster

import (
	"os"
	"testing"

	fakecache "github.com/maycommit/circlerr/internal/k8s/controller/cache/fake"
	circleApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"
	projectApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/project/v1alpha1"
	clientset "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned"

	custominformers "github.com/maycommit/circlerr/pkg/k8s/controller/client/informers/externalversions"

	clientsetFake "github.com/maycommit/circlerr/pkg/k8s/controller/client/clientset/versioned/fake"

	"github.com/maycommit/circlerr/internal/k8s/controller/cache"

	"github.com/argoproj/gitops-engine/pkg/utils/kube/kubetest"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicFake "k8s.io/client-go/dynamic/fake"
	kubeFake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

type ClusterTestSuite struct {
	suite.Suite

	namespace             string
	dynamicClient         *dynamicFake.FakeDynamicClient
	customClientsetFake   *clientsetFake.Clientset
	kubeClient            *kubeFake.Clientset
	syncOpts              SyncOpts
	clientset             *clientset.Clientset
	controllerCache       *cache.Cache
	customInformerFactory custominformers.SharedInformerFactory
}

func (s *ClusterTestSuite) SetupSuite() {
	os.Setenv("KUBECONFIG_PATH", "../../../test/kubeconfig.yaml")
	os.Setenv("K8S_CONN_TYPE", "out")

	namespace := "default"

	config := &rest.Config{}

	s.customClientsetFake = clientsetFake.NewSimpleClientset()
	s.customInformerFactory = custominformers.NewSharedInformerFactory(s.customClientsetFake, 0)
	s.dynamicClient = dynamicFake.NewSimpleDynamicClient(runtime.NewScheme())
	s.kubeClient = kubeFake.NewSimpleClientset()
	s.controllerCache = fakecache.NewCache(config, namespace)
	s.syncOpts = New(config, s.dynamicClient, namespace, s.controllerCache, nil, &kubetest.MockKubectlCmd{})
	s.clientset = clientset.NewForConfigOrDie(config)
}

func TestClusterTestSuite(t *testing.T) {
	suite.Run(t, new(ClusterTestSuite))
}

var project1 = projectApi.Project{
	TypeMeta: v1.TypeMeta{
		APIVersion: "circlerr.io/v1alpha1",
		Kind:       "Project",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "project-1",
	},
	Spec: projectApi.ProjectSpec{
		RepoURL: "https://github.com/maycommit/argocd-example-apps",
		Path:    "guestbook",
		Template: &projectApi.ProjectTemplate{
			TemplateType: "puremanifest",
		},
	},
}

var circle1 = circleApi.Circle{
	TypeMeta: v1.TypeMeta{
		APIVersion: "circlerr.io/v1alpha1",
		Kind:       "Circle",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: "circle-1",
	},
	Spec: circleApi.CircleSpec{
		Segments: []circleApi.Segment{
			{Key: "username", Condition: "=", Value: "test@mail.com"},
		},
		Release: &circleApi.CircleRelease{
			Name: "release-1",
			Projects: []circleApi.CircleProject{
				{
					Name:  project1.GetName(),
					Image: "latest",
				},
			},
		},
		Destination: circleApi.CircleDestination{
			Namespace: "default",
		},
	},
}

var gv = schema.GroupVersion{
	Group:   "circlerr.io",
	Version: "v1alpha",
}

var deployment = &unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": "demo-deployment",
		},
		"spec": map[string]interface{}{},
	},
}

func (s *ClusterTestSuite) TestCreateResourcesByCircle() {
	return
}
