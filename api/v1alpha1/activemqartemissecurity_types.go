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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ActiveMQArtemisSecuritySpec defines the desired state of ActiveMQArtemisSecurity
type ActiveMQArtemisSecuritySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Specifies the login modules (deprecated in favour of ActiveMQArtemisSpec.DeploymentPlan.ExtraMounts.Secrets -jaas-config)
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Login Modules"
	LoginModules LoginModulesType `json:"loginModules,omitempty"`
	// Specifies the security domains (deprecated in favour of ActiveMQArtemisSpec.DeploymentPlan.ExtraMounts.Secrets -jaas-config)
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Security Domains"
	SecurityDomains SecurityDomainsType `json:"securityDomains,omitempty"`
	// Specifies the security settings
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Security Settings"
	SecuritySettings SecuritySettingsType `json:"securitySettings,omitempty"`
	// Apply this security config to the broker crs in the current namespace. A value of * or empty string means applying to all broker crs. Default apply to all broker crs
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Apply to Broker CR Names"
	ApplyToCrNames []string `json:"applyToCrNames,omitempty"`
}

type LoginModulesType struct {
	// Specifies the properties login modules
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Properties Login Modules"
	PropertiesLoginModules []PropertiesLoginModuleType `json:"propertiesLoginModules,omitempty"`
	// Specifies the guest login modules
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Guest Login Modules"
	GuestLoginModules []GuestLoginModuleType `json:"guestLoginModules,omitempty"`
	// Specifies the Keycloak login modules
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Keycloak Login Modules"
	KeycloakLoginModules []KeycloakLoginModuleType `json:"keycloakLoginModules,omitempty"`
}

type PropertiesLoginModuleType struct {
	// Name for PropertiesLoginModule
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name string `json:"name,omitempty"`
	// Specifies the users
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Users"
	Users []UserType `json:"users,omitempty"`
}

type UserType struct {
	// User name to be defined in properties login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name string `json:"name,omitempty"`
	// Password to be defined in properties login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:password"}
	Password *string `json:"password,omitempty"`
	// Roles to be defined in properties login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Roles"
	Roles []string `json:"roles,omitempty"`
}

type GuestLoginModuleType struct {
	// Name for GuestLoginModule
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name string `json:"name,omitempty"`
	// The guest user name
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Guest User",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	GuestUser *string `json:"guestUser,omitempty"`
	// The guest user role
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Guest Role",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	GuestRole *string `json:"guestRole,omitempty"`
}

type KeycloakLoginModuleType struct {
	// Name for KeycloakLoginModule
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name string `json:"name,omitempty"`
	// Type of KeycloakLoginModule directAccess or bearerToken
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Module Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ModuleType *string `json:"moduleType,omitempty"`
	// Specifies the Keycloak module configuration
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Keycloa Module Configuration"
	Configuration KeycloakModuleConfigurationType `json:"configuration,omitempty"`
}

