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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CircleProject struct {
	Name    string `json:"name"`
	RepoURL string `json:"repoUrl"`
	Path    string `json:"path"`
}

type CircleRelease struct {
	Name     string          `json:"name"`
	Tag      string          `json:"tag"`
	Projects []CircleProject `json:"projects"`
}

type CircleDestination struct {
	Namespace string `json:"namespace"`
}

type CircleSpec struct {
	Release     CircleRelease     `json:"release"`
	Destination CircleDestination `json:"destination"`
}

type ResourceStatus struct {
	Group   string               `json:"group,omitempty"`
	Version string               `json:"version,omitempty"`
	Kind    string               `json:"kind,omitempty"`
	Name    string               `json:"name,omitempty"`
	Status  string               `json:"status,omitempty"`
	Health  *health.HealthStatus `json:"health,omitempty"`
}

type ProjectStatus struct {
	Name      string           `json:"name,omitempty"`
	Status    string           `json:"status,omitempty"`
	Resources []ResourceStatus `json:"resources,omitempty"`
}

type CircleStatus struct {
	Projects []ProjectStatus `json:"projects,omitempty"`
}

type Circle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CircleSpec   `json:"spec"`
	Status            CircleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type CircleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Circle `json:"items"`
}
