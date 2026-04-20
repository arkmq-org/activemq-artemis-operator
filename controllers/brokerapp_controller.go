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

package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	broker "github.com/arkmq-org/activemq-artemis-operator/api/v1beta2"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/appselector"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/resources"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/resources/secrets"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type BrokerAppReconciler struct {
	*ReconcilerLoop
}

type BrokerAppInstanceReconciler struct {
	*BrokerAppReconciler
	instance *broker.BrokerApp
	service  *broker.BrokerService
	status   *broker.BrokerAppStatus
}

func (reconciler BrokerAppInstanceReconciler) validateSpec() error {
	// Validate resource name
	if err := ValidateResourceNameAndSetCondition(reconciler.instance.Name, &reconciler.status.Conditions); err != nil {
		return err
	}

	// Validate capability address types
	if err := reconciler.verifyCapabilityAddressType(); err != nil {
		return err
	}

	// Add additional spec validations here as needed
	// Future validations would go here

	return nil
}

func (reconciler BrokerAppInstanceReconciler) processBindingSecret() error {

	// Only manage binding secret if app has been bound to a service (annotation exists)
	serviceAnnotation, hasAnnotation := reconciler.instance.Annotations[common.AppServiceAnnotation]
	if !hasAnnotation {
		return nil
	}

	bindingSecretNsName := types.NamespacedName{
		Namespace: reconciler.instance.Namespace,
		Name:      BindingsSecretName(reconciler.instance.Name),
	}

	var desired *corev1.Secret

	obj := reconciler.CloneOfDeployed(reflect.TypeOf(corev1.Secret{}), bindingSecretNsName.Name)
	if obj != nil {
		desired = obj.(*corev1.Secret)
	} else {
		desired = secrets.NewSecret(bindingSecretNsName, nil, nil)
	}

	// Parse annotation to get service namespace and name
	serviceNamespace, serviceName, ok := parseServiceAnnotation(serviceAnnotation)
	if ok {
		desired.Data = map[string][]byte{
			// host as FQQN to work everywhere in the cluster
			"host": []byte(fmt.Sprintf("%s.%s.svc.%s", serviceName, serviceNamespace, common.GetClusterDomain())),
			"port": []byte(fmt.Sprintf("%d", reconciler.instance.Spec.Acceptor.Port)),
			"uri":  []byte(fmt.Sprintf("amqps://%s.%s.svc.%s:%d", serviceName, serviceNamespace, common.GetClusterDomain(), reconciler.instance.Spec.Acceptor.Port)),
		}
	}

	reconciler.status.Binding = &corev1.LocalObjectReference{
		Name: bindingSecretNsName.Name,
	}
	reconciler.TrackDesired(desired)
	return nil
}

func NewBrokerAppReconciler(client client.Client, scheme *runtime.Scheme, config *rest.Config, logger logr.Logger) *BrokerAppReconciler {
	reconciler := BrokerAppReconciler{ReconcilerLoop: &ReconcilerLoop{KubeBits: &KubeBits{
		Client: client, Scheme: scheme, Config: config, log: logger}}}
	reconciler.ReconcilerLoopType = &reconciler
	return &reconciler
}

//+kubebuilder:rbac:groups=broker.arkmq.org,namespace=arkmq-org-broker-operator,resources=brokerapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=broker.arkmq.org,namespace=arkmq-org-broker-operator,resources=brokerapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=broker.arkmq.org,namespace=arkmq-org-broker-operator,resources=brokerservices,verbs=get;list;watch;update

func (reconciler *BrokerAppReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := reconciler.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "Reconciling", "BrokerApp")

	instance := &broker.BrokerApp{}
	var err = reconciler.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	localLoop := &ReconcilerLoop{
		KubeBits:           reconciler.KubeBits,
		ReconcilerLoopType: reconciler,
	}

	processor := BrokerAppInstanceReconciler{
		BrokerAppReconciler: &BrokerAppReconciler{ReconcilerLoop: localLoop},
		instance:            instance,
		status:              instance.Status.DeepCopy(),
	}

	reqLogger.V(2).Info("Reconciler Processing...", "CRD.Name", instance.Name, "CRD ver", instance.ObjectMeta.ResourceVersion, "CRD Gen", instance.ObjectMeta.Generation)
	if err = processor.validateSpec(); err == nil {
		if err = processor.resolveBrokerService(); err == nil {
			if err = processor.InitDeployed(instance, processor.getOwned()...); err == nil {
				if err = processor.processBindingSecret(); err == nil {
					err = processor.SyncDesiredWithDeployed(processor.instance)
				}
			}
		}
	}

	statusErr, retry := processor.processStatus(err)
	if err == nil {
		err = statusErr
	}
	reqLogger.V(2).Info("Reconciler Processed...", "CRD.Name", instance.Name, "CRD ver", instance.ObjectMeta.ResourceVersion, "CRD Gen", instance.ObjectMeta.Generation, "error", err)
	if err == nil && retry {
		return ctrl.Result{Requeue: true, RequeueAfter: common.GetReconcileResyncPeriod()}, nil
	}
	return ctrl.Result{}, err
}

