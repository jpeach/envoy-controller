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

// ListenerSpec defines the desired state of Listener.
type ListenerSpec struct {
	// +required
	Listener Message `json:"listener"`
}

// ListenerStatus defines the observed state of Listener.
type ListenerStatus struct {
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Listener is the Schema for the listeners API.
//
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration.html#lds
type Listener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ListenerSpec   `json:"spec,omitempty"`
	Status ListenerStatus `json:"status,omitempty"`
}

// GetStatusConditions ...
func (l *Listener) GetStatusConditions() []Condition {
	return l.Status.Conditions
}

// SetStatusConditions ...
func (l *Listener) SetStatusConditions(conditions []Condition) {
	l.Status.Conditions = conditions
}

// GetSpecMessage ...
func (l *Listener) GetSpecMessage() *Message {
	return &l.Spec.Listener
}

var _ Object = &Listener{}

// +kubebuilder:object:root=true

// ListenerList contains a list of Listener.
type ListenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Listener `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Listener{}, &ListenerList{})
}
