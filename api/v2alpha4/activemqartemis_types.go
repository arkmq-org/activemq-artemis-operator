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

package v2alpha4

import (
	"github.com/RHsyseng/operator-utils/pkg/olm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisSpec defines the desired state of ActiveMQArtemis
type ActiveMQArtemisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// User name for standard broker user. It is required for connecting to the broker and the web console. If left empty, it will be generated.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Admin User",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	AdminUser string `json:"adminUser,omitempty"`
	// Password for standard broker user. It is required for connecting to the broker and the web console. If left empty, it will be generated.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Admin Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:password"}
	AdminPassword string `json:"adminPassword,omitempty"`
	// Specifies the deployment plan
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Deployment Plan"
	DeploymentPlan DeploymentPlanType `json:"deploymentPlan,omitempty"`
	// Specifies the acceptor configuration
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Acceptors"
	Acceptors []AcceptorType `json:"acceptors,omitempty"`
	// Specifies connectors and connector configuration
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Connectors"
	Connectors []ConnectorType `json:"connectors,omitempty"`
	// Specifies the console configuration
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Console Configurations"
	Console ConsoleType `json:"console,omitempty"`
	// The desired version of the broker. Can be x, or x.y or x.y.z to configure upgrades
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Version",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Version string `json:"version,omitempty"`
	// Specifies the upgrades (deprecated in favour of Version)
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Upgrades"
	Upgrades ActiveMQArtemisUpgrades `json:"upgrades,omitempty"`
	// Specifies the address configurations (deprecated in favour of BrokerProperties)
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Address Configurations"
	AddressSettings AddressSettingsType `json:"addressSettings,omitempty"`
}

type AddressSettingsType struct {
	// How to merge the address settings to broker configuration
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Apply Rule",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ApplyRule *string `json:"applyRule,omitempty"`
	// Specifies the address settings
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Address Settings"
	AddressSetting []AddressSettingType `json:"addressSetting,omitempty"`
}

