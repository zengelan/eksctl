/*
Copyright 2018 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// kustomize "sigs.k8s.io/kustomize/pkg/types"
)

type ModuleSrource struct {
	URL string `json:"url,omitempty"`
	//Kustomizations *kustomize.Kustomization `json:"kustomizations"`
	HelmTemplate *HelmTemplate `json:"helmTemplate"`
}

type HelmTemplate struct {
	ChartDir string `json:"chartDir"`
}

// ModuleSpec defines the desired state of Module
type ModuleSpec struct {
	// Important: Run "make" to regenerate code after modifying this file
	Source   ModuleSrource `json:"source"`
	External bool          `json:"external,omitempty"`
}

// ModuleStatus defines the observed state of Module
type ModuleStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// Module is the Schema for the modules API
// +k8s:openapi-gen=true
type Module struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModuleSpec   `json:"spec,omitempty"`
	Status ModuleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// ModuleList contains a list of Module
type ModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Module `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Module{}, &ModuleList{})
}