type KeycloakModuleConfigurationType struct {
	// Realm for KeycloakLoginModule
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Realm",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Realm *string `json:"realm,omitempty"`
	// Public key for the realm
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Realm PublicKey",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	RealmPublicKey *string `json:"realmPublicKey,omitempty"`
	// URL of the keycloak authentication server
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auth Server Url",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	AuthServerUrl *string `json:"authServerUrl,omitempty"`
	// How SSL is required
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="SSL Required",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	SslRequired *string `json:"sslRequired,omitempty"`
	// Resource Name
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resource",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Resource *string `json:"resource,omitempty"`
	// If it is public client
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Public Client",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	PublicClient *bool `json:"publicClient,omitempty"`
	// Specify the credentials
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Credentials"
	Credentials []KeyValueType `json:"credentials,omitempty"`
	// If to use resource role mappings
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Use Resource Role Mappings",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	UseResourceRoleMappings *bool `json:"useResourceRoleMappings,omitempty"`
	// If to enable CORS
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enable Cors",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	EnableCors *bool `json:"enableCors,omitempty"`
	// CORS max age
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Cors Max Age",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	CorsMaxAge *int64 `json:"corsMaxAge,omitempty"`
	// CORS allowed methods
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Cors Allowed Methods",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	CorsAllowedMethods *string `json:"corsAllowedMethods,omitempty"`
	// CORS allowed headers
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Cors Allowed Headers",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	CorsAllowedHeaders *string `json:"corsAllowedHeaders,omitempty"`
	// CORS exposed headers
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Cors Exposed Headers",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	CorsExposedHeaders *string `json:"corsExposedHeaders,omitempty"`
	// If to expose access token
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expose Token",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	ExposeToken *bool `json:"exposeToken,omitempty"`
	// If only verify bearer token
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Bearer Only",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	BearerOnly *bool `json:"bearerOnly,omitempty"`
	// If auto-detect bearer token only
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Auto Detect Bearer Only",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AutoDetectBearerOnly *bool `json:"autoDetectBearerOnly,omitempty"`
	// Size of the connection pool
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Connection Pool Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ConnectionPoolSize *int64 `json:"connectionPoolSize,omitempty"`
	// If to allow any host name
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Allow Any Host Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AllowAnyHostName *bool `json:"allowAnyHostName,omitempty"`
	// If to disable trust manager
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Disable Trust Manager",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	DisableTrustManager *bool `json:"disableTrustManager,omitempty"`
	// Path of a trust store
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="TrustStore",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	TrustStore *string `json:"trustStore,omitempty"`
	// Truststore password
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="TrustStore Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:password"}
	TrustStorePassword *string `json:"trustStorePassword,omitempty"`
	// Path of a client keystore
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Client KeyStore",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ClientKeyStore *string `json:"clientKeyStore,omitempty"`
	// Client keystore password
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Client KeyStore Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:password"}
	ClientKeyStorePassword *string `json:"clientKeyStorePassword,omitempty"`
	// Client key password
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Client Key Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:password"}
	ClientKeyPassword *string `json:"clientKeyPassword,omitempty"`
	// If always refresh token
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Always Refresh Token",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	AlwaysRefreshToken *bool `json:"alwaysRefreshToken,omitempty"`
	// If register node at startup
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Register Node At Startup",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	RegisterNodeAtStartup *bool `json:"registerNodeAtStartup,omitempty"`
	// Period for re-registering node
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Register Node Period",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	RegisterNodePeriod *int64 `json:"registerNodePeriod,omitempty"`
	// Type of token store. session or cookie
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Token Store",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	TokenStore *string `json:"tokenStore,omitempty"`
	// Cookie path for a cookie store
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Token Cookie Path",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	TokenCookiePath *string `json:"tokenCookiePath,omitempty"`
	// OpenID Connect ID Token attribute to populate the UserPrincipal name with
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Principal Attribute",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	PrincipalAttribute *string `json:"principalAttribute,omitempty"`
	// The proxy URL
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Proxy Url",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ProxyUrl *string `json:"proxyUrl,omitempty"`
	// If not to change session id on a successful login
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Turn Off Change SessionId On Login",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	TurnOffChangeSessionIdOnLogin *bool `json:"turnOffChangeSessionIdOnLogin,omitempty"`
	// Minimum time to refresh an active access token
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Token Minimum Time To Live",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	TokenMinimumTimeToLive *int64 `json:"tokenMinimumTimeToLive,omitempty"`
	// Minimum interval between two requests to Keycloak to retrieve new public keys
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Min Time Between Jwks Requests",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	MinTimeBetweenJwksRequests *int64 `json:"minTimeBetweenJwksRequests,omitempty"`
	// Maximum interval between two requests to Keycloak to retrieve new public keys
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Public Key Cache Ttl",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	PublicKeyCacheTtl *int64 `json:"publicKeyCacheTtl,omitempty"`
	// Whether to turn off processing of the access_token query parameter for bearer token processing
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ignore Oauth Query Parameter",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	IgnoreOauthQueryParameter *bool `json:"ignoreOauthQueryParameter,omitempty"`
	// Verify whether the token contains this client name (resource) as an audience
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Verify Token Audience",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	VerifyTokenAudience *bool `json:"verifyTokenAudience,omitempty"`
	// Whether to support basic authentication
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enable Basic Auth",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	EnableBasicAuth *bool `json:"enableBasicAuth,omitempty"`
	// The confidential port used by the Keycloak server for secure connections over SSL/TLS
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Confidential Port",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	ConfidentialPort *int32 `json:"confidentialPort,omitempty"`
	// Specify the redirect rewrite rules
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Redirect Rewrite Rules"
	RedirectRewriteRules []KeyValueType `json:"redirectRewriteRules,omitempty"`
	// The OAuth2 scope parameter for DirectAccessGrantsLoginModule
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Scope",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Scope *string `json:"scope,omitempty"`
}

type KeyValueType struct {
	// The regular expression to match the Redirect URI
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Key string `json:"key,omitempty"`
	// The replacement value
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Value",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Value *string `json:"value,omitempty"`
}

type SecurityDomainsType struct {
	// Specify the broker domain
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Broker Domain"
	BrokerDomain BrokerDomainType `json:"brokerDomain,omitempty"`
	// Specify the console domain
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Console Domain"
	ConsoleDomain BrokerDomainType `json:"consoleDomain,omitempty"`
}

type BrokerDomainType struct {
	// Name for the broker/console domain
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name *string `json:"name,omitempty"`
	// Specify the login modules
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Login Modules"
	LoginModules []LoginModuleReferenceType `json:"loginModules,omitempty"`
}

type LoginModuleReferenceType struct {
	// Name of the login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Name *string `json:"name,omitempty"`
	// Flag of the login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Flag",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Flag *string `json:"flag,omitempty"`
	// Debug option of the login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Debug",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Debug *bool `json:"debug,omitempty"`
	// Reload option of the login module
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Reload",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Reload *bool `json:"reload,omitempty"`
}

type SecuritySettingsType struct {
	// Specify the broker security settings
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Broker Security Settings"
	Broker []BrokerSecuritySettingType `json:"broker,omitempty"`
	// Specify the management security settings
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Management Security Settings"
	Management ManagementSecuritySettingsType `json:"management,omitempty"`
}

