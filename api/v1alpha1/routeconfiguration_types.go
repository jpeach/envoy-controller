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

// RouteConfigurationSpec defines the desired state of RouteConfiguration.
type RouteConfigurationSpec struct {
	RouteConfiguration Message `json:"routeConfiguration"`
}

// RouteConfigurationStatus defines the observed state of RouteConfiguration.
type RouteConfigurationStatus struct {
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RouteConfiguration is the Schema for the routeconfigurations API.
//
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration.html#rds
type RouteConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteConfigurationSpec   `json:"spec,omitempty"`
	Status RouteConfigurationStatus `json:"status,omitempty"`
}

// GetStatusConditions ...
func (c *RouteConfiguration) GetStatusConditions() []Condition {
	return c.Status.Conditions
}

// SetStatusConditions ...
func (c *RouteConfiguration) SetStatusConditions(conditions []Condition) {
	c.Status.Conditions = conditions
}

// GetSpecMessage ...
func (c *RouteConfiguration) GetSpecMessage() *Message {
	return &c.Spec.RouteConfiguration
}

var _ Object = &RouteConfiguration{}

// +kubebuilder:object:root=true

// RouteConfigurationList contains a list of RouteConfiguration.
type RouteConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RouteConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RouteConfiguration{}, &RouteConfigurationList{})
}
