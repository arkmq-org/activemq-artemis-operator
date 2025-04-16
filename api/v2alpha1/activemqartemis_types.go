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
	"github.com/RHsyseng/operator-utils/pkg/olm"
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
}

type DeploymentPlanType struct {
	//The image used for the broker, all upgrades are disabled. Needs a corresponding initImage
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Image",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Image string `json:"image,omitempty"`
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

// ActiveMQArtemisStatus defines the observed state of ActiveMQArtemis
type ActiveMQArtemisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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