// instance specifics for a reconciler loop
func (r *BrokerAppReconciler) getOwned() []client.ObjectList {
	return []client.ObjectList{&corev1.SecretList{}}
}

func (r *BrokerAppReconciler) getOrderedTypeList() []reflect.Type {
	types := make([]reflect.Type, 1)
	types[0] = reflect.TypeOf(corev1.Secret{})
	return types
}

func (reconciler *BrokerAppInstanceReconciler) resolveBrokerService() error {

	// find the matching service to find the matching brokers
	var list = &broker.BrokerServiceList{}
	var opts, err = metav1.LabelSelectorAsSelector(reconciler.instance.Spec.ServiceSelector)
	if err != nil {
		err = fmt.Errorf("failed to evaluate Spec.Selector %v", err)
		meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  broker.ValidConditionSpecSelectorError,
			Message: err.Error(),
		})
		return err
	}
	err = reconciler.Client.List(context.TODO(), list, &client.ListOptions{LabelSelector: opts})
	if err != nil {
		// API error, not a CR validation issue - let processStatus handle it in Deployed condition
		return fmt.Errorf("Spec.Selector list error %v", err)
	}

	// CR spec is valid (selector syntax was valid)
	// Set this early so it's always set even if we return errors below
	// Runtime issues (no matching services, no capacity, API errors) are handled in Deployed condition
	meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
		Type:   broker.ValidConditionType,
		Status: metav1.ConditionTrue,
		Reason: broker.ValidConditionSuccessReason,
	})

	var service *broker.BrokerService
	needsServiceAssignment := false

	// Check if we have an existing annotation
	deployedTo, hasAnnotation := reconciler.instance.Annotations[common.AppServiceAnnotation]

	if hasAnnotation {
		// Try to find the annotated service in the list
		for index, candidate := range list.Items {
			if annotationNameFromService(&candidate) == deployedTo {
				service = &list.Items[index]
				break
			}
		}

		// Check if we found the service and if it still matches
		if service != nil {
			matches, matchErr := reconciler.matchesServiceSelector(service)
			if matchErr != nil {
				// CEL evaluation error - surface to user via status
				return NewConditionError(broker.DeployedConditionSelectorEvaluationError,
					"error evaluating service selector: %v", matchErr)
			}
			if !matches {
				reconciler.log.V(1).Info("App no longer matches annotated service selector, removing binding",
					"app", reconciler.instance.Name,
					"service", deployedTo)
				// Remove the annotation since we no longer match
				delete(reconciler.instance.Annotations, common.AppServiceAnnotation)
				if err := resources.Update(reconciler.Client, reconciler.instance); err != nil {
					return err
				}
				service = nil
				needsServiceAssignment = true
			}
		}

		// If annotated service not found, selector may have changed
		if service == nil && !needsServiceAssignment {
			if len(list.Items) > 0 {
				reconciler.log.V(1).Info("Annotated service not found in selector matches, reassigning",
					"app", reconciler.instance.Name,
					"old-annotation", deployedTo,
					"matching-services", len(list.Items))
				needsServiceAssignment = true
			}
			// else: no services match current selector, processStatus will handle it
		}
	} else {
		// No annotation yet, need initial assignment
		needsServiceAssignment = true
	}

	// Assign service if needed (initial assignment or reassignment due to selector change)
	if needsServiceAssignment {

		if len(list.Items) == 0 {
			// No matching services is a runtime issue, not a CR validation issue
			// Let processStatus handle it in Deployed condition
			return NewConditionError(broker.DeployedConditionNoMatchingServiceReason,
				"no matching services available for selector %v", opts)
		}

		service, err = reconciler.findServiceWithCapacity(list)
		if service != nil {
			// Update annotation to bind to this service
			common.ApplyAnnotations(&reconciler.instance.ObjectMeta, map[string]string{common.AppServiceAnnotation: annotationNameFromService(service)})
			err = resources.Update(reconciler.Client, reconciler.instance)
		} else {
			// Check if error is already a ConditionError (e.g., AppSelectorNoMatch)
			if _, isCondErr := AsConditionError(err); isCondErr {
				// Already a ConditionError, return as-is
				return err
			}
			// No service with capacity is a runtime issue, not a CR validation issue
			// Let processStatus handle it in Deployed condition
			err = NewConditionError(broker.DeployedConditionNoServiceCapacityReason,
				"no service with capacity available for selector %v, %v", opts, err)
		}
	}

	reconciler.service = service
	return err
}

