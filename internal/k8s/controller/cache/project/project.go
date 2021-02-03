package project

import (
	"sync"

	projectApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/project/v1alpha1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ProjectRoute struct {
	CircleID    string `json:"circleId"`
	ReleaseName string `json:"releaseName"`
}

type ProjectCache struct {
	lock sync.RWMutex
	projectApi.ProjectSpec
	revision  string
	routes    map[string]ProjectRoute
	manifests []*unstructured.Unstructured
}

func (c *ProjectCache) SetRoute(circleName string, route ProjectRoute) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.routes[circleName] = route
}

func (c *ProjectCache) DeleteRoute(circleName string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.routes, circleName)
}

func (c *ProjectCache) GetRoutes() map[string]ProjectRoute {
	return c.routes
}

func (c *ProjectCache) SetRevision(revision string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.revision = revision
}

func (c *ProjectCache) GetRevision() string {
	return c.revision
}

func (c *ProjectCache) SetManifests(manifests []*unstructured.Unstructured) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests = manifests
}

func (c *ProjectCache) GetManifests() []*unstructured.Unstructured {
	return c.manifests
}

type ProjectsCache struct {
	lock     sync.RWMutex
	projects map[string]*ProjectCache
}

func NewProjectCache() *ProjectsCache {
	return &ProjectsCache{
		lock:     sync.RWMutex{},
		projects: map[string]*ProjectCache{},
	}
}

func (c *ProjectsCache) List() map[string]*ProjectCache {
	return c.projects
}

func (c *ProjectsCache) Get(projectName string) *ProjectCache {
	return c.projects[projectName]
}

func (m *ProjectsCache) Put(projectName string, project projectApi.Project) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.projects[projectName] = &ProjectCache{
		lock:        sync.RWMutex{},
		routes:      map[string]ProjectRoute{},
		ProjectSpec: project.Spec,
	}
}

func (m *ProjectsCache) Delete(projectName string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.projects, projectName)
}
