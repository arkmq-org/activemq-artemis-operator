/*
Copyright 2024.

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

type ActiveMQArtemisAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Selector"
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Authentication Types"
	Auth []AppAuthType `json:"auth,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Messaging Capabilities"
	Capabilities []AppCapabilityType `json:"capabilities,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resources"
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type AppAddressType struct {
	// one of Name or QueueName is required
	Name      string `json:"name,omitempty"`
	QueueName string `json:"queuename,omitempty"`
}

type AppCapabilityType struct {
	// defaults to cr.Name
	Role string `json:"role,omitempty"`

	// only apps that split producer/consumer across roles will be able to shard for elastic queue
	ProducerOf []AppAddressType `json:"producerof,omitempty"` // ProducerTo

	ConsumerOf []AppAddressType `json:"consumerof,omitempty"` // ConsumerFrom/ConsumerOf

	ProducerAndConsumerOf []AppAddressType `json:"producerandconsumerof,omitempty"`
}

type ActiveMQArtemisAppStatus struct {
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
//+kubebuilder:resource:path=activemqartemisapps,shortName=aaapp
//+operator-sdk:csv:customresourcedefinitions:resources={{"Secret", "v1"}}

// Describes the messaging requirements of an application
// +operator-sdk:csv:customresourapplicationcedefinitions:displayName="ActiveMQ Artemis Messaging Application"
type ActiveMQArtemisApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisAppSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type ActiveMQArtemisAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisApp{}, &ActiveMQArtemisAppList{})
}