func BindingsSecretName(crName string) string {
	return fmt.Sprintf("%s-binding-secret", crName)
}

func annotationNameFromService(service *broker.BrokerService) string {
	return fmt.Sprintf("%s:%s", service.Namespace, service.Name)
}

func parseServiceAnnotation(annotation string) (namespace string, name string, ok bool) {
	parts := strings.Split(annotation, ":")
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return "", "", false
}

// matchesServiceSelector checks if the app matches the service's appSelectorExpression.
// Returns (matches, nil) if evaluation succeeds.
// Returns (false, error) if CEL evaluation fails - caller should surface this to user.
func (reconciler *BrokerAppInstanceReconciler) matchesServiceSelector(service *broker.BrokerService) (bool, error) {
	matches, err := appselector.Matches(reconciler.instance, service, reconciler.Client)
	if err != nil {
		reconciler.log.Error(err, "Failed to evaluate appSelectorExpression",
			"app", reconciler.instance.Namespace+"/"+reconciler.instance.Name,
			"service", service.Namespace+"/"+service.Name,
			"expression", service.Spec.AppSelectorExpression)
		return false, fmt.Errorf("failed to evaluate appSelectorExpression on service %s/%s: %w",
			service.Namespace, service.Name, err)
	}

	// Additional debug logging for BrokerApp controller
	if reconciler.log.V(2).Enabled() {
		expression := service.Spec.AppSelectorExpression
		if expression == "" {
			expression = appselector.DefaultExpression
		}
		reconciler.log.V(2).Info("App selector check result",
			"expression", expression,
			"app", reconciler.instance.Namespace+"/"+reconciler.instance.Name,
			"service", service.Namespace+"/"+service.Name,
			"matches", matches)
	}

	return matches, nil
}

