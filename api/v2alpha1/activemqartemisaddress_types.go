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

package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisAddressSpec defines the desired state of ActiveMQArtemisAddress
type ActiveMQArtemisAddressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The Address Name
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Address Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	AddressName string `json:"addressName,omitempty"`
	// The Queue Name
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Queue Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	QueueName string `json:"queueName,omitempty"`
	// The Routing Type
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Routing Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	RoutingType string `json:"routingType,omitempty"`
}

// ActiveMQArtemisAddressStatus defines the observed state of ActiveMQArtemisAddress
type ActiveMQArtemisAddressStatus struct {
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
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=activemqartemisaddresses,shortName=aaa
//+operator-sdk:csv:customresourcedefinitions:resources={{"Secret", "v1"}}

// +kubebuilder:deprecatedversion:warning="The ActiveMQArtemisAddress CRD is deprecated. Use the spec.brokerProperties attribute in the ActiveMQArtemis CR to create addresses and queues instead"
// Adding and removing addresses using custom resource definitions
// +operator-sdk:csv:customresourcedefinitions:displayName="ActiveMQ Artemis Address"
type ActiveMQArtemisAddress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisAddressSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisAddressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ActiveMQArtemisAddressList contains a list of ActiveMQArtemisAddress
type ActiveMQArtemisAddressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisAddress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisAddress{}, &ActiveMQArtemisAddressList{})
}
