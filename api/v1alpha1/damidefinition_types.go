/*
Copyright 2022.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DamiDefinitionSpec defines the desired state of DamiDefinition
type DamiDefinitionSpec struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64

	// Resp corresponds to response that dami api going to return.
	// For example, if resp field of the DamiDefinition is "resp from k8s",
	// dami api will return JSON response with 'document' field equals to
	// given resp field, "resp from k8s".
	Resp string `json:"resp"`

	// Port specifies the port that dami server listens. Currently, it is ineffective.
	// +optional
	// +kubebuilder:default=8001
	Port int32 `json:"port,omitempty"`
}

// DamiDefinitionStatus defines the observed state of DamiDefinition
type DamiDefinitionStatus struct {
	DefinitionID string `json:"definition_id"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DamiDefinition is the Schema for the damidefinitions API
type DamiDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DamiDefinitionSpec   `json:"spec,omitempty"`
	Status DamiDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DamiDefinitionList contains a list of DamiDefinition
type DamiDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DamiDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DamiDefinition{}, &DamiDefinitionList{})
}