func (reconciler *BrokerAppInstanceReconciler) findServiceWithCapacity(list *broker.BrokerServiceList) (chosen *broker.BrokerService, err error) {
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("no services in list")
	}

	// Get the app's resource requirements
	appMemoryRequest := reconciler.instance.Spec.Resources.Requests.Memory()
	appPort := reconciler.instance.Spec.Acceptor.Port

	var bestService *broker.BrokerService
	var maxAvailable int64 = -1

	// Track why services were rejected for better error messages
	rejectionReasons := make(map[string]string)
	selectorRejectionCount := 0
	var celEvaluationError error // Track first CEL error encountered

	for i := range list.Items {
		service := &list.Items[i]

		// Check selector match first
		matches, matchErr := reconciler.matchesServiceSelector(service)
		if matchErr != nil {
			// CEL evaluation error - track it and continue to check other services
			reconciler.log.V(1).Info("Failed to evaluate selector for service",
				"service", service.Name,
				"error", matchErr)
			rejectionReasons[service.Name] = fmt.Sprintf("selector evaluation failed: %v", matchErr)
			if celEvaluationError == nil {
				celEvaluationError = matchErr // Keep first error for reporting
			}
			selectorRejectionCount++
			continue
		}
		if !matches {
			reconciler.log.V(1).Info("App does not match service selector",
				"service", service.Name,
				"app-namespace", reconciler.instance.Namespace)
			rejectionReasons[service.Name] = "does not match selector"
			selectorRejectionCount++
			continue
		}

		// Check memory capacity
		available, checkErr := reconciler.getAvailableMemory(service)
		if checkErr != nil {
			reconciler.log.V(1).Info("Failed to check capacity for service",
				"service", service.Name,
				"error", checkErr)
			rejectionReasons[service.Name] = fmt.Sprintf("error checking capacity: %v", checkErr)
			continue
		}

		if appMemoryRequest != nil && available < appMemoryRequest.Value() {
			reconciler.log.V(1).Info("Service has insufficient memory capacity",
				"service", service.Name,
				"available", available,
				"required", appMemoryRequest.Value())
			rejectionReasons[service.Name] = fmt.Sprintf("insufficient memory (available: %d, required: %d)",
				available, appMemoryRequest.Value())
			continue
		}

		// Check port availability
		conflictingApp, hasConflict, portErr := reconciler.getPortConflict(service, appPort)
		if portErr != nil {
			// Error checking for conflicts - treat as service unavailable
			reconciler.log.Error(portErr, "Failed to check port conflicts for service",
				"service", service.Name)
			rejectionReasons[service.Name] = fmt.Sprintf("error checking ports: %v", portErr)
			continue
		}
		if hasConflict {
			reconciler.log.V(1).Info("Service has port conflict",
				"service", service.Name,
				"port", appPort,
				"conflicting-app", conflictingApp)
			rejectionReasons[service.Name] = fmt.Sprintf("port %d conflict with %s", appPort, conflictingApp)
			continue
		}

		// Track service with most available capacity
		if available > maxAvailable {
			maxAvailable = available
			bestService = service
		}
	}

	if bestService == nil {
		// Build detailed error message with reasons
		reasons := []string{}
		for svcName, reason := range rejectionReasons {
			reasons = append(reasons, fmt.Sprintf("%s: %s", svcName, reason))
		}

		// If we encountered CEL evaluation errors, surface them to user
		if celEvaluationError != nil {
			return nil, NewConditionError(broker.DeployedConditionSelectorEvaluationError,
				"failed to evaluate service selector: %v", celEvaluationError)
		}

		// If all services were rejected due to selector mismatch, return specific error
		if selectorRejectionCount > 0 && selectorRejectionCount == len(rejectionReasons) {
			return nil, NewConditionError(broker.DeployedConditionDoesNotMatchReason,
				"app in namespace %s does not match selector for any service", reconciler.instance.Namespace)
		}

		if len(reasons) > 0 {
			return nil, fmt.Errorf("no service with capacity for port %d and memory %v: %s",
				appPort, appMemoryRequest, strings.Join(reasons, "; "))
		}
		return nil, fmt.Errorf("no service with capacity for port %d and memory %v", appPort, appMemoryRequest)
	}

	reconciler.log.V(1).Info("Selected service with capacity",
		"service", bestService.Name,
		"available-memory", maxAvailable,
		"port", appPort)
	return bestService, nil
}

func (reconciler *BrokerAppInstanceReconciler) listOtherAppsForService(service *broker.BrokerService) ([]broker.BrokerApp, error) {
	apps := &broker.BrokerAppList{}
	serviceKey := annotationNameFromService(service)
	if err := reconciler.Client.List(context.TODO(), apps, client.MatchingFields{common.AppServiceAnnotation: serviceKey}); err != nil {
		return nil, err
	}

	// Filter out ourselves
	result := make([]broker.BrokerApp, 0, len(apps.Items))
	for _, app := range apps.Items {
		if app.Namespace == reconciler.instance.Namespace && app.Name == reconciler.instance.Name {
			continue
		}
		result = append(result, app)
	}
	return result, nil
}

func (reconciler *BrokerAppInstanceReconciler) getAvailableMemory(service *broker.BrokerService) (int64, error) {
	// Get service's total memory limit (0 if not specified means unlimited)
	serviceMemory := service.Spec.Resources.Limits.Memory()
	if serviceMemory == nil || serviceMemory.IsZero() {
		// No limit specified, treat as unlimited
		return int64(^uint64(0) >> 1), nil // max int64
	}
	totalMemory := serviceMemory.Value()

	// Find all other apps currently provisioned on this service
	apps, err := reconciler.listOtherAppsForService(service)
	if err != nil {
		return 0, err
	}

	// Sum up memory requests of all provisioned apps
	var usedMemory int64 = 0
	for _, app := range apps {
		appMemory := app.Spec.Resources.Requests.Memory()
		if appMemory != nil {
			usedMemory += appMemory.Value()
		}
	}

	available := totalMemory - usedMemory
	if available < 0 {
		available = 0
	}

	return available, nil
}

func (reconciler *BrokerAppInstanceReconciler) verifyCapabilityAddressType() (err error) {

	for _, capability := range reconciler.instance.Spec.Capabilities {
		for index, address := range capability.SubscriberOf {
			if !strings.Contains(address.Address, "::") {
				err = fmt.Errorf("Spec.Capability.SubscriberOf[%d] address must specify a FQQN, %v", index, err)
				break
			}
		}
	}
	if err != nil {
		meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  broker.ValidConditionAddressTypeError,
			Message: err.Error(),
		})
	}
	return err
}

