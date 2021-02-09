package circle

import (
	"sync"

	circleApi "github.com/maycommit/circlerr/pkg/k8s/controller/apis/circle/v1alpha1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type CircleCache struct {
	lock       sync.RWMutex
	circle     circleApi.Circle
	manifests  []*unstructured.Unstructured
	isDeletion bool
}

func (c *CircleCache) Circle() circleApi.Circle {
	return c.circle
}

func (c *CircleCache) SetCircle(circle circleApi.Circle) circleApi.Circle {
	c.circle = circle
	return c.circle
}

func (c *CircleCache) Manifests() []*unstructured.Unstructured {
	return c.manifests
}

func (c *CircleCache) SetManifests(manifests []*unstructured.Unstructured) []*unstructured.Unstructured {
	c.manifests = manifests
	return c.manifests
}

func (c *CircleCache) SetDeletion() bool {
	c.isDeletion = true
	return c.isDeletion
}

func (c *CircleCache) IsDeletion() bool {
	return c.isDeletion
}

type CirclesCache struct {
	lock    sync.RWMutex
	circles map[string]*CircleCache
}

func (c *CirclesCache) Circles() map[string]*CircleCache {
	return c.circles
}

func (c *CirclesCache) Add(circleName string, circle circleApi.Circle) *CircleCache {
	c.circles[circleName] = &CircleCache{
		circle:     circle,
		isDeletion: false,
		manifests:  []*unstructured.Unstructured{},
	}

	return c.circles[circleName]
}

func (c *CirclesCache) Set(circleName string, circle circleApi.Circle) *CircleCache {
	c.circles[circleName].circle = circle

	return c.circles[circleName]
}

func (c *CirclesCache) Get(circleName string) *CircleCache {
	return c.circles[circleName]
}

func (c *CirclesCache) Delete(circleName string) {
	delete(c.circles, circleName)
}

func (c *CirclesCache) IterateAllCircles(cb func(circleName string, circle *CircleCache)) {
	for circleName, circle := range c.Circles() {
		cb(circleName, circle)
	}
}

func NewCirclesCache() *CirclesCache {
	return &CirclesCache{
		lock:    sync.RWMutex{},
		circles: map[string]*CircleCache{},
	}
}