type AddressSettingType struct {
	// the address to send dead messages to
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dead Letter Address",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DeadLetterAddress *string `json:"deadLetterAddress,omitempty"`
	// whether or not to automatically create the dead-letter-address and/or a corresponding queue on that address when a message found to be undeliverable
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="AutoCreateDeadLetterResources",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateDeadLetterResources *bool `json:"autoCreateDeadLetterResources,omitempty"`
	// the prefix to use for auto-created dead letter queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dead Letter Queue Prefix",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DeadLetterQueuePrefix *string `json:"deadLetterQueuePrefix,omitempty"`
	// the suffix to use for auto-created dead letter queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dead Letter Queue Suffix",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DeadLetterQueueSuffix *string `json:"deadLetterQueueSuffix,omitempty"`
	// the address to send expired messages to
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expiry Address",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ExpiryAddress *string `json:"expiryAddress,omitempty"`
	// whether or not to automatically create the expiry-address and/or a corresponding queue on that address when a message is sent to a matching queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Create Expiry Resources",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateExpiryResources *bool `json:"autoCreateExpiryResources,omitempty"`
	// the prefix to use for auto-created expiry queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expiry Queue Prefix",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ExpiryQueuePrefix *string `json:"expiryQueuePrefix,omitempty"`
	// the suffix to use for auto-created expiry queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expiry Queue Suffix",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ExpiryQueueSuffix *string `json:"expiryQueueSuffix,omitempty"`
	// Overrides the expiration time for messages using the default value for expiration time. "-1" disables this setting.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expiry Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ExpiryDelay *int32 `json:"expiryDelay,omitempty"`
	// Overrides the expiration time for messages using a lower value. "-1" disables this setting.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Min Expiry Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MinExpiryDelay *int32 `json:"minExpiryDelay,omitempty"`
	// Overrides the expiration time for messages using a higher value. "-1" disables this setting.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Max Expiry Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MaxExpiryDelay *int32 `json:"maxExpiryDelay,omitempty"`
	// the time (in ms) to wait before redelivering a cancelled message.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Redelivery Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	RedeliveryDelay           *int32 `json:"redeliveryDelay,omitempty"`
	RedeliveryDelayMultiplier *int32 `json:"redeliveryDelayMultiplier,omitempty"`
	//
	RedeliveryCollisionAvoidanceFactor *float32 `json:"redeliveryCollisionAvoidanceFactor,omitempty"` // controller-gen requires crd:allowDangerousTypes=true to allow support float types
	//
	// Maximum value for the redelivery-delay
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Max Redelivery Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MaxRedeliveryDelay *int32 `json:"maxRedeliveryDelay,omitempty"`
	// how many times to attempt to deliver a message before sending to dead letter address
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Max Delivery Attempts",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MaxDeliveryAttempts *int32 `json:"maxDeliveryAttempts,omitempty"`
	// the maximum size in bytes for an address. -1 means no limits. This is used in PAGING, BLOCK and FAIL policies. Supports byte notation like K, Mb, GB, etc.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Max Size Bytes",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	MaxSizeBytes *string `json:"maxSizeBytes,omitempty"`
	// used with the address full BLOCK policy, the maximum size in bytes an address can reach before messages start getting rejected. Works in combination with max-size-bytes for AMQP protocol only.  Default = -1 (no limit).
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Max Size Bytes Reject Threshold",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MaxSizeBytesRejectThreshold *int32 `json:"maxSizeBytesRejectThreshold,omitempty"`
	// The page size in bytes to use for an address. Supports byte notation like K, Mb, GB, etc.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Page Size Bytes",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	PageSizeBytes *string `json:"pageSizeBytes,omitempty"`
	// Number of paging files to cache in memory to avoid IO during paging navigation
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Page Max Cache Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	PageMaxCacheSize *int32 `json:"pageMaxCacheSize,omitempty"`
	// what happens when an address where maxSizeBytes is specified becomes full
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Address Full Policy",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	AddressFullPolicy *string `json:"addressFullPolicy,omitempty"`
	// how many days to keep message counter history for this address
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Message Counter History Day Limit",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MessageCounterHistoryDayLimit *int32 `json:"messageCounterHistoryDayLimit,omitempty"`
	// This is deprecated please use default-last-value-queue instead.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Last Value Queue",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	LastValueQueue *bool `json:"lastValueQueue,omitempty"`
	// whether to treat the queues under the address as a last value queues by default
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Last Value Queue",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DefaultLastValueQueue *bool `json:"defaultLastValueQueue,omitempty"`
	// the property to use as the key for a last value queue by default
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Last Value Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DefaultLastValueKey *string `json:"defaultLastValueKey,omitempty"`
	// whether the queue should be non-destructive by default
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Non Destructive",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DefaultNonDestructive *bool `json:"defaultNonDestructive,omitempty"`
	// whether to treat the queues under the address as exclusive queues by default
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Exclusive Queue",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DefaultExclusiveQueue *bool `json:"defaultExclusiveQueue,omitempty"`
	// whether to rebalance groups when a consumer is added
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Group Rebalance",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DefaultGroupRebalance *bool `json:"defaultGroupRebalance,omitempty"`
	// whether to pause dispatch when rebalancing groups
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Group Rebalance Pause Dispatch",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DefaultGroupRebalancePauseDispatch *bool `json:"defaultGroupRebalancePauseDispatch,omitempty"`
	// number of buckets to use for grouping, -1 (default) is unlimited and uses the raw group, 0 disables message groups.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Group Buckets",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DefaultGroupBuckets *int32 `json:"defaultGroupBuckets,omitempty"`
	// key used to mark a message is first in a group for a consumer
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Group First Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DefaultGroupFirstKey *string `json:"defaultGroupFirstKey,omitempty"`
	// the default number of consumers needed before dispatch can start for queues under the address.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Consumers Before Dispatch",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DefaultConsumersBeforeDispatch *int32 `json:"defaultConsumersBeforeDispatch,omitempty"`
	// the default delay (in milliseconds) to wait before dispatching if number of consumers before dispatch is not met for queues under the address.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Delay Before Dispatch",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DefaultDelayBeforeDispatch *int32 `json:"defaultDelayBeforeDispatch,omitempty"`
	// how long (in ms) to wait after the last consumer is closed on a queue before redistributing messages.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Redistribution Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	RedistributionDelay *int32 `json:"redistributionDelay,omitempty"`
	// if there are no queues matching this address, whether to forward message to DLA (if it exists for this address)
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Send To DLA On No Route",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	SendToDlaOnNoRoute *bool `json:"sendToDlaOnNoRoute,omitempty"`
	// The minimum rate of message consumption allowed before a consumer is considered "slow." Measured in messages-per-second.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Slow Consumer Threshold",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	SlowConsumerThreshold *int32 `json:"slowConsumerThreshold,omitempty"`
	// what happens when a slow consumer is identified
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Slow Consumer Policy",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SlowConsumerPolicy *string `json:"slowConsumerPolicy,omitempty"`
	// How often to check for slow consumers on a particular queue. Measured in seconds.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Slow Consumer Check Period",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	SlowConsumerCheckPeriod *int32 `json:"slowConsumerCheckPeriod,omitempty"`
	// DEPRECATED. whether or not to automatically create JMS queues when a producer sends or a consumer connects to a queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Create Jms Queues",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateJmsQueues *bool `json:"autoCreateJmsQueues,omitempty"`
	// DEPRECATED. whether or not to delete auto-created JMS queues when the queue has 0 consumers and 0 messages
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Jms Queues",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDeleteJmsQueues *bool `json:"autoDeleteJmsQueues,omitempty"`
	// DEPRECATED. whether or not to automatically create JMS topics when a producer sends or a consumer subscribes to a topic
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Create Jms Topics",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateJmsTopics *bool `json:"autoCreateJmsTopics,omitempty"`
	// DEPRECATED. whether or not to delete auto-created JMS topics when the last subscription is closed
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Jms Topics",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDeleteJmsTopics *bool `json:"autoDeleteJmsTopics,omitempty"`
	// whether or not to automatically create a queue when a client sends a message to or attempts to consume a message from a queue
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Create Queues",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateQueues *bool `json:"autoCreateQueues,omitempty"`
	// whether or not to delete auto-created queues when the queue has 0 consumers and 0 messages
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Queues",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDeleteQueues *bool `json:"autoDeleteQueues,omitempty"`
	// whether or not to delete created queues when the queue has 0 consumers and 0 messages
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Created Queues",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDeleteCreatedQueues *bool `json:"autoDeleteCreatedQueues,omitempty"`
	// how long to wait (in milliseconds) before deleting auto-created queues after the queue has 0 consumers.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Queues Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	AutoDeleteQueuesDelay *int32 `json:"autoDeleteQueuesDelay,omitempty"`
	// the message count the queue must be at or below before it can be evaluated to be auto deleted, 0 waits until empty queue (default) and -1 disables this check.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Queues Message Count",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	AutoDeleteQueuesMessageCount *int32 `json:"autoDeleteQueuesMessageCount,omitempty"`
	//What to do when a queue is no longer in broker.xml.  OFF = will do nothing queues will remain, FORCE = delete queues even if messages remaining.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Config Delete Queues",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ConfigDeleteQueues *string `json:"configDeleteQueues,omitempty"`
	// whether or not to automatically create addresses when a client sends a message to or attempts to consume a message from a queue mapped to an address that doesnt exist
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Create Addresses",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoCreateAddresses *bool `json:"autoCreateAddresses,omitempty"`
	// whether or not to delete auto-created addresses when it no longer has any queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Addresses",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDeleteAddresses *bool `json:"autoDeleteAddresses,omitempty"`
	// how long to wait (in milliseconds) before deleting auto-created addresses after they no longer have any queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Delete Addresses Delay",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	AutoDeleteAddressesDelay *int32 `json:"autoDeleteAddressesDelay,omitempty"`
	// What to do when an address is no longer in broker.xml.  OFF = will do nothing addresses will remain, FORCE = delete address and its queues even if messages remaining.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Config Delete Addresses",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ConfigDeleteAddresses *string `json:"configDeleteAddresses,omitempty"`
	// how many message a management resource can browse
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Management Browse Page Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ManagementBrowsePageSize *int32 `json:"managementBrowsePageSize,omitempty"`
	// purge the contents of the queue once there are no consumers
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Purge On No Consumers",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DefaultPurgeOnNoConsumers *bool `json:"defaultPurgeOnNoConsumers,omitempty"`
	// the maximum number of consumers allowed on this queue at any one time
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Max Consumers",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DefaultMaxConsumers *int32 `json:"defaultMaxConsumers,omitempty"`
	// the routing-type used on auto-created queues
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Queue Routing Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DefaultQueueRoutingType *string `json:"defaultQueueRoutingType,omitempty"`
	// the routing-type used on auto-created addresses
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Address Routing Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	DefaultAddressRoutingType *string `json:"defaultAddressRoutingType,omitempty"`
	// the default window size for a consumer
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Consumer Window Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DefaultConsumerWindowSize *int32 `json:"defaultConsumerWindowSize,omitempty"`
	// the default ring-size value for any matching queue which doesnt have ring-size explicitly defined
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Ring Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	DefaultRingSize *int32 `json:"defaultRingSize,omitempty"`
	// the number of messages to preserve for future queues created on the matching address
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Retroactive Message Count",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	RetroactiveMessageCount *int32 `json:"retroactiveMessageCount,omitempty"`
	// whether or not to enable metrics for metrics plugins on the matching address
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enable Metrics",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	EnableMetrics *bool `json:"enableMetrics,omitempty"`
	// pattern for matching settings against addresses; can use wildards
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Match",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Match string `json:"match,omitempty"`
}

