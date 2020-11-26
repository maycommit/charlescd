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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CircleProject struct {
	Name    string `json:"name" protobuf:"bytes,1,opt,name=name"`
	RepoURL string `json:"repoUrl" protobuf:"bytes,2,opt,name=repoUrl"`
	Path    string `json:"path" protobuf:"bytes,3,opt,name=path"`
}

type CircleRelease struct {
	Name     string          `json:"name" protobuf:"bytes,1,opt,name=name"`
	Tag      string          `json:"tag" protobuf:"bytes,2,opt,name=tag"`
	Projects []CircleProject `json:"projects" protobuf:"bytes,3,opt,name=projects"`
}

type CircleDestination struct {
	Namespace string `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
}

type CircleSpec struct {
	Release     CircleRelease     `json:"release" protobuf:"bytes,1,opt,name=release"`
	Destination CircleDestination `json:"destination" protobuf:"bytes,2,opt,name=destination"`
}

type ResourceHealth struct {
	Status  string `json:"status,omitempty" protobuf:"bytes,1,opt,name=string"`
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
}

type ResourceStatus struct {
	Group   string         `json:"group,omitempty" protobuf:"bytes,1,opt,name=group"`
	Version string         `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
	Kind    string         `json:"kind,omitempty" protobuf:"bytes,3,opt,name=kind"`
	Name    string         `json:"name,omitempty" protobuf:"bytes,4,opt,name=name"`
	Status  string         `json:"status,omitempty" protobuf:"bytes,5,opt,name=status"`
	Health  *ResourceHealth `json:"health,omitempty" protobuf:"bytes,5,opt,name=health"`
}

type CircleStatus struct {
	Resources []ResourceStatus `json:"resources,omitempty" protobuf:"bytes,1,opt,name=resources"`
}

type Circle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              CircleSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            CircleStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type CircleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Circle `json:"items" protobuf:"bytes,2,rep,name=items"`
}
