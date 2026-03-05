/*
Copyright 2021.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisScaledownSpec defines the desired state of ActiveMQArtemisScaledown
type ActiveMQArtemisScaledownSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Triggered by main ActiveMQArtemis CRD messageMigration entry
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Temporary",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	LocalOnly bool `json:"localOnly"`
	// Specifies the minimum/maximum amount of compute resources required/allowed
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resource Requirements",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:resourceRequirements"}
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// ActiveMQArtemisScaledownStatus defines the observed state of ActiveMQArtemisScaledown
type ActiveMQArtemisScaledownStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Current state of the resource
	// Conditions represent the latest available observations of an object's state
	//+optional
	//+patchMergeKey=type
	//+patchStrategy=merge
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Conditions",xDescriptors="urn:alm:descriptor:io.kubernetes.conditions"
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:storageversion
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=activemqartemisscaledowns,shortName=aad
//+operator-sdk:csv:customresourcedefinitions:resources={{"Secret", "v1"}}

// +kubebuilder:deprecatedversion:warning="The ActiveMQArtemisScaledown CRD is deprecated, it is an internal only api"
// Provides internal message migration on clustered broker scaledown
// +operator-sdk:csv:customresourcedefinitions:displayName="ArkMQ Scaledown"
type ActiveMQArtemisScaledown struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisScaledownSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisScaledownStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ActiveMQArtemisScaledownList contains a list of ActiveMQArtemisScaledown
type ActiveMQArtemisScaledownList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisScaledown `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisScaledown{}, &ActiveMQArtemisScaledownList{})
}
