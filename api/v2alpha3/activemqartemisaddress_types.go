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

package v2alpha3

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
	QueueName *string `json:"queueName,omitempty"`
	// The Routing Type
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Routing Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	RoutingType *string `json:"routingType,omitempty"`
	// Whether or not delete the queue from broker when CR is undeployed(default false)
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Remove From Broker On Delete",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	RemoveFromBrokerOnDelete bool `json:"removeFromBrokerOnDelete,omitempty"`
	// User name for creating the queue or address
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="User",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	User *string `json:"user,omitempty"`
	// The password for the user
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:password"}
	Password *string `json:"password,omitempty"`
	// Specify the queue configuration
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Queue Configuration"
	QueueConfiguration *QueueConfigurationType `json:"queueConfiguration,omitempty"`
	// Apply to the broker crs in the current namespace. A value of * or empty string means applying to all broker crs. Default apply to all broker crs
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Apply To Broker CR Names"
	ApplyToCrNames []string `json:"applyToCrNames,omitempty"`
}

type QueueConfigurationType struct {
	// If ignore if the target queue already exists
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ignore If Exists",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	IgnoreIfExists *bool `json:"ignoreIfExists,omitempty"`
	// The routing type of the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Routing Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	RoutingType *string `json:"routingType,omitempty"`
	// The filter string for the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Filter String",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	FilterString *string `json:"filterString,omitempty"`
	// If the queue is durable or not
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Durable",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Durable *bool `json:"durable,omitempty"`
	// The user associated with the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="User",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	User *string `json:"user,omitempty"`
	// Max number of consumers allowed on this queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Max Consumers",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MaxConsumers *int32 `json:"maxConsumers,omitempty"`
	// If the queue is exclusive
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Exclusive",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Exclusive *bool `json:"exclusive,omitempty"`
	// If rebalance the message group
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Group Rebalance",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	GroupRebalance *bool `json:"groupRebalance,omitempty"`
	// If pause message dispatch when rebalancing groups
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Group Rebalance Pause Dispatch",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	GroupRebalancePauseDispatch *bool `json:"groupRebalancePauseDispatch,omitempty"`
	// Number of messaging group buckets
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Group Buckets",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	GroupBuckets *int32 `json:"groupBuckets,omitempty"`
	// Header set on the first group message
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Group First Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	GroupFirstKey *string `json:"groupFirstKey,omitempty"`
	// If it is a last value queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Last Value",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	LastValue *bool `json:"lastValue,omitempty"`
	// The property used for last value queue to identify last values
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Last Value Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	LastValueKey *string `json:"lastValueKey,omitempty"`
	// If force non-destructive consumers on the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Non Destructive",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	NonDestructive *bool `json:"nonDestructive,omitempty"`
	// Whether to delete all messages when no consumers connected to the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Purge On No Consumers",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	PurgeOnNoConsumers *bool `json:"purgeOnNoConsumers,omitempty"`
	// If the queue is enabled
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Enabled *bool `json:"enabled,omitempty"`
	// Number of consumers required before dispatching messages
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Consumers Before Dispatch",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ConsumersBeforeDispatch *int32 `json:"consumersBeforeDispatch,omitempty"`
	// Milliseconds to wait for `consumers-before-dispatch` to be met before dispatching messages anyway
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Delay Before Dispatch",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DelayBeforeDispatch *int64 `json:"delayBeforeDispatch,omitempty"`
	// Consumer Priority
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Consumer Priority",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ConsumerPriority *int32 `json:"consumerPriority,omitempty"`
	// Auto-delete the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDelete *bool `json:"autoDelete,omitempty"`
	// Delay (Milliseconds) before auto-delete the queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	AutoDeleteDelay *int64 `json:"autoDeleteDelay,omitempty"`
	// Message count of the queue to allow auto delete
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Message Count",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	AutoDeleteMessageCount *int64 `json:"autoDeleteMessageCount,omitempty"`
	// The size the queue should maintain according to ring semantics
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ring Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	RingSize *int64 `json:"ringSize,omitempty"`
	//  If the queue is configuration managed
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Configuration Managed",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	ConfigurationManaged *bool `json:"configurationManaged,omitempty"`
	// If the queue is temporary
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Temporary",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Temporary *bool `json:"temporary,omitempty"`
	// Whether auto create address
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Create Address",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateAddress *bool `json:"autoCreateAddress,omitempty"`
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
