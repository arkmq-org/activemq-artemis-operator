This guide provides practical information for developers who want to contribute to the ActiveMQ Artemis Operator by adding new features, debugging issues, or understanding the development workflow.

> **Note**: This guide is validated against actual project history. The patterns, examples, and workflows described here reflect real development practices used in this project.

**Table of Contents**

*   [Feature Development Guide](#feature-development-guide)
    *   [Adding New API Fields](#adding-new-api-fields)
    *   [Extending the Reconciler](#extending-the-reconciler)
    *   [Creating New Resource Types](#creating-new-resource-types)
    *   [Testing New Features](#testing-new-features)
    *   [Development Workflow](#development-workflow)
        *   [Prerequisites and Initial Setup](#1-prerequisites-and-initial-setup)
        *   [Development Cycle](#2-development-cycle-tdd-approach)
        *   [Local Testing Options](#3-local-testing-options)
        *   [Required File Changes Pattern](#4-required-file-changes-pattern)
        *   [Commit Pattern](#5-commit-pattern)
        *   [Testing Strategy](#6-testing-strategy)
        *   [Common Integration Points](#7-common-integration-points)
        *   [Documentation Requirements](#8-documentation-requirements)
*   [Development and Debugging Guide](#development-and-debugging-guide)
    *   [Running Tests](#running-tests)
    *   [Debugging Reconciliation Issues](#debugging-reconciliation-issues)
    *   [Code Navigation Tips](#code-navigation-tips)

## Feature Development Guide

This section provides practical guidance for developers who want to add new functionality to the operator.

> **Philosophy**: The operator follows **strict patterns and conventions**. New features must integrate seamlessly with existing architecture while maintaining backward compatibility and comprehensive test coverage.

### Adding New API Fields

Follow this established pattern for adding new API fields:

> **Reference**: See the [Operator SDK Go tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/) for comprehensive guidance on defining APIs and [kubebuilder markers reference](https://book.kubebuilder.io/reference/markers.html) for annotation details.

#### **1. Add Field to API Type**

Add your field to api/v1beta1/activemqartemis_types.go::ActiveMQArtemisSpec with proper annotations:

**Example** - The `Restricted` field (api/v1beta1/activemqartemis_types.go::Restricted):

```go
// Restricted deployment, mtls jolokia agent with RBAC
//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Restricted"
Restricted *bool `json:"restricted,omitempty"`
```

**Pattern**:
```go
// MyNewFeature enables new functionality
//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="My New Feature"
MyNewFeature *bool `json:"myNewFeature,omitempty"`
```

**Required Annotations**:
- `//+operator-sdk:csv:customresourcedefinitions:type=spec` - OLM integration
- `displayName` - Human-readable name for UI
- `omitempty` - Makes field optional (convention)

#### **2. Add Validation Logic**

Add validation function to `controllers/activemqartemis_controller.go` validation chain:

> **Reference**: For webhook-based validation, see the [Operator SDK Admission Webhooks guide](https://sdk.operatorframework.io/docs/building-operators/golang/webhook/).

```go
func validateMyNewFeature(customResource *brokerv1beta1.ActiveMQArtemis) (*metav1.Condition, bool) {
    if customResource.Spec.MyNewFeature != nil && *customResource.Spec.MyNewFeature {
        // Add your validation logic here
        if /* validation fails */ {
            return &metav1.Condition{
                Type:    brokerv1beta1.ValidConditionType,
                Status:  metav1.ConditionFalse,
                Reason:  "ValidConditionMyNewFeatureInvalid",
                Message: "MyNewFeature validation failed: reason",
            }, false
        }
    }
    return nil, false
}
```

Integrate into validation chain in controllers/activemqartemis_controller.go::validate():

```go
if validationCondition.Status != metav1.ConditionFalse {
    condition, retry = validateMyNewFeature(customResource)
    if condition != nil {
        validationCondition = *condition
    }
}
```

#### **3. Regenerate Code**

```bash
# Regenerate CRDs and DeepCopy methods
make manifests generate
```

This updates:
- `config/crd/bases/broker.amq.io_activemqartemises.yaml`
- `api/v1beta1/zz_generated.deepcopy.go`
- Bundle manifests for OLM

> **Reference**: See the [Operator SDK Go tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/) for details on CRD generation workflow and the [controller-gen documentation](https://book.kubebuilder.io/reference/controller-gen.html) for generation options.

### Extending the Reconciler

Follow this pattern to add new processing logic:

> **Reference**: See the [Operator SDK Go tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/) for best practices on implementing reconciliation logic and the [Operator SDK Advanced Topics](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/) for complex scenarios.

#### **1. Create ProcessX Function**

Add new process function to `controllers/activemqartemis_reconciler.go` following the established signature pattern:

**Example** - The `ProcessConsole()` function (controllers/activemqartemis_reconciler.go::ProcessConsole):

```go
func (reconciler *ActiveMQArtemisReconcilerImpl) ProcessConsole(
    customResource *brokerv1beta1.ActiveMQArtemis,
    namer common.Namers,
    client rtclient.Client,
    scheme *runtime.Scheme,
    currentStatefulSet *appsv1.StatefulSet) error {

    reqLogger := reconciler.log.WithName("ProcessConsole")

    if customResource.Spec.Console.Expose {
        // Console exposure logic
        labels := namer.LabelBuilder.Labels()
        serviceName := namer.SvcHeadlessNameBuilder.Name()

        // Create service for console
        consoleService := svc.NewServiceDefinitionForCR(...)
        reconciler.trackDesired(consoleService)
    }

    return nil
}
```

**Pattern**:
```go
func (reconciler *ActiveMQArtemisReconcilerImpl) ProcessMyNewFeature(
    customResource *brokerv1beta1.ActiveMQArtemis,
    namer common.Namers,
    client rtclient.Client,
    scheme *runtime.Scheme,
    currentStatefulSet *appsv1.StatefulSet) error {

    reqLogger := reconciler.log.WithName("ProcessMyNewFeature")

    if customResource.Spec.MyNewFeature != nil && *customResource.Spec.MyNewFeature {
        // Your feature logic here

        // Create any needed resources
        myResource := mybuilders.MakeMyResource(...)
        reconciler.trackDesired(myResource)

        // Modify StatefulSet if needed
        // currentStatefulSet.Spec.Template.Spec.Containers[0].Env = append(...)
    }

    return nil
}
```

#### **2. Integrate into Main Flow**

Add call in controllers/activemqartemis_reconciler.go::Process() function:

```go
err = reconciler.ProcessConsole(customResource, namer, client, scheme, desiredStatefulSet)
if err != nil {
    reconciler.log.Error(err, "Error processing console")
    return err
}

// Add your new process call here
err = reconciler.ProcessMyNewFeature(customResource, namer, client, scheme, desiredStatefulSet)
if err != nil {
    reconciler.log.Error(err, "Error processing my new feature")
    return err
}
```

#### **3. Handle Feature Detection**

Follow the pkg/utils/common/common.go for feature detection:

**Example** - The `IsRestricted()` function (pkg/utils/common/common.go):

```go
func IsRestricted(customResource *brokerv1beta1.ActiveMQArtemis) bool {
    return customResource.Spec.Restricted != nil && *customResource.Spec.Restricted
}
```

**Pattern**:
```go
// In pkg/utils/common/common.go
func IsMyNewFeatureEnabled(customResource *brokerv1beta1.ActiveMQArtemis) bool {
    return customResource.Spec.MyNewFeature != nil && *customResource.Spec.MyNewFeature
}
```

### Creating New Resource Types

Follow the established resource builder patterns:

#### **1. Create Resource Builder Package**

Create `pkg/resources/mynewresource/mynewresource.go`:

```go
package mynewresource

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    // Import the Kubernetes resource type you're creating
)

func MakeMyNewResource(existing *MyResourceType, namespacedName types.NamespacedName,
                      labels map[string]string, customResource *brokerv1beta1.ActiveMQArtemis) *MyResourceType {

    if existing == nil {
        existing = &MyResourceType{
            TypeMeta: metav1.TypeMeta{
                Kind:       "MyResourceType",
                APIVersion: "...",
            },
            ObjectMeta: metav1.ObjectMeta{},
            Spec:       MyResourceTypeSpec{},
        }
    }

    // Apply desired state
    existing.ObjectMeta.Labels = labels
    existing.ObjectMeta.Name = namespacedName.Name
    existing.ObjectMeta.Namespace = namespacedName.Namespace

    // Configure spec based on customResource
    // existing.Spec.Field = customResource.Spec.MyNewFeature

    return existing
}
```

#### **2. Add Custom Comparator (if needed)**

If your resource needs special comparison logic, add to controllers/activemqartemis_reconciler.go::ProcessResources():

```go
comparator.Comparator.SetComparator(reflect.TypeOf(MyResourceType{}), reconciler.CompareMyResource)
```

#### **3. Integration Pattern**

Use in your `ProcessX()` function:

```go
myResourceName := types.NamespacedName{
    Name:      namer.MyResourceNameBuilder.Name(),
    Namespace: customResource.Namespace,
}

var myResourceDefinition *MyResourceType
obj := reconciler.cloneOfDeployed(reflect.TypeOf(MyResourceType{}), myResourceName.Name)
if obj != nil {
    myResourceDefinition = obj.(*MyResourceType)
}

myResourceDefinition = mynewresource.MakeMyNewResource(myResourceDefinition, myResourceName, labels, customResource)
reconciler.trackDesired(myResourceDefinition)
```

### Testing New Features

Follow the established testing patterns:

> **Reference**: See the [Operator SDK Testing guide](https://sdk.operatorframework.io/docs/building-operators/golang/testing/) for comprehensive testing strategies including envtest, integration tests, and end-to-end testing with OLM.

#### **1. Create Feature Test File**

Create `controllers/activemqartemis_mynewfeature_test.go`:

**Example** - The block-reconcile feature test (`controllers/activemqartemis_controller_block_reconcile_test.go`):

```go
var _ = Describe("block reconcile annotation", func() {
    BeforeEach(func() {
        BeforeEachSpec()
    })

    AfterEach(func() {
        AfterEachSpec()
    })

    Context("annotation test", Label("block-reconcile"), func() {
        It("should pause reconciliation when annotation is set", func() {
            By("deploying a broker")
            brokerCr, createdBrokerCr := DeployCustomBroker(defaultNamespace, func(candidate *brokerv1beta1.ActiveMQArtemis) {
                candidate.Spec.DeploymentPlan.Size = common.Int32ToPtr(1)
            })

            By("setting block reconcile annotation")
            Eventually(func(g Gomega) {
                g.Expect(k8sClient.Get(ctx, brokerKey, createdBrokerCr)).Should(Succeed())
                createdBrokerCr.Annotations = map[string]string{
                    common.BlockReconcileAnnotation: "true",
                }
                g.Expect(k8sClient.Update(ctx, createdBrokerCr)).Should(Succeed())
            }, timeout, interval).Should(Succeed())

            By("cleanup")
            CleanResource(createdBrokerCr, brokerCr.Name, defaultNamespace)
        })
    })
})
```

**Pattern**:
```go
package controllers

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    // Other imports...
)

var _ = Describe("my new feature", func() {

    BeforeEach(func() {
        BeforeEachSpec()
    })

    AfterEach(func() {
        AfterEachSpec()
    })

    Context("feature scenario", Label("mynewfeature"), func() {
        It("should enable new functionality", func() {
            if os.Getenv("USE_EXISTING_CLUSTER") == "true" {

                By("deploying broker with new feature enabled")
                brokerCr, createdBrokerCr := DeployCustomBroker(defaultNamespace, func(candidate *brokerv1beta1.ActiveMQArtemis) {
                    candidate.Spec.MyNewFeature = &[]bool{true}[0]
                    // Other customizations...
                })

                By("verifying feature is active")
                Eventually(func(g Gomega) {
                    // Verify your feature works
                }, existingClusterTimeout, existingClusterInterval).Should(Succeed())

                By("cleanup")
                CleanResource(createdBrokerCr, brokerCr.Name, defaultNamespace)
            } else {
                fmt.Println("Test skipped as it requires an existing cluster")
            }
        })
    })
})
```

#### **2. Add Unit Tests**

Add unit tests to `controllers/activemqartemis_reconciler_test.go`:

```go
func TestMyNewFeatureProcessing(t *testing.T) {
    // Test your processing logic
    assert := assert.New(t)

    // Test cases...
    assert.True(condition, "Expected behavior")
}
```

#### **3. Test Both Modes**

Ensure your feature works in both restricted and legacy modes:

```go
Context("restricted mode", func() {
    It("should work with restricted: true", func() {
        brokerCr, createdBrokerCr := DeployCustomBroker(defaultNamespace, func(candidate *brokerv1beta1.ActiveMQArtemis) {
            candidate.Spec.Restricted = &[]bool{true}[0]
            candidate.Spec.MyNewFeature = &[]bool{true}[0]
        })
        // Test implementation...
    })
})

Context("legacy mode", func() {
    It("should work with restricted: false", func() {
        // Similar test without restricted mode
    })
})
```

### Development Workflow

Standard development workflow for contributing to the operator:

#### **1. Prerequisites and Initial Setup**

> **📖 See**: [Building Documentation](../help/building.md) for complete prerequisites and build instructions.

**Quick Setup**:
```bash
# Prerequisites: Go 1.23.9, Operator SDK v1.28.0, Docker
# See docs/help/building.md for installation instructions

# Clone repository
git clone https://github.com/arkmq-org/activemq-artemis-operator
cd activemq-artemis-operator

# Build operator locally
make build

# Install development tools (for CRD generation and testing)
make controller-gen envtest
```

> **Reference**: See the [Operator SDK Installation guide](https://sdk.operatorframework.io/docs/installation/) for setting up your development environment and the [Operator SDK Quickstart](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/) for an introduction to operator development.

#### **2. Development Cycle (TDD Approach)**

```bash
# 1. Create feature branch
git checkout -b feature/my-new-feature

# 2. Write tests first (TDD approach)
# Create test file: controllers/activemqartemis_mynewfeature_test.go

# 3. Run tests (should fail initially)
make test

# 4. Implement feature
# - Add API field (api/v1beta1/activemqartemis_types.go)
# - Add validation (controllers/activemqartemis_controller.go)
# - Add processing logic (controllers/activemqartemis_reconciler.go)
# - Add resource builders (pkg/resources/*)

# 5. Regenerate code
make manifests generate

# 6. Run tests again
make test

# 7. Test locally (operator runs on your machine, connects to cluster)
export OPERATOR_LOG_LEVEL=debug
make run
# In another terminal: kubectl apply -f config/samples/...

# 8. Test against real cluster (integration tests)
USE_EXISTING_CLUSTER=true make test
```

#### **3. Local Testing Options**

Choose the testing approach that fits your development needs:

**Option A: Quick Iteration (Recommended for development)**

Run operator locally connected to a Kubernetes cluster:

```bash
# Requires: Running Kubernetes cluster (minikube, kind, or remote)
# Operator runs on your machine (not in cluster)

# 1. Ensure you have a cluster
kubectl cluster-info

# 2. Run operator with debug logging
export OPERATOR_LOG_LEVEL=debug
make run

# 3. In another terminal, deploy test broker
kubectl apply -f config/samples/broker_v1beta1_activemqartemis.yaml

# 4. Make changes to code and restart (Ctrl+C then make run again)
```

**Benefits**: Fast iteration, easy debugging, see logs directly
**Limitations**: Not testing deployed operator image, some cluster-specific features may behave differently

> **📖 See**: [Quick Start Guide](../getting-started/quick-start.md) for cluster setup options including Minikube and OpenShift Local (CRC).

**Option B: Full Integration Testing**

Test against a real cluster with deployed operator:

```bash
# Requires: Running Kubernetes cluster with operator deployed

# Run integration tests against existing cluster
USE_EXISTING_CLUSTER=true make test

# Run specific test suite
USE_EXISTING_CLUSTER=true go test -v ./controllers -run TestActiveMQArtemisController
```

> **📖 See**: [Quick Start Guide](../getting-started/quick-start.md#install-the-operator) for operator installation methods.

**Option C: Building and Deploying Custom Operator Image**

For testing operator image changes, webhooks, or OLM integration:

> **📖 See**: [Building Documentation](../help/building.md#building-the-operator-image) for detailed image build and push instructions.

```bash
# 1. Build custom operator image
make docker-build OPERATOR_IMAGE_REPO=quay.io/your-org/activemq-artemis-operator OPERATOR_VERSION=dev

# 2. Push to registry (or use local registry for minikube/kind)
make docker-push OPERATOR_IMAGE_REPO=quay.io/your-org/activemq-artemis-operator OPERATOR_VERSION=dev

# 3. Deploy operator with custom image
make deploy IMG=quay.io/your-org/activemq-artemis-operator:dev

# 4. Verify deployment
kubectl get pods -n activemq-artemis-operator-system

# 5. Check operator logs
kubectl logs -n activemq-artemis-operator-system \
  deployment/activemq-artemis-controller-manager -f

# 6. Test your changes
kubectl apply -f config/samples/broker_v1beta1_activemqartemis.yaml

# 7. Cleanup
make undeploy
```

**Local Registry Setup** (for minikube/kind):

For minikube:
```bash
# Start minikube with sufficient resources
minikube start --cpus=4 --memory=8192

# Use minikube's Docker daemon
eval $(minikube docker-env)

# Build image (automatically available in minikube)
make docker-build IMG=localhost:5000/activemq-artemis-operator:dev
```

For kind with local registry:
```bash
# See: https://kind.sigs.k8s.io/docs/user/local-registry/
# Create kind cluster with local registry configured

# Build and push to local registry
make docker-build IMG=localhost:5000/activemq-artemis-operator:dev
docker push localhost:5000/activemq-artemis-operator:dev
```

#### **4. Required File Changes Pattern**

Every feature addition must touch these files:

**For API field additions** (example: [#1038](https://github.com/arkmq-org/activemq-artemis-operator/issues/1038)/[PR#1037](https://github.com/arkmq-org/activemq-artemis-operator/pull/1037) restricted mode, [#685](https://github.com/arkmq-org/activemq-artemis-operator/issues/685)/[PR#1109](https://github.com/arkmq-org/activemq-artemis-operator/pull/1109) block-reconcile):
```bash
# Core implementation
api/v1beta1/activemqartemis_types.go           # API field definition
controllers/activemqartemis_controller.go      # Validation logic
controllers/activemqartemis_reconciler.go      # Feature implementation
controllers/activemqartemis_mynewfeature_test.go # Feature tests

# Auto-generated (via make manifests generate)
api/v1beta1/zz_generated.deepcopy.go           # DeepCopy methods
config/crd/bases/broker.amq.io_activemqartemises.yaml  # CRD definition
bundle/manifests/*.yaml                        # OLM bundle manifests
deploy/crds/*.yaml                             # Deploy CRDs
deploy/activemq-artemis-operator.yaml          # Deploy operator

# Documentation
docs/help/operator.md                          # User-facing docs
```

**For reconciler-only changes** (example: [#993](https://github.com/arkmq-org/activemq-artemis-operator/issues/993)/[PR#994](https://github.com/arkmq-org/activemq-artemis-operator/pull/994) JSON format support):
```bash
controllers/activemqartemis_reconciler.go      # Implementation
controllers/activemqartemis_controller_test.go # Tests
docs/help/operator.md                          # Documentation
```

**For utility/infrastructure** (example: [#623](https://github.com/arkmq-org/activemq-artemis-operator/issues/623)/[PR#673](https://github.com/arkmq-org/activemq-artemis-operator/pull/673) cert-manager support):
```bash
api/v1beta1/activemqartemis_types.go           # New fields
controllers/activemqartemis_controller.go      # Validation
controllers/activemqartemis_reconciler.go      # Integration
controllers/activemqartemis_controller_cert_manager_test.go # Dedicated tests
pkg/utils/certutil/certutil.go                 # New utility package
pkg/resources/secrets/secret.go                # Resource changes
go.mod, go.sum                                 # New dependencies
```

#### **5. Commit Pattern**

Follow the established commit patterns observed in the codebase:

**Single-line commits** (most common for focused changes):
```bash
# For issue-tracked work
git commit -m "[#123] Add support for my new feature"

# For non-issue work (docs, tests, CI, cleanup)
git commit -m "[NO-ISSUE] Update tutorial examples"

# For multiple related issues
git commit -m "[#123] [#456] Fix validation and add tests"
```

**Multi-line commits** (for complex features or when additional context is valuable):
```bash
git commit -m "[#123] Add support for my new feature

This commit introduces a new feature that...

- Add MyNewFeature field to ActiveMQArtemisSpec
- Add validation for new feature
- Add ProcessMyNewFeature to reconciliation flow
- Add comprehensive test coverage
- Update user documentation"
```

**Common [NO-ISSUE] use cases** (from codebase analysis):
- Documentation updates: `[NO-ISSUE] Update operator documentation`
- Tutorial additions: `[NO-ISSUE] tutorials: locked down scrapping`
- Test improvements: `[NO-ISSUE] Improve test reliability`
- CI/workflow changes: `[NO-ISSUE] Update CI workflow`
- Code cleanup: `[NO-ISSUE] Code clean up - unused parameters`
- Dependency bumps: Auto-generated by dependabot

**Version updates** (no issue number):
```bash
git commit -m "Update version to 2.0.6"
git commit -m "Update related images to 2.42.0"
```

#### **6. Testing Strategy**

```bash
# Test progression
1. Unit tests first:
   go test -v ./controllers -run TestMyNewFeatureProcessing

2. Integration tests:
   go test -v ./controllers -run TestMyNewFeature

3. Full test suite:
   make test

4. Real cluster testing:
   USE_EXISTING_CLUSTER=true make test

5. Specific scenario testing:
   go test -v ./controllers -ginkgo.label-filter="mynewfeature"
```

#### **7. Common Integration Points**

Common integration points for new features:

**Environment Variables** (from controllers/activemqartemis_reconciler.go):
```go
// Pattern: Define as constants in activemqartemis_reconciler.go
const (
    javaArgsAppendEnvVarName = "JAVA_ARGS_APPEND"
    jdkJavaOptionsEnvVarName = "JDK_JAVA_OPTIONS"
    // Add your new env var following the pattern
    myFeatureEnvVarName = "MY_FEATURE_CONFIG"
)
```

**Constants and Suffixes** (from controllers/activemqartemis_reconciler.go):
```go
const (
    // Suffix-based magic behavior (auto-detected by operator)
    jaasConfigSuffix    = "-jaas-config"
    loggingConfigSuffix = "-logging-config"
    brokerPropsSuffix   = "-bp"

    // Resource naming postfixes
    ServiceTypePostfix  = "svc"
    RouteTypePostfix    = "rte"
    IngressTypePostfix  = "ing"

    // Add your new suffix/constant
    myFeatureSuffix = "-myfeature"
)
```

**Resource Naming** (from controllers/activemqartemis_controller.go::MakeNamers):
```go
// In MakeNamers() function
newNamers := common.Namers{
    SsNameBuilder:                 namer.NamerData{},
    SvcHeadlessNameBuilder:        namer.NamerData{},
    SecretsCredentialsNameBuilder: namer.NamerData{},
    // Add your new name builder...
}

// Pattern: Prefix(cr-name).Base(feature).Suffix(type).Generate()
newNamers.SvcHeadlessNameBuilder.Prefix(customResource.Name).Base("hdls").Suffix("svc").Generate()
newNamers.MyFeatureNameBuilder.Prefix(customResource.Name).Base("myfeature").Suffix("svc").Generate()
```

**Validation Chain** (from controllers/activemqartemis_controller.go::validate):
```go
// In validate() function - actual pattern from codebase
func (r *ActiveMQArtemisReconcilerImpl) validate(customResource *brokerv1beta1.ActiveMQArtemis, client rtclient.Client, namer common.Namers) (bool, retry bool) {
    validationCondition := metav1.Condition{
        Type:   brokerv1beta1.ValidConditionType,
        Status: metav1.ConditionTrue,
        Reason: brokerv1beta1.ValidConditionSuccessReason,
    }

    // Existing validations
    condition, retry := validateExtraMounts(customResource, client)
    if condition != nil {
        validationCondition = *condition
    }

    // Add your validation following the chain pattern
    if validationCondition.Status != metav1.ConditionFalse {
        condition, retry = validateMyNewFeature(customResource)
        if condition != nil {
            validationCondition = *condition
        }
    }

    // More validations...

    meta.SetStatusCondition(&customResource.Status.Conditions, validationCondition)
    return validationCondition.Status != metav1.ConditionFalse, retry
}
```

#### **8. Documentation Requirements**

> **📖 See**: [User Documentation](../help/operator.md) for the main user-facing documentation and [Quick Start Guide](../getting-started/quick-start.md) for examples of broker deployment and management.

**User Documentation** (`docs/help/operator.md`) - **FREQUENTLY UPDATED**:
- Updated in ~10% of issue-tracked feature commits
- Examples: [#1187](https://github.com/arkmq-org/activemq-artemis-operator/issues/1187)/[PR#1189](https://github.com/arkmq-org/activemq-artemis-operator/pull/1189) Acceptors docs, [#1113](https://github.com/arkmq-org/activemq-artemis-operator/issues/1113)/[PR#1119](https://github.com/arkmq-org/activemq-artemis-operator/pull/1119) credentials secret docs, [#976](https://github.com/arkmq-org/activemq-artemis-operator/issues/976)/[PR#977](https://github.com/arkmq-org/activemq-artemis-operator/pull/977) Prometheus metrics
- Pattern: Always update for user-facing features

**Examples Directory** (`examples/`) - **RARELY UPDATED**:
- Updated in <2% of feature commits
- Mainly updated for catalog/subscription examples ([#1167](https://github.com/arkmq-org/activemq-artemis-operator/issues/1167)/[PR#1168](https://github.com/arkmq-org/activemq-artemis-operator/pull/1168))
- Not required for most internal features

**Technical Documentation** - **Use codebase_search for technical details**:
- Use `codebase_search` to discover implementation patterns
- Most implementation details are documented via code comments and tests
- Update only for significant architectural changes

**Required for ALL features**:
1. ✅ **User Documentation**: Update docs/help/operator.md::docs/help/operator.md
   - Document new API fields
   - Provide usage examples
   - Explain behavior and limitations

2. ✅ **API Documentation**: Proper kubebuilder annotations
   ```go
   //+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="My Feature"
   MyFeature *bool `json:"myFeature,omitempty"`
   ```

**Optional (as needed)**:
3. ⚠️ **Examples**: Add to examples/::examples/ if:
   - Feature requires complex YAML configuration
   - Multiple resources must be coordinated
   - Common use case needs demonstration

4. ⚠️ **Architecture Changes**: For major structural changes:
   - Document via comprehensive code comments
   - Update tests to demonstrate new patterns
   - Introducing new design patterns

**Tutorials** (`docs/tutorials/`) - frequently added for complex scenarios ([NO-ISSUE] pattern)

Following these patterns ensures new features integrate seamlessly with the operator's established architecture while maintaining high quality standards.

## Development and Debugging Guide

This section provides essential information for developers working with the operator codebase.

> **Reference**: For additional debugging techniques, see the [Operator SDK Go tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/) and [Advanced Topics](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/).

### Running Tests

The operator follows **Test-Driven Development (TDD)** with **396 test scenarios** across 23 test files using [Ginkgo v2](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/):

```bash
# Run all tests
make test

# Run with verbose output
TEST_VERBOSE=true make test

# Run specific test files
go test -v ./controllers -run TestActiveMQArtemisController

# Run tests with coverage
go test -v -cover ./controllers/...
```

**Test Environment Modes** (controllers/suite_test.go):
- **EnvTest Mode**: controllers/suite_test.go with local Kubernetes API server
- **Real Cluster Mode**: Tests against actual clusters when `DEPLOY_OPERATOR=true`
- **Existing Cluster Mode**: Reuses existing clusters when `USE_EXISTING_CLUSTER=true`

### Running Full E2E Test Suite on Minikube

This section provides a complete procedure for running the entire E2E test suite against minikube, which is essential for validating changes before submitting PRs.

#### Complete Setup Procedure

**1. Start Minikube with Sufficient Resources**

```bash
# Start minikube with adequate resources for running all tests
minikube start
```

**2. Enable Required Addons**

The E2E tests require ingress functionality with SSL passthrough:

```bash
# Enable ingress addon
minikube addons enable ingress

# Wait for ingress controller to be ready
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

# Enable SSL passthrough (required for TLS proxy used by tests)
kubectl patch deployment ingress-nginx-controller \
  -n ingress-nginx \
  --type='json' \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--enable-ssl-passthrough"}]'

# Wait for the controller to restart with new configuration
kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx --timeout=120s

# Verify SSL passthrough is enabled
kubectl get deployment ingress-nginx-controller -n ingress-nginx \
  -o jsonpath='{.spec.template.spec.containers[0].args}' | grep "enable-ssl-passthrough"
```

**3. Install Development Tools**

```bash
# Install helm (required for cert-manager installation during tests)
make helm

# Install controller-gen and envtest
make controller-gen envtest
```

**4. Generate Manifests and Install CRDs**

```bash
# Generate CRDs, DeepCopy methods, and format code
make manifests generate fmt vet

# Install CRDs (using server-side apply to handle large CRDs)
kubectl apply -f config/crd/bases/broker.amq.io_activemqartemises.yaml --server-side
kubectl apply -f config/crd/bases/broker.amq.io_activemqartemisaddresses.yaml
kubectl apply -f config/crd/bases/broker.amq.io_activemqartemisscaledowns.yaml
kubectl apply -f config/crd/bases/broker.amq.io_activemqartemissecurities.yaml

# Verify CRDs are installed
kubectl get crds | grep broker.amq.io
```

**5. Run the Full Test Suite**

Run all tests excluding those requiring deployed operator (`do` label):

```bash
# Run with make target (includes fail-fast)
make test-mk-v

# OR run without fail-fast to see all test results
USE_EXISTING_CLUSTER=true \
RECONCILE_RESYNC_PERIOD=5s \
HELM=/home/tlavocat/dev/activemq-artemis-operator/bin/helm \
go test ./controllers -p 1 \
  -test.timeout=120m \
  -ginkgo.label-filter='!do' \
  -v \
  -ginkgo.poll-progress-after=150s \
  -coverprofile cover-mk.out
```

**Expected Results:**
- **~213 tests** should pass (out of 226 total)
- **~13 tests** skipped (labeled with `do` - for deployed operator mode)
- **Coverage**: ~91% of statements
- **Duration**: ~45-50 minutes depending on hardware

#### Test Targets Reference

The Makefile provides several test targets for minikube:

```bash
# Local operator mode (operator runs on your machine, excludes 'do' tests)
make test-mk          # Without verbose output
make test-mk-v        # With verbose output

# Deployed operator mode (operator deployed in cluster, only 'do' tests)
make test-mk-do       # Without verbose output
make test-mk-do-v     # With verbose output

# Fast deployed operator mode (excludes 'slow' tests)
make test-mk-do-fast
make test-mk-do-fast-v
```

**Makefile Configuration** (from `Makefile:149-168`):
- `test-mk`: 120-minute timeout, excludes `do` labeled tests
- `test-mk-do`: 60-minute timeout, only runs `do` labeled tests
- `test-mk-do-fast`: 30-minute timeout, runs `do && !slow` labeled tests

#### Automatic Features

**Cert-Manager Auto-Installation**

The E2E test suite automatically handles cert-manager installation:

- Tests check if cert-manager is installed (`CertManagerInstalled()`)
- If not present, tests install it automatically via Helm
- After tests complete, auto-installed cert-manager is cleaned up
- If cert-manager was already present, tests leave it unchanged

See controllers/common_util_test.go::controllers/common_util_test.go:857-892 for implementation.

#### Troubleshooting

**Issue: Tests timeout waiting for broker pods**

```bash
# Check if pods are pulling images
kubectl get pods -n test -w

# Check events for issues
kubectl get events -n test --sort-by='.lastTimestamp'

# Check if images are available
kubectl describe pod <pod-name> -n test
```

**Issue: Connection refused to ingress**

```bash
# Verify ingress controller is running
kubectl get pods -n ingress-nginx

# Check SSL passthrough is enabled
kubectl get deployment ingress-nginx-controller -n ingress-nginx \
  -o jsonpath='{.spec.template.spec.containers[0].args}' | grep ssl-passthrough

# Check ingress resources
kubectl get ingress -A
```

**Issue: CRD annotation size limit exceeded**

The `activemqartemises.broker.amq.io` CRD is very large. Use server-side apply:

```bash
kubectl apply -f config/crd/bases/broker.amq.io_activemqartemises.yaml --server-side
```

**Issue: Tests fail with "not enough resources"**

```bash
# Increase minikube resources
minikube stop
minikube delete
minikube start --cpus=6 --memory=12288
```

#### Running Specific Test Suites

Run specific tests by label or name:

```bash
# Run only cert-manager tests
USE_EXISTING_CLUSTER=true go test -v ./controllers \
  -ginkgo.label-filter="controller-cert-mgr-test" \
  -timeout 120m

# Run only control plane tests
USE_EXISTING_CLUSTER=true go test -v ./controllers \
  -ginkgo.focus="minimal restricted rbac" \
  -timeout 120m

# Skip specific failing tests
USE_EXISTING_CLUSTER=true go test -v ./controllers \
  -ginkgo.skip="operator role access" \
  -ginkgo.label-filter='!do' \
  -timeout 120m

# Run tests for a specific file
USE_EXISTING_CLUSTER=true go test -v ./controllers \
  -run TestAPIs \
  -ginkgo.focus="control plane override" \
  -timeout 30m
```

#### Clean Up

After testing, you can clean up resources:

```bash
# Delete test resources
kubectl delete namespace test --ignore-not-found=true
kubectl delete namespace other --ignore-not-found=true
kubectl delete namespace restricted --ignore-not-found=true

# Stop minikube (preserves state)
minikube stop

# Or delete minikube completely
minikube delete
```

### Debugging Reconciliation Issues

**Enable verbose logging**:

Set the `OPERATOR_LOG_LEVEL` environment variable to `debug` when running the operator:

```bash
export OPERATOR_LOG_LEVEL=debug
make run
```

**Common debugging steps**:

1. **controllers/activemqartemis_controller.go**: The CR's `.status` field contains conditions that reflect the current state of reconciliation. Look for error messages in condition reasons.

   ```bash
   kubectl get activemqartemis <name> -o jsonpath='{.status.conditions}' | jq
   ```

2. **controllers/activemqartemis_reconciler.go**: Compare the generated `StatefulSet`, `Service`, and `Secret` resources against expectations.

   ```bash
   kubectl get statefulset <name>-ss -o yaml
   kubectl get secret <name>-props -o yaml
   ```

3. **Check Broker Logs**: The broker's stdout/stderr often contains configuration errors that aren't surfaced to the operator.

   ```bash
   kubectl logs <broker-pod-name> -c <broker-container-name>
   ```

4. **controllers/activemqartemis_reconciler.go**: Verify the properties secret contains the expected configuration and formatting.

   ```bash
   kubectl get secret <name>-props-0 -o jsonpath='{.data.broker\.properties}' | base64 -d
   ```

### Code Navigation Tips

For developers new to the codebase, these are the key entry points:

- **Reconciliation Entry Point**: Start with `controllers/activemqartemis_reconciler.go::Process()` to understand the main reconciliation flow
- **StatefulSet Generation**: See `pkg/resources/statefulsets/statefulset.go` for StatefulSet creation patterns
- **Broker Properties**: Check `controllers/activemqartemis_reconciler.go::ProcessBrokerProperties()` for configuration management
- **Service Creation**: Look at `pkg/resources/services/service.go` for headless, ping, and acceptor services
- **Certificate Management**: Search for `cert-manager.io` usage in `controllers/activemqartemis_reconciler.go`
- **Drain Controller**: See `pkg/draincontroller/` for message migration on scale-down

## Related Documentation

### Internal Documentation
- **Test Coverage**: Use `codebase_search` to find test patterns and examples in `controllers/*_test.go`
- **User Documentation**: See `docs/help/operator.md` for end-user configuration guide
- **Implementation Patterns**: Use `codebase_search` to discover patterns and conventions in the codebase
- **Tutorials**: See `docs/tutorials/` for step-by-step guides for common scenarios

### External Resources
- **[Operator SDK Documentation](https://sdk.operatorframework.io/)**: Official Operator SDK documentation
- **[Kubebuilder Book](https://book.kubebuilder.io/)**: Comprehensive guide to building Kubernetes operators
- **[Operator Best Practices](https://sdk.operatorframework.io/docs/best-practices/)**: Industry best practices for operator development
- **[OLM Integration](https://sdk.operatorframework.io/docs/olm-integration/)**: Guide for Operator Lifecycle Manager integration

