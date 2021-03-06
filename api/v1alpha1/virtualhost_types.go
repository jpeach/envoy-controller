/*
Copyright 2020 VMware, Inc.

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
)

// VirtualHostSpec defines the desired state of VirtualHost.
type VirtualHostSpec struct {
	// +required
	VirtualHost Message `json:"virtualHost"`
}

// VirtualHostStatus defines the observed state of VirtualHost.
type VirtualHostStatus struct {
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VirtualHost is the Schema for the virtualhosts API.
//
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration.html#vhds
type VirtualHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualHostSpec   `json:"spec,omitempty"`
	Status VirtualHostStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualHostList contains a list of VirtualHost.
type VirtualHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualHost{}, &VirtualHostList{})
}
