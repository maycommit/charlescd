package circle

import (
	"sync"

	circleApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CircleCache struct {
	lock sync.RWMutex
	circleApi.Circle
	manifests  []*unstructured.Unstructured
	isDeletion bool
}

func (c *CircleCache) SetDeletion(isDeletion bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.isDeletion = isDeletion
}

func (c *CircleCache) GetDeletion() bool {
	return c.isDeletion
}

func (c *CircleCache) SetManifests(manifests []*unstructured.Unstructured) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests = manifests
}

func (c *CircleCache) GetManifests() []*unstructured.Unstructured {
	return c.manifests
}

func (c *CircleCache) GetRelease() *circleApi.CircleRelease {
	return c.Spec.Release
}

func (c *CircleCache) GetProject(projectName string) circleApi.CircleProject {
	project := circleApi.CircleProject{}

	for _, circleProject := range c.Spec.Release.Projects {
		if circleProject.Name == projectName {
			project = circleProject
		}
	}

	return project
}

type CirclesCache struct {
	lock    sync.RWMutex
	circles map[string]*CircleCache
}

func NewCircleCache() *CirclesCache {
	return &CirclesCache{
		lock:    sync.RWMutex{},
		circles: map[string]*CircleCache{},
	}
}

func (m *CirclesCache) List() map[string]*CircleCache {
	return m.circles
}

func (m *CirclesCache) Get(circleName string) *CircleCache {
	return m.circles[circleName]
}

func (m *CirclesCache) Put(circleName string, circle circleApi.Circle) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.circles[circleName] = &CircleCache{
		lock:   sync.RWMutex{},
		Circle: circle,
	}
}

func (m *CirclesCache) Delete(circleName string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.circles, circleName)
}

func (m *CirclesCache) LogicDeletion(circleName string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Get(circleName).SetDeletion(true)
}