type DeploymentPlanType struct {
	//The image used for the broker, all upgrades are disabled. Needs a corresponding initImage
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Image string `json:"image,omitempty"`
	// The init container image used to configure broker, all upgrades are disabled. Needs a corresponding image
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Init Image",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	InitImage string `json:"initImage,omitempty"`
	// The number of broker pods to deploy
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:podCount"}
	Size int32 `json:"size,omitempty"`
	// If true require user password login credentials for broker protocol ports
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Require Login",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	RequireLogin bool `json:"requireLogin,omitempty"`
	// If true use persistent volume via persistent volume claim for journal storage
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Persistence Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	PersistenceEnabled bool `json:"persistenceEnabled,omitempty"`
	// If aio use ASYNCIO, if nio use NIO for journal IO
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Journal Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	JournalType string `json:"journalType,omitempty"`
	//If true migrate messages on scaledown
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Message Migration",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	MessageMigration *bool `json:"messageMigration,omitempty"`
	// Specifies the minimum/maximum amount of compute resources required/allowed
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resource Requirements",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:resourceRequirements"}
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Specifies the storage configurations
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage Configurations"
	Storage StorageType `json:"storage,omitempty"`
	// If true enable the Jolokia JVM Agent
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Jolokia Agent Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	JolokiaAgentEnabled bool `json:"jolokiaAgentEnabled,omitempty"`
	// If true enable the management role based access control
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Management RBAC Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	ManagementRBACEnabled bool `json:"managementRBACEnabled,omitempty"`
}

