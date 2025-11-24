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

type ActiveMQArtemisServiceSpec struct {

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resources"
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Environment"
	Env []corev1.EnvVar `json:"env,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auth"
	Auth []AppAuthType `json:"auth,omitempty"`

	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Acceptors"
	Acceptors []AppAcceptor `json:"acceptors,omitempty"`
}

type AppAuthType string

const (
	MTLS  AppAuthType = "mtls"
	Token AppAuthType = "token" // oath2
)

type AppAcceptor struct {
	Name string `json:"name"`
	//	Protocols []string `json:"protocols,omitempty"` // atribute of the matched service, defaults to capitalised name
	//	Port      *int32   `json:"port,omitempty"`      // optional inferred from protocol iana reserved port
	TlsSecret *string `json:"tlsSecret,omitempty"` // defaults to service broker-cert

	//MaxConnections *int32 `json:"makconnections,omitempty"` //sensible to bound, default 1k, but needs app to aquire chunks to allow sizing
}

type ActiveMQArtemisServiceStatus struct {
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
//+kubebuilder:resource:path=activemqartemisservices,shortName=aasvc
//+operator-sdk:csv:customresourcedefinitions:resources={{"Secret", "v1"}}
//+operator-sdk:csv:customresourcedefinitions:resources={{"Service", "v1"}}
//+operator-sdk:csv:customresourcedefinitions:resources={{"ActiveMQArtemis", "v1beta1"}}

// Provides a broker service
// +operator-sdk:csv:customresourcedefinitions:displayName="ActiveMQ Artemis Service"
type ActiveMQArtemisService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisServiceSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type ActiveMQArtemisServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisService{}, &ActiveMQArtemisServiceList{})
}
