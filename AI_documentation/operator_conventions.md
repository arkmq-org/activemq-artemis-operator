---
ai_metadata:
  doc_type: "conventions_reference"
  primary_audience: ["developers", "users"]
  key_topics: ["naming", "defaults", "magic_behavior", "configuration_precedence"]
  answers_questions:
    - "What are the naming conventions?"
    - "What defaults does the operator use?"
    - "How does convention over configuration work?"
  concepts_defined: 45
  code_links: 44
  complexity: "reference"
  related_docs: ["operator_defaults_reference.md", "operator_architecture.md"]
---

> **Code Reference Notice**
> 
> Source code references in this document are accurate as of **October 9, 2025** (commit [`4acadb95`](https://github.com/arkmq-org/activemq-artemis-operator/commit/4acadb95603c38b82d5d7f63fb538c37ed855662)).
> 
> All code links use GitHub permalinks. Conventions are generally stable, but to verify current behavior:
> - Replace `/blob/4acadb95...` with `/blob/main` in any link
> - Function/constant names in link text help locate current implementation
> - Core conventions rarely change, but check current code for critical decisions

This document provides a comprehensive reference of all conventions, defaults, and automatic behaviors implemented by the ActiveMQ Artemis Operator. Understanding these conventions is essential for both using and extending the operator.

**Table of Contents**

*   [Convention Over Configuration Philosophy](#convention-over-configuration-philosophy)
*   [Resource Naming Conventions](#resource-naming-conventions)
*   [Platform Detection and Adaptation Conventions](#platform-detection-and-adaptation-conventions)
*   [Persistence and Storage Conventions](#persistence-and-storage-conventions)
*   [Service Discovery and Networking Conventions](#service-discovery-and-networking-conventions)
*   [Configuration Precedence Conventions](#configuration-precedence-conventions)
*   [Zero-Configuration Deployment Convention](#zero-configuration-deployment-convention)
*   [Suffix-Based Magic Behavior](#suffix-based-magic-behavior)
*   [Mount Path Conventions](#mount-path-conventions)
*   [Environment Variable Patterns](#environment-variable-patterns)
*   [Default Values and Configurations](#default-values-and-configurations)
*   [Broker Version Dependencies](#broker-version-dependencies)
*   [Default Retry and Timeout Conventions](#default-retry-and-timeout-conventions)
*   [Certificate Management Conventions](#certificate-management-conventions)
*   [Annotation and Label Conventions](#annotation-and-label-conventions)
*   [Default Security Conventions](#default-security-conventions)
*   [Constants and Magic Values](#constants-and-magic-values)
*   [Regex Patterns](#regex-patterns)
*   [Hash-Based Resource Management](#hash-based-resource-management)
*   [Volume Mount Ordering](#volume-mount-ordering)

## Convention Over Configuration Philosophy

The ActiveMQ Artemis Operator embodies the **"Convention Over Configuration"** philosophy, providing sensible defaults for everything while allowing customization when needed.

> **Philosophy**: The operator should work out-of-the-box with minimal configuration, following established patterns and conventions that users can rely on without explicit configuration.

**🧪 Tested in**: 
- [`broker_name_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/broker_name_test.go) - Naming convention tests
- [`common_util_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/common_util_test.go) - Utility function and convention testing
- [`suite_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/suite_test.go) - Platform detection and environment setup testing

## Resource Naming Conventions

The operator [follows strict naming patterns](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/namer/namer.go) that are [implemented through](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L175) the namer package:

### **StatefulSet and Pod Naming**
```go
// Pattern: <cr-name>-ss follows convention
StatefulSet: "ex-aao-ss"            // Generated via CrToSS()
Pods: "ex-aao-ss-0", "ex-aao-ss-1"  // Kubernetes StatefulSet convention
```
[Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/namer/namer.go#L65): `CrToSS()` function  
[Usage](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L783): Applied in `MakeNamers()`

### **Service Naming**
```go
// Services follow <cr-name>-<type>-<suffix> pattern
HeadlessService: "ex-aao-hdls-svc"  // For StatefulSet DNS
PingService: "ex-aao-ping-svc"      // For cluster discovery
AcceptorService: "ex-aao-artemis-svc" // Per-acceptor exposure
```
[Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L785): Services are [named in MakeNamers()](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L786)

### **Secret Naming**
```go
// Secrets follow <cr-name>-<purpose>-secret pattern
BrokerProperties: "ex-aao-props"           // Properties storage
Credentials: "ex-aao-credentials-secret"   // Auth credentials
Console: "ex-aao-console-secret"           // Console TLS
Netty: "ex-aao-netty-secret"              // Netty TLS
```
[Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2782): Properties naming via `getPropertiesResourceNsName()`  
[Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L788): Other secrets [named in MakeNamers()](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L794)

### **Label Conventions**
```go
// Standard labels applied to all resources
Labels: {
    "application": "<cr-name>-app",     // Application grouping
    "ActiveMQArtemis": "<cr-name>"      // Resource association
}
```
[Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/selectors/label.go#L35): Labels are [generated](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/selectors/label.go#L4) with [standard keys](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/selectors/label.go#L5)

## Platform Detection and Adaptation Conventions

The operator [automatically detects and adapts](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/main.go#L248) to different Kubernetes platforms:

### **OpenShift vs Kubernetes Detection**
```go
// Auto-detection via API server inspection
isOpenshift = DetectOpenshiftWith(restConfig)

// Platform-specific behavior:
// OpenShift: Uses Routes for external exposure
// Kubernetes: Uses Ingress for external exposure  
// OpenShift: Detects FIPS mode from cluster config
// Kubernetes: Assumes standard compliance
```
[Detection Logic](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go)

### **Architecture Support**
```go
// Multi-architecture image selection
// Convention: Try arch-specific image first, fall back to generic
RELATED_IMAGE_ActiveMQ_Artemis_Broker_Kubernetes_<version>_<arch>
RELATED_IMAGE_ActiveMQ_Artemis_Broker_Kubernetes_<version>

// Supported architectures: amd64, arm64, s390x, ppc64le
```
[Image Selection Logic](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L327): Architecture-specific images are [tried first](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L330), then [falls back to generic](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L335)

## Persistence and Storage Conventions

```go
// Default storage behavior
persistenceEnabled: false             // Ephemeral by default
storage.size: "2Gi"                  // Default PVC size when enabled
storage.accessMode: "ReadWriteOnce"   // Default access mode

// Convention: No storage class specified = use cluster default
// Convention: PVCs retained on scale-down when messageMigration enabled
// Convention: Each broker pod gets its own PVC in StatefulSet
```
[Storage Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3066): Persistence is [enabled when configured](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3066), [PVC templates are created](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3061) for each broker

## Service Discovery and Networking Conventions

```go
// Service naming patterns
HeadlessService: "<cr-name>-hdls-svc"           // Always created
PingService: "<cr-name>-ping-svc"               // Cluster discovery
AcceptorService: "<cr-name>-<acceptor>-<ordinal>-svc"  // Per-acceptor, per-pod

// Route naming (OpenShift)
Route: "<service-name>-rte"                     // Route per service
// Convention: Uses passthrough TLS termination for SSL acceptors

// Ingress naming (Kubernetes)  
Ingress: "<service-name>-ing"                   // Ingress per service
// Convention: Requires ingress controller with SSL passthrough support
```
[Service Creation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/services/service.go#L11): Services are [created with standard patterns](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/services/service.go#L38)  
[Route Creation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/routes/route.go#L10): Routes [use passthrough TLS](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/routes/route.go#L41) when [SSL is enabled](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/routes/route.go#L41)  
[Ingress Creation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/ingresses/ingress.go#L10): Ingresses are [created with SSL support](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/ingresses/ingress.go#L75)

## Configuration Precedence Conventions

The operator [follows strict precedence rules](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L1952) for configuration sources. This section demonstrates how the operator consistently prioritizes **explicit user configuration** over **environment-based overrides** over **defaults** across different aspects of the system.

Two key examples illustrate this precedence philosophy:

**1. Image Selection Precedence** - How the operator chooses which container image to use for the broker (infrastructure concern):

```go
// Precedence order (highest to lowest):
1. Explicit CR field values (e.g., spec.deploymentPlan.image)
2. Environment variable overrides (e.g., DEFAULT_BROKER_KUBE_IMAGE)  
3. Architecture-specific related images (RELATED_IMAGE_*_<arch>)
4. Generic related images (RELATED_IMAGE_*)
5. Built-in defaults (LatestKubeImage, LatestInitImage)
```

**2. Broker Properties Precedence** - How the operator merges broker configuration properties (application concern):

```go
// Precedence order (highest to lowest):
1. CR spec.brokerProperties (applied first)
2. External -bp secrets (alphabetical by secret name)
3. Files within -bp secrets (alphabetical by key name)
4. Operator-injected properties (restricted mode, metrics, etc.)
```
[Image Selection](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L317): Images are [determined using precedence](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L330) in `DetermineImageToUse()`  
[Properties Precedence](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2569): Properties are [applied in order](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2572) during JVM configuration

## Zero-Configuration Deployment Convention

A minimal CR demonstrates the convention over configuration philosophy:

```yaml
# Minimal working configuration - everything else uses conventions
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemis
metadata:
  name: my-broker
spec: {}  # Everything else is defaulted!

# This creates:
# - 1 broker pod (DefaultDeploymentSize)
# - All standard ports (61616, 8161, 8778, 7800)
# - Headless and ping services
# - Latest stable broker image
# - Ephemeral storage
# - Standard labels and naming
# - CORE protocol enabled for clustering
```

## Suffix-Based Magic Behavior

The operator uses suffix-based conventions for automatic behavior:

### **ExtraMount Secret Suffixes**
- **`-bp`**: [Auto-detected as broker properties source](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2997) when secrets are [processed](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2997)
- **`-jaas-config`**: [JAAS configuration](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L561) (must be Secret, not ConfigMap) - [validated during setup](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L561)
- **`-logging-config`**: [Logging configuration](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L558) (must be ConfigMap) - [processed during validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller.go#L558)

### **File Extension Processing**
- **`.properties`**: [Java properties format](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L85) (default handling)
- **`.json`**: [JSON format for `-bp` secrets](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L86) (broker parses internally)
- **`.config`**: [JAAS login configuration](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L88) files

## Mount Path Conventions

```bash
# Base paths (defined as constants)
ConfigMaps: "/amq/extra/configmaps/<name>/"  # cfgMapPathBase at line 79
Secrets: "/amq/extra/secrets/<name>/"        # secretPathBase at line 80

# Specific paths
BrokerProperties: "/amq/extra/secrets/<name>/broker.properties"
JaasConfig: "/amq/extra/secrets/<name>/login.config"           # JaasConfigKey
LoggingConfig: "/amq/extra/configmaps/<name>/logging.properties" # LoggingConfigKey
```
[Path Usage](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2980): Base paths are [used to construct mount points](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3020) for all configuration files

## Environment Variable Patterns

The operator recognizes different environment variables based on deployment mode and operational context. All environment variables default to empty/unset unless explicitly configured.

### **JVM Configuration Variables (Mode-Dependent)**

The operator [uses different environment variables](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2539) based on deployment mode:

**Legacy Mode** (`restricted: false`):
- **JAVA_ARGS_APPEND**: [Used for additional JVM arguments](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L95) in [legacy deployments](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2542)
- **JAVA_OPTS**: [Standard Java options](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L97) for JVM configuration
- **DEBUG_ARGS**: Debug configuration options

**Restricted Mode** (`restricted: true`):
- **JDK_JAVA_OPTIONS**: [Modern JVM options](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L98) (preferred) [used in restricted mode](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2545)
- **STATEFUL_SET_ORDINAL**: [Extracted from hostname](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2502) for [ordinal-specific config](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2566)

### **Operator Behavior Variables**

These environment variables control operator-level behavior:

- **OPERATOR_WATCH_NAMESPACE**: Limits the operator to watch specific namespace(s). Empty = all namespaces (cluster-wide)
- **OPERATOR_NAMESPACE**: The namespace where the operator itself is deployed
- **OPERATOR_OPENSHIFT**: Override automatic OpenShift detection (`"true"` or `"false"`)
- **OPERATOR_LOG_LEVEL**: Set operator logging level (`debug`, `info`, `warn`, `error`)

### **Image Selection Variables**

See [Configuration Precedence Conventions](#configuration-precedence-conventions) for the complete list of image-related environment variables (`DEFAULT_BROKER_VERSION`, `DEFAULT_BROKER_KUBE_IMAGE`, `RELATED_IMAGE_*`, etc.)

## Constants and Magic Values

Important constants that [drive operator behavior](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L72):

```go
// Suffixes (lines 75-77)
brokerPropsSuffix = "-bp"
jaasConfigSuffix = "-jaas-config"
loggingConfigSuffix = "-logging-config"

// Prefixes and Separators (lines 82-84)
OrdinalPrefix = "broker-"
OrdinalPrefixSep = "."
UncheckedPrefix = "_"

// File Extensions (lines 85-87)
PropertiesSuffix = ".properties"
JsonSuffix = ".json"
BrokerPropertiesName = "broker.properties"

// Special Values (line 94)
RemoveKeySpecialValue = "-"  // Used to remove/disable features
```
[Constant Usage](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go): These constants are [used throughout the reconciler](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2997) to [maintain consistency](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3031)

## Regex Patterns

### **Ordinal Property Parsing**
```go
// Matches: broker-0.property=value, broker-123.another=value
Pattern: ^(broker-[0-9]+)\.(.*)$
```
[Pattern Implementation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3033): Regex is [used to parse](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3031) ordinal-specific properties

### **JAAS Config Validation**
```go
// Complex regex for JAAS syntax validation (customizable via env var)
JAAS_CONFIG_SYNTAX_MATCH_REGEX: ^(?:(\s*|(?://.*)|(?s:/\*.*\*/))*\S+\s*{...
```
[JAAS Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L60): Regex [validates JAAS syntax](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L122) and is [customizable via environment](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L115)

## Hash-Based Resource Management

The operator [uses Adler-32 checksums](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2810) for efficient change detection:

```go
// Legacy immutable ConfigMaps with hash in name
ConfigMapName: "<cr-name>-props-<adler32-hash>"

// Modern mutable Secrets with fixed name  
SecretName: "<cr-name>-props"
```
[Hash Generation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2764): Checksums are [computed for change detection](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2768) to [trigger updates](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L188)

## Volume Mount Ordering

Properties are [loaded in strict order](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2569):
1. [CR `brokerProperties`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2563) (global first, then ordinal-specific)
2. [`-bp` secrets](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2570) ([alphabetical by secret name](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L2572))
3. [Within each `-bp` secret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3008): [alphabetical by key name](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L3008)

## Default Values Reference

> **Complete Default Values Documentation**: For a comprehensive reference of all default values (deployment sizing, container images, ports, paths, timeouts, JAAS configuration, RBAC, probes, etc.), see the **[Operator Defaults Reference](operator_defaults_reference.md)**.

The operator follows a "convention over configuration" philosophy, providing sensible defaults for:
- Single-broker deployment by default
- NIO journal type
- Message migration enabled for safe scale-down
- Management RBAC enabled
- Jolokia agent enabled
- Standard protocol ports in legacy mode
- No exposed ports in restricted mode (explicit configuration required)
- 30-second reconciliation period
- 5-second liveness probe delay

All defaults can be overridden through CR specification, environment variables, or operator configuration according to [configuration precedence conventions](#configuration-precedence-conventions).

## Broker Version Dependencies

- **Ordinal-specific properties**: Requires broker >= 2.27.1
- **Directory-based property loading**: Requires broker >= 2.27.1
- **JSON format support**: Available in all supported versions

## Certificate Management Conventions

The operator integrates with `cert-manager` for automated certificate lifecycle management:

**Certificate Names**: See [Certificate Management Defaults](operator_defaults_reference.md#certificate-management-defaults) for default secret names

**Key Conventions**:
- All certificates use the same CA for trust relationships
- Certificates are auto-renewed by cert-manager before expiration
- mTLS is **required** in restricted mode
- mTLS is **optional** in legacy mode (can use plain-text acceptors)
- Certificates are mounted as secrets in the broker pods

**Source**: [`pkg/utils/common/common.go#L69`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L69)

## Annotation and Label Conventions

```go
// Control annotations
BlockReconcileAnnotation = "arkmq.org/block-reconcile"  // Pause reconciliation

// Standard Kubernetes labels
"statefulset.kubernetes.io/pod-name"  // Pod identification
"app.kubernetes.io/name"             // Application name
"app.kubernetes.io/instance"         // Instance identifier
"app.kubernetes.io/version"          // Version tracking
"app.kubernetes.io/component"        // Component type
"app.kubernetes.io/part-of"          // Application suite
"app.kubernetes.io/managed-by"       // Management tool

// Operator-specific labels (../../pkg/utils/selectors/label.go)
"application": "<cr-name>-app"        // Application grouping
"ActiveMQArtemis": "<cr-name>"        // Resource association
```

## Security Conventions

**JAAS Configuration**: See [JAAS Configuration Defaults](operator_defaults_reference.md#jaas-configuration-defaults) for complete JAAS defaults, including:
- Syntax validation regex (customizable via environment variable)
- Automatic JAAS generation in restricted mode
- User-provided JAAS in legacy mode

**Address Settings**: Default apply rule is `merge_all` for merging address setting configurations

**Key Security Conventions**:
- Restricted mode **requires** mTLS for all communication (broker, management, metrics)
- Legacy mode **allows** plain-text connections for development/testing
- Management RBAC always enabled by default
- Operator has full access to all broker management functions
- In restricted mode: no web console, no init container, certificate-based authentication only

**Source**: [`pkg/utils/common/common.go#L60`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common.go#L60) (JAAS validation), [`controllers/activemqartemis_reconciler.go#L110`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler.go#L110) (apply rule)