func (reconciler *BrokerAppInstanceReconciler) getPortConflict(service *broker.BrokerService, port int32) (conflictingApp string, hasConflict bool, err error) {
	// Query all other apps assigned to this service
	apps, err := reconciler.listOtherAppsForService(service)
	if err != nil {
		return "", false, fmt.Errorf("failed to list apps for service %s: %w", annotationNameFromService(service), err)
	}

	// Check for port conflicts with other apps
	for _, app := range apps {
		if app.Spec.Acceptor.Port == port {
			return fmt.Sprintf("%s/%s", app.Namespace, app.Name), true, nil
		}
	}

	return "", false, nil
}

func (reconciler *BrokerAppInstanceReconciler) processStatus(reconcilerError error) (err error, retry bool) {

	var deployedCondition metav1.Condition = metav1.Condition{
		Type: broker.DeployedConditionType,
	}
	if serviceName, found := reconciler.instance.Annotations[common.AppServiceAnnotation]; found {

		deployedCondition.Status = metav1.ConditionFalse
		deployedCondition.Reason = broker.DeployedConditionProvisioningPendingReason
		deployedCondition.Message = "Waiting for broker to apply application properties"

		if reconciler.service != nil {
			appIdentity := AppIdentity(reconciler.instance)
			for _, appliedApp := range reconciler.service.Status.ProvisionedApps {
				if appliedApp == appIdentity {
					deployedCondition.Status = metav1.ConditionTrue
					deployedCondition.Reason = broker.DeployedConditionProvisionedReason
					deployedCondition.Message = ""
					break
				}
			}
		} else {
			deployedCondition.Status = metav1.ConditionFalse
			deployedCondition.Reason = broker.DeployedConditionMatchedServiceNotFoundReason
			deployedCondition.Message = fmt.Sprintf("matching service from annotation %s not found", serviceName)
		}

	} else if reconcilerError != nil {
		// Check if error is a ConditionError with a specific reason
		if condErr, ok := AsConditionError(reconcilerError); ok {
			deployedCondition.Status = metav1.ConditionFalse
			deployedCondition.Reason = condErr.Reason
			deployedCondition.Message = condErr.Message
		} else {
			// Generic error (API errors, update errors, etc.)
			deployedCondition.Status = metav1.ConditionUnknown
			deployedCondition.Reason = broker.DeployedConditionCrudKindErrorReason
			deployedCondition.Message = fmt.Sprintf("error on resource crud %v", reconcilerError)
		}
	} else {
		// No annotation and no error means we haven't tried to assign yet
		deployedCondition.Status = metav1.ConditionFalse
		deployedCondition.Reason = broker.DeployedConditionNoMatchingServiceReason
	}
	meta.SetStatusCondition(&reconciler.status.Conditions, deployedCondition)
	common.SetReadyCondition(&reconciler.status.Conditions)

	if !reflect.DeepEqual(reconciler.instance.Status, *reconciler.status) {
		reconciler.instance.Status = *reconciler.status
		err = resources.UpdateStatus(reconciler.Client, reconciler.instance)
	}
	retry = meta.IsStatusConditionTrue(reconciler.instance.Status.Conditions, broker.ValidConditionType) &&
		meta.IsStatusConditionFalse(reconciler.instance.Status.Conditions, broker.DeployedConditionType)
	return err, retry
}

func (r *BrokerAppReconciler) enqueueAppsForService() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		service := obj.(*broker.BrokerService)
		serviceAnnotation := fmt.Sprintf("%s:%s", service.Namespace, service.Name)

		// Find all BrokerApps that reference this service via annotation using field index
		appList := &broker.BrokerAppList{}
		if err := r.Client.List(ctx, appList, client.MatchingFields{common.AppServiceAnnotation: serviceAnnotation}); err != nil {
			r.log.Error(err, "Failed to list BrokerApps for service watch", "service", serviceAnnotation)
			return nil
		}

		requests := make([]reconcile.Request, 0, len(appList.Items))
		for _, app := range appList.Items {
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: app.Namespace,
					Name:      app.Name,
				},
			})
		}
		return requests
	})
}

func (r *BrokerAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Note: Namespace informer is set up in main.go for CEL evaluation
	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.BrokerApp{}).
		Owns(&corev1.Secret{}).
		Watches(&broker.BrokerService{}, r.enqueueAppsForService()).
		Complete(r)
}
