package gitops

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/maycommit/circlerr/internal/k8s/controller/utils/git"

	"github.com/maycommit/circlerr/internal/k8s/controller/utils/kube"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

func init() {
	config := flag.String("config", "", "Path to config repository list.")
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	masterUrl := flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	gitDir := flag.String("gitdir", "./tmp/git", "")
	flag.Parse()

	os.Setenv("CONFIG", *gitDir)
	os.Setenv("GIT_DIR", *gitDir)
	os.Setenv("KUBECONFIG", *kubeconfig)
	os.Setenv("MASTER_URL", *masterUrl)
}

type Repository struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

type Repositories struct {
	Repositories []Repository `yaml:"repositories"`
}

func loadRepositories() (*Repositories, error) {
	conf := &Repositories{}
	configData, err := ioutil.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configData, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func main() {

	repositories, err := loadRepositories()
	if err != nil {
		klog.Fatalf("Load repositories by config file failed: %s\n", err.Error())
	}

	config, err := kube.GetClusterConfig()
	if err != nil {
		klog.Fatalf("Get cluster config failed: %s\n", err.Error())
	}

	for _, r := range repositories.Repositories {
		gitRepository, err := git.CloneAndOpenRepository(r.Url)
		if err != nil {
			klog.Fatalf("Clone or open git repository failed: %s\n", err.Error())
		}

		revision, err := git.SyncRepository()
	}
}