type BrokerSecuritySettingType struct {
	// The address match pattern of a security setting
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Match",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Match string `json:"match,omitempty"`
	// Specify the permissions
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Permissions"
	Permissions []PermissionType `json:"permissions,omitempty"`
}

type PermissionType struct {
	// The operation type of a security setting
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Operation Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	OperationType string `json:"operationType"`
	// The roles of a security setting
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Roles"
	Roles []string `json:"roles,omitempty"`
}

type ManagementSecuritySettingsType struct {
	// The roles allowed to login hawtio
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Hawtio Roles"
	HawtioRoles []string `json:"hawtioRoles,omitempty"`
	// Specify connector configurations
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Connector Configurations"
	Connector ConnectorConfigType `json:"connector,omitempty"`
	// Specify the authorisation configurations
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Authorisation Configurations"
	Authorisation AuthorisationConfigType `json:"authorisation,omitempty"`
}

type ConnectorConfigType struct {
	// The connector host for connecting to management
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Host",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Host *string `json:"host,omitempty"`
	// The connector port for connecting to management
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Port",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	Port *int32 `json:"port,omitempty"`
	// The RMI registry port for management
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Rmi Registry Port",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:number"}
	RmiRegistryPort *int32 `json:"rmiRegistryPort,omitempty"`
	// The JMX realm of management
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Jmx Realm",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	JmxRealm *string `json:"jmxRealm,omitempty"`
	// The JMX object name of management
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Object Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	ObjectName *string `json:"objectName,omitempty"`
	// The management authentication type
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Authenticator Type",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	AuthenticatorType *string `json:"authenticatorType,omitempty"`
	// Whether management connection is secured
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Secured",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:booleanSwitch"}
	Secured *bool `json:"secured,omitempty"`
	// The keystore provider for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="KeyStore Provider",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	KeyStoreProvider *string `json:"keyStoreProvider,omitempty"`
	// The keystore path for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="KeyStore Path",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	KeyStorePath *string `json:"keyStorePath,omitempty"`
	// The keystore password for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="KeyStore Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	KeyStorePassword *string `json:"keyStorePassword,omitempty"`
	// The truststore provider for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="TrustStore Provider",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	TrustStoreProvider *string `json:"trustStoreProvider,omitempty"`
	// The truststore path for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="TrustStore Path",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	TrustStorePath *string `json:"trustStorePath,omitempty"`
	// The truststore password for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="TrustStore Password",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	TrustStorePassword *string `json:"trustStorePassword,omitempty"`
	// The password codec for management connector
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Password Codec",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	PasswordCodec *string `json:"passwordCodec,omitempty"`
}

type AuthorisationConfigType struct {
	// Specify the allowed entries
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Allowed Entries"
	AllowedList []AllowedListEntryType `json:"allowedList,omitempty"`
	// Specify the default accesses
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Accesses"
	DefaultAccess []DefaultAccessType `json:"defaultAccess,omitempty"`
	// Specify the role accesses
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Role Accesses"
	RoleAccess []RoleAccessType `json:"roleAccess,omitempty"`
}

type AllowedListEntryType struct {
	// The domain of allowedList
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Domain",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Domain *string `json:"domain,omitempty"`
	// The key of allowedList
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Key *string `json:"key,omitempty"`
}

type DefaultAccessType struct {
	// Specifies the access entry method
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Method",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Method *string `json:"method,omitempty"`
	// Specifies the access entry roles
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Roles"
	Roles []string `json:"roles,omitempty"`
}

type RoleAccessType struct {
	// The domain of the role access
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Domain",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Domain *string `json:"domain,omitempty"`
	// The key of the role access
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Key",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:text"}
	Key *string `json:"key,omitempty"`
	// Specify the default accesses
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Default Accesses"
	AccessList []DefaultAccessType `json:"accessList,omitempty"`
}

// ActiveMQArtemisSecurityStatus defines the observed state of ActiveMQArtemisSecurity
type ActiveMQArtemisSecurityStatus struct {
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
//+kubebuilder:resource:path=activemqartemissecurities,shortName=aas
//+operator-sdk:csv:customresourcedefinitions:resources={{"Secret", "v1"}}

// +kubebuilder:deprecatedversion:warning="The ActiveMQArtemisSecurity CRD is deprecated. Use the spec.brokerProperties attribute in the ActiveMQArtemis CR and -jaas-config extraMount to configure security instead"
// Security configuration for the broker
// +operator-sdk:csv:customresourcedefinitions:displayName="ArkMQ Security"
type ActiveMQArtemisSecurity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ActiveMQArtemisSecuritySpec `json:"spec,omitempty"`
	// Specifies the security status modules
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="ArkMQ Security Status"
	Status ActiveMQArtemisSecurityStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ActiveMQArtemisSecurityList contains a list of ActiveMQArtemisSecurity
type ActiveMQArtemisSecurityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ActiveMQArtemisSecurity `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ActiveMQArtemisSecurity{}, &ActiveMQArtemisSecurityList{})
}
