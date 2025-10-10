---
ai_metadata:
  doc_type: "test_catalog"
  primary_audience: ["developers", "testers", "ai_assistants"]
  key_topics: ["testing", "tdd", "test_scenarios", "feature_coverage"]
  answers_questions:
    - "What tests exist?"
    - "How is feature X tested?"
    - "What test coverage exists?"
  test_scenarios: 396
  test_files: 39
  code_links: 369
  complexity: "comprehensive"
  related_docs: ["contribution_guide.md", "operator_architecture.md"]
---

# ActiveMQ Artemis Operator Test Index

**Generated**: 2025-10-07  
**Last Updated**: 2025-10-08  
**Total Test Files**: 39  
**Estimated Test Scenarios**: 396+
**Linked Scenarios**: 322/396 (Sections 1-35 complete)

This document provides a comprehensive index of all tests in the ActiveMQ Artemis Operator codebase. Since the project follows **Test-Driven Development (TDD)**, this index serves as a functional inventory of all implemented features. Individual test scenarios are systematically linked to their source code implementation.

> **Code Reference Notice**
> 
> All test links in this document are accurate as of **October 9, 2025** (commit [`4acadb95`](https://github.com/arkmq-org/activemq-artemis-operator/commit/4acadb95603c38b82d5d7f63fb538c37ed855662)).
> 
> Test file links use GitHub permalinks. Test files are generally stable, but if line numbers drift:
> - Use the test name (in link text) to locate it in current code
> - Replace `/blob/4acadb95...` with `/blob/main` to see current implementation
> - Test context and description help identify the test even if moved

---

## Table of Contents

- [Detailed Test Catalog](#detailed-test-catalog)
   - [Controller Tests](#controller-tests)
    - [1. `activemqartemis_controller_test.go`](#1-activemqartemiscontrollertestgo-143-test-scenarios)
      - [TLS/SSL Configuration](#tlsssl-configuration)
        - [TLS Secret Reuse](#tls-secret-reuse)
      - [StatefulSet Management](#statefulset-management)
        - [Statefulset Options](#statefulset-options)
      - [High Availability](#high-availability)
        - [Pod Disruption Budget](#pod-disruption-budget)
      - [Error Handling](#error-handling)
        - [Broker Status on Resource Error](#broker-status-on-resource-error)
      - [Version Management](#version-management)
        - [Broker Versions Map](#broker-versions-map)
        - [Broker Versions](#broker-versions)
      - [Resource Tracking](#resource-tracking)
        - [Broker Resource Tracking](#broker-resource-tracking)
      - [Address Settings](#address-settings)
        - [New Address Settings Options](#new-address-settings-options)
      - [Versions Test](#versions-test)
        - [Versions Test](#versions-test)
      - [Version Conversion](#version-conversion)
        - [CR Version Conversion Test](#cr-version-conversion-test)
      - [Console Configuration](#console-configuration)
        - [Console Config Test](#console-config-test)
      - [Expose Modes](#expose-modes)
        - [Expose Mode Test](#expose-mode-test)
      - [Console Secret Validation](#console-secret-validation)
        - [SSL Console Secret Validation](#ssl-console-secret-validation)
      - [Image Management](#image-management)
        - [Image Update Test](#image-update-test)
      - [Client ID Auto-Sharding](#client-id-auto-sharding)
        - [ClientID Autoshard Test](#clientid-autoshard-test)
      - [Probe Configuration](#probe-configuration)
        - [Probe Defaults Reconcile](#probe-defaults-reconcile)
      - [StatefulSet Recovery](#statefulset-recovery)
        - [SS Delete Recreate Test](#ss-delete-recreate-test)
      - [PVC Management](#pvc-management)
        - [PVC No GC Test](#pvc-no-gc-test)
        - [PVC Upgrade Owner Reference Test](#pvc-upgrade-owner-reference-test)
      - [Tolerations](#tolerations)
        - [Tolerations Existing Cluster](#tolerations-existing-cluster)
      - [Console Secrets](#console-secrets)
        - [Console Secret Test](#console-secret-test)
      - [Security Contexts](#security-contexts)
        - [PodSecurityContext Test](#podsecuritycontext-test)
      - [Affinity Configuration](#affinity-configuration)
        - [Affinity Test](#affinity-test)
      - [Node Selector](#node-selector)
        - [Node Selector Test](#node-selector-test)
      - [Annotations](#annotations)
        - [Annotations Test](#annotations-test)
      - [Labels](#labels)
        - [Labels Test](#labels-test)
      - [Namespace Isolation](#namespace-isolation)
        - [Different Namespace, Deployed Before Start](#different-namespace-deployed-before-start)
      - [Tolerations Configuration](#tolerations-configuration)
        - [Tolerations Test](#tolerations-test)
      - [Liveness Probes](#liveness-probes)
        - [Liveness Probe Tests](#liveness-probe-tests)
      - [Readiness Probes](#readiness-probes)
        - [Readiness Probe Tests](#readiness-probe-tests)
      - [Startup Probes](#startup-probes)
        - [Startup Probe Tests](#startup-probe-tests)
      - [Status Tracking](#status-tracking)
        - [Status](#status)
      - [Validation](#validation)
        - [Validation](#validation)
      - [Rolling Updates](#rolling-updates)
        - [Env Var Updates TRIGGERED_ROLL_COUNT Checksum](#env-var-updates-triggeredrollcount-checksum)
      - [Broker Properties](#broker-properties)
        - [BrokerProperties](#brokerproperties)
      - [Broker Version](#broker-version)
        - [BrokerVersion](#brokerversion)
      - [Logger Properties](#logger-properties)
        - [LoggerProperties](#loggerproperties)
      - [Address Settings](#address-settings)
        - [With Address Settings via Updated CR](#with-address-settings-via-updated-cr)
      - [Address CR](#address-cr)
        - [With Address CR](#with-address-cr)
      - [Persistence Toggle](#persistence-toggle)
        - [With Toggle Persistence=True](#with-toggle-persistencetrue)
        - [With Update Persistence=True Single Reconcile](#with-update-persistencetrue-single-reconcile)
      - [Version Toggle](#version-toggle)
        - [Toggle spec.Version](#toggle-specversion)
      - [Performance](#performance)
        - [Time to Ready](#time-to-ready)
      - [Resource Templates](#resource-templates)
        - [Template Tests](#template-tests)
      - [Acceptor Management](#acceptor-management)
        - [With Deployed Controller - Acceptor](#with-deployed-controller-acceptor)
      - [Acceptor/Connector Updates](#acceptorconnector-updates)
        - [Acceptor/Connector Update No Owner Ref](#acceptorconnector-update-no-owner-ref)
      - [API Version Compatibility](#api-version-compatibility)
        - [With Deployed Controller](#with-deployed-controller)
      - [Storage Configuration](#storage-configuration)
        - [With Deployed Controller](#with-deployed-controller)
      - [User Management](#user-management)
        - [User Management Tests](#user-management-tests)
      - [Lifecycle](#lifecycle)
        - [Lifecycle Tests](#lifecycle-tests)
      - [Management RBAC](#management-rbac)
        - [Management RBAC Tests](#management-rbac-tests)
      - [Environment Variables](#environment-variables)
        - [Environment Variable Tests](#environment-variable-tests)
      - [Metrics](#metrics)
        - [Metrics Tests](#metrics-tests)
      - [Config Projection](#config-projection)
        - [Config Projection](#config-projection)
    - [2. `activemqartemis_controller2_test.go`](#2-activemqartemiscontroller2testgo-4-test-scenarios)
      - [Persistent Volumes](#persistent-volumes)
        - [Persistent Volumes Tests](#persistent-volumes-tests)
      - [Route Reconciliation](#route-reconciliation)
        - [Route Reconcile](#route-reconcile)
    - [3. `activemqartemis_controller_unit_test.go`](#3-activemqartemiscontrollerunittestgo-6-test-functions)
      - [Validation Tests](#validation-tests)
        - [Validation Tests](#validation-tests)
      - [Status Caching Tests](#status-caching-tests)
        - [Status Caching Tests](#status-caching-tests)
    - [4. `activemqartemis_controller_deploy_operator_test.go`](#4-activemqartemiscontrollerdeployoperatortestgo-4-test-scenarios)
      - [TLS Jolokia Access](#tls-jolokia-access)
        - [TLS Jolokia Access](#tls-jolokia-access)
      - [Operator Logging Configuration](#operator-logging-configuration)
        - [Operator Logging Config](#operator-logging-config)
      - [Custom Related Images](#custom-related-images)
        - [Operator Deployment in Default Namespace](#operator-deployment-in-default-namespace)
      - [Restricted Namespace](#restricted-namespace)
        - [Operator Deployment in Restricted Namespace](#operator-deployment-in-restricted-namespace)
    - [5. `activemqartemis_controller_cert_manager_test.go`](#5-activemqartemiscontrollercertmanagertestgo-8-test-scenarios)
      - [Java Keystore Integration](#java-keystore-integration)
        - [Cert-Manager Cert with Java Store](#cert-manager-cert-with-java-store)
      - [TLS Exposure](#tls-exposure)
        - [TLS Exposure with Cert Manager](#tls-exposure-with-cert-manager)
      - [Certificate Utilities](#certificate-utilities)
        - [Certutil Functions](#certutil-functions)
      - [Certificate Annotations](#certificate-annotations)
        - [Certificate from Annotations](#certificate-from-annotations)
      - [Certificate Rotation](#certificate-rotation)
        - [Certificate Rotation](#certificate-rotation)
      - [Certificate Bundle](#certificate-bundle)
        - [Certificate Bundle](#certificate-bundle)
    - [6. `activemqartemis_controller_block_reconcile_test.go`](#6-activemqartemiscontrollerblockreconciletestgo-2-test-scenarios)
      - [Block Reconcile Annotation](#block-reconcile-annotation)
        - [Test](#test)
   - [Reconciler Tests](#reconciler-tests)
    - [7. `activemqartemis_reconciler_test.go`](#7-activemqartemisreconcilertestgo-38-test-functions)
      - [Hashing and Comparison](#hashing-and-comparison)
      - [Configuration Management](#configuration-management)
      - [Broker Properties Processing](#broker-properties-processing)
      - [JAAS Configuration](#jaas-configuration)
      - [Pod Template Generation](#pod-template-generation)
      - [Resource Template Processing](#resource-template-processing)
      - [Status and Utilities](#status-and-utilities)
  - [Specialized Feature Tests](#specialized-feature-tests)
    - [8. `activemqartemis_pub_sub_scale_test.go`](#8-activemqartemispubsubscaletestgo-2-test-scenarios)
      - [Pub/Sub HA with Partitioned Subscribers](#pubsub-ha-with-partitioned-subscribers)
        - [Pub n Sub, HA Pub, Partitioned Sub](#pub-n-sub-ha-pub-partitioned-sub)
    - [9. `activemqartemis_scale_zero_test.go`](#9-activemqartemisscalezerotestgo-2-test-scenarios)
      - [Scale to Zero via Subresource](#scale-to-zero-via-subresource)
        - [Scale to Zero via Scale Subresource](#scale-to-zero-via-scale-subresource)
    - [10. `activemqartemis_work_queue_test.go`](#10-activemqartemisworkqueuetestgo-1-test-scenario)
      - [HA Work Queue with Competing Subscribers](#ha-work-queue-with-competing-subscribers)
        - [HA Pub and HA Competing Sub, Compromised Total Message Order](#ha-pub-and-ha-competing-sub-compromised-total-message-order)
    - [11. `activemqartemis_rwm_pvc_ha_test.go`](#11-activemqartemisrwmpvchatestgo-1-test-scenario)
      - [Shared Store HA with Two Peer CRs](#shared-store-ha-with-two-peer-crs)
        - [Two Peer CRs](#two-peer-crs)
    - [12. `activemqartemis_jdbc_ha_test.go`](#12-activemqartemisjdbchatestgo-1-test-scenario)
      - [JDBC Store HA](#jdbc-store-ha)
        - [Artemis](#artemis)
  - [Address Controller Tests](#address-controller-tests)
    - [13. `activemqartemisaddress_controller_test.go`](#13-activemqartemisaddresscontrollertestgo-20-test-scenarios)
      - [Configuration Defaults](#configuration-defaults)
        - [Address Queue Config Defaults](#address-queue-config-defaults)
      - [Address Management](#address-management)
        - [Broker with Address Custom Resources](#broker-with-address-custom-resources)
      - [Reconciliation](#reconciliation)
        - [Address Controller Test with Reconcile](#address-controller-test-with-reconcile)
      - [Scale Down](#scale-down)
        - [Address Delete and Scale Down](#address-delete-and-scale-down)
      - [Address Creation](#address-creation)
        - [Address Creation Test](#address-creation-test)
      - [Console Agent](#console-agent)
        - [Address CR with Console Agent](#address-cr-with-console-agent)
      - [Conversion and Migration](#conversion-and-migration)
        - [Address Conversion Test](#address-conversion-test)
        - [Address CR Migration Test](#address-cr-migration-test)
    - [14. `activemqartemisaddress_controller_unit_test.go`](#14-activemqartemisaddresscontrollerunittestgo-4-test-functions)
      - [Address Deletion](#address-deletion)
    - [15. `activemqartemisaddress_controller_deploy_operator_test.go`](#15-activemqartemisaddresscontrollerdeployoperatortestgo-7-test-scenarios)
      - [Address Management](#address-management)
        - [Address Test](#address-test)
      - [Scale Down](#scale-down)
        - [Address Delete and Scale Down](#address-delete-and-scale-down)
      - [Address Creation](#address-creation)
        - [Address Creation Test](#address-creation-test)
      - [Console Agent](#console-agent)
        - [Address CR with Console Agent](#address-cr-with-console-agent)
    - [16. `activemqartemis_address_broker_properties_test.go`](#16-activemqartemisaddressbrokerpropertiestestgo-6-test-scenarios)
      - [Configuration Defaults](#configuration-defaults)
        - [BP Address Queue Config Defaults](#bp-address-queue-config-defaults)
      - [Broker with Queue](#broker-with-queue)
        - [BP Broker with Queue](#bp-broker-with-queue)
      - [Duplicate Queue Handling](#duplicate-queue-handling)
        - [BP Broker with Duplicate Queue](#bp-broker-with-duplicate-queue)
      - [Reconciliation](#reconciliation)
        - [BP Address Controller Test with Reconcile](#bp-address-controller-test-with-reconcile)
      - [Address Deletion](#address-deletion)
        - [BP Address Delete](#bp-address-delete)
      - [Address Creation](#address-creation)
        - [BP Address Creation Test](#bp-address-creation-test)
  - [Security Controller Tests](#security-controller-tests)
    - [17. `activemqartemissecurity_controller_test.go`](#17-activemqartemissecuritycontrollertestgo-11-test-scenarios)
      - [Security Resources](#security-resources)
        - [Broker with Security Custom Resources](#broker-with-security-custom-resources)
      - [Reconciliation](#reconciliation)
        - [Reconcile Test](#reconcile-test)
      - [Force Reconcile](#force-reconcile)
    - [18. `activemqartemissecurity_broker_properties_test.go`](#18-activemqartemissecuritybrokerpropertiestestgo-1-test-scenario)
      - [BrokerProperties RBAC](#brokerproperties-rbac)
        - [BrokerProperties RBAC](#brokerproperties-rbac)
  - [Scale Down Controller Tests](#scale-down-controller-tests)
    - [19. `activemqartemisscaledown_controller_test.go`](#19-activemqartemisscaledowncontrollertestgo-2-test-scenarios)
      - [Scale Down Test](#scale-down-test)
        - [Scale Down Test](#scale-down-test)
      - [Toleration During Scale Down](#toleration-during-scale-down)
  - [Controller Manager and Infrastructure Tests](#controller-manager-and-infrastructure-tests)
    - [20. `controllermanager_test.go`](#20-controllermanagertestgo-7-test-scenarios)
      - [Operator Namespaces Test](#operator-namespaces-test)
        - [Operator Namespaces Test](#operator-namespaces-test)
    - [21. `controll_plane_test.go`](#21-controllplanetestgo-1-test-scenario)
      - [Restricted RBAC](#restricted-rbac)
        - [Restricted RBAC](#restricted-rbac)
    - [22. `broker_name_test.go`](#22-brokernametestgo-2-test-scenarios)
      - [Non-Default Broker Names](#non-default-broker-names)
        - [Set Non Default Value](#set-non-default-value)
    - [23. `common_util_test.go`](#23-commonutiltestgo-0-test-scenarios-support-file)
    - [24. `suite_test.go`](#24-suitetestgo-1-test-function-setup-file)
      - [Test Suite Setup](#test-suite-setup)
   - [Component Tests](#component-tests)
  - [Utility Package Tests](#utility-package-tests)
    - [25. `pkg/utils/common/common_test.go`](#25-pkgutilscommoncommontestgo-1-test-context)
      - [Resync Period Configuration](#resync-period-configuration)
        - [Default Resync Period](#default-resync-period)
    - [26. `pkg/utils/common/conditions_test.go`](#26-pkgutilscommonconditionstestgo-15-test-scenarios)
      - [SetReadyCondition](#setreadycondition)
        - [SetReadyCondition](#setreadycondition)
      - [IsConditionPresentAndEquals](#isconditionpresentandequals)
        - [IsConditionPresentAndEquals](#isconditionpresentandequals)
      - [IsConditionPresentAndEqualsIgnoreMessage](#isconditionpresentandequalsignoremessage)
        - [IsConditionPresentAndEqualsIgnoreMessage](#isconditionpresentandequalsignoremessage)
      - [PodStartingStatusDigestMessage](#podstartingstatusdigestmessage)
        - [PodStartingStatusDigestMessage](#podstartingstatusdigestmessage)
    - [27. `pkg/utils/jolokia_client/jolokia_client_test.go`](#27-pkgutilsjolokiaclientjolokiaclienttestgo-7-test-scenarios)
      - [GetBrokers](#getbrokers)
    - [28. `pkg/utils/artemis/artemis_test.go`](#28-pkgutilsartemisartemistestgo-4-test-functions)
      - [Status Retrieval](#status-retrieval)
    - [29. `pkg/utils/cr2jinja2/cr2jinja2_test.go`](#29-pkgutilscr2jinja2cr2jinja2testgo-1-test-context)
      - [Type Validation](#type-validation)
        - [Testing Special Key Must Be of String Type](#testing-special-key-must-be-of-string-type)
  - [Resource Builder Tests](#resource-builder-tests)
    - [30. `pkg/resources/statefulsets/statefulset_test.go`](#30-pkgresourcesstatefulsetsstatefulsettestgo-4-test-scenarios)
      - [GetDeployedStatefulSetNames](#getdeployedstatefulsetnames)
  - [Log Package Tests](#log-package-tests)
    - [31. `pkg/log/filtered_log_sink_test.go`](#31-pkglogfilteredlogsinktestgo-1-test-function)
      - [Filtered Log Sink](#filtered-log-sink)
  - [DrainController Tests](#draincontroller-tests)
    - [32. `pkg/draincontroller/draincontroller_test.go`](#32-pkgdraincontrollerdraincontrollertestgo-1-test-scenario)
      - [Global Template Test](#global-template-test)
        - [Global Template Test](#global-template-test)
  - [Test Utilities](#test-utilities)
    - [33-35. Test Utilities](#33-35-test-utilities-10-test-scenarios-total)
      - [33. `test/utils/namer/namer_test.go`](#33-testutilsnamernamertestgo-3-test-scenarios)
        - [Pod/StatefulSet Name Matching](#podstatefulset-name-matching)
        - [Testing Pod Belonging to Statefulset](#testing-pod-belonging-to-statefulset)
        - [Testing Different Namespaces](#testing-different-namespaces)
      - [34. `test/utils/config/config_test.go`](#34-testutilsconfigconfigtestgo-1-test-scenario)
        - [Address Settings Comparison](#address-settings-comparison)
      - [35. `test/utils/common/common_test.go`](#35-testutilscommoncommontestgo-6-test-scenarios)
        - [Resource Comparison](#resource-comparison)
        - [Testing 2 Resources Are Equal](#testing-2-resources-are-equal)
        - [Testing 2 Resources Are Different](#testing-2-resources-are-different)
        - [Testing 2 Resources Are Different](#testing-2-resources-are-different)
        - [Testing 2 Resources Are Different](#testing-2-resources-are-different)
- [Test Labels and Categories](#test-labels-and-categories)
  - [Ginkgo Labels Used](#ginkgo-labels-used)
  - [Test Categorization by Type](#test-categorization-by-type)
- [Test Patterns and Conventions](#test-patterns-and-conventions)
  - [Ginkgo Test Structure](#ginkgo-test-structure)
  - [Common Test Utilities](#common-test-utilities)
  - [Test Modes](#test-modes)

---

## Detailed Test Catalog

### Controller Tests

#### 1. [`activemqartemis_controller_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L77) (143 test scenarios)



##### TLS/SSL Configuration:

Tests validating TLS/SSL secret management and configuration across broker components. Ensures secure communication channels can be established using shared or dedicated certificates for console, acceptor, and connector endpoints.

###### TLS Secret Reuse

Validates that a single TLS secret can be shared across multiple broker components (console, acceptors, connectors) while maintaining proper service, route, and ingress configurations.

<u>[Console and Acceptor Share One Secret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L90)</u>

This test deploys a broker configured with a common TLS secret shared across the console (with SSL enabled on port 8443), an AMQP acceptor (SSL on port 5671), and a TCP connector (SSL on port 5671). The test creates the shared TLS secret, deploys the broker CR with all three components referencing the same secret, and validates that services are properly created for each exposed endpoint. It verifies the operator correctly handles SSL secret mounting and configuration across different broker components without conflicts, and confirms both Kubernetes ingress and OpenShift route resources are created appropriately based on the platform.

##### StatefulSet Management:

Tests controlling StatefulSet-level configurations that affect how Kubernetes manages broker pod deployments, including revision history tracking and update strategies.

###### Statefulset Options

Validates StatefulSet-specific configuration options that control deployment behavior and history management.

<u>[Revision History Limit](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L314)</u>

This test deploys a broker with a custom revision history limit of 1000 set via `Spec.DeploymentPlan.RevisionHistoryLimit`. It retrieves the created StatefulSet and verifies that the `RevisionHistoryLimit` field is properly propagated from the ActiveMQArtemis CR to the underlying StatefulSet resource. This ensures the operator correctly translates this Kubernetes-level configuration setting, which controls how many old ReplicaSets are retained for rollback purposes.

##### High Availability:

Tests ensuring broker deployments can maintain availability during voluntary disruptions through PodDisruptionBudget configurations. Validates various PDB settings including numeric and percentage-based availability requirements.

###### Pod Disruption Budget

Validates PodDisruptionBudget creation, validation, and updates to ensure high availability during voluntary disruptions like node drains or rolling updates.

<u>[Pod Disruption Budget Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L337)</u>

This test verifies that when a PodDisruptionBudget is configured with a non-nil selector, the operator detects this as invalid and updates the broker CR status with a validation error. The test deploys a broker with a PDB containing custom match labels and verifies the status condition shows `ValidConditionType` as false with reason `ValidConditionPDBNonNilSelectorReason`. This ensures the operator prevents invalid PDB configurations where users specify custom selectors that could bypass the operator's label management.

<u>[Pod Disruption Apply Number](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L374)</u>

This test validates PodDisruptionBudget creation with a numeric MinAvailable value. It deploys a 2-replica broker with MinAvailable set to 1, then verifies the PDB resource is created with the correct specification including MinAvailable=1 and proper label selectors matching the broker CR name. In existing cluster environments, it also validates the PDB status reflects the correct health metrics (DisruptionsAllowed=1, CurrentHealthy=2, DesiredHealthy=1, ExpectedPods=2), confirming Kubernetes properly tracks availability constraints.

<u>[Pod Disruption Apply Percentage](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L410)</u>

This test validates PodDisruptionBudget creation using percentage-based availability constraints. It deploys a 2-replica broker with MaxUnavailable set to "50%", then verifies the PDB is created with the percentage value properly preserved and label selectors correctly configured. This ensures the operator supports both numeric and percentage-based PDB specifications, allowing flexible high-availability policies.

<u>[Pod Disruption Update](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L439)</u>

This test validates that PodDisruptionBudget updates are properly propagated when the broker CR is modified. It initially deploys a broker with MinAvailable=1, verifies the PDB creation and status, then updates the CR to change MinAvailable to 2. The test confirms the PDB resource is updated (ResourceVersion changes) and the new MinAvailable value is applied, ensuring the operator properly reconciles PDB changes and maintains synchronization between the CR and the actual PDB resource.

##### Error Handling:

Tests validating that the operator properly detects, reports, and handles invalid configurations and resource creation failures. Ensures error conditions are surfaced in CR status with meaningful messages.

###### Broker Status on Resource Error

Validates that when resource creation fails due to invalid configurations, the broker CR status is updated with appropriate error conditions and messages.

<u>[Two Acceptors with Names Too Long for Kube](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L494)</u>

This test deploys a broker with two exposed acceptors that have names exceeding Kubernetes resource name length limits (63 characters when combined with suffixes for services/routes/ingress). On Kubernetes, it verifies that the DeployedCondition status becomes False with reason `DeployedConditionCrudKindErrorReason` and a message containing the problematic acceptor name and "failed". It confirms the operator detects and reports resource creation failures in the CR status, setting both Deployed and Ready conditions to false. On OpenShift, routes have no such length limits, so deployment succeeds.

<u>[Two Acceptors with Port Clash](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L544)</u>

This test validates port conflict detection between acceptors. It deploys a broker with two acceptors ("first" and "second") both configured to use port 61636 and both exposed. The test verifies that the operator detects this port clash and updates the CR status with ValidConditionType set to False and reason `ValidConditionFailedDuplicateAcceptorPort`. This ensures the operator prevents invalid configurations where multiple acceptors would conflict at the network level, catching configuration errors before deployment.

##### Version Management:

Tests validating broker version resolution, image selection, and version compatibility checks. Ensures the operator correctly handles various version specifications (major, minor, patch levels) and validates image pairings.

###### Broker Versions Map

Tests the version resolution algorithm that selects the appropriate broker version based on user input.

<u>[Resolve Various Versions](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L608)</u>

This test validates the `ResolveBrokerVersion` function's ability to match version specifications to available versions. It tests resolution with empty strings (defaults to latest 7.10.2), major-only versions ("7" → 7.10.2), major.minor versions ("7.8" → 7.8.3, "7.9" → 7.9.4), and exact versions ("7.9.0" → 7.9.0). It also confirms nil is returned for unavailable versions ("7.11", "8"). This ensures the operator's version resolution logic correctly implements semver matching with partial version specifications.

###### Broker Versions

Tests version validation logic for different combinations of version strings and explicit image specifications.

<u>[Version Validation When Both Version and Images Are Explicitly Specified](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L648)</u>

This test validates that when both a specific version and custom images are provided, the CR is marked as valid. It deploys a broker with `Spec.Version` set to the default version and custom images for both Image and InitImage fields. The test confirms the ValidConditionType status is True, indicating the operator accepts this configuration where users override images while still specifying a version.

<u>[Version Validation When Version Is Loosly Specified and Images Are Explicitly Specified](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L667)</u>

This test checks validation behavior when a major-only version (e.g., "7") is specified alongside custom images. It deploys with a major version extracted from the default version and custom Image/InitImage paths. The test verifies the status shows ValidConditionType as Unknown with reason `ValidConditionUnknownReason` and message containing "NotSupportedImageVersionMessage", indicating the operator cannot verify compatibility when loose versions are combined with custom images.

<u>[Version Validation When Version and Images Are Loosly Specified](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L692)</u>

This test validates behavior when only a major version is specified without any explicit image overrides. The broker is deployed with just the major version number, allowing the operator to select appropriate images. The test confirms ValidConditionType is Unknown with the unsupported version message, as the operator requires more specific version information for validation without explicit images.

<u>[Images Need to Be in Pairs](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L712)</u>

This test verifies that when custom images are specified, both the broker image and init image must be provided as a pair. It deploys a broker with only Spec.DeploymentPlan.Image set (no InitImage). The test validates that the CR status shows ValidConditionType as False with reason `ValidConditionFailedImagePairReason`, ensuring the operator prevents incomplete image configurations that would lead to runtime failures.

<u>[Images Need to Be in Pairs, Image Placeholder OK](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L745)</u>

This test confirms that the special placeholder value "placeholder" is acceptable for the image field when a version is specified. It deploys with Image="placeholder" and InitImage empty. The test verifies ValidConditionType is True, showing the operator treats "placeholder" as a special marker that doesn't require a paired init image, allowing flexible configurations during template-based deployments.

<u>[Version Validation Invalid](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L773)</u>

This test validates rejection of syntactically invalid version strings. It attempts to deploy with Version="x.y.z" (not semver-compliant). The test confirms ValidConditionType becomes False with reason `ValidConditionFailedVersionInvalidReason`, ensuring the operator rejects malformed version specifications before attempting broker deployment.

<u>[Version Validation Not Found](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L796)</u>

This test checks handling of valid but unavailable version strings. It deploys with Version="9.9.9" (syntactically valid but not in supported versions). The test verifies ValidConditionType is False with reason `ValidConditionFailedVersionNotFoundReason`, confirming the operator detects when specified versions don't exist in the supported versions catalog.

<u>[Version Handling When Images Are Explicitly Specified](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L818)</u>

This test validates that when both Image and InitImage are explicitly provided (paired), no version is required. It deploys a broker with only custom images, no version string. The test confirms ValidConditionType is True, allowing users to fully control image selection without version-based resolution.

<u>[Specify Only Major Version](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L863)</u>

This test validates broker deployment with only a major version specified (e.g., "7"). It deploys with just the major version number, allowing the operator to select the latest patch within that major version. The test confirms the broker deploys successfully and the StatefulSet is created, verifying the operator's version resolution selects appropriate images for major-only version specifications.

<u>[Specify Only Major.Minor Version](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L911)</u>

This test validates deployment with major.minor version specifications. It deploys with a version like "7.10", allowing the operator to select the latest patch version. The test verifies successful broker and StatefulSet creation, confirming the operator correctly resolves major.minor versions to specific patch versions with corresponding images.

<u>[Default Broker Versions](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L960)</u>

This test validates default version behavior when no version is specified. It deploys a broker CR without setting Spec.Version, allowing the operator to use its default version. The test confirms the broker deploys successfully with a StatefulSet created, verifying the operator's default version fallback mechanism works correctly.

<u>[OK With Valid=Unknown, Relaxed Version Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1006)</u>

This test validates that brokers can deploy even when ValidConditionType is Unknown. It deploys with a major-only version and custom images (triggering Unknown validation status). The test confirms that despite the Unknown validity status, the broker still deploys successfully with a StatefulSet created, showing the operator's relaxed validation allows deployment when validation cannot be definitively determined.

##### Resource Tracking:

Tests validating that operator-managed resources (secrets, services, etc.) maintain proper ownership references and are tracked correctly for updates and garbage collection.

###### Broker Resource Tracking

Validates credential secret generation, ownership tracking, and proper mounting into broker pods.

<u>[Default User Credential Secret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1070)</u>

This test validates automatic credential secret generation when no credentials are provided. It deploys a broker without explicit credentials, then verifies a secret named `{broker-name}-credentials-secret` is created containing four keys: AMQ_USER, AMQ_PASSWORD, AMQ_CLUSTER_USER, and AMQ_CLUSTER_PASSWORD. The test confirms the secret has proper owner references pointing to the broker CR and that the StatefulSet's init container references these secrets via environment variables. After updates, it verifies the secrets remain stable (same values) across reconciliation cycles.

<u>[Default User Credential Secret with Values in CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1231)</u>

This test validates credential secret generation when explicit admin username and password are provided in the CR. It deploys a broker with Spec.AdminUser and Spec.AdminPassword set, then verifies the credential secret is created with these values for AMQ_USER and AMQ_PASSWORD, while AMQ_CLUSTER_USER and AMQ_CLUSTER_PASSWORD are still auto-generated. The test confirms the specified credentials are properly propagated to the secret and mounted in the StatefulSet.

<u>[User Credential Secret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1300)</u>

This test validates using a pre-existing external credential secret. It creates a custom secret before deploying the broker, then sets Spec.DeploymentPlan.CredentialsSecret to reference this external secret. The test verifies the operator does not create a default credential secret and instead uses the specified external secret, confirming the StatefulSet's init container references the correct external secret name for all credential environment variables.

<u>[User Credential Secret with CR Values Provided](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1462)</u>

This test validates the combination of an external credential secret with explicit admin values in the CR. It creates an external secret, deploys a broker referencing that secret while also specifying AdminUser/AdminPassword in the CR. The test verifies the operator respects the external secret reference and does not attempt to create or manage credential secrets, deferring entirely to the user-provided secret configuration.

##### Address Settings:

Tests validating the configuration of address-level policies that control message handling behavior such as size limits, expiry, and redistribution.

###### New Address Settings Options

Validates that address settings configuration options are properly converted to broker configuration.

<u>[Deploy Broker with New Address Settings](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1618)</u>

This test validates comprehensive address settings configuration including newer options. It deploys a broker with AddressSettings containing two address patterns ("abc#" and "#") with multiple settings: DeadLetterAddress, MaxSizeBytes ("10m"), MaxSizeMessages (5000), ConfigDeleteDiverts ("OFF"), and RedistributionDelay (5000ms). The test verifies these settings are properly converted to YAML configuration in the init container args, specifically checking for "max_size_messages: 5000" and the shell-safe substitution of "config_delete_diverts: OFF". In existing cluster mode, it confirms the broker starts successfully with the settings applied.

##### Versions Test:

Tests validating the image and version resolution logic for determining which broker images to use.

###### Versions Test

Validates default image and version resolution when no explicit values are specified.

<u>[Default Image to Use Latest](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1702)</u>

This test validates the image selection logic when no explicit images are specified in the CR. It calls `DetermineImageToUse` for both "Kubernetes" and "Init" image types without setting Spec.DeploymentPlan.Image or Spec.DeploymentPlan.InitImage, verifying the functions return the default Kubernetes and init images from the version package. It also validates `DetermineCompactVersionToUse` returns the default compact version string. This ensures the operator has sensible defaults for image selection.

##### Version Conversion:

Tests validating that the operator can handle CRs from different API versions and migrate configurations to the current version format.

###### CR Version Conversion Test

Validates cross-version reconciliation and configuration migration through brokerProperties.

<u>[Can Reconcile Different Version](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1717)</u>

This test validates that a broker deployed with an older API version can be successfully reconciled with the current operator version. It deploys a broker CR, waits for it to be ready, then verifies reconciliation continues to work successfully. The test confirms the operator maintains backward compatibility with CRs created using earlier API versions, ensuring smooth upgrades.

<u>[Can Reconcile Latest Version with Config in brokerProperties](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1776)</u>

This test validates that advanced configurations specified via brokerProperties work correctly with the latest API version. It deploys a broker with configuration specified in brokerProperties (not higher-level CR fields), verifies the broker becomes ready, and confirms the operator properly handles configuration defined at the low-level brokerProperties level rather than through typed CR fields.

##### Console Configuration:

Tests validating web console deployment, exposure via ingress/routes, SSL configuration, and Jolokia metrics integration.

###### Console Config Test

Validates console service configuration, Jolokia integration, and SSL exposure mechanisms.

<u>[Checking Console Target Service Port Name for Metrics](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1820)</u>

This test validates that console services are configured with correct ports for metrics scraping. It deploys a 2-replica broker with console exposed and Jolokia enabled, then verifies the StatefulSet containers expose two ports: "jolokia" (8778) and "wconsj" (8161). For each broker replica, it validates that individual services are created with two service ports: one targeting port 8161 for console access and another targeting port 8162, both configured for TCP. This ensures Prometheus or other monitoring systems can scrape Jolokia metrics through properly named service ports.

<u>[Exposing Secured Console](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L1890)</u>

This test validates SSL-secured console exposure on actual clusters (OpenShift and Kubernetes). It creates a TLS secret with the expected name pattern `{broker-name}-console-secret`, deploys a broker with Console.Expose=true and Console.SSLEnabled=true, then verifies platform-specific exposure: On OpenShift, it checks a Route is created with TLS passthrough termination and proper host configuration. On Kubernetes, it verifies an Ingress with the "nginx" class is created with TLS configuration referencing the console secret. The test confirms SSL configuration is properly mounted and the console is accessible via HTTPS.

<u>[Exposing Secured Broker with Custom Ingress Hosts](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2013)</u>

This test validates custom hostname configuration for exposed console and acceptors. It deploys a broker with custom IngressDomain, exposed SSL-enabled console, and an exposed SSL acceptor, then verifies that the generated route/ingress resources use the specified custom ingress domain in their hostnames rather than defaulting to cluster defaults. This ensures users can control the external DNS names used for accessing broker services.

<u>[Exposing Secured Console with User Specified Secret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2197)</u>

This test validates using a custom-named TLS secret for console SSL instead of the default naming convention. It creates a TLS secret with a custom name, deploys a broker with Console.SSLSecret explicitly specifying this custom secret name, and verifies the console deployment succeeds with SSL properly configured using the specified secret. This allows users to manage SSL certificates independently and reuse existing secrets.

##### Expose Modes:

Tests validating different exposure strategies (Ingress, Route) for making broker services accessible outside the cluster, with and without TLS.

###### Expose Mode Test

Validates different exposure methods across Kubernetes (Ingress) and OpenShift (Route) platforms.

<u>[Expose with Ingress Mode](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2333)</u>

This test validates HTTP (non-SSL) exposure via Kubernetes Ingress. It deploys a broker with an exposed acceptor without SSL, then verifies an Ingress resource is created with the "nginx" ingress class, proper hostname based on the IngressDomain, and HTTP-only configuration (no TLS block). The test confirms the Ingress routes traffic to the correct service and port for the acceptor.

<u>[Expose with Secure Ingress Mode](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2491)</u>

This test validates HTTPS/SSL exposure via Kubernetes Ingress. It creates a TLS secret, deploys a broker with an SSL-enabled exposed acceptor, and verifies the Ingress is created with a TLS configuration block referencing the acceptor's TLS secret and proper secure hostname. This ensures SSL-secured acceptors are correctly exposed through Ingress with certificate termination.

<u>[Expose with Route Mode](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2648)</u>

This test validates OpenShift Route-based exposure. It deploys a broker with an exposed acceptor on OpenShift and verifies a Route resource is created with proper configuration including the target service, port, and hostname. This confirms the operator correctly uses OpenShift Routes instead of Ingress on OpenShift clusters.

##### Console Secret Validation:

Tests validating SSL secret validation for console configuration.

###### SSL Console Secret Validation

Validates that invalid or missing SSL secrets are detected and reported.

<u>[SSL Console Secret Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2764)</u>

This test validates that when Console.SSLEnabled is true but no valid SSL secret exists, the operator properly detects and reports this configuration error. It deploys a broker with SSL enabled but without creating the required TLS secret, then verifies the broker CR status is updated with a validation error condition indicating the missing secret. This ensures users are notified of invalid SSL configurations before attempting deployment.

##### Image Management:

Tests validating operator behavior when broker image pulls fail and how recovery is handled through CR updates or deletions.

###### Image Update Test

Validates that broken deployments due to image pull failures can be recovered through CR modifications.

<u>[Deploy ImagePullBackOff Update Delete OK](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L2861)</u>

This test validates recovery from ImagePullBackOff errors. It deploys a broker with an invalid image reference that will fail to pull, verifies the StatefulSet enters ImagePullBackOff state, then updates the broker CR with a valid image and confirms the broker recovers and deploys successfully. Finally, it verifies clean deletion works even after recovery. This ensures the operator can handle and recover from image pull failures through CR updates.

##### Client ID Auto-Sharding:

Tests validating automatic distribution of client IDs across broker instances for load balancing connections in a cluster.

###### ClientID Autoshard Test

Validates the client ID auto-sharding feature that distributes MQTT/AMQP client connections across broker replicas.

<u>[Deploy 2 with ClientID Auto Sharding](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3055)</u>

This test validates the client ID auto-sharding feature for multi-replica deployments. It deploys a 2-replica broker and verifies that the broker configuration includes client ID auto-sharding settings, which distribute client connections across broker instances based on client ID hashes. The test confirms the operator properly configures this load-balancing mechanism for MQTT and AMQP protocols in clustered deployments.

##### Probe Configuration:

Tests validating Kubernetes liveness, readiness, and startup probe configurations with default values and override behavior.

###### Probe Defaults Reconcile

Validates default probe configuration and minimal update detection logic.

<u>[Deploy](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3179)</u>

This test validates that brokers deploy with sensible default probe configurations when none are explicitly specified. It deploys a broker without probe overrides and verifies the StatefulSet includes default liveness, readiness, and startup probes with reasonable timeout, period, and threshold values. This ensures brokers have health checking configured by default for Kubernetes pod lifecycle management.

<u>[Verify Minimal Updates via Unequal](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3307)</u>

This test validates the operator's ability to detect meaningful changes vs cosmetic differences in probe configurations. It compares probe configurations that are functionally equivalent but represented differently, verifying the operator correctly identifies when actual updates are needed versus when configurations are semantically identical. This prevents unnecessary pod restarts from reconciliation loops.

##### StatefulSet Recovery:

Tests validating that the operator can detect and recover from StatefulSet deletion scenarios.

###### SS Delete Recreate Test

Validates automatic StatefulSet recreation when the underlying StatefulSet is deleted.

<u>[Deploy, Delete SS, Verify](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3351)</u>

This test validates the operator's self-healing capability for StatefulSet resources. It deploys a broker, waits for the StatefulSet to be created and running, manually deletes the StatefulSet, then verifies the operator detects the deletion and automatically recreates the StatefulSet with the correct configuration. This ensures the operator maintains desired state even when critical resources are accidentally or intentionally deleted.

##### PVC Management:

Tests validating PersistentVolumeClaim lifecycle, persistence across broker deletions, and owner reference management for upgrades.

###### PVC No GC Test

Validates that PVCs persist across broker deletion and can be reused on redeployment.

<u>[Deploy, Verify, Undeploy, Verify, Redeploy, Verify](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3412)</u>

This test validates PVC persistence through broker lifecycle operations. It deploys a broker with persistence enabled, verifies PVCs are created, deletes the broker CR, confirms PVCs remain (not garbage collected with the broker), redeploys the same broker, and verifies it reattaches to the existing PVCs with data intact. This ensures PVCs survive broker deletions for data preservation and enables safe broker recreation.

###### PVC Upgrade Owner Reference Test

Validates automatic adoption of PVCs that lost owner references during operator upgrades.

<u>[Faking a Broker Deployment with Owned PVC](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3522)</u>

This test simulates an upgrade scenario where PVCs exist without proper owner references (common in older operator versions). It creates a PVC without owner refs, deploys a broker that should use that PVC, and verifies the operator detects the orphaned PVC and automatically adds owner references to establish proper ownership. This ensures smooth upgrades from older operator versions that didn't set owner references on PVCs.

##### Tolerations:

Tests validating Kubernetes pod toleration configurations for scheduling on tainted nodes.

###### Tolerations Existing Cluster

Validates toleration configuration and propagation to StatefulSet pod specs.

<u>[Toleration of Artemis](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3585)</u>

This test validates basic toleration configuration. It deploys a broker with tolerations specified in the CR, then verifies the StatefulSet's pod template includes the specified tolerations with correct keys, operators, effects, and optional toleration seconds. This ensures brokers can be scheduled on nodes with matching taints.

<u>[Toleration of Artemis Required Add/Remove Verify Status](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3685)</u>

This test validates dynamic toleration updates and status tracking. It deploys a broker, adds tolerations via CR update, verifies they're applied to the StatefulSet and broker status reflects readiness, then removes tolerations and confirms the StatefulSet is updated accordingly. This ensures toleration changes trigger proper reconciliation and pod recreation.

##### Console Secrets:

Tests validating internal console SSL secret management, owner references, and adoption of orphaned secrets.

###### Console Secret Test

Validates automatic secret generation for SSL consoles and secret adoption mechanisms.

<u>[Deploy Broker with SSL Enabled Console](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3799)</u>

This test validates automatic internal secret generation for SSL consoles when no external secret is specified. It deploys a broker with Console.SSLEnabled=true without specifying Console.SSLSecret, verifies the operator generates an internal secret with the pattern `{broker-name}-console-secret-{hash}`, and confirms the secret contains proper TLS certificate and key data.

<u>[Reconcile Verify Internal Secret Owner Ref](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3874)</u>

This test validates that internally-generated console secrets have proper owner references pointing to the broker CR. It deploys a broker with SSL console, retrieves the generated secret, and verifies the owner references are correctly set with the broker CR as the controller. This ensures proper garbage collection when the broker is deleted.

<u>[Reconcile Verify Adopt Internal Secret That Has Lost Owner Ref](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L3976)</u>

This test validates the operator's ability to adopt orphaned console secrets. It simulates a scenario where an internal console secret exists but has lost its owner reference (due to upgrade or manual edit), triggers reconciliation, and verifies the operator detects and re-establishes the owner reference. This prevents secret leaks and ensures proper resource tracking.

##### Security Contexts:

Tests validating PodSecurityContext configuration for controlling pod-level security settings.

###### PodSecurityContext Test

Validates pod security context propagation to StatefulSet specifications.

<u>[Setting the Pods PodSecurityContext](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4057)</u>

This test validates pod security context configuration. It deploys a broker with Spec.DeploymentPlan.PodSecurityContext specifying RunAsUser, FSGroup, and other security settings, then verifies the StatefulSet's pod template includes the exact security context specification. This ensures security policies like user/group IDs and filesystem permissions are properly enforced on broker pods.

##### Affinity Configuration:

Tests validating pod affinity, anti-affinity, and node anti-affinity rules for controlling pod placement and distribution.

###### Affinity Test

Validates various affinity and anti-affinity configurations for pod scheduling.

<u>[Setting Pod Affinity](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4172)</u>

This test validates pod affinity configuration that attracts pods to nodes with specific labels. It deploys a broker with PodAffinity rules specifying preferred pod affinity terms with topology keys and label selectors, then verifies the StatefulSet's pod template includes the affinity specification. This enables co-location of broker pods with other specific workloads.

<u>[Setting Pod AntiAffinity](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4254)</u>

This test validates pod anti-affinity configuration that spreads pods across nodes. It deploys a broker with PodAntiAffinity rules specifying required anti-affinity terms to prevent multiple broker pods on the same node, then verifies the StatefulSet includes these rules. This ensures high availability by preventing single-point-of-failure node issues.

<u>[Setting Node AntiAffinity](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4333)</u>

This test validates node anti-affinity rules that prevent pods from scheduling on nodes with specific labels. It deploys a broker with NodeAntiAffinity configuration, then verifies the StatefulSet's affinity rules include node selector terms that exclude nodes matching certain labels. This allows steering broker pods away from unsuitable nodes.

##### Node Selector:

Tests validating node selector configuration for constraining pods to specific nodes.

###### Node Selector Test

Validates node selector label matching for pod placement.

<u>[Passing in 2 Labels](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4419)</u>

This test validates node selector configuration with multiple labels. It deploys a broker with Spec.DeploymentPlan.NodeSelector containing two key-value pairs, then verifies the StatefulSet's pod template includes the nodeSelector map with both labels. This ensures broker pods only schedule on nodes matching all specified label criteria.

##### Annotations:

Tests validating custom annotation propagation to broker pods.

###### Annotations Test

Validates pod annotation configuration and propagation.

<u>[Add Some Annotations](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4491)</u>

This test validates custom annotation application to broker pods. It deploys a broker with Spec.DeploymentPlan.Annotations containing custom key-value pairs, then verifies the StatefulSet's pod template metadata includes these annotations. This allows users to add metadata for monitoring, security policies, or other pod-level integrations that rely on annotations.

##### Labels:

Tests validating custom label application to broker pods beyond operator-managed labels.

###### Labels Test

Validates pod label configuration and limits.

<u>[Passing in 2 Labels](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4595)</u>

This test validates custom label application to broker pods. It deploys a broker with Spec.DeploymentPlan.Labels containing two custom key-value pairs, then verifies the StatefulSet's pod template includes these labels in addition to operator-managed labels. This enables pod selection and organization through custom labeling schemes.

<u>[Passing in 8 Labels](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4636)</u>

This test validates handling of multiple custom labels. It deploys a broker with eight custom labels, verifying all are properly applied to the StatefulSet pod template. This confirms the operator doesn't have artificial limits on the number of custom labels users can apply to broker pods.

##### Namespace Isolation:

Tests validating that the operator correctly isolates resources to their respective namespaces and doesn't interfere with other namespaces.

###### Different Namespace, Deployed Before Start

Validates namespace-scoped reconciliation and resource isolation.

<u>[Verify Reconcile in Own Namespace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4669)</u>

This test validates namespace isolation by deploying a broker in a different namespace before the test suite starts, then verifying the operator only reconciles brokers in its watched namespace and doesn't interfere with brokers in other namespaces. This ensures multi-tenant clusters can run multiple operator instances safely without cross-namespace conflicts.

##### Tolerations Configuration:

Additional tests for toleration configuration patterns.

###### Tolerations Test

Validates multiple toleration configurations with different operators and effects.

<u>[Passing in 2 Tolerations](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4738)</u>

This test validates multiple toleration configurations. It deploys a broker with two different tolerations (one with operator=Exists, another with operator=Equal and a specific value), then verifies both are correctly applied to the StatefulSet's pod template. This ensures brokers can tolerate multiple different node taints simultaneously for flexible pod placement.

##### Liveness Probes:

Tests validating Kubernetes liveness probe configurations that determine when to restart unhealthy containers, with support for Exec, TCPSocket, HTTPGet, and gRPC probe types.

###### Liveness Probe Tests

Validates liveness probe override mechanisms and default configurations across different probe handler types.

<u>[Override Liveness Probe No Exec](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4831)</u>

This test validates overriding the default liveness probe configuration without using Exec handler. It deploys a broker with custom LivenessProbe settings specifying InitialDelaySeconds, PeriodSeconds, TimeoutSeconds, SuccessThreshold, and FailureThreshold, then verifies the StatefulSet's liveness probe matches these custom values. This ensures users can fine-tune probe timing parameters for their specific deployment requirements.

<u>[Override Liveness Probe Exec](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4921)</u>

This test validates Exec-based liveness probes. It deploys a broker with LivenessProbe.Exec.Command specified (custom shell command to check broker health), then verifies the StatefulSet's liveness probe uses the Exec handler with the provided command. This allows users to implement custom health checks through shell scripts.

<u>[Default Liveness Probe](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L4973)</u>

This test validates the operator's default liveness probe when no override is specified. It deploys a broker without LivenessProbe configuration and verifies the StatefulSet includes a default liveness probe with reasonable timeout, period, and threshold values. This ensures brokers have basic health checking even without explicit configuration.

<u>[Override Liveness Probe Default TCPSocket](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5011)</u>

This test validates TCPSocket-based liveness probes with default settings. It deploys a broker with LivenessProbe.TCPSocket specified (checking a specific port), then verifies the StatefulSet's liveness probe uses the TCPSocket handler targeting the correct port. This enables simple TCP connectivity checks for broker health.

<u>[Override Liveness Probe Default HTTPGet](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5055)</u>

This test validates HTTPGet-based liveness probes. It deploys a broker with LivenessProbe.HTTPGet specified (with path and port), then verifies the StatefulSet's liveness probe uses the HTTPGet handler targeting the specified endpoint. This enables HTTP-based health check endpoints for determining broker liveness.

<u>[Override Liveness Probe Default GRPC](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5097)</u>

This test validates gRPC-based liveness probes (Kubernetes 1.24+). It deploys a broker with LivenessProbe.GRPC specified (with port and optional service name), then verifies the StatefulSet's liveness probe uses the GRPC handler. This enables native gRPC health checks for brokers exposing gRPC interfaces.

##### Readiness Probes:

Tests validating Kubernetes readiness probe configurations that determine when pods are ready to receive traffic, with support for various probe handler types.

###### Readiness Probe Tests

Validates readiness probe configurations including overrides and default behavior.

<u>[Override Readiness Probe No Exec](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5141)</u>

This test validates custom readiness probe timing configuration without Exec handler. It deploys a broker with ReadinessProbe specifying custom InitialDelaySeconds, PeriodSeconds, and other timing parameters, then verifies the StatefulSet's readiness probe uses these values. This allows users to adjust when pods are marked ready based on their startup characteristics.

<u>[Override Readiness Probe Exec](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5186)</u>

This test validates Exec-based readiness probes. It deploys a broker with ReadinessProbe.Exec.Command specified, then verifies the StatefulSet's readiness probe uses the Exec handler with the provided command for determining pod readiness through custom scripts.

<u>[Override Readiness Probe GRPC](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5228)</u>

This test validates gRPC-based readiness probes. It deploys a broker with ReadinessProbe.GRPC configuration, then verifies the StatefulSet applies the GRPC handler for readiness checks. This allows brokers with gRPC interfaces to use native gRPC health checking protocol for readiness.

<u>[Default Readiness Probe](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5268)</u>

This test validates default readiness probe configuration. It deploys a broker without ReadinessProbe override and verifies the StatefulSet includes a default readiness probe with appropriate timing and handler configuration. This ensures pods have readiness checking configured automatically.

##### Startup Probes:

Tests validating Kubernetes startup probe configurations that provide extra time for slow-starting containers before liveness checks begin.

###### Startup Probe Tests

Validates startup probe configuration for handling slow broker initialization.

<u>[Startup Probe with Exec](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5306)</u>

This test validates Exec-based startup probes. It deploys a broker with StartupProbe.Exec.Command specified with custom timing parameters including longer FailureThreshold for slow startups, then verifies the StatefulSet's startup probe uses the Exec handler and timing values. This allows accommodating brokers with long initialization times through custom startup checks.

<u>[Default Startup Probe](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5374)</u>

This test validates default startup probe configuration. It deploys a broker without StartupProbe override and verifies the StatefulSet includes a default startup probe with reasonable timing for typical broker startup scenarios. This ensures proper startup detection without explicit configuration.

##### Status Tracking:

Tests validating that the operator correctly tracks and reports broker status including pod states and conditions.

###### Status

Validates status reporting and pod descriptor information in CR status.

<u>[Expect Pod Desc](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5406)</u>

This test validates that broker CR status correctly reports pod descriptor information. It deploys a broker and verifies the status includes PodStatus with detailed information about each broker pod including names, states, and other descriptors. This ensures users can monitor broker pod states through the CR status without directly querying pods.

##### Validation:

Tests validating CR configuration validation logic that prevents invalid or incompatible configurations.

###### Validation

Validates various configuration validation scenarios and platform-specific constraints.

<u>[With Test Labels](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5474)</u>

This test validates that label-based test configurations work correctly with validation logic. It deploys a broker with test-specific labels and verifies the validation logic doesn't incorrectly flag valid test configurations as errors. This ensures the validation system is robust and doesn't produce false positives during testing.

<u>[With Ingress in OpenShift](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5569)</u>

This test validates that attempting to use Ingress resources on OpenShift platforms is properly detected and reported as invalid. It deploys a broker on OpenShift with Ingress-specific configuration and verifies the operator reports a validation error indicating Routes should be used instead of Ingress on OpenShift. This guides users to use platform-appropriate exposure mechanisms.

##### Rolling Updates:

Tests validating triggered rolling updates through environment variable checksums and configuration changes.

###### Env Var Updates TRIGGERED_ROLL_COUNT Checksum

Validates that configuration changes trigger pod rolling updates via checksum environment variables.

<u>[Expect TRIGGERED_ROLL_COUNT Count Non 0](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5611)</u>

This test validates the rolling update trigger mechanism via environment variable checksums. It deploys a broker, makes a configuration change, and verifies that the TRIGGERED_ROLL_COUNT environment variable or checksum is updated to trigger a rolling restart of broker pods. This confirms the operator properly propagates configuration changes to running brokers through pod recreation.

##### Broker Properties:

Tests validating low-level broker configuration via brokerProperties that bypass higher-level CR fields for direct artemis.xml configuration.

###### BrokerProperties

Validates brokerProperties configuration management including validation, updates, and coexistence of multiple brokers.

<u>[Expect Vol Mount via Config Map](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5692)</u>

This test validates that brokerProperties are stored in a ConfigMap and properly mounted into broker pods. It deploys a broker with brokerProperties specified, verifies a ConfigMap is created containing the properties, and confirms the StatefulSet mounts this ConfigMap at the correct path for the broker to consume.

<u>[Expect Updated Secret on Update to BrokerProperties](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5787)</u>

This test validates that changes to brokerProperties trigger ConfigMap updates and pod restarts. It deploys a broker, updates the brokerProperties in the CR, and verifies the ConfigMap is updated and pods are restarted to apply the new configuration. This ensures configuration changes are propagated correctly.

<u>[Upgrade BrokerProps Respect Existing Immutable ConfigMap](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5849)</u>

This test validates handling of immutable ConfigMaps during upgrades. It simulates an older operator version that created immutable ConfigMaps, then upgrades and verifies the operator handles the immutable ConfigMap correctly without errors. This ensures smooth upgrades across operator versions.

<u>[Expect Two CRs to Coexist](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5936)</u>

This test validates that multiple broker CRs can coexist with different brokerProperties without conflicts. It deploys two separate brokers with different configurations and verifies each maintains its own ConfigMap and configuration independently. This ensures multi-tenant broker deployments work correctly.

<u>[Expect Error Message on Invalid Property Value](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L5972)</u>

This test validates that invalid brokerProperties are detected and reported with clear error messages. It deploys a broker with a malformed property value and verifies the CR status shows a validation error with details about the invalid configuration. This guides users to correct configuration mistakes.

##### Broker Version:

Tests validating broker version resolution and compatibility checks between specified versions and deployed images.

###### BrokerVersion

Validates version matching logic and error reporting for version-image incompatibilities.

<u>[Expect Version Match When Version Is Loosly Specified](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6011)</u>

This test validates loose version matching where only major or major.minor is specified. It deploys with a partial version string and verifies the operator resolves it to the latest matching patch version and deploys with the correct images. This allows flexible version specifications.

<u>[Expect Version Match on Latest Version](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6044)</u>

This test validates deployment with the latest available version. It deploys without specifying a version and verifies the operator selects and deploys the latest supported broker version with corresponding images. This ensures users get current versions by default.

<u>[Expect Error Message on Wrong Image Version](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6077)</u>

This test validates detection of version-image mismatches. It deploys with a specified version but incompatible custom images and verifies the operator reports a validation error explaining the mismatch. This prevents deployment of incompatible version-image combinations.

##### Logger Properties:

Tests validating custom logging configuration via external ConfigMaps/Secrets and validation of protected environment variables.

###### LoggerProperties

Validates logging configuration management, environment variable protection, and volume mounting for custom logging configs.

<u>[Test Validate Can Pick Up All Internal Vars Misusage](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6124)</u>

This test validates that the operator detects when users inadvertently reference internal/protected environment variables in their configurations. It attempts to use reserved variable names and verifies validation errors are raised listing all the misused variables. This protects operator-managed environment variables from user override.

<u>[Validate User Directly Using Internal Env Vars](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6180)</u>

This test validates detection of direct use of internal environment variables in the Env field. It deploys a broker with Env variables that conflict with operator-managed variables and verifies validation errors are reported. This ensures users can't bypass protection by directly setting internal variables.

<u>[Custom Logging Not to Override JAVA_ARGS_APPEND](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6241)</u>

This test validates that custom logging configurations don't override the JAVA_ARGS_APPEND environment variable. It deploys with a custom logging ConfigMap and verifies the operator properly merges logging configuration without losing other Java arguments. This ensures logging customization doesn't break other broker functionality.

<u>[Logging ConfigMap Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6328)</u>

This test validates that missing logging ConfigMaps are detected. It references a non-existent logging ConfigMap and verifies the operator reports a validation error. This catches configuration errors before deployment.

<u>[Logging Secret and ConfigMap Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6368)</u>

This test validates that only one logging source (ConfigMap or Secret, not both) can be specified. It attempts to configure both simultaneously and verifies a validation error is raised. This prevents ambiguous logging configuration.

<u>[Expect Vol Mount for Logging ConfigMap Deployed](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6465)</u>

This test validates that custom logging ConfigMaps are properly mounted. It deploys with a logging ConfigMap reference and verifies the StatefulSet mounts the ConfigMap at the correct path for the broker's logging system to consume.

<u>[Expect Vol Mount for Logging Secret Deployed](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6532)</u>

This test validates that custom logging Secrets are properly mounted. It deploys with a logging Secret reference and verifies the StatefulSet mounts the Secret at the correct path, allowing sensitive logging configurations to be stored securely.

##### Address Settings:

Additional tests for address settings configuration updates.

###### With Address Settings via Updated CR

Validates address settings application through CR updates.

<u>[Expect OK Deploy](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6601)</u>

This test validates that address settings can be added or updated after initial deployment. It deploys a broker, updates the CR to add address settings, and verifies the configuration is applied and the broker becomes ready with the new settings. This confirms dynamic address settings management.

##### Address CR:

Tests validating interaction between broker CR and Address CR resources.

###### With Address CR

Validates that Address CRs work correctly alongside broker deployments.

<u>[Expect OK Deploy](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6720)</u>

This test validates that Address CRs can be deployed alongside a broker and properly managed. It deploys a broker and creates Address CRs targeting that broker, verifying the addresses are created in the broker and the Address CR statuses reflect successful creation. This confirms the Address controller integration.

##### Persistence Toggle:

Tests validating dynamic persistence configuration changes and reconciliation behavior.

###### With Toggle Persistence=True

Validates enabling persistence on an existing broker.

<u>[Expect OK Redeploy](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6800)</u>

This test validates that persistence can be enabled on a running broker. It deploys a broker without persistence, updates the CR to enable persistence (Spec.DeploymentPlan.PersistenceEnabled=true), and verifies PVCs are created and the broker is redeployed with persistent storage. This confirms dynamic persistence configuration.

###### With Update Persistence=True Single Reconcile

Validates that persistence changes don't cause excessive reconciliation loops.

<u>[Check Reconcile Loop](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6859)</u>

This test validates that enabling persistence doesn't trigger infinite reconciliation loops. It monitors reconciliation events during a persistence toggle and confirms the operator reaches a stable state after a reasonable number of reconciliations. This ensures configuration changes converge efficiently.

##### Version Toggle:

Tests validating broker version updates and StatefulSet generation management.

###### Toggle spec.Version

Validates changing broker versions through CR updates.

<u>[Expect OK Update and New SS Generation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L6953)</u>

This test validates broker version updates. It deploys a broker with one version, updates Spec.Version to a different version, and verifies the StatefulSet is updated with new images corresponding to the new version, and the StatefulSet generation number increments. This confirms the operator properly handles version upgrades and downgrades.

##### Performance:

Tests validating operator performance characteristics and timing.

###### Time to Ready

Validates broker startup time and readiness detection performance.

<u>[Expect OK Ready Is Fast](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7064)</u>

This test validates that brokers reach ready state within acceptable time limits. It deploys a broker and measures the time until the Ready condition becomes true, asserting it's within performance thresholds. This ensures the operator doesn't introduce unnecessary delays in broker startup.

##### Resource Templates:

Tests validating custom resource templates that allow advanced customization of operator-generated resources.

###### Template Tests

Validates resource template customization capabilities.

<u>[Expect Custom Service](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7117)</u>

This test validates that custom Service templates can override operator-generated Services. It deploys a broker with a custom Service template specified in Spec.ResourceTemplates and verifies the created Service matches the custom template rather than the operator's default. This allows advanced users to customize Kubernetes resources beyond standard CR fields.

##### Acceptor Management:

Tests validating acceptor configuration including dynamic updates, exposure toggling, SSL configuration, and bind address options.

###### With Deployed Controller - Acceptor

Validates comprehensive acceptor management scenarios.

<u>[Add Acceptor via Update](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7215)</u>

This test validates dynamic acceptor addition. It deploys a broker, updates the CR to add new acceptors, and verifies the StatefulSet configuration is updated to include the new acceptors with proper Service and Route/Ingress resources created. This confirms acceptors can be added to running brokers.

<u>[Checking Acceptor Service and Route/Ingress While Toggle Expose](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7289)</u>

This test validates expose toggling for acceptors. It deploys a broker with exposed acceptors, verifies Services and Routes/Ingress are created, toggles Expose to false, verifies external resources are deleted, then toggles back to true and confirms resources are recreated. This ensures dynamic exposure control works correctly.

<u>[Testing Acceptor bindToAllInterfaces Default](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7585)</u>

This test validates default bindToAllInterfaces behavior. It deploys an acceptor without explicitly setting bindToAllInterfaces and verifies the broker configuration uses the default value (0.0.0.0 bind address). This ensures sensible defaults for network binding.

<u>[Testing Acceptor bindToAllInterfaces Being False](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7637)</u>

This test validates explicit bindToAllInterfaces=false configuration. It deploys an acceptor with bindToAllInterfaces explicitly set to false and verifies the broker binds to specific interfaces rather than all interfaces. This allows fine-grained network configuration.

<u>[Testing Acceptor bindToAllInterfaces Being True](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7692)</u>

This test validates explicit bindToAllInterfaces=true configuration. It deploys an acceptor with bindToAllInterfaces=true and verifies the broker binds to all network interfaces (0.0.0.0). This ensures explicit all-interface binding works correctly.

<u>[Testing Acceptor keyStoreProvider Being Set](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7764)</u>

This test validates custom keyStoreProvider configuration for SSL acceptors. It deploys an SSL-enabled acceptor with a custom keyStoreProvider (e.g., "PKCS11") and verifies the broker configuration includes this provider setting. This allows integration with hardware security modules and custom keystore providers.

<u>[Testing Acceptor trustStoreType Being Set and Unset](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7819)</u>

This test validates trustStoreType configuration and removal. It deploys an SSL acceptor with a custom trustStoreType, verifies it's applied, then updates to remove the trustStoreType and confirms the broker configuration is updated accordingly. This ensures truststore type customization and defaults work correctly.

<u>[Testing Acceptor trustStoreProvider Being Set](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7927)</u>

This test validates custom trustStoreProvider configuration. It deploys an SSL acceptor with a custom trustStoreProvider specified and verifies the broker includes this setting in its SSL configuration. This allows advanced SSL/TLS configuration with custom trust providers.

##### Acceptor/Connector Updates:

Tests validating upgrade scenarios where acceptor/connector services lack owner references.

###### Acceptor/Connector Update No Owner Ref

Validates adoption of orphaned acceptor and connector services.

<u>[With Existing Acceptor Service](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7471)</u>

This test validates adoption of pre-existing acceptor services without owner references. It creates an acceptor service manually (simulating an upgrade from older operator versions), deploys a broker that should use that service, and verifies the operator adopts the service by adding proper owner references. This ensures smooth upgrades.

<u>[With Existing Connector Service](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7527)</u>

This test validates adoption of pre-existing connector services without owner references. It creates a connector service manually, deploys a broker with matching connector configuration, and verifies the operator adopts and manages the service properly. This handles upgrade scenarios where services existed before proper ownership tracking.

##### API Version Compatibility:

Tests validating cross-API-version compatibility and conversion webhooks.

###### With Deployed Controller

Validates backward compatibility with older API versions.

<u>[Verify Old Ver Support](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L7984)</u>

This test validates v2alpha4 to v1beta1 API version conversion. It creates a broker CR using the v2alpha4 API version and verifies the operator correctly converts and reconciles it as v1beta1, ensuring all functionality works correctly across API version boundaries. This ensures long-term backward compatibility.

##### Storage Configuration:

Tests validating PersistentVolumeClaim storage configuration including storage class and size specifications.

###### With Deployed Controller

Validates storage configuration and validation.

<u>[Checking storageClassName Is Configured](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8021)</u>

This test validates custom storageClassName configuration. It deploys a broker with Spec.DeploymentPlan.Storage.StorageClassName specified and verifies the created PVCs use the specified storage class. This allows users to select appropriate storage classes for their environments.

<u>[Checking Storage Size Mal Configured](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8084)</u>

This test validates detection of malformed storage size specifications. It deploys a broker with an invalid storage size value and verifies the operator reports a validation error. This catches configuration errors in storage sizing before deployment.

##### User Management:

Tests validating user credential management, validation, and security CR integration.

###### User Management Tests

Validates user credential handling across various scenarios.

<u>[populateValidatedUser](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8128)</u>

This test validates the populateValidatedUser function's handling of explicit admin credentials. It tests the logic that processes admin username/password from the CR and ensures proper validation and secret population. This validates the core user credential management logic.

<u>[populateValidatedUser as Auto Generated Guest](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8226)</u>

This test validates automatic guest user generation when no admin credentials are provided. It tests the populateValidatedUser function with empty credentials and verifies a guest user is automatically generated with appropriate permissions. This ensures brokers have usable authentication even without explicit configuration.

<u>[Credential Secret Manually Created](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8282)</u>

This test validates using manually pre-created credential secrets. It creates a credential secret before deploying the broker, references it in the CR, and verifies the operator uses the existing secret rather than generating a new one. This allows external credential management.

<u>[Deploy Security CR While Broker Is Not Yet Ready](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8331)</u>

This test validates Security CR handling when the target broker isn't ready. It deploys a Security CR before the broker is fully ready and verifies the Security controller properly waits and applies security configuration once the broker becomes available. This ensures robust security CR lifecycle management.

##### Lifecycle:

Tests validating broker lifecycle management including cascade deletion behavior.

###### Lifecycle Tests

Validates resource cleanup and cascade deletion.

<u>[Cascade Delete Foreground Test](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8482)</u>

This test validates foreground cascade deletion of broker resources. It deploys a broker, initiates deletion with foreground cascading enabled, and verifies all owned resources (StatefulSet, Services, PVCs, Secrets, etc.) are deleted in the proper order before the broker CR is removed. This ensures clean resource cleanup without orphaned resources.

##### Management RBAC:

Tests validating management RBAC configuration that controls broker management console access.

###### Management RBAC Tests

Validates management RBAC feature toggling.

<u>[managementRBACEnabled Is False](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8527)</u>

This test validates disabling management RBAC. It deploys a broker with Spec.DeploymentPlan.ManagementRBACEnabled=false and verifies the broker configuration doesn't include RBAC restrictions on management operations. This allows open management access for development or trusted environments.

##### Environment Variables:

Tests validating custom environment variable injection into broker pods.

###### Environment Variable Tests

Validates custom environment variable configuration.

<u>[Env Var](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8583)</u>

This test validates custom environment variable injection. It deploys a broker with Spec.Env containing custom environment variables and verifies the StatefulSet's containers include these variables with correct values. This allows users to configure broker behavior through environment variables beyond the CR's typed fields.

##### Metrics:

Tests validating metrics exporter configuration including JVM and broker metrics.

###### Metrics Tests

Validates metrics configuration via broker properties.

<u>[Enable JVM Metrics by Using Broker Properties](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8691)</u>

This test validates JVM metrics enablement through brokerProperties. It configures metrics collection using broker property settings and verifies the broker exposes JVM and broker metrics for Prometheus scraping. This confirms metrics integration works correctly via low-level broker configuration.

##### Config Projection:

Tests validating advanced configuration projection mechanisms including ordinal-specific properties, extraMounts for custom configurations, JAAS security integration, and per-broker customization in StatefulSet deployments.

###### Config Projection

Validates complex configuration projection scenarios for StatefulSet-based broker deployments.

<u>[Ordinal Broker Properties](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8747)</u>

This test validates ordinal-specific broker properties that apply to individual broker pods in a StatefulSet. It deploys a multi-replica broker with properties suffixed with `-broker-0`, `-broker-1`, etc., and verifies each pod receives its specific configuration. This enables per-broker customization within a single broker cluster.

<u>[Ordinal Broker Properties with Other Secret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8803)</u>

This test validates combining ordinal broker properties from multiple secrets. It creates separate secrets with ordinal-specific properties, references them in extraMounts, and verifies each broker pod receives the correct merged configuration from all sources. This allows complex multi-source configuration scenarios.

<u>[Invalid Ordinal Prefix Broker Properties](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8850)</u>

This test validates detection of malformed ordinal property names. It provides properties with invalid ordinal suffixes (e.g., wrong format, out-of-range ordinals) and verifies validation errors are reported. This catches configuration mistakes in ordinal property naming.

<u>[BrokerProperties with Escapes](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8882)</u>

This test validates proper escaping and handling of special characters in broker properties. It configures properties containing quotes, backslashes, and other special characters, and verifies they're correctly escaped and applied in the broker configuration without syntax errors.

<u>[Mod Ordinal Broker Properties with Error and Update](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L8931)</u>

This test validates error recovery in ordinal property updates. It deploys with invalid ordinal properties causing errors, updates the CR to fix the errors, and verifies the operator recovers and applies the corrected configuration. This ensures robustness in configuration management.

<u>[ExtraMount.ConfigMap Projection Update](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9017)</u>

This test validates ConfigMap-based extraMount updates. It deploys with an extraMount ConfigMap, updates the ConfigMap content, and verifies the broker pods detect and apply the updated configuration. This confirms dynamic configuration updates through ConfigMaps.

<u>[ExtraMount.Secret jaas-config Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9112)</u>

This test validates JAAS configuration via extraMount secrets. It mounts a secret containing jaas-config for authentication and verifies the operator properly validates and applies the JAAS configuration, enabling custom authentication realms.

<u>[ExtraMount jaas-config Once Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9194)</u>

This test validates that jaas-config can only be specified once. It attempts to mount multiple jaas-config files and verifies a validation error is raised. This prevents ambiguous JAAS configuration.

<u>[ExtraMount.Secret x-jaas-config Single Realm Update Status](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9259)</u>

This test validates single JAAS realm configuration and status updates. It deploys with a JAAS config defining one security realm, verifies the broker applies it, and confirms the CR status reflects successful JAAS configuration. This tests basic JAAS integration.

<u>[ExtraMount.Secret y-jaas-config Mgmt Realm OK Connect and Status Check](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9438)</u>

This test validates JAAS management realm configuration. It configures JAAS with a management realm for broker console/JMX access, attempts connections using JAAS-authenticated users, and verifies both connection success and CR status reporting. This confirms end-to-end JAAS management authentication.

<u>[CR.brokerProperties and -bp Duplicate Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9557)</u>

This test validates detection of duplicate properties between CR.brokerProperties and extraMount secrets with `-bp` suffix. It specifies the same property in both locations and verifies a validation error is raised to prevent conflicting configurations.

<u>[Onboarding - jaas-config New User Queue RBAC](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9668)</u>

This test validates complete JAAS-based RBAC scenario. It configures JAAS authentication, creates users with specific queue permissions, attempts operations with those users, and verifies RBAC is properly enforced. This demonstrates full JAAS security integration for user onboarding scenarios.

<u>[jaas-config Not Allowed in Config Map](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9782)</u>

This test validates that JAAS configuration must be in secrets, not ConfigMaps. It attempts to provide jaas-config via ConfigMap and verifies a validation error is raised. This enforces security best practices for sensitive authentication configuration.

<u>[jaas-config Syntax Check](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9821)</u>

This test validates JAAS configuration syntax checking. It provides malformed JAAS configuration and verifies the operator detects syntax errors and reports them in CR status. This catches JAAS configuration mistakes before deployment.

<u>[Ordinal Status - jaas-config](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9861)</u>

This test validates per-broker ordinal status reporting with JAAS configuration. It deploys with ordinal-specific JAAS configs and verifies the CR status reports each broker's JAAS configuration state individually. This provides fine-grained status visibility for StatefulSet brokers.

<u>[ExtraSecret with Broker Properties -bp Suffix](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L9985)</u>

This test validates the `-bp` suffix convention for broker properties in extraMount secrets. It creates secrets with keys ending in `-bp`, mounts them, and verifies the operator extracts and applies the properties to broker configuration. This enables property injection via secrets.

<u>[ExtraSecret with JSON Broker Properties -bp Suffix](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10065)</u>

This test validates JSON-formatted broker properties in extraMount secrets. It provides properties as JSON in secrets with `-bp` suffix and verifies proper parsing and application. This supports structured property configuration.

<u>[Complex ExtraSecret with JSON Broker Properties with Quote and -bp Suffix](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10176)</u>

This test validates complex JSON property values with special characters. It configures JSON properties containing quotes, nested structures, and escapes in `-bp` secrets, and verifies correct parsing and application. This ensures robust JSON property handling.

<u>[-bp Suffix Secret broker-n Support](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10306)</u>

This test validates ordinal-specific properties via `-bp` suffix in secrets. It creates secrets with keys like `property-bp-broker-0`, `property-bp-broker-1`, and verifies each broker pod receives its ordinal-specific properties. This combines ordinal targeting with secret-based injection.

<u>[-bp Suffix Secret broker-n JSON Support](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10411)</u>

This test validates ordinal-specific JSON properties in secrets. It provides JSON-formatted, ordinal-specific properties via `-bp` secrets and verifies proper parsing and ordinal-targeted application. This supports structured, per-broker configuration.

<u>[ExtraMount.ConfigMap Logging Config Manually](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10516)</u>

This test validates manual logging configuration via extraMount ConfigMaps as an alternative to the Spec.Logging field. It mounts a ConfigMap with logging configuration, verifies it's applied, and confirms logging behavior matches the custom configuration.

<u>[Secure Connections with Wildcard DNS Name](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10598)</u>

This test validates TLS configuration with wildcard DNS names in certificates. It creates TLS secrets with wildcard SANs (e.g., `*.example.com`), configures SSL acceptors, and verifies secure connections work correctly with wildcard certificate matching.

<u>[Secure Connections with OpenShift Serving Certificate](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10642)</u>

This test validates integration with OpenShift service serving certificates. It uses OpenShift's automatic certificate injection for services, configures SSL with these auto-generated certificates, and verifies secure broker connections work with platform-managed certificates.

<u>[Secure Connections with Multiple DNS Names](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10691)</u>

This test validates TLS with multiple Subject Alternative Names (SANs). It creates certificates with multiple DNS names, configures SSL acceptors, and verifies connections succeed using any of the certificate's DNS names. This supports multi-hostname broker access.

<u>[Secure Connections with Verify Host Disabled](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_test.go#L10736)</u>

This test validates disabling hostname verification for SSL connections. It configures SSL with verifyHost=false, uses certificates that don't match the connection hostname, and verifies connections succeed without hostname validation. This supports development and test environments with non-matching certificates.

**Test Labels in this file**:
- `tls-secret-reuse`
- `statefulset-options`
- `pod-disruption-budget`
- `broker-resource-tracking-context`
- `console-expose-metrics`
- `console-expose-ssl`
- `console`, `acceptor`, `connector`
- `ingress`, `route`, `ssl`
- `pvc-owner-reference-upgrade`
- `annotations-test`
- `with-ing-in-openshift`
- `LoggerProperties-test`
- `slow`

---

#### 2. [`activemqartemis_controller2_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller2_test.go#L44) (4 test scenarios)


Additional integration tests for broker controller, focusing on persistent volumes and external resources


##### Persistent Volumes:

Tests validating PVC management, external volume attachment, and controller resource recovery scenarios.

###### Persistent Volumes Tests

Validates PVC handling and controller recovery mechanisms.

<u>[Controller Resource Recover Test](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller2_test.go#L55)</u>

This test validates the controller's ability to recover and re-adopt resources after controller restarts or failures. It deploys a broker, simulates controller failure/restart, and verifies the controller rediscovers and resumes managing existing broker resources without disruption. This ensures the operator is resilient to controller pod failures.

<u>[External Volumes Attach](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller2_test.go#L175)</u>

This test validates attaching pre-existing external volumes to broker pods. It creates PVCs independently, references them in the broker CR's volume configuration, and verifies the StatefulSet mounts these external PVCs rather than creating new ones. This allows reusing existing storage or migrating data from external sources.

<u>[External PVC with PVC Templates](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller2_test.go#L276)</u>

This test validates using PVC templates to customize PVC specifications while still allowing external PVC references. It configures PVC templates with custom storage classes and sizes alongside external PVC references, verifying the operator respects both configurations appropriately. This enables fine-grained PVC customization.

##### Route Reconciliation:

Tests validating OpenShift Route management and reconciliation.

###### Route Reconcile

Validates Route resource lifecycle management on OpenShift.

<u>[Route Reconcile](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller2_test.go#L410)</u>

This test validates OpenShift Route reconciliation logic. It deploys a broker with exposed services on OpenShift, verifies Routes are created, manually modifies a Route to simulate drift, and confirms the operator detects and reconciles the Route back to the desired state. This ensures Routes maintain correct configuration despite manual changes.

---

#### 3. [`activemqartemis_controller_unit_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go) (6 test functions)

Unit tests for controller validation logic and status caching


##### Validation Tests:

Unit tests validating CR configuration validation logic including reserved labels and broker property duplicate detection.

###### Validation Tests

Validates configuration validation functions for catching user errors.

<u>[TestValidate](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go#L41)</u>

This unit test validates the detection of reserved operator-managed labels in user-provided resource templates. It attempts to specify labels like `ActiveMQArtemis` or `application` in Spec.ResourceTemplates and verifies the validation function rejects these reserved labels with appropriate error messages. This prevents users from interfering with operator-managed resource selection.

<u>[TestValidateBrokerPropsDuplicate](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go#L71)</u>

This unit test validates duplicate broker property key detection. It provides brokerProperties with duplicate keys and verifies the validation function detects and reports the duplicates. This catches configuration mistakes where the same property is specified multiple times with potentially conflicting values.

<u>[TestValidateBrokerPropsDuplicateOnFirstEquals](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go#L99)</u>

This unit test validates duplicate detection when property names contain escaped equals signs. It tests properties like `key\\=foo=value1` and `key\\=foo=value2` and verifies the validation correctly identifies them as duplicates despite the escaped equals in the name. This ensures robust property name normalization.

<u>[TestValidateBrokerPropsDuplicateOnFirstEqualsIncorrectButUnrealisticForOurBrokerConfigUsecase](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go#L127)</u>

This unit test validates edge case handling in property parsing. It tests an unrealistic but technically possible scenario with escaped equals followed by different suffixes, verifying the validation logic handles these edge cases without false positives or crashes. This ensures robustness in the validation code.

##### Status Caching Tests:

Unit tests validating that status checks are properly cached to avoid excessive API calls and improve performance.

###### Status Caching Tests

Validates caching mechanisms for status checks.

<u>[TestStatusPodsCheckCached](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go#L155)</u>

This unit test validates pod status check caching. It calls the status check function multiple times and verifies that subsequent calls use cached data rather than making redundant Kubernetes API calls. This ensures the operator is efficient and doesn't overload the API server with repeated identical queries.

<u>[TestJolokiaStatusCached](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_unit_test.go#L203)</u>

This unit test validates Jolokia status check caching. It calls Jolokia status functions multiple times and verifies caching prevents redundant Jolokia API calls. It also validates that errors are properly cached to avoid retry loops. This ensures efficient broker status monitoring.

**Coverage**:
- Resource template validation (reserved labels)
- BrokerProperties duplicate key detection
- Escaped character handling in property names
- Status check caching optimization
- Jolokia client error handling

---

#### 4. [`activemqartemis_controller_deploy_operator_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_deploy_operator_test.go#L40) (4 test scenarios)


Integration tests requiring full operator deployment in cluster


##### TLS Jolokia Access:

Tests validating secure Jolokia access with TLS and SNI configuration.

###### TLS Jolokia Access

Validates Jolokia client functionality over TLS connections.

<u>[Check the Util Works in Test Env](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_deploy_operator_test.go#L51)</u>

This test validates that the Jolokia utility functions work correctly in the test environment. It sets up basic Jolokia connectivity without TLS and verifies the utility can successfully connect and retrieve data. This ensures the test infrastructure is properly configured before testing TLS scenarios.

<u>[Get Status from Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_deploy_operator_test.go#L55)</u>

This test validates broker status retrieval via Jolokia over TLS with SNI (Server Name Indication). It deploys a broker with SSL-enabled console, configures Jolokia client for TLS with proper hostname verification, and successfully retrieves broker status through the encrypted connection. This confirms the operator's Jolokia integration works with production-grade secure configurations.

##### Operator Logging Configuration:

Tests validating operator logging level configuration via environment variables.

###### Operator Logging Config

Validates operator logging customization through deployment configuration.

<u>[Test Operator with Env Var](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_deploy_operator_test.go#L111)</u>

This test validates operator logging configuration via environment variables. It deploys the operator with custom LOG_LEVEL environment variable, triggers operator actions, and verifies log output matches the configured level. This ensures operators can adjust logging verbosity for debugging or production environments.

##### Custom Related Images:

Tests validating custom broker image configuration through related images.

###### Operator Deployment in Default Namespace

Validates broker version resolution with custom related images.

<u>[Default Broker Versions](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_deploy_operator_test.go#L175)</u>

This test validates broker version resolution when custom related images are configured via operator environment variables. It configures RELATED_IMAGE_* environment variables pointing to custom broker and init images, deploys a broker without explicit image specification, and verifies the operator uses the custom related images. This supports air-gapped environments and custom image registries.

##### Restricted Namespace:

Tests validating operator behavior in restricted namespace mode with minimal RBAC permissions.

###### Operator Deployment in Restricted Namespace

Validates restricted namespace operation mode.

<u>[Test in a Restricted Namespace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_deploy_operator_test.go#L228)</u>

This test validates the operator's restricted namespace mode where the operator has minimal permissions and watches only its own namespace. It deploys the operator in restricted mode, creates brokers, and verifies functionality works with reduced RBAC. This ensures the operator can run in security-conscious environments with least-privilege principles.

**Coverage**:
- TLS/SNI configuration for Jolokia
- Operator logging configuration
- Custom image resolution
- Restricted namespace mode

---

#### 5. [`activemqartemis_controller_cert_manager_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L68) (8 test scenarios)


Integration tests for cert-manager integration and certificate lifecycle management


##### Java Keystore Integration:

Tests validating Java keystore (JKS) format support for legacy broker configurations.

###### Cert-Manager Cert with Java Store

Validates cert-manager certificate conversion to Java keystore format.

<u>[Test Configured with Cert Secret as Legacy One](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L194)</u>

This test validates that cert-manager-generated certificates in PEM format can be converted to legacy Java keystore format for brokers requiring JKS. It uses cert-manager to issue a certificate, configures the broker to use Java keystore, and verifies the operator converts the PEM certificate to JKS format and properly mounts it for broker consumption. This supports legacy broker versions and Java-specific SSL configurations.

##### TLS Exposure:

Tests validating TLS configuration with cert-manager certificates for secure broker exposure.

###### TLS Exposure with Cert Manager

Validates cert-manager integration for broker TLS configuration.

<u>[Test Configured with Cert and CA Bundle](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L262)</u>

This test validates broker TLS configuration using cert-manager certificates with CA bundle. It provisions certificates with complete certificate chains, configures brokers to use them, and verifies the entire certificate chain is properly mounted and configured. This ensures proper certificate validation in production environments.

<u>[Test Console Cert Broker Status Access](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L267)</u>

This test validates accessing broker status through cert-manager-secured console. It deploys a broker with console TLS using cert-manager certificates, connects via Jolokia using the certificate, and successfully retrieves broker status. This confirms end-to-end cert-manager integration for secure management access.

<u>[Test SSL Args with Keystore Secrets Only](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L272)</u>

This test validates SSL configuration using only keystore secrets without PEM certificates. It configures brokers with keystore-only TLS (no PEM format), and verifies the broker accepts this configuration and establishes secure connections. This supports environments where certificates are managed exclusively in keystore format.

##### Certificate Utilities:

Tests validating certificate utility functions for detecting and handling cert-manager certificates.

###### Certutil Functions

Validates certificate detection and handling utilities.

<u>[Certutil - Is Secret from Cert](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L293)</u>

This unit test validates the certificate secret detection utility function. It creates secrets with cert-manager annotations/labels and secrets without, then calls the detection function to verify it correctly identifies cert-manager-managed secrets. This ensures the operator properly distinguishes cert-manager certificates from manually-created secrets.

##### Certificate Annotations:

Tests validating certificate provisioning via Kubernetes annotations.

###### Certificate from Annotations

Validates annotation-driven certificate provisioning for Ingress resources.

<u>[Ingress Certificate Annotations](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L321)</u>

This test validates automatic certificate provisioning for Ingress via cert-manager annotations. It creates Ingress resources with cert-manager annotations (`cert-manager.io/issuer`), verifies cert-manager automatically provisions certificates, and confirms the broker uses these auto-provisioned certificates for TLS. This demonstrates fully automated certificate lifecycle management.

##### Certificate Rotation:

Tests validating certificate rotation and automatic pod updates when certificates change.

###### Certificate Rotation

Validates certificate lifecycle including rotation and renewal.

<u>[Broker Certificate Rotation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L465)</u>

This test validates broker certificate rotation handling. It deploys a broker with a cert-manager certificate, triggers certificate renewal (simulating expiration), and verifies the operator detects the certificate change and triggers broker pod rolling restart to load the new certificate. This ensures certificates can be rotated without manual intervention.

<u>[Broker Issuer Certificate Rotation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L605)</u>

This test validates CA issuer certificate rotation. It uses a cert-manager CA issuer, rotates the CA certificate, and verifies brokers using certificates from that CA are updated with new certificates signed by the new CA. This ensures complete certificate chain rotation including CA changes.

##### Certificate Bundle:

Tests validating mutual TLS (mTLS) with client certificate authentication.

###### Certificate Bundle

Validates mTLS configuration with certificate bundles for client authentication.

<u>[Mutual Authentication](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_cert_manager_test.go#L817)</u>

This test validates mutual TLS (mTLS) configuration using cert-manager. It provisions server and client certificates, configures the broker to require client certificates, and verifies that connections succeed with valid client certificates and fail without them. This confirms full mTLS support for high-security broker deployments.

**Coverage**:
- cert-manager integration
- Java keystore format support
- Certificate lifecycle management
- Certificate rotation (broker + CA)
- mTLS with certificate bundles
- Ingress certificate annotations

---

#### 6. [`activemqartemis_controller_block_reconcile_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_block_reconcile_test.go#L36) (2 test scenarios)


Tests for annotation-based reconciliation blocking mechanism


##### Block Reconcile Annotation:

Tests validating the manual reconciliation pause mechanism via annotations for maintenance or debugging.

###### Test

Validates annotation-based reconciliation control.

<u>[Deploy, Annotate, Verify](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_block_reconcile_test.go#L47)</u>

This test validates pausing reconciliation on a running broker. It deploys a broker, waits for it to become ready, adds the `operator.activemqartemis.apache.org/block-reconcile` annotation, makes configuration changes, and verifies the operator does NOT reconcile the changes while the annotation is present. This allows administrators to perform manual maintenance without operator interference.

<u>[Annotate, Deploy, Verify](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_controller_block_reconcile_test.go#L152)</u>

This test validates blocking reconciliation from initial deployment. It creates a broker CR with the block-reconcile annotation already present, and verifies the operator does not proceed with resource creation. Removing the annotation allows reconciliation to proceed. This enables staged deployments where CRs are created but activation is deferred.

**Coverage**:
- Annotation-based reconciliation control (`operator.activemqartemis.apache.org/block-reconcile`)
- Reconciliation pause after deployment
- Reconciliation block before deployment
- Manual intervention support

---

### Reconciler Tests

#### 7. [`activemqartemis_reconciler_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go) (38 test functions)

Unit tests for reconciler logic and utility functions


##### Hashing and Comparison:

Unit tests for configuration hashing and resource comparison algorithms used to detect changes.

<u>[TestHexShaHashOfMap](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L29)</u>

Tests SHA-256 hashing of map structures. Validates that configuration maps produce consistent hashes for change detection, ensuring order-independent hashing.

<u>[TestMapComparatorForStatefulSet](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L67)</u>

Tests StatefulSet comparison logic. Validates the comparator correctly identifies meaningful differences while ignoring runtime metadata changes like timestamps and resource versions.

<u>[TestComparatorMetaAndSpec](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L203)</u>

Tests metadata and spec field comparison. Validates detection of changes in labels, annotations, and spec to determine when resources need updates.

##### Configuration Management:

Unit tests for configuration retrieval, naming conventions, and checksum extraction.

<u>[TestGetSingleStatefulSetStatus](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L226)</u>

Tests StatefulSet status retrieval function. Validates correct extraction of replica counts, ready status, and other StatefulSet state information.

<u>[TestGetConfigAppliedConfigMapName](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L268)</u>

Tests ConfigMap name generation for applied configurations. Validates naming convention produces valid, predictable Kubernetes resource names.

<u>[TestExtractSha](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L280)</u>

Tests SHA checksum extraction from resource annotations. Validates parsing of config-applied annotations used to track configuration versions.

##### Broker Properties Processing:

Unit tests for Adler-32 checksum generation and broker property parsing with various edge cases.

<u>[TestAlder32Gen](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L312)</u>

Tests basic Adler-32 checksum generation. Validates the checksum function produces consistent values for configuration strings.

<u>[TestAlder32GenSpace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L322)</u>

Tests checksum with trailing spaces. Validates whitespace handling in checksum generation.

<u>[TestAlder32GenWithEmptyLine](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L330)</u>

Tests checksum with empty lines. Validates handling of blank lines in property configurations.

<u>[TestAlder32GenWithSpace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L341)</u>

Tests checksum with embedded whitespace. Validates consistent checksum despite whitespace variations.

<u>[TestAlder32GenWithDots](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L351)</u>

Tests checksum with dotted property keys. Validates handling of hierarchical property names like `broker.config.setting`.

<u>[TestAlder32GenBrokerProps](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L361)</u>

Tests checksum for actual broker properties. Validates realistic broker configuration checksum generation.

<u>[TestAlder32RolesProps](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L369)</u>

Tests checksum for JAAS roles properties. Validates security configuration checksum generation.

<u>[TestAlder32PropsWithFF](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L402)</u>

Tests checksum with form feed characters. Validates handling of special whitespace in configuration files.

<u>[TestExtractErrors](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L415)</u>

Tests error message extraction from broker logs. Validates parsing of broker error output for CR status reporting.

##### JAAS Configuration:

Unit tests for JAAS mount path resolution and syntax validation.

<u>[TestGetJaasConfigExtraMountPath](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L453)</u>

Tests JAAS config mount path resolution. Validates identification of jaas-config files in extraMounts and correct mount path generation.

<u>[TestGetJaasConfigExtraMountPathNotPresent](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L494)</u>

Tests behavior when JAAS config is absent. Validates the function returns appropriate values when no jaas-config is configured.

<u>[TestLoginConfigSyntaxCheck](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1210)</u>

Tests JAAS login configuration syntax validation. Validates detection of malformed JAAS config and error reporting.

##### Pod Template Generation:

Unit tests for pod template generation functions with various configuration features.

<u>[TestNewPodTemplateSpecForCR_IncludesDebugArgs](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L511)</u>

Tests debug argument inclusion in pod templates. Validates Java debug options are added when debug mode is enabled.

<u>[TestNewPodTemplateSpecForCR_AppendsDebugArgs](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1080)</u>

Tests debug argument appending logic. Validates debug options are appended to existing args without replacement.

<u>[TestNewPodTemplateSpecForCR_IncludesImagePullSecret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1118)</u>

Tests ImagePullSecrets inclusion. Validates configured pull secrets are properly referenced in generated pod specs.

<u>[TestNewPodTemplateSpecForCR_IncludesTopologySpreadConstraints](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1145)</u>

Tests topology spread constraints inclusion. Validates topology configuration is correctly applied to pod templates.

<u>[TestNewPodTemplateSpecForCR_IncludesContainerSecurityContext](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1184)</u>

Tests container security context inclusion. Validates security context settings are applied to container specs.

##### Resource Template Processing:

Unit tests for resource template parsing, customization, and variable substitution.

<u>[TestProcess_TemplateIncludesLabelsServiceAndSecret](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L545)</u>

Tests label application via resource templates. Validates custom labels are merged correctly with operator-managed labels in Services and Secrets.

<u>[TestProcess_TemplateIncludesLabelsSecretRegexp](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L619)</u>

Tests regex-based secret template matching. Validates template patterns correctly match target secret resources.

<u>[TestProcess_TemplateDuplicateKeyReplacesOk](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L682)</u>

Tests duplicate key handling in templates. Validates later values override earlier ones in template processing.

<u>[Test_Respect_existing_JAVA_OPTS_properties_def](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L719)</u>

Tests JAVA_OPTS preservation during template processing. Validates critical environment variables aren't overwritten.

<u>[TestProcess_TemplateKeyValue](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L766)</u>

Tests key-value pair processing in templates. Validates variable substitution and value replacement.

<u>[TestProcess_TemplateCustomAttributeIngress](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L860)</u>

Tests Ingress resource template customization. Validates custom attributes can be applied to Ingress resources.

<u>[TestProcess_TemplateCustomAttributeMisSpellingIngress](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L921)</u>

Tests misspelled attribute handling. Validates graceful handling of typos in template attributes.

<u>[TestProcess_TemplateCustomAttributeContainerSecurityContext](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L969)</u>

Tests container security context template customization. Validates security context can be modified via templates.

<u>[TestProcess_TemplateCustomAttributePriorityClassName](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1029)</u>

Tests priority class name template application. Validates pod priority can be set through resource templates.

##### Status and Utilities:

Unit tests for status handling, host resolution, and ordinal property management.

<u>[TestStatusMarshall](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1381)</u>

Tests broker status JSON marshalling. Validates status structures serialize correctly for API storage and responses.

<u>[TestGetBrokerHost](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1402)</u>

Tests broker hostname resolution. Validates hostname generation for broker service discovery and connections.

<u>[TestParseBrokerPropertyWithOrdinal](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1426)</u>

Tests ordinal property parsing. Validates extraction of ordinal-specific properties from keys like `setting-broker-0=value`.

<u>[TestBrokerPropertiesData](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1454)</u>

Tests broker properties data processing. Validates conversion of property maps to broker configuration format.

<u>[TestBrokerPropertiesDataWithOrdinal](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1467)</u>

Tests ordinal-specific property extraction. Validates filtering and application of per-broker properties.

<u>[TestBrokerPropertiesDataWithAndWithoutOrdinal](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_reconciler_test.go#L1490)</u>

Tests mixed ordinal and global properties. Validates merging of shared and broker-specific properties for StatefulSet deployments.

**Coverage**:
- Configuration hashing and change detection
- StatefulSet comparison and status management
- Broker properties parsing and checksum validation
- JAAS configuration handling and syntax validation
- Pod template generation with various features
- Resource template processing and customization
- Status marshalling and host resolution
- Ordinal-based property management

---

### Specialized Feature Tests

#### 8. [`activemqartemis_pub_sub_scale_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_pub_sub_scale_test.go#L45) (2 test scenarios)


Integration tests for pub/sub messaging patterns with cluster scaling


##### Pub/Sub HA with Partitioned Subscribers:

Tests validating publish/subscribe messaging patterns with high availability and subscriber partitioning across a scaling broker cluster.

###### Pub n Sub, HA Pub, Partitioned Sub

Validates pub/sub messaging under cluster scaling scenarios.

<u>[Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_pub_sub_scale_test.go#L56)</u>

This test validates the initial configuration for pub/sub HA testing. It sets up topic addresses, publisher configurations, and partitioned subscriber settings, verifying the broker cluster is properly configured for pub/sub messaging with high availability before load testing begins.

<u>[From Zero to Hero](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_pub_sub_scale_test.go#L367)</u>

This test validates pub/sub behavior during cluster scaling. It starts with a minimal broker deployment, publishes messages to topics, scales the cluster up to full capacity while maintaining active publishers and subscribers, and verifies messages are correctly distributed across all subscribers with partitioning preserved. This confirms pub/sub HA works correctly during dynamic scaling operations.

**Coverage**:
- Pub/sub messaging patterns with HA
- Partitioned subscriber configuration
- Cluster scaling under pub/sub workload
- Message distribution across scaled brokers

---

#### 9. [`activemqartemis_scale_zero_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_scale_zero_test.go#L43) (2 test scenarios)


Integration tests for Kubernetes scale subresource API and zero-replica deployments


##### Scale to Zero via Subresource:

Tests validating Kubernetes scale subresource functionality for scaling brokers to zero replicas and autoscaling integration.

###### Scale to Zero via Scale Subresource

Validates scale subresource API usage and zero-replica support.

<u>[Deploy Plan 1](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_scale_zero_test.go#L54)</u>

This test validates scaling a broker to zero replicas using the Kubernetes scale subresource. It deploys a single-replica broker, uses the scale subresource API to set replicas=0, and verifies the StatefulSet scales down to zero pods. This confirms the operator properly implements the scale subresource and supports zero-replica deployments for cost savings or maintenance.

<u>[Deploy Plan with Label for Autoscale](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_scale_zero_test.go#L153)</u>

This test validates integration with autoscaling systems. It deploys a broker with labels suitable for HorizontalPodAutoscaler (HPA) configuration, verifies the labels are applied, and confirms the broker can be scaled via the scale subresource. This enables integration with Kubernetes autoscaling mechanisms.

**Coverage**:
- Kubernetes scale subresource API
- Zero-replica deployment support
- Autoscaling label integration
- StatefulSet scaling behavior

---

#### 10. [`activemqartemis_work_queue_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_work_queue_test.go#L45) (1 test scenario)


Integration tests for work queue messaging pattern with HA and competing consumers


##### HA Work Queue with Competing Subscribers:

Tests validating work queue messaging pattern where multiple consumers compete for messages with high availability.

###### HA Pub and HA Competing Sub, Compromised Total Message Order

Validates work queue pattern with HA publishers and competing consumers.

<u>[Validation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_work_queue_test.go#L56)</u>

This test validates the work queue messaging pattern with high availability. It deploys a multi-replica broker cluster, creates queue addresses, starts multiple competing consumers, publishes messages, and verifies that messages are distributed across consumers with each message delivered to exactly one consumer (competing pattern). The test confirms HA failover works during message processing and that total message ordering may be compromised (acceptable for work queues) while individual consumer ordering is preserved.

**Coverage**:
- Work queue messaging pattern
- High availability with competing consumers
- Message distribution across multiple subscribers
- Message ordering behavior in HA scenarios

---

#### 11. [`activemqartemis_rwm_pvc_ha_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_rwm_pvc_ha_test.go#L37) (1 test scenario)


Integration tests for ReadWriteMany PVC-based shared store HA configuration


##### Shared Store HA with Two Peer CRs:

Tests validating shared storage high availability where multiple broker instances access a common ReadWriteMany PVC for fast failover.

###### Two Peer CRs

Validates shared store configuration with multiple peer broker CRs.

<u>[With Shared Store](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_rwm_pvc_ha_test.go#L48)</u>

This test validates shared store HA configuration using ReadWriteMany PVCs. It deploys two separate broker CRs configured to share a common ReadWriteMany PVC for message persistence, simulates failover by stopping one broker, and verifies the peer broker can immediately access the shared messages without data loss or delay. This confirms fast failover capability with shared storage, enabling active-passive HA deployments.

**Coverage**:
- ReadWriteMany PVC support
- Shared storage HA configuration
- Fast failover scenarios
- Peer broker coordination

---

#### 12. [`activemqartemis_jdbc_ha_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_jdbc_ha_test.go#L42) (1 test scenario)


Integration tests for JDBC-based persistence and HA configuration


##### JDBC Store HA:

Tests validating JDBC database-backed message persistence for high availability scenarios.

###### Artemis

Validates JDBC store configuration and failover.

<u>[CR with DB Store](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_jdbc_ha_test.go#L53)</u>

This test validates JDBC-based persistence for HA. It deploys a broker configured with JDBC persistence (database backend instead of file-based journal), configures database connection parameters via brokerProperties, stores messages in the database, simulates broker restart, and verifies messages are recovered from the database. This confirms JDBC persistence enables shared-nothing HA architectures using database-backed message storage.

**Coverage**:
- JDBC persistence configuration
- Database-backed message store
- JDBC-based HA setup
- Fast failover with database backend

---

### Address Controller Tests

#### 13. [`activemqartemisaddress_controller_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L55) (20 test scenarios)


Integration tests for ActiveMQArtemisAddress controller and address lifecycle management


##### Configuration Defaults:

Tests validating default address configuration behavior and secret management for Address CRs.

###### Address Queue Config Defaults

Validates default configurationManaged behavior and secret lifecycle.

<u>[configurationManaged Default Should Be True (Without Queue Configuration)](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L70)</u>

This test validates that configurationManaged defaults to true when creating Address CRs without queue configuration. It creates an Address CR specifying only address name and routing type, and verifies the operator sets configurationManaged=true automatically. This ensures addresses are operator-managed by default.

<u>[Verifying the lscrs Secret Is Gone After Deleting Address CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L94)</u>

This test validates cleanup of address configuration secrets when Address CRs are deleted. It creates an Address CR, verifies a secret is created for storing address configuration, deletes the CR, and confirms the secret is garbage collected. This ensures no secret leaks occur.

<u>[configurationManaged Default Should Be True (With Queue Configuration)](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L132)</u>

This test validates configurationManaged defaults even with queue configuration specified. It creates an Address CR with queue settings and verifies configurationManaged still defaults to true. This confirms operator management is the default regardless of configuration complexity.

##### Address Management:

Tests validating Address CR lifecycle when broker CRs are recreated.

###### Broker with Address Custom Resources

Validates address persistence across broker recreation.

<u>[Address After Recreating Broker CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L165)</u>

This test validates that Address CRs survive broker deletion and recreation. It creates a broker and Address CR, deletes the broker CR, recreates an identical broker CR, and verifies the Address CR automatically reconnects to the new broker instance and the address is recreated. This confirms Address CRs are broker-independent resources.

##### Reconciliation:

Tests validating address reconciliation across multi-pod broker deployments.

###### Address Controller Test with Reconcile

Validates address creation across StatefulSet pods.

<u>[Deploy CR with Size 5 (Pods)](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L203)</u>

This test validates address creation across a 5-replica broker cluster. It deploys a multi-pod broker, creates an Address CR, and verifies the address is created on all broker pods consistently. This confirms the Address controller correctly handles multi-replica brokers.

##### Scale Down:

Tests validating address deletion behavior during broker scale-down operations.

###### Address Delete and Scale Down

Validates RemoveFromBrokerOnDelete functionality during scaling.

<u>[Scale Down, Verify RemoveFromBrokerOnDelete](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L349)</u>

This test validates address cleanup during broker scale-down. It creates addresses, scales down the broker cluster, and verifies that addresses are removed from brokers being terminated when RemoveFromBrokerOnDelete=true. This prevents orphaned addresses on scaled-down brokers.

##### Address Creation:

Tests validating various address and queue creation scenarios including auto-create behavior.

###### Address Creation Test

Validates address and queue creation with different configurations.

<u>[Create a Queue with Non-Existing Address, No Auto-Create](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L419)</u>

This test validates that queue creation fails when the parent address doesn't exist and auto-create is disabled. It attempts to create a queue without creating the address first, with autoCreateAddresses=false, and verifies proper error handling. This ensures address existence is enforced when auto-create is disabled.

<u>[Create Address with Name Only](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L514)</u>

This test validates minimal address creation with only name specified. It creates an Address CR with just the address name (no routing type, queue config, etc.) and verifies the address is created with default settings. This confirms simple address creation scenarios work correctly.

##### Console Agent:

Tests validating the console agent mechanism for managing addresses and cross-namespace/cross-broker scenarios.

###### Address CR with Console Agent

Validates console agent integration for address management across brokers and namespaces.

<u>[Address Creation via Agent](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L563)</u>

This test validates address creation through the console agent. It creates an Address CR, and verifies the operator uses the console agent to create the address in the running broker. This confirms the agent-based address management mechanism works correctly.

<u>[Address Creation Before Controller Manager Restart](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L616)</u>

This test validates address state persistence across controller restarts. It creates addresses, simulates controller restart, and verifies addresses remain correctly configured after the controller recovers. This ensures controller failures don't lose address state.

<u>[Address Creation via with ApplyToCrNames Before Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L680)</u>

This test validates Address CRs created with ApplyToCrNames before the target broker exists. It creates an Address CR specifying ApplyToCrNames targeting a not-yet-created broker, then deploys the broker, and verifies the address is automatically created once the broker becomes available. This supports declarative address provisioning.

<u>[Address Creation via with ApplyToCrNames After Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L735)</u>

This test validates Address CRs with ApplyToCrNames targeting existing brokers. It deploys a broker first, then creates an Address CR with ApplyToCrNames specifying that broker, and verifies the address is created. This confirms the standard ApplyToCrNames workflow.

<u>[Address Creation with Multiple Namespaces Before Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L790)</u>

This test validates cross-namespace address management. It creates an Address CR in one namespace with ApplyToCrNames targeting a broker in a different namespace (not yet created), deploys the broker in its namespace, and verifies the address is created cross-namespace. This confirms multi-namespace address management capabilities.

<u>[Address Creation with Multiple ApplyToCrNames](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L896)</u>

This test validates a single Address CR applying to multiple broker CRs. It creates an Address CR with ApplyToCrNames listing multiple broker names, deploys those brokers, and verifies the address is created in all specified brokers. This enables centralized address configuration for multiple broker instances.

##### Conversion and Migration:

Tests validating API version conversion and migration from Address CRs to brokerProperties-based configuration.

###### Address Conversion Test

Validates v1alpha1 to v1beta1 conversion for Address CRs.

<u>[Convert Empty Address CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L1010)</u>

This test validates conversion of minimal v1alpha1 Address CRs. It creates an empty or minimal v1alpha1 Address CR and verifies proper conversion to v1beta1 with appropriate defaults applied.

###### Address CR Migration Test

Validates migration from Address CR management to brokerProperties-based configuration.

<u>[Migrate to Broker Properties When configDelete Is FORCE](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L1057)</u>

This test validates migration when configDelete is set to FORCE. It creates Address CRs, triggers migration to brokerProperties approach, and verifies addresses are forcefully removed from the broker and recreated via brokerProperties. This supports migration strategies for deprecating Address CRs.

<u>[Migrate to Broker Properties When RemoveFromBrokerOnDelete Is False](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L1133)</u>

This test validates migration with RemoveFromBrokerOnDelete=false. It migrates addresses to brokerProperties while preserving them in the broker even after Address CR deletion. This enables graceful migration without address downtime.

<u>[Migrate to Broker Properties When RemoveFromBrokerOnDelete Is True](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_test.go#L1192)</u>

This test validates migration with RemoveFromBrokerOnDelete=true. It migrates to brokerProperties and verifies addresses are removed from the broker when the Address CR is deleted. This confirms cleanup behavior during migration.

**Coverage**:
- Address CR lifecycle management
- Queue configuration defaults
- Address settings and routing
- Migration from Address CR to brokerProperties
- Multi-broker address management
- ApplyToCrNames functionality
- Console agent integration
- RemoveFromBrokerOnDelete behavior
- Cross-namespace address creation

---

#### 14. [`activemqartemisaddress_controller_unit_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_unit_test.go) (4 test functions)

Unit tests for address controller deletion logic and error handling


##### Address Deletion:

Unit tests for address deletion logic including error handling scenarios.

<u>[TestDeleteExistingAddress](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_unit_test.go#L40)</u>

Tests standard address deletion. Validates the deleteAddress function successfully removes an address when it exists in the broker.

<u>[TestDeleteExistingAddressWithRemoveFromBrokerOnDelete](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_unit_test.go#L44)</u>

Tests address deletion with RemoveFromBrokerOnDelete flag. Validates that when this flag is set, addresses are actively removed from the broker rather than just unmarked.

<u>[TestDeleteAddressWithNotFoundError](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_unit_test.go#L90)</u>

Tests error handling when address doesn't exist. Validates the deletion function handles "not found" errors gracefully without failing reconciliation.

<u>[TestDeleteAddressWithInternalError](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_unit_test.go#L108)</u>

Tests error handling for internal failures during deletion. Validates appropriate error propagation when address deletion fails due to internal broker errors.

**Coverage**:
- Address deletion logic
- RemoveFromBrokerOnDelete flag behavior
- Error handling for not found scenarios
- Error handling for internal failures

---

#### 15. [`activemqartemisaddress_controller_deploy_operator_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L52) (7 test scenarios)


Address controller integration tests with full operator deployment


##### Address Management:

Tests for address management with deployed operator.

###### Address Test

Validates address creation across multi-pod brokers.

<u>[Deploy CR with Size 5 (Pods)](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L64)</u>

This test validates address creation across a 5-replica broker cluster with full operator deployment. It deploys a large broker cluster, creates Address CRs, and verifies addresses are consistently created across all 5 broker pods. This confirms the Address controller works correctly in real cluster scenarios.

##### Scale Down:

Tests for address behavior during scaling operations.

###### Address Delete and Scale Down

Validates address cleanup during broker scaling.

<u>[Scale Down, Verify RemoveFromBrokerOnDelete](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L208)</u>

This test validates address removal during broker scale-down in real clusters. It creates addresses, scales down brokers, and verifies RemoveFromBrokerOnDelete behavior removes addresses from terminating brokers.

##### Address Creation:

Tests for various address creation scenarios.

###### Address Creation Test

Validates address and queue creation patterns.

<u>[Create a Queue with Non-Existing Address, No Auto-Create](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L277)</u>

This test validates queue creation failure handling when parent addresses don't exist. It attempts queue creation with autoCreateAddresses=false and verifies proper error reporting in real cluster conditions.

##### Console Agent:

Tests for console agent-based address management.

###### Address CR with Console Agent

Validates console agent integration with full operator deployment.

<u>[Address Creation via Agent](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L375)</u>

This test validates console agent-based address creation in deployed operator scenarios. It uses the console agent to create addresses and verifies end-to-end functionality.

<u>[Address Creation via with ApplyToCrNames Before Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L428)</u>

This test validates ApplyToCrNames with addresses created before broker deployment. It ensures addresses automatically apply once target brokers become available.

<u>[Address Creation via with ApplyToCrNames After Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L483)</u>

This test validates ApplyToCrNames with addresses created after broker deployment. It confirms standard ApplyToCrNames workflow in deployed operator scenarios.

<u>[Address Creation with Multiple ApplyToCrNames](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisaddress_controller_deploy_operator_test.go#L538)</u>

This test validates applying addresses to multiple brokers simultaneously. It creates Address CRs targeting multiple broker names and verifies addresses are created in all specified brokers.

**Coverage**:
- Address controller with deployed operator (Label: `do`)
- Real cluster scenarios
- Console agent integration
- ApplyToCrNames functionality
- RemoveFromBrokerOnDelete behavior

---

#### 16. [`activemqartemis_address_broker_properties_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L41) (6 test scenarios)


Integration tests for address configuration via brokerProperties (alternative to Address CR)


##### Configuration Defaults:

Tests for address configuration defaults via brokerProperties.

###### BP Address Queue Config Defaults

Validates default behavior when managing addresses through brokerProperties.

<u>[configurationManaged Default Should Be True (Without Queue Configuration)](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L55)</u>

This test validates default configurationManaged behavior when addresses are defined via brokerProperties instead of Address CRs. It configures addresses through broker properties and verifies configurationManaged defaults correctly.

##### Broker with Queue:

Tests for queue configuration via brokerProperties.

###### BP Broker with Queue

Validates queue persistence and management through brokerProperties.

<u>[Address After Recreating Broker CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L84)</u>

This test validates that addresses configured via brokerProperties persist across broker recreation. It deploys a broker with addresses in brokerProperties, deletes and recreates the broker, and verifies addresses are restored.

##### Duplicate Queue Handling:

Tests for duplicate queue detection in brokerProperties.

###### BP Broker with Duplicate Queue

Validates handling of duplicate queue definitions.

<u>[Address After Recreating Broker CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L120)</u>

This test validates detection and handling of duplicate queue definitions in brokerProperties. It configures the same queue multiple times and verifies the operator detects and handles duplicates appropriately.

##### Reconciliation:

Tests for address reconciliation via brokerProperties across multi-pod brokers.

###### BP Address Controller Test with Reconcile

Validates address synchronization across broker replicas.

<u>[BP Deploy CR with Size 5 (Pods)](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L175)</u>

This test validates address reconciliation across a 5-replica broker cluster when configured via brokerProperties. It verifies addresses are consistently configured across all broker pods.

##### Address Deletion:

Tests for address cleanup via brokerProperties.

###### BP Address Delete

Validates address removal from brokerProperties configuration.

<u>[BP Verify Remove from Shared Secret After Scale Up](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L225)</u>

This test validates address cleanup from shared configuration when addresses are removed from brokerProperties during scale-up scenarios. It ensures configuration consistency across scaling operations.

##### Address Creation:

Tests for direct address creation via brokerProperties.

###### BP Address Creation Test

Validates queue creation directly through broker properties.

<u>[BP Create a Queue, Direct](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemis_address_broker_properties_test.go#L334)</u>

This test validates direct queue creation via brokerProperties without using Address CRs. It specifies queue configuration in broker properties and verifies the queue is created in the broker, demonstrating the alternative configuration approach.

**Coverage**:
- Alternative address management pattern via brokerProperties
- Address configuration without Address CR
- Duplicate queue detection
- Multi-pod address reconciliation
- Address deletion and cleanup

---

### Security Controller Tests

#### 17. [`activemqartemissecurity_controller_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L54) (11 test scenarios)


Integration tests for ActiveMQArtemisSecurity controller (DEPRECATED)


##### Security Resources:

Tests validating Security CR functionality including management connectors and security persistence (DEPRECATED feature).

###### Broker with Security Custom Resources

Validates Security CR integration with brokers.

<u>[Management Connector Config](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L70)</u>

This test validates management connector configuration via Security CR. It creates a Security CR with management connector settings and verifies the broker applies the security configuration for management interfaces.

<u>[No Password in Security Log Test](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L149)</u>

This test validates password protection in logging. It configures security with passwords and verifies that passwords are not exposed in operator or broker logs, ensuring sensitive credential information is properly protected.

<u>[Security After Recreating Broker CR](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L186)</u>

This test validates security configuration persistence across broker recreation. It creates Security CR and broker, deletes the broker, recreates it, and verifies security configuration is automatically reapplied.

<u>[Security with Console Domain Name Specified](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L255)</u>

This test validates console-specific security domain configuration. It creates a Security CR with console domain settings and verifies the broker applies different security policies for console access versus broker messaging.

##### Reconciliation:

Tests validating Security CR reconciliation logic.

###### Reconcile Test

Validates security application, idempotency, and multi-broker scenarios.

<u>[Testing Security Applied After Broker](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L350)</u>

This test validates applying security after broker deployment. It deploys a broker first, then creates a Security CR, and verifies the broker pods restart to apply security configuration.

<u>[Reconcile Twice with Nothing Changed](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L642)</u>

This test validates reconciliation idempotency. It applies security configuration, triggers reconciliation again without changes, and verifies no unnecessary updates or pod restarts occur.

<u>[Testing applyToCrNames Working Properly](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L730)</u>

This test validates ApplyToCrNames functionality for Security CRs. It creates a Security CR targeting specific broker names and verifies security is applied only to the specified brokers.

<u>[Reconcile Security on Multiple Broker CRs](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L857)</u>

This test validates a single Security CR applying to multiple brokers. It creates multiple brokers and one Security CR targeting all of them, verifying security configuration is consistently applied across all brokers.

<u>[Reconcile Security on Broker with Non Shell Safe Annotations](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L936)</u>

This test validates handling of special characters in security configuration. It uses annotations containing shell-unsafe characters and verifies proper escaping and handling in security configuration processing.

<u>[Reconcile Security with Management Role Access](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L1008)</u>

This test validates management RBAC configuration via Security CR. It configures management roles and permissions, and verifies appropriate access control for broker management interfaces.

##### Force Reconcile:

Tests validating forced reconciliation of Security CRs.

<u>[Reconcile After Broker CR Deployed, Verify Force Reconcile](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_controller_test.go#L1146)</u>

This test validates force reconciliation of security configuration. It deploys security, triggers a force reconcile, and verifies security configuration is reapplied even without changes.

**Coverage**:
- Security CR lifecycle
- JAAS configuration
- Login modules
- Security domains
- Management RBAC
- ApplyToCrNames functionality
- ⚠️ **DEPRECATED**: Security CR is deprecated in favor of brokerProperties

---

#### 18. [`activemqartemissecurity_broker_properties_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_broker_properties_test.go#L41) (1 test scenario)


Integration tests for security configuration via brokerProperties (modern approach)


##### BrokerProperties RBAC:

Tests validating security configuration through brokerProperties without Security CR (modern approach).

###### BrokerProperties RBAC

Validates modern security configuration approach.

<u>[Security with Management Role Access](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemissecurity_broker_properties_test.go#L57)</u>

This test validates security configuration entirely through brokerProperties. It deploys a broker with RBAC, user definitions, and role mappings specified in brokerProperties (no Security CR), and verifies security is properly applied. This demonstrates the modern recommended approach for security configuration, replacing the deprecated Security CR.

**Coverage**:
- Security configuration via brokerProperties (modern approach)
- RBAC configuration without Security CR
- Management role access configuration

---

### Scale Down Controller Tests

#### 19. [`activemqartemisscaledown_controller_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisscaledown_controller_test.go#L39) (2 test scenarios)


Integration tests for scale-down message migration (DEPRECATED)


##### Scale Down Test:

Tests validating message migration during broker scale-down operations (DEPRECATED).

###### Scale Down Test

Validates scale-down orchestration and message migration.

<u>[Deploy Plan 2 Clustered](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisscaledown_controller_test.go#L50)</u>

This test validates message migration during scale-down. It deploys a 2-broker cluster with persistence, stores messages, scales down to 1 broker, and verifies messages are migrated from the terminating broker to the surviving broker via the drain controller. This ensures no message loss during scale-down operations.

##### Toleration During Scale Down:

Tests validating toleration handling during scale-down.

<u>[Toleration OK, Verify Scaledown](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/activemqartemisscaledown_controller_test.go#L167)</u>

This test validates that drain pods inherit toleration configuration during scale-down. It scales down a broker with tolerations and verifies the drain pod created for message migration also has the required tolerations to run on tainted nodes.

**Coverage**:
- Message migration (DrainController)
- Scale-down orchestration
- Persistence requirements for migration
- Toleration configuration during scale down
- ⚠️ **DEPRECATED**: Scaledown CR is deprecated

---

### Controller Manager and Infrastructure Tests

#### 20. [`controllermanager_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L36) (7 test scenarios)


Integration tests for controller manager namespace watching configuration


##### Operator Namespaces Test:

Tests validating namespace watching configurations including single-namespace, restricted, multi-namespace, and cluster-wide modes.

###### Operator Namespaces Test

Validates WATCH_NAMESPACE environment variable processing and multi-tenant scenarios.

<u>[Test Resolving Watching Namespace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L48)</u>

This test validates namespace resolution logic. It tests various WATCH_NAMESPACE values and verifies the controller manager correctly parses and resolves target namespaces for watching.

<u>[Test Watching Single Default Namespace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L68)</u>

This test validates single-namespace watching mode. It configures the operator to watch only one specific namespace and verifies it only reconciles resources in that namespace, ignoring others.

<u>[Test Watching Restricted Namespace](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L117)</u>

This test validates restricted namespace mode. It configures the operator with restricted RBAC watching only its deployment namespace and verifies proper resource isolation and management within that single namespace.

<u>[Test Watching All Namespaces](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L181)</u>

This test validates cluster-wide watching mode. It sets WATCH_NAMESPACE="" (empty) and verifies the operator watches and reconciles resources across all namespaces in the cluster.

<u>[Test Watching All Namespaces with Same CR Name and Address Creation](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L224)</u>

This test validates cross-namespace address management in all-namespace mode. It creates brokers with identical names in different namespaces, creates addresses, and verifies each namespace's resources are properly isolated and managed.

<u>[Test Watching All Namespaces with Same CR Name and Exposed Console](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L289)</u>

This test validates console exposure across namespaces. It deploys brokers with identical names in different namespaces with exposed consoles and verifies each namespace gets its own distinct Routes/Ingress without conflicts.

<u>[Test Watching Multiple Namespaces](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controllermanager_test.go#L334)</u>

This test validates multi-namespace watching mode. It configures WATCH_NAMESPACE with a comma-separated list of specific namespaces and verifies the operator watches only those namespaces while ignoring others. This enables targeted multi-tenant deployments.

**Coverage**:
- Namespace watching configuration (`WATCH_NAMESPACE`)
- Multi-namespace operator deployment
- Namespace isolation
- Cross-namespace resource management
- Restricted namespace mode
- Cluster-wide deployment

---

#### 21. [`controll_plane_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controll_plane_test.go#L45) (1 test scenario)


Integration tests for restricted RBAC mode with minimal permissions


##### Restricted RBAC:

Tests validating restricted control plane mode for minimal-permission operator deployments.

###### Restricted RBAC

Validates restricted RBAC control plane functionality.

<u>[Operator Role Access](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/controll_plane_test.go#L105)</u>

This test validates operator functionality in restricted RBAC mode. It deploys the operator with minimal permissions (Role instead of ClusterRole), creates brokers in the restricted namespace, and verifies full broker lifecycle works with reduced permissions. This enables secure multi-tenant Kubernetes environments where operators cannot have cluster-wide privileges.

**Coverage**:
- Restricted mode RBAC
- Minimal permissions configuration
- Operator functionality with reduced privileges

---

#### 22. [`broker_name_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/broker_name_test.go#L41) (2 test scenarios)


Integration tests for broker naming conventions in different modes


##### Non-Default Broker Names:

Tests validating broker naming conventions in restricted and non-restricted modes.

###### Set Non Default Value

Validates custom broker name configuration across different operational modes.

<u>[Non Restricted](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/broker_name_test.go#L101)</u>

This test validates broker naming in non-restricted mode. It creates a broker with custom name configuration and verifies all generated resources (StatefulSet, Services, Routes/Ingress) use the correct naming pattern derived from the broker name. This ensures consistent resource naming.

<u>[Restricted](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/broker_name_test.go#L145)</u>

This test validates broker naming in restricted mode. It deploys in restricted namespace mode and verifies broker naming follows restricted patterns with appropriate resource name generation. This confirms naming conventions adapt to operational mode.

**Coverage**:
- Broker naming conventions
- Resource naming patterns
- Mode-specific naming (restricted vs non-restricted)
- Custom broker name configuration

---

#### 23. [`common_util_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/common_util_test.go) (0 test scenarios - Support File)

Common utility and helper functions for test suites

**Coverage**:
- Test utility functions used across controller tests
- Helper methods for test setup and assertions
- Shared test infrastructure
- No direct test scenarios (support file only)

---

#### 24. [`suite_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/suite_test.go#L171) (1 test function - Setup File) 

Test suite setup and configuration for controller tests

##### Test Suite Setup:

Test infrastructure initialization and configuration.

<u>[TestAPIs](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/controllers/suite_test.go#L171)</u>

This test function serves as the Ginkgo test suite entry point. It initializes the test environment including Kubernetes client configuration, EnvTest control plane setup, CRD installation, webhook configuration, cluster detection (OpenShift vs Kubernetes), and SSH tunnel support for remote cluster testing. This provides the complete test infrastructure for all controller integration tests.

**Provides**:
- Test environment setup and teardown
- Kubernetes client configuration
- EnvTest configuration (test control plane)
- Test utilities (BeforeEachSpec, AfterEachSpec)
- SSH tunnel support for remote cluster testing
- Cluster detection (OpenShift vs Kubernetes)
- CRD installation and cleanup
- Webhook configuration

---

## Component Tests

### Utility Package Tests

#### 25. [`pkg/utils/common/common_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common_test.go#L27) (1 test context)



Unit tests for common utility functions


##### Resync Period Configuration:

Tests validating reconciliation resync period configuration from environment variables or defaults.

###### Default Resync Period

Validates resync period resolution logic.

<u>[Default Resync Period](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/common_test.go#L32)</u>

This test validates reconcile resync period configuration. It tests the GetReconcileResyncPeriod function with various RECONCILE_RESYNC_PERIOD environment variable values and verifies correct parsing and default fallback to 30 seconds. This ensures efficient reconciliation timing.

**Coverage**:
- Resync period configuration
- Environment variable parsing (`RECONCILE_RESYNC_PERIOD`)
- Default value handling
- Duration parsing

---

#### 26. [`pkg/utils/common/conditions_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L29) (15 test scenarios)



Unit tests for status condition management


##### SetReadyCondition:

Tests for ready condition calculation based on other conditions.

###### SetReadyCondition

Validates ReadyConditionType calculation based on other status conditions.

<u>[Is Ready When There Are No Other Conditions](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L36)</u>

Tests that Ready condition is True when no other conditions exist, establishing the default ready state.

<u>[Changes Back to Ready When There Are No Other Conditions and Ready Was False](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L41)</u>

Tests that Ready condition transitions back to True when problematic conditions are resolved.

<u>[Is Ready When All Conditions Are True](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L57)</u>

Tests that Ready condition is True when all other conditions (Deployed, Valid, etc.) are True.

<u>[Is Not Ready When One Condition Is False](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L79)</u>

Tests that Ready condition becomes False when any other condition is False, implementing fail-fast status logic.

<u>[Changes to Not Ready When One Condition Is False and Was Ready](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L107)</u>

Tests ready-to-not-ready transitions when conditions degrade during runtime.

<u>[Ignores Condition in Unknown State](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L136)</u>

Tests that Unknown conditions don't affect Ready status, allowing permissive validation.

##### IsConditionPresentAndEquals:

Tests for precise condition comparison including all fields.

###### IsConditionPresentAndEquals

Validates exact condition matching logic.

<u>[Returns False When There Are No Conditions to Compare With](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L164)</u>

Tests comparison returns false when condition list is empty.

<u>[Returns False When One Field Is Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L167)</u>

Tests that any field difference (Status, Reason, Message) causes comparison to fail.

<u>[Returns True When Only Time Is Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L194)</u>

Tests that timestamp differences are ignored in condition comparison, allowing idempotent status updates.

##### IsConditionPresentAndEqualsIgnoreMessage:

Tests for condition comparison ignoring message field.

###### IsConditionPresentAndEqualsIgnoreMessage

Validates condition matching without message comparison.

<u>[Returns False When There Are No Conditions to Compare With](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L225)</u>

Tests comparison with empty condition list.

<u>[Returns True When Message Field Is Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L228)</u>

Tests that message differences are ignored while other fields must match.

<u>[Returns True When Only Time Is Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L250)</u>

Tests that both message and timestamp differences are ignored in this comparison mode.

##### PodStartingStatusDigestMessage:

Tests for pod status digest message generation from container states.

###### PodStartingStatusDigestMessage

Validates digest message generation from various pod states.

<u>[OK on Empty](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L281)</u>

Tests digest generation with empty pod status.

<u>[OK on No Conditions](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L284)</u>

Tests digest with pod status but no container conditions.

<u>[OK on Simple Conditions No Status](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L287)</u>

Tests digest with conditions but no status details.

<u>[OK on Simple Condition Status Reason](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L297)</u>

Tests digest including status reason information.

<u>[OK on Simple Condition Status Reason Message](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L309)</u>

Tests digest including both reason and message fields.

<u>[OK on Simple Condition Status Message](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L322)</u>

Tests digest with status message only.

<u>[OK on Two Condition Status Message](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/common/conditions_test.go#L334)</u>

Tests digest aggregation from multiple container conditions.

**Coverage**:
- Status condition creation and updates
- Ready condition calculation
- Condition comparison logic
- Message-aware vs message-ignoring comparisons
- Pod status digest message generation

---

#### 27. [`pkg/utils/jolokia_client/jolokia_client_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L17) (7 test scenarios)

**Test Entry Point**: [`TestJolokia`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/suite_test.go#L10)


Unit/integration tests for Jolokia client library


##### GetBrokers:

Tests for broker discovery and Jolokia client initialization.

**GetBrokers:**

Validates broker discovery from StatefulSets/Pods and Jolokia client creation.

**Without Any Target**

**[Should Not Return Any Brokers](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L21):**

Tests that GetBrokers returns empty when no broker target is configured.

**Without a Valid Deployment**

**[Should Return an Empty Array If No StatefulSets Exist with the Given Name](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L34):**

Tests error handling when referenced StatefulSet doesn't exist.

**[Should Return an Empty Array If No Pods Are Available](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L56):**

Tests handling of StatefulSet without ready pods.

**With a Valid Target**

**[Should Return a Jolokia Client](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L91):**

Tests successful Jolokia client creation for HTTP connections.

**With PEM TLS Console**

**[Should Return a Jolokia Client](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L137):**

Tests Jolokia client creation with PEM-format TLS certificates.

**With JKS TLS Console**

**[Should Return a Jolokia Client](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L194):**

Tests Jolokia client with Java KeyStore TLS configuration.

**With HTTP Web URI**

**[Should Return a Jolokia Client](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/jolokia_client/jolokia_client_test.go#L251):**

Tests Jolokia client with HTTP web URI connection.

**Coverage**:
- Jolokia client initialization
- Broker discovery from StatefulSets and Pods
- mTLS connections (PEM and JKS formats)
- HTTP and HTTPS connections
- Error handling for missing resources

---

#### 28. [`pkg/utils/artemis/artemis_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/artemis/artemis_test.go#L12) (4 test functions)

**Test Entry Point**: [`TestArtemis`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/artemis/suite_test.go#L10)

Unit tests for Artemis broker status retrieval utilities


##### Status Retrieval:

Tests for broker status retrieval via Jolokia.

<u>[TestGetStatus](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/artemis/artemis_test.go#L12)</u>

Tests successful status retrieval from a running broker via Jolokia API.

<u>[TestGetStatusWithError](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/artemis/artemis_test.go#L39)</u>

Tests error handling when Jolokia API call fails, ensuring graceful error propagation.

<u>[TestGetStatusWithErrorStatusOnly](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/artemis/artemis_test.go#L65)</u>

Tests handling of partial status responses where some data is returned alongside errors.

<u>[TestGetStatusWithNilStatus](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/artemis/artemis_test.go#L91)</u>

Tests handling of nil/empty status responses from Jolokia.

**Coverage**:
- Artemis broker status retrieval via Jolokia
- Error handling for various failure scenarios
- Nil/partial response handling
- Status parsing and validation

---

#### 29. [`pkg/utils/cr2jinja2/cr2jinja2_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/cr2jinja2/cr2jinja2_test.go#L11) (1 test context)



Unit tests for CR to Jinja2 template conversion (legacy feature)


##### Type Validation:

Tests for template variable type validation.

###### Testing Special Key Must Be of String Type

Validates type constraints for Jinja2 template variables.

<u>[Testing Special Key Must Be of String Type](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/utils/cr2jinja2/cr2jinja2_test.go#L17)</u>

This test validates that special keys in CR configuration must be string type for Jinja2 template conversion. It attempts to use non-string values and verifies type validation errors are raised. This ensures template variable compatibility.

**Coverage**:
- Template generation from CR
- Type validation for template variables
- Legacy support for init container Jinja2 templates
- String type enforcement for special keys

---

### Resource Builder Tests

#### 30. [`pkg/resources/statefulsets/statefulset_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/statefulsets/statefulset_test.go#L14) (4 test scenarios)

**Test Entry Point**: [`TestStatefulsets`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/statefulsets/suite_test.go#L10)


Unit tests for StatefulSet resource builder utilities


##### GetDeployedStatefulSetNames:

Tests for StatefulSet discovery and filtering utilities.

**GetDeployedStatefulSetNames:**

Validates StatefulSet discovery logic.

**Without Existing Stafulsets**

**[Returns an Empty Collection If the Filter Is Not Empty](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/statefulsets/statefulset_test.go#L20):**

Tests empty cluster state handling when filtering is active.

**With Deployed Statefulsets**

**[Returns the Same Collection When the Filter Is Empty](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/statefulsets/statefulset_test.go#L62):**

Tests unfiltered StatefulSet listing returns all deployed StatefulSets.

**[Returns the Statefulsets That Match the Namespace and Name](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/statefulsets/statefulset_test.go#L84):**

Tests namespace and name filtering correctly selects matching StatefulSets.

**[Returns an Empty Collection If the Filter Doesn't Match Any Statefulset](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/resources/statefulsets/statefulset_test.go#L106):**

Tests that non-matching filters return empty collections.

**Coverage**:
- StatefulSet discovery and filtering
- Namespace and name matching
- Empty collection handling
- Filter logic validation

---

### Log Package Tests

#### 31. [`pkg/log/filtered_log_sink_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/log/filtered_log_sink_test.go#L16) (1 test function)

**Test Entry Point**: [`TestArtemis`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/log/suite_test.go#L10)

Unit tests for filtered logging infrastructure


##### Filtered Log Sink:

Tests for log filtering and level management.

<u>[TestFilteredLogSink](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/log/filtered_log_sink_test.go#L16)</u>

Tests filtered log sink implementation. Validates log filtering logic and level management for operator logging infrastructure.

**Coverage**:
- Log sink filtering logic
- Log level management
- Custom log sink implementation
- Log message routing

---

### DrainController Tests

#### 32. [`pkg/draincontroller/draincontroller_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/draincontroller/draincontroller_test.go#L27) (1 test scenario)



Unit tests for drain controller message migration


##### Global Template Test:

Tests for drain pod template generation.

###### Global Template Test

Validates drain pod configuration templates.

<u>[Testing Global Template Has Correct OpenShift Ping DNS Service Port](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/pkg/draincontroller/draincontroller_test.go#L34)</u>

Tests that drain pod templates include correct OpenShift DNS service port configuration for cluster discovery during message migration.

**Coverage**:
- Drain pod template generation
- OpenShift DNS service port configuration
- Message migration infrastructure
- Scale-down coordination setup

---

### Test Utilities

#### 33-35. Test Utilities (10 test scenarios total)

Test helper utilities for common operations

---

##### 33. [`test/utils/namer/namer_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L12) (3 test scenarios)



Unit tests for pod and StatefulSet name matching utilities


###### Pod/StatefulSet Name Matching:

Tests for name matching logic.

<u>[TestMatchPodNameAndStatefulSetName: match](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L18)</u>

###### [Testing Pod Belonging to Statefulset](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L24)

Tests valid pod name matches StatefulSet naming convention.

<u>[TestMatchPodNameAndStatefulSetName: namespace mismatch](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L39)</u>

###### [Testing Different Namespaces](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L45)

Tests namespace mismatch detection in name matching.

<u>[TestMatchPodNameAndStatefulSetName: mismatch-different ss names](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L60)</u>

**[Testing Different Statefulset Name](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/namer/namer_test.go#L66):**

Tests name mismatch detection when pod doesn't belong to StatefulSet.

**Coverage**:
- Pod to StatefulSet name matching
- Namespace validation
- Name format validation

---

##### 34. [`test/utils/config/config_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/config/config_test.go#L12) (1 test scenario)



Unit tests for address settings equality comparison


###### Address Settings Comparison:

Tests for address settings deep equality.

<u>[TestAddressSettingsEqual](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/config/config_test.go#L18)</u>

**[Testing Equal](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/config/config_test.go#L83):**

Tests deep equality validation of address settings structures, ensuring all fields are compared correctly.

**Coverage**:
- Address settings comparison
- Deep equality validation
- Configuration matching logic

---

##### 35. [`test/utils/common/common_test.go`](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L14) (6 test scenarios)



Unit tests for resource comparison utilities


###### Resource Comparison:

Tests for resource equality and difference detection.

<u>[TestCompareResources: Equals](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L20)</u>

###### [Testing 2 Resources Are Equal](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L25)

Tests equal resources validation.

<u>[TestCompareResources: Different1](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L30)</u>

###### [Testing 2 Resources Are Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L35)

Tests different resource detection (case 1).

<u>[TestCompareResources: Different2](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L40)</u>

###### [Testing 2 Resources Are Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L45)

Tests different resource detection (case 2).

<u>[TestCompareResources: Different3](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L50)</u>

###### [Testing 2 Resources Are Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L55)

Tests different resource detection (case 3).

<u>[TestCompareResources: Different3](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L60)</u>

**[Testing 2 Resources Are Different](https://github.com/arkmq-org/activemq-artemis-operator/blob/4acadb95603c38b82d5d7f63fb538c37ed855662/test/utils/common/common_test.go#L65):**

Tests different resource detection (case 4 - duplicate context name in source).

**Coverage**:
- Resource equality comparison
- Multiple difference detection scenarios
- Resource matching logic

---

## Test Labels and Categories

### Ginkgo Labels Used

| Label | Purpose | Test Count |
|-------|---------|------------|
| `slow` | Long-running tests | 10+ |
| `do` | Deploy operator tests | 5+ |
| `tls-secret-reuse` | TLS/SSL tests | 3 |
| `broker-security-res` | Security tests | 10+ |
| `broker-address-res` | Address management | 8+ |
| `address-creation-test` | Address creation | 5+ |
| `queue-config-defaults` | Queue configuration | 3 |
| `console-expose-ssl` | Console SSL | 5+ |
| `ingress` | Ingress tests | 8+ |
| `route` | Route tests (OpenShift) | 8+ |
| `console` | Console tests | 12+ |
| `acceptor` | Acceptor tests | 15+ |
| `connector` | Connector tests | 10+ |
| `annotations-test` | Annotations | 2 |
| `block-reconcile` | Block reconcile | 2 |
| `pvc-owner-reference-upgrade` | PVC upgrade | 2 |
| `LoggerProperties-test` | Logging | 5+ |

### Test Categorization by Type

**Integration Tests** (require cluster):
- All tests with `USE_EXISTING_CLUSTER` check
- All tests with `DEPLOY_OPERATOR=true` check
- ~200 scenarios

**Unit Tests** (no cluster required):
- All `func Test*` in reconciler_test.go
- All controller_unit_test.go tests
- All pkg/* tests
- ~150 scenarios

**E2E Tests** (full deployment):
- Tests with Label `do` (deploy operator)
- Tests with Label `slow`
- ~50 scenarios

---

## Test Patterns and Conventions

### Ginkgo Test Structure

```go
var _ = Describe("feature name", func() {
    BeforeEach(func() {
        // Setup
    })

    AfterEach(func() {
        // Cleanup
    })

    Context("scenario group", Label("label-name"), func() {
        It("should do something", func() {
            // Test implementation
            By("step description")
            // Assertions
        })
    })
})
```

### Common Test Utilities

**Broker Deployment**:
```go
DeployCustomBroker(namespace, customizer)
DeployBroker(namespace)
WaitForPod(brokerName)
CleanResource(cr, name, namespace)
```

**Assertions**:
```go
Eventually(func(g Gomega) {
    g.Expect(...).Should(Succeed())
}, timeout, interval).Should(Succeed())
```

### Test Modes

**EnvTest Mode** (default):
- Uses local Kubernetes API server
- Fast, isolated testing
- No actual pods/containers

**Existing Cluster Mode**:
```bash
USE_EXISTING_CLUSTER=true make test
```
- Tests against real cluster
- Actual pod deployment
- Full integration testing

**Deploy Operator Mode**:
```bash
DEPLOY_OPERATOR=true USE_EXISTING_CLUSTER=true make test
```
- Deploys operator in cluster
- Full E2E testing
- Tests operator image
