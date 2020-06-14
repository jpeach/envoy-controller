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

// ScopedRouteConfigurationSpec defines the desired state of ScopedRouteConfiguration
type ScopedRouteConfigurationSpec struct {
	// +required
	ScopedRouteConfiguration Message `json:"scopedRouteConfiguration"`
}

// ScopedRouteConfigurationStatus defines the observed state of ScopedRouteConfiguration
type ScopedRouteConfigurationStatus struct {
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ScopedRouteConfiguration is the Schema for the scopedrouteconfigurations API
//
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration.html#srds
type ScopedRouteConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScopedRouteConfigurationSpec   `json:"spec,omitempty"`
	Status ScopedRouteConfigurationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScopedRouteConfigurationList contains a list of ScopedRouteConfiguration
type ScopedRouteConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScopedRouteConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScopedRouteConfiguration{}, &ScopedRouteConfigurationList{})
}