type StorageType struct {
	// The storage size
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Size string `json:"size,omitempty"`
}

type AcceptorType struct {
	// The acceptor name
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name string `json:"name"`
	// Port number
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Port",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	Port int32 `json:"port,omitempty"`
	// The protocols to enable for this acceptor
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Protocols",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Protocols string `json:"protocols,omitempty"`
	// Whether or not to enable SSL on this port
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	SSLEnabled bool `json:"sslEnabled,omitempty"`
	// Name of the secret to use for ssl information
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Secret",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SSLSecret string `json:"sslSecret,omitempty"`
	// Comma separated list of cipher suites used for SSL communication.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled Cipher Suites",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	EnabledCipherSuites string `json:"enabledCipherSuites,omitempty"`
	// Comma separated list of protocols used for SSL communication.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled Protocols",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	EnabledProtocols string `json:"enabledProtocols,omitempty"`
	// Tells a client connecting to this acceptor that 2-way SSL is required. This property takes precedence over wantClientAuth.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Need Client Auth",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	NeedClientAuth bool `json:"needClientAuth,omitempty"`
	// Tells a client connecting to this acceptor that 2-way SSL is requested but not required. Overridden by needClientAuth.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Want Client Auth",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	WantClientAuth bool `json:"wantClientAuth,omitempty"`
	// The CN of the connecting client's SSL certificate will be compared to its hostname to verify they match. This is useful only for 2-way SSL.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Verify Host",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	VerifyHost bool `json:"verifyHost,omitempty"`
	// Used to change the SSL Provider between JDK and OPENSSL. The default is JDK.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Provider",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SSLProvider string `json:"sslProvider,omitempty"`
	// A regular expression used to match the server_name extension on incoming SSL connections. If the name doesn't match then the connection to the acceptor will be rejected.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SNI Host",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SNIHost string `json:"sniHost,omitempty"`
	// Whether or not to expose this acceptor
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expose",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Expose bool `json:"expose,omitempty"`
	// To indicate which kind of routing type to use.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Anycast Prefix",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	AnycastPrefix string `json:"anycastPrefix,omitempty"`
	// To indicate which kind of routing type to use
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Multicast Prefix",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	MulticastPrefix string `json:"multicastPrefix,omitempty"`
	// Max number of connections allowed to make
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Connections Allowed",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ConnectionsAllowed int `json:"connectionsAllowed,omitempty"`
	// AMQP Minimum Large Message Size
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="AMQP Min Large Message Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	AMQPMinLargeMessageSize int `json:"amqpMinLargeMessageSize,omitempty"`
}

