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

// ClusterLoadAssignmentSpec defines the desired state of ClusterLoadAssignment.
type ClusterLoadAssignmentSpec struct {
	// +required
	ClusterLoadAssignment Message `json:"clusterLoadAssignment"`
}

// ClusterLoadAssignmentStatus defines the observed state of ClusterLoadAssignment.
type ClusterLoadAssignmentStatus struct {
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterLoadAssignment is the Schema for the clusterloadassignments API.
//
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration.html#eds
type ClusterLoadAssignment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterLoadAssignmentSpec   `json:"spec,omitempty"`
	Status ClusterLoadAssignmentStatus `json:"status,omitempty"`
}

// GetStatusConditions ...
func (c *ClusterLoadAssignment) GetStatusConditions() []Condition {
	return c.Status.Conditions
}

// SetStatusConditions ...
func (c *ClusterLoadAssignment) SetStatusConditions(conditions []Condition) {
	c.Status.Conditions = conditions
}

// GetSpecMessage ...
func (c *ClusterLoadAssignment) GetSpecMessage() *Message {
	return &c.Spec.ClusterLoadAssignment
}

var _ Object = &ClusterLoadAssignment{}

// +kubebuilder:object:root=true

// ClusterLoadAssignmentList contains a list of ClusterLoadAssignment.
type ClusterLoadAssignmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterLoadAssignment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterLoadAssignment{}, &ClusterLoadAssignmentList{})
}
