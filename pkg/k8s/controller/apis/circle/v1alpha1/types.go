/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/argoproj/gitops-engine/pkg/health"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Helm struct {
}

type CircleTemplate struct {
	Helm *Helm `json:"helm"`
}

type CircleProject struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type CircleRelease struct {
	Name     string          `json:"name"`
	Projects []CircleProject `json:"projects"`
}

type CircleDestination struct {
	Namespace string `json:"namespace"`
}

type Segment struct {
	Key       string `json:"key"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
}

type Environment struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Router struct {
	RouterType string `json:"routerType"`
}

type CircleSpec struct {
	Release      *CircleRelease    `json:"release,omitempty"`
	Destination  CircleDestination `json:"destination,omitempty"`
	Environments []Environment     `json:"environments,omitempty"`
	Segments     []Segment         `json:"segments,omitempty"`
}

type ResourceStatus struct {
	Group             string               `json:"group,omitempty" protobuf:"bytes,1,opt,name=group"`
	Kind              string               `json:"kind,omitempty" protobuf:"bytes,3,opt,name=kind"`
	Name              string               `json:"name,omitempty" protobuf:"bytes,5,opt,name=name"`
	Namespace         string               `json:"namespace"`
	Status            string               `json:"status"`
	Health            *health.HealthStatus `json:"health,omitempty" protobuf:"bytes,7,opt,name=health"`
	CreationTimestamp metav1.Time          `json:"creationTimestamp"`
}

type ProjectStatus struct {
	Name      string              `json:"name,omitempty"`
	Health    health.HealthStatus `json:"status,omitempty"`
	Resources []ResourceStatus    `json:"resources,omitempty"`
}

type CircleStatus struct {
	Health        health.HealthStatusCode `json:"status,omitempty"`
	Projects      []ProjectStatus         `json:"projects,omitempty"`
	BreakingError string                  `json:"error"`
}

type CircleError struct {
}

type Circle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CircleSpec   `json:"spec"`
	Status            CircleStatus `json:"status,omitempty"`
}

func (c *Circle) ToUnstructured() (*unstructured.Unstructured, error) {
	return kube.ToUnstructured(c)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type CircleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Circle `json:"items"`
}