type ConnectorType struct {
	// The name of the connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name string `json:"name"`
	// The type either tcp or vm
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Type string `json:"type,omitempty"`
	// Hostname or IP to connect to
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Host",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Host string `json:"host"`
	// Port number
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Port",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	Port int32 `json:"port"`
	//  Whether or not to enable SSL on this port
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	SSLEnabled bool `json:"sslEnabled,omitempty"`
	// Name of the secret to use for ssl information
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Secret",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SSLSecret string `json:"sslSecret,omitempty"`
	// Comma separated list of cipher suites used for SSL communication.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled Cipher Suites",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	EnabledCipherSuites string `json:"enabledCipherSuites,omitempty"`
	// Comma separated list of protocols used for SSL communication.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled Protocols",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	EnabledProtocols string `json:"enabledProtocols,omitempty"`
	// Tells a client connecting to this connector that 2-way SSL is required. This property takes precedence over wantClientAuth.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Need Client Auth",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	NeedClientAuth bool `json:"needClientAuth,omitempty"`
	// Tells a client connecting to this connector that 2-way SSL is requested but not required. Overridden by needClientAuth.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Want Client Auth",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	WantClientAuth bool `json:"wantClientAuth,omitempty"`
	// The CN of the connecting client's SSL certificate will be compared to its hostname to verify they match. This is useful only for 2-way SSL.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Verify Host",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	VerifyHost bool `json:"verifyHost,omitempty"`
	// Used to change the SSL Provider between JDK and OPENSSL. The default is JDK.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Provider",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SSLProvider string `json:"sslProvider,omitempty"`
	// A regular expression used to match the server_name extension on incoming SSL connections. If the name doesn't match then the connection to the acceptor will be rejected.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SNI Host",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SNIHost string `json:"sniHost,omitempty"`
	// Whether or not to expose this connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expose",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Expose bool `json:"expose,omitempty"`
}

type ConsoleType struct {
	// Whether or not to expose this port
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expose",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Expose bool `json:"expose,omitempty"`
	// Whether or not to enable SSL on this port
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Enabled",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	SSLEnabled bool `json:"sslEnabled,omitempty"`
	// Name of the secret to use for ssl information
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Secret",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SSLSecret string `json:"sslSecret,omitempty"`
	// If the embedded server requires client authentication
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Use Client Auth",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	UseClientAuth bool `json:"useClientAuth,omitempty"`
}

// ActiveMQArtemis App product upgrade flags
type ActiveMQArtemisUpgrades struct {
	// Set true to enable automatic micro version product upgrades, it is disabled by default.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enable Upgrades",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:ui:booleanSwitch"}
	Enabled bool `json:"enabled"`
	// Set true to enable automatic minor product version upgrades, it is disabled by default. Requires spec.upgrades.enabled to be true.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Include minor version upgrades",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:fieldDependency:upgrades.enabled:true","urn:alm:descriptor:com.tectonic.ui:ui:booleanSwitch"}
	Minor bool `json:"minor"`
}

// ActiveMQArtemisStatus defines the observed state of ActiveMQArtemis
type ActiveMQArtemisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The current pods
	//+operator-sdk:csv:customresourcedefinitions:type=status,displayName="Pods Status",xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses"
	PodStatus olm.DeploymentStatus `json:"podStatus"`

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
//+kubebuilder:resource:path=activemqartemises
//+kubebuilder:resource:path=activemqartemises,shortName=aa
//+operator-sdk:csv:customresourcedefinitions:resources={{"Secret", "v1"}}

// A stateful deployment of one or more brokers
// +operator-sdk:csv:customresourcedefinitions:displayName="ActiveMQ Artemis"
type ActiveMQArtemis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ActiveMQArtemisSpec   `json:"spec,omitempty"`
	Status ActiveMQArtemisStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ActiveMQArtemisList contains a list of ActiveMQArtemis
type ActiveMQArtemisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemis{}, &ActiveMQArtemisList{})
}
