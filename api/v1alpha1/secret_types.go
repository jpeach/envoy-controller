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

// SecretSpec defines the desired state of Secret.
type SecretSpec struct {
	// +required
	Secret Message `json:"secret"`
}

// SecretStatus defines the observed state of Secret.
type SecretStatus struct {
	Conditions []Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Secret is the Schema for the secrets API.
//
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration.html#sds
type Secret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretSpec   `json:"spec,omitempty"`
	Status SecretStatus `json:"status,omitempty"`
}

// GetStatusConditions ...
func (s *Secret) GetStatusConditions() []Condition {
	return s.Status.Conditions
}

// SetStatusConditions ...
func (s *Secret) SetStatusConditions(conditions []Condition) {
	s.Status.Conditions = conditions
}

// GetSpecMessage ...
func (s *Secret) GetSpecMessage() *Message {
	return &s.Spec.Secret
}

var _ Object = &Secret{}

// +kubebuilder:object:root=true

// SecretList contains a list of Secret.
type SecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Secret{}, &SecretList{})
}
