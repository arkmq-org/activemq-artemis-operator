---
ai_metadata:
  doc_type: "defaults_reference"
  primary_audience: ["developers", "users", "troubleshooters"]
  key_topics: ["default_values", "configuration", "ports", "timeouts"]
  answers_questions:
    - "What are the default values?"
    - "What ports does the operator use?"
    - "What are the timeout defaults?"
  concepts_defined: 80
  code_links: 34
  complexity: "reference"
  related_docs: ["operator_conventions.md", "operator_architecture.md"]
---

> **Code Reference Notice**
> 
> Source code references in this document are accurate as of **October 9, 2025** (commit [`4acadb95`](https://github.com/arkmq-org/activemq-artemis-operator/commit/4acadb95603c38b82d5d7f63fb538c37ed855662)).
> 
> All code links use GitHub permalinks. Default values are generally stable, but to verify current defaults:
> - Replace `/blob/4acadb95...` with `/blob/main` in any link
> - Constant/variable names in link text help locate current values
> - Check the current code when default values are critical to your use case

This document provides a comprehensive reference of all default values used by the ActiveMQ Artemis Operator. Understanding these defaults is essential for configuring deployments and troubleshooting behavior.

> **Note**: All defaults are subject to override through CR fields, environment variables, or operator configuration. See the [operator conventions documentation](operator_conventions.md#configuration-precedence-conventions) for precedence rules.

> **Environment Variables**: All environment variables used by the operator default to empty/unset. For information about which environment variables are recognized and their usage patterns in different deployment modes, see the [environment variable patterns section](operator_conventions.md#environment-variable-patterns) in the conventions documentation.

**Table of Contents**

*   [API Spec Defaults](#api-spec-defaults)
*   [Deployment Configuration Defaults](#deployment-configuration-defaults)
*   [Container Image Defaults](#container-image-defaults)
*   [Port and Networking Defaults](#port-and-networking-defaults)
*   [Acceptor Configuration Defaults](#acceptor-configuration-defaults)
*   [Storage and Persistence Defaults](#storage-and-persistence-defaults)
*   [Path and Mount Defaults](#path-and-mount-defaults)
*   [Timeout and Retry Defaults](#timeout-and-retry-defaults)
*   [Certificate Management Defaults](#certificate-management-defaults)
*   [Security and RBAC Defaults](#security-and-rbac-defaults)
*   [Probe Configuration Defaults](#probe-configuration-defaults)
*   [Broker Configuration Defaults](#broker-configuration-defaults)
*   [Address Settings Defaults](#address-settings-defaults)
*   [Console Configuration Defaults](#console-configuration-defaults)
*   [Constants and Magic Values](#constants-and-magic-values)

## API Spec Defaults

All fields in the `ActiveMQArtemis` custom resource have implicit defaults when not specified.

### ActiveMQArtemisSpec Defaults

| Field | Default Value | Description | Source |
|-------|---------------|-------------|--------|
| `adminUser` | `""` (generated if empty) | Admin username for broker and console | Auto-generated in reconciler |
| `adminPassword` | `""` (generated if empty) | Admin password for broker and console | Auto-generated in reconciler |
| `deploymentPlan.size` | `1` | Number of broker replicas | [`common.go#L54`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L54) |
| `deploymentPlan.image` | Latest stable image | Broker container image | [`version.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/version/version.go) |
| `deploymentPlan.initImage` | Latest stable init image | Init container image | [`version.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/version/version.go) |
| `deploymentPlan.persistenceEnabled` | `false` | Ephemeral storage by default | Kubernetes default for `omitempty` |
| `deploymentPlan.messageMigration` | `true` | Enable safe message migration on scale-down | [`activemqartemis_reconciler.go#L101`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L101) |
| `deploymentPlan.resources` | `{}` (unlimited) | Resource requests/limits | Kubernetes default |
| `deploymentPlan.storage.size` | `"2Gi"` | PVC size when persistence enabled | Operator default |
| `deploymentPlan.storage.storageClassName` | `""` (cluster default) | Storage class for PVCs | Uses cluster default StorageClass |
| `deploymentPlan.livenessProbe.initialDelaySeconds` | `5` | Liveness probe startup delay | [`activemqartemis_reconciler.go#L73`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L73) |
| `version` | Latest stable version | Broker version to use | [`version.go#L19`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/version/version.go#L19) |
| `restricted` | `false` | Standard deployment mode | Kubernetes default for `*bool` |
| `addressSettings.applyRule` | `"merge_all"` | Address settings merge strategy | [`activemqartemis_reconciler.go#L110`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L110) |
| `console.expose` | `false` | Do not expose console by default | Kubernetes default |
| `brokerProperties` | `[]` | No additional broker properties | Empty array default |
| `acceptors` | `[]` (auto-generated) | Default acceptors created automatically | Generated in reconciler |
| `connectors` | `[]` | No connectors by default | Empty array default |

### DeploymentPlanType Defaults

| Field | Default Value | Description |
|-------|---------------|-------------|
| `image` | `""` (uses version-based resolution) | If empty, resolved from version |
| `initImage` | `""` (uses version-based resolution) | If empty, resolved from version |
| `size` | `1` | Single broker instance |
| `persistenceEnabled` | `false` | Ephemeral storage |
| `requireLogin` | `true` | Authentication required by broker |
| `messageMigration` | `true` | Safe scale-down enabled |
| `managementRBACEnabled` | `true` | RBAC enabled for management |
| `jolokiaAgentEnabled` | `true` | Jolokia agent enabled |
| `journalType` | `"nio"` | NIO journal by default |

### AddressSettingType Defaults

All address setting fields are pointers (`*int32`, `*string`, `*bool`) and default to `nil` (not set), which means the broker's internal defaults apply. When `nil`, the ActiveMQ Artemis broker uses its own default values from its configuration.

| Field | Default When Set | Broker Internal Default (when nil) |
|-------|------------------|-------------------------------------|
| `match` | Required field | N/A (must be specified) |
| `deadLetterAddress` | `nil` | Broker default (typically `DLQ`) |
| `expiryAddress` | `nil` | Broker default (typically `ExpiryQueue`) |
| `redeliveryDelay` | `nil` | Broker default (`0` ms) |
| `maxDeliveryAttempts` | `nil` | Broker default (`10`) |
| `maxSizeBytes` | `nil` | Broker default (`-1` = unlimited) |
| `pageSizeBytes` | `nil` | Broker default (`10485760` = 10MB) |
| `pageMaxCacheSize` | `nil` | Broker default (`5`) |
| `addressFullPolicy` | `nil` | Broker default (`PAGE`) |
| `autoCreateQueues` | `nil` | Broker default (`true`) |
| `autoDeleteQueues` | `nil` | Broker default (`false`) |
| `autoCreateAddresses` | `nil` | Broker default (`true`) |
| `autoDeleteAddresses` | `nil` | Broker default (`false`) |

## Deployment Configuration Defaults

### Replica and Sizing Defaults

```go
// Default deployment size
DefaultDeploymentSize = int32(1)  // Single broker instance

// Default message migration
defaultMessageMigration = true    // Safe scale-down enabled
```

**Source**: [`pkg/utils/common/common.go#L54`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L54), [`controllers/activemqartemis_reconciler.go#L101`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L101)

### Reconciliation Defaults

```go
// Controller resync period
DEFAULT_RESYNC_PERIOD = 30 * time.Second

// Resource discovery retry settings
defaultRetries = 10                     // Number of API server retry attempts
defaultRetryInterval = 3 * time.Second  // Wait period between retries
```

**Source**: [`pkg/utils/common/common.go#L57`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L57)

## Container Image Defaults

### Image Version Defaults

```go
// Latest supported Artemis version (as of last update)
LatestVersion = "2.42.0"

// Default container images
LatestKubeImage = "quay.io/arkmq-org/activemq-artemis-broker-kubernetes:artemis.2.42.0"
LatestInitImage = "quay.io/arkmq-org/activemq-artemis-broker-init:artemis.2.42.0"
```

**Source**: [`version/version.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/version/version.go)

### Image Selection Environment Variables

```bash
# Environment variable patterns (in precedence order):
DEFAULT_BROKER_VERSION              # Override version
DEFAULT_BROKER_KUBE_IMAGE          # Override broker image
DEFAULT_BROKER_INIT_IMAGE          # Override init image
RELATED_IMAGE_ActiveMQ_Artemis_Broker_Kubernetes_<version>_<arch>  # Architecture-specific
RELATED_IMAGE_ActiveMQ_Artemis_Broker_Kubernetes_<version>         # Generic
RELATED_IMAGE_ActiveMQ_Artemis_Broker_Init_<version>_<arch>       # Init arch-specific
RELATED_IMAGE_ActiveMQ_Artemis_Broker_Init_<version>              # Init generic
```

**Supported Architectures**: `amd64`, `arm64`, `s390x`, `ppc64le`

**Source**: [`pkg/utils/common/common.go#L51`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L51)

## Port and Networking Defaults

### Legacy Mode Ports (restricted: false)

```go
// Default protocol ports
61616  // "all" - All protocols (AMQP,CORE,HORNETQ,MQTT,OPENWIRE,STOMP)
8161   // "console-jolokia" - Web console and Jolokia management
8778   // "jolokia" - Dedicated Jolokia endpoint
7800   // "jgroups" - Cluster communication (JGroups)

// Protocol-specific ports
1883   // MQTT (plain)
8883   // MQTT (SSL)
5672   // AMQP (plain)
5671   // AMQP (SSL)
61613  // STOMP (plain)
61612  // STOMP (SSL)
61616  // CORE (always includes CORE for clustering)
```

**Source**: [`pkg/resources/serviceports/service_port.go#L12`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/serviceports/service_port.go#L12)

### Restricted Mode Ports (restricted: true)

```go
// No default ports exposed
// Only explicitly configured acceptors are created
// Metrics exposed on pod-local port 8888 (requires manual Service creation)
```

**Source**: [`pkg/resources/serviceports/service_port.go#L15`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/serviceports/service_port.go#L15)

### Service Naming Patterns

```go
// Service names follow convention
HeadlessService: "<cr-name>-hdls-svc"              // Always created
PingService: "<cr-name>-ping-svc"                  // Cluster discovery
AcceptorService: "<cr-name>-<acceptor>-<ordinal>-svc"  // Per-acceptor, per-pod

// OpenShift Routes
Route: "<service-name>-rte"

// Kubernetes Ingress
Ingress: "<service-name>-ing"
```

## Acceptor Configuration Defaults

### Auto-Generated Acceptor Defaults

```go
// Default acceptor parameters
defaultArgs = "tcpSendBufferSize=1048576;tcpReceiveBufferSize=1048576;useEpoll=true;amqpCredits=1000;amqpMinCredits=300"

// Default protocols when "all" is specified
defaultProtocols = "AMQP,CORE,HORNETQ,MQTT,OPENWIRE,STOMP"

// Auto-assigned port behavior
defaultPortIncrement = 10  // Ports start at 61626, increment by 10
```

**Conventions**:
- Port `61616` always includes CORE protocol for clustering
- Acceptors without explicit ports get auto-assigned starting at `61626`
- `bindAddress` defaults to `"ACCEPTOR_IP"` unless `bindToAllInterfaces=true`

**Source**: [`controllers/activemqartemis_reconciler.go#L750`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L750)

## Storage and Persistence Defaults

### Persistence Defaults

```go
// Storage disabled by default
persistenceEnabled: false

// When enabled:
storage.size: "2Gi"                    // Default PVC size
storage.accessMode: "ReadWriteOnce"    // Default access mode
storage.storageClassName: ""           // Uses cluster default StorageClass
```

**Conventions**:
- No storage class specified = use cluster default
- PVCs retained on scale-down when `messageMigration` enabled
- Each broker pod gets its own PVC in StatefulSet

**Source**: [`controllers/activemqartemis_reconciler.go#L3066`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3066)

> **Note**: For resource naming patterns (StatefulSet names, Service names, Secret names, labels), see the [operator conventions documentation](operator_conventions.md#resource-naming-conventions). These follow strict conventions rather than configurable defaults.

## Path and Mount Defaults

### Container Path Defaults

**Broker Configuration Paths**:
```bash
brokerConfigRoot = "/amq/init/config"              # XML configuration
initHelperScript = "/opt/amq-broker/script/default.sh"  # Startup script
configCmd = "/opt/amq/bin/launch.sh"               # Main broker command
```
**Source**: [`controllers/activemqartemis_reconciler.go#L105`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L105)

**Mount Base Paths**:
```bash
cfgMapPathBase = "/amq/extra/configmaps/"          # ConfigMap mounts
secretPathBase = "/amq/extra/secrets/"             # Secret mounts
```
**Source**: [`controllers/activemqartemis_reconciler.go#L79`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L79)

**Data Paths**:
```bash
GLOBAL_DATA_PATH = "/opt/<cr-name>/data"           # Broker data directory
```
**Source**: Generated dynamically in [`activemqartemis_controller.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L781) via `MakeNamers()`

**Specific Configuration File Paths** (derived from base paths):
```bash
BrokerProperties: "/amq/extra/secrets/<name>/broker.properties"
BrokerPropertiesOrdinal: "/amq/extra/secrets/<name>/broker-N/broker.properties"
JaasConfig: "/amq/extra/secrets/<name>/login.config"
LoggingConfig: "/amq/extra/configmaps/<name>/logging.properties"
```
These paths are constructed using the base paths above combined with specific filenames defined as constants ([`BrokerPropertiesName`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L87), [`JaasConfigKey`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L88), [`LoggingConfigKey`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L89)).

## Timeout and Retry Defaults

### API Server Interaction Defaults

```go
// Resource discovery and reconciliation
defaultRetries = 10                     // API server retry attempts
defaultRetryInterval = 3 * time.Second  // Retry wait period
DEFAULT_RESYNC_PERIOD = 30 * time.Second // Controller resync interval
```

**Source**: [`pkg/utils/common/common.go#L62`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L62)

## Certificate Management Defaults

### cert-manager Integration Defaults

```go
// Default certificate secret names
DefaultOperatorCertSecretName = "activemq-artemis-manager-cert"  // Operator client cert
DefaultOperatorCASecretName = "activemq-artemis-manager-ca"      // CA certificate
DefaultOperandCertSecretName = "broker-cert"                     // Broker server cert (or <cr-name>-broker-cert)
```

**Conventions**:
- All certificates use same CA for trust relationships
- Certificates auto-renewed by cert-manager
- mTLS required in restricted mode, optional in legacy mode

**Source**: [`pkg/utils/common/common.go#L69`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L69)

## Security and RBAC Defaults

### JAAS Configuration Defaults

#### JAAS Syntax Validation

The operator validates JAAS configuration syntax using a regex pattern that can be customized via environment variable:

```go
// Default JAAS syntax validation regex
JaasConfigSyntaxMatchRegExDefault = `^(?:(\s*|(?://.*)|(?s:/\*.*\*/))*\S+\s*{...`

// Customizable via environment variable
JAAS_CONFIG_SYNTAX_MATCH_REGEX  // Override the default validation pattern
```

**Source**: [`pkg/utils/common/common.go#L60`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L60), environment variable check at [`common.go#L115`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L115)

#### Default JAAS Login Module

> **Important**: The operator ONLY generates a default JAAS configuration when `restricted: true`. In legacy mode (`restricted: false`), users must provide their own JAAS configuration via a secret with the `-jaas-config` suffix if JAAS authentication is needed.

When `restricted: true`, the operator automatically generates this JAAS login configuration:

```java
// File: login.config
http_server_authenticator {
  org.apache.activemq.artemis.spi.core.security.jaas.TextFileCertificateLoginModule required
   reload=true
   debug=false
   org.apache.activemq.jaas.textfiledn.user=_cert-users
   org.apache.activemq.jaas.textfiledn.role=_cert-roles
   baseDir="/amq/extra/secrets/<cr-name>-props"
  ;
};
```

**Configuration Details**:
- **Login Module**: `TextFileCertificateLoginModule` for certificate-based authentication
- **Reload**: `true` - hot-reload of user/role files
- **Debug**: `false` - debug logging disabled by default
- **User File**: `_cert-users` (generated by operator)
- **Role File**: `_cert-roles` (generated by operator)
- **Base Directory**: Points to the broker properties secret mount

**Source**: [`controllers/activemqartemis_reconciler.go#L1981`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L1981) (restricted mode check), [`L1993`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L1993) (login configuration generation)

#### Legacy Mode JAAS Configuration

In legacy mode (`restricted: false` or not set), the operator does **not** generate any JAAS configuration. Instead:

- Users must provide their own `login.config` file via a Secret or ConfigMap with the **`-jaas-config` suffix**
- The operator [detects the `-jaas-config` mount](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2578) and automatically sets the JVM argument:
  ```bash
  -Djava.security.auth.login.config=/amq/extra/secrets/<name>/login.config
  ```
- If no `-jaas-config` secret/configmap is provided, JAAS authentication is not configured
- See the [suffix-based magic behavior](operator_conventions.md#suffix-based-magic-behavior) documentation for details

**Source**: [`controllers/activemqartemis_reconciler.go#L2578`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2578) (JAAS config detection), [`L2247`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2247) (JVM argument injection)

#### Default Certificate-to-User Mappings (Restricted Mode Only)

When `restricted: true`, the operator generates default certificate-to-user mappings:

```properties
# File: _cert-users (Java properties format with regex values)
hawtio=/CN = hawtio-online\.hawtio\.svc.*/
operator=/.*<operator-cert-cn>.*/
probe=/.*<broker-cert-cn>.*/
```

**Mapping Details**:
- **`hawtio` user**: Matches certificates with CN `hawtio-online.hawtio.svc.*`
- **`operator` user**: Matches certificates with the operator's client certificate CN (dynamically extracted)
- **`probe` user**: Matches certificates with the broker's server certificate CN (dynamically extracted)
- All patterns use regex syntax delimited by `/` (Java properties file format requirement)

**Source**: [`controllers/activemqartemis_reconciler.go#L2034`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2034)

#### Default Certificate-to-Role Mappings (Restricted Mode Only)

```properties
# File: _cert-roles (Java properties format)
status=operator,probe
metrics=operator
hawtio=hawtio
```

**Role Assignments**:
- **`status` role**: Assigned to `operator` and `probe` users (status check permissions)
- **`metrics` role**: Assigned to `operator` user only (metrics collection permissions)
- **`hawtio` role**: Assigned to `hawtio` user only (Hawtio console permissions)

**Source**: [`controllers/activemqartemis_reconciler.go#L2041`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2041)

### Management RBAC Defaults

```go
// Management RBAC enabled by default
managementRBACEnabled: true

// Applies to all deployment modes
// Controls access to management operations via JMX and Jolokia
```

**Source**: API spec default (boolean field defaults to `true` when `omitempty`)

## Probe Configuration Defaults

### Liveness Probe Defaults

```go
// Liveness probe settings
defaultLivenessProbeInitialDelay = 5  // 5 second startup delay
TCPLivenessPort = 8161                // Health check port
```

**Source**: [`controllers/activemqartemis_reconciler.go#L73`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L73)

## Broker Configuration Defaults

### Journal and Persistence Defaults

```go
// Default journal type
journalType: "nio"  // NIO journal (alternative: "aio" for better performance on Linux)
```

**Source**: [`controllers/activemqartemis_reconciler.go#L3133`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3133)

### Message Migration Defaults

```go
// Safe scale-down enabled by default
defaultMessageMigration = true

// PVC retention on scale-down
// When messageMigration is enabled, PVCs are retained for recovery
```

**Source**: [`controllers/activemqartemis_reconciler.go#L101`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L101)

## Address Settings Defaults

### Apply Rule Default

```go
// Default merge strategy for address settings
defApplyRule = "merge_all"  // Merge all address settings
```

**Source**: [`controllers/activemqartemis_reconciler.go#L110`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L110)

## Console Configuration Defaults

### Web Console Defaults

```go
// Console exposure
console.expose: false  // Not exposed by default

// Console SSLEnabled
console.sslEnabled: false  // Plain HTTP by default

// In restricted mode:
// Console is completely disabled (no Jetty, no web console)
```

**Sources**: 
- `console.expose` defaults to `false`: [`api/v1beta1/activemqartemis_types.go#L663`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/api/v1beta1/activemqartemis_types.go#L663) (boolean field defaults to `false`)
- `console.sslEnabled` defaults to `false`: [`api/v1beta1/activemqartemis_types.go#L669`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/api/v1beta1/activemqartemis_types.go#L669) (boolean field defaults to `false`)
- Console disabled in restricted mode: [`controllers/activemqartemis_reconciler.go#L459`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L459) (early return when `IsRestricted(customResource)`)

## Non-Defaults: Conventions and Magic Values

> **Note**: For magic values, special syntaxes, and constants (like the `-` removal value, file suffixes `-bp`/`-jaas-config`, prefixes `broker-`/`_`, file extensions `.properties`/`.json`, and resource naming postfixes `-svc`/`-rte`/`-ing`), see the [operator conventions documentation](operator_conventions.md#constants-and-magic-values). These are conventions with specific semantic meaning, not configurable defaults.

## Related Documentation

- **[Operator Architecture](operator_architecture.md)**: Technical implementation details
- **[Operator Conventions](operator_conventions.md)**: Conventions and patterns used
- **[Contribution Guide](contribution_guide.md)**: How to extend the operator
- **[User Documentation](../help/operator.md)**: End-user configuration guide
