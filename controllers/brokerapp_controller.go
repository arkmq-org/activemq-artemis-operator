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

	"github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	broker "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
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

func (reconciler BrokerAppInstanceReconciler) validateResourceName() error {
	err := common.ValidateResourceName(reconciler.instance.Name)
	if err != nil {
		meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "InvalidResourceName",
			Message: err.Error(),
		})
	}
	return err
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

//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerservices,verbs=get;list;watch;update

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
	if err = processor.validateResourceName(); err == nil {
		if err = processor.verifyCapabilityAddressType(); err == nil {
			if err = processor.resolveBrokerService(); err == nil {
				if err = processor.InitDeployed(instance, processor.getOwned()...); err == nil {
					if err = processor.processBindingSecret(); err == nil {
						err = processor.SyncDesiredWithDeployed(processor.instance)
					}
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
			Reason:  "SpecSelectorError",
			Message: err.Error(),
		})
		return err
	}
	err = reconciler.Client.List(context.TODO(), list, &client.ListOptions{LabelSelector: opts})
	if err != nil {
		err = fmt.Errorf("Spec.Selector list error %v", err)
		meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "SpecSelectorListError",
			Message: err.Error(),
		})

		return err
	}

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

		// If annotated service not found, selector may have changed
		if service == nil {
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
			err = fmt.Errorf("no matching services available for selector %v", opts)
			meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
				Type:    broker.ValidConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "SpecSelectorNoMatch",
				Message: err.Error(),
			})
			return err
		}

		service, err = reconciler.findServiceWithCapacity(list)
		if service != nil {
			// Update annotation to bind to this service
			common.ApplyAnnotations(&reconciler.instance.ObjectMeta, map[string]string{common.AppServiceAnnotation: annotationNameFromService(service)})
			err = resources.Update(reconciler.Client, reconciler.instance)
		} else {
			// No service with capacity available
			err = fmt.Errorf("no service with capacity available for selector %v, %v", opts, err)
			meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
				Type:    broker.ValidConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "NoServiceCapacity",
				Message: err.Error(),
			})
		}
	}

	reconciler.service = service
	if err == nil {
		// CR is valid regardless of whether service is currently available
		meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
			Type:   broker.ValidConditionType,
			Status: metav1.ConditionTrue,
			Reason: "ServiceResolved",
		})
	}
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

func (reconciler *BrokerAppInstanceReconciler) findServiceWithCapacity(list *broker.BrokerServiceList) (chosen *broker.BrokerService, err error) {
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("no services in list")
	}

	// Get the app's memory request (0 if not specified)
	appMemoryRequest := reconciler.instance.Spec.Resources.Requests.Memory()

	var bestService *broker.BrokerService
	var maxAvailable int64 = -1

	for i := range list.Items {
		service := &list.Items[i]
		available, checkErr := reconciler.getAvailableMemory(service)

		if checkErr != nil {
			reconciler.log.V(1).Info("Failed to check capacity for service",
				"service", service.Name,
				"error", checkErr)
			continue
		}

		// Check if this service has enough capacity
		if appMemoryRequest != nil && available < appMemoryRequest.Value() {
			reconciler.log.V(1).Info("Service has insufficient memory capacity",
				"service", service.Name,
				"available", available,
				"required", appMemoryRequest.Value())
			continue
		}

		// Track service with most available capacity
		if available > maxAvailable {
			maxAvailable = available
			bestService = service
		}
	}

	if bestService == nil {
		if appMemoryRequest != nil && !appMemoryRequest.IsZero() {
			return nil, fmt.Errorf("no service with sufficient memory capacity (required: %s)", appMemoryRequest.String())
		}
		// No memory constraints, pick first
		return &list.Items[0], nil
	}

	reconciler.log.V(1).Info("Selected service with capacity",
		"service", bestService.Name,
		"available-memory", maxAvailable)
	return bestService, nil
}

func (reconciler *BrokerAppInstanceReconciler) getAvailableMemory(service *broker.BrokerService) (int64, error) {
	// Get service's total memory limit (0 if not specified means unlimited)
	serviceMemory := service.Spec.Resources.Limits.Memory()
	if serviceMemory == nil || serviceMemory.IsZero() {
		// No limit specified, treat as unlimited
		return int64(^uint64(0) >> 1), nil // max int64
	}
	totalMemory := serviceMemory.Value()

	// Find all apps currently provisioned on this service
	apps := &broker.BrokerAppList{}
	serviceKey := annotationNameFromService(service)
	if err := reconciler.Client.List(context.TODO(), apps, client.MatchingFields{common.AppServiceAnnotation: serviceKey}); err != nil {
		return 0, err
	}

	// Sum up memory requests of all provisioned apps
	var usedMemory int64 = 0
	for _, app := range apps.Items {
		// Skip ourselves if we're already in the list
		if app.Namespace == reconciler.instance.Namespace && app.Name == reconciler.instance.Name {
			continue
		}

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
			Reason:  "AddressTypeError",
			Message: err.Error(),
		})
	}
	return err
}

func (reconciler *BrokerAppInstanceReconciler) processStatus(reconcilerError error) (err error, retry bool) {

	var deployedCondition metav1.Condition = metav1.Condition{
		Type: v1beta1.DeployedConditionType,
	}
	if serviceName, found := reconciler.instance.Annotations[common.AppServiceAnnotation]; found {

		deployedCondition.Status = metav1.ConditionFalse
		deployedCondition.Reason = "ProvisioningPending"
		deployedCondition.Message = "Waiting for broker to apply application properties"

		if reconciler.service != nil {
			appIdentity := AppIdentity(reconciler.instance)
			for _, appliedApp := range reconciler.service.Status.ProvisionedApps {
				if appliedApp == appIdentity {
					deployedCondition.Status = metav1.ConditionTrue
					deployedCondition.Reason = "Provisioned"
					deployedCondition.Message = ""
					break
				}
			}
		} else {
			deployedCondition.Status = metav1.ConditionFalse
			deployedCondition.Reason = "MatchedServiceNotFound"
			deployedCondition.Message = fmt.Sprintf("matching service from annotation %s not found", serviceName)
		}

	} else if reconcilerError != nil {
		deployedCondition.Status = metav1.ConditionUnknown
		deployedCondition.Reason = broker.DeployedConditionCrudKindErrorReason
		deployedCondition.Message = fmt.Sprintf("error on resource crud %v", reconcilerError)
	} else {
		deployedCondition.Status = metav1.ConditionFalse
		deployedCondition.Reason = "NoMatchingService"
	}
	meta.SetStatusCondition(&reconciler.status.Conditions, deployedCondition)
	common.SetReadyCondition(&reconciler.status.Conditions)

	if !reflect.DeepEqual(reconciler.instance.Status, *reconciler.status) {
		reconciler.instance.Status = *reconciler.status
		err = resources.UpdateStatus(reconciler.Client, reconciler.instance)
	}
	retry = meta.IsStatusConditionTrue(reconciler.instance.Status.Conditions, v1beta1.ValidConditionType) &&
		meta.IsStatusConditionFalse(reconciler.instance.Status.Conditions, v1beta1.DeployedConditionType)
	return err, retry
}

func (r *BrokerAppReconciler) enqueueAppsForService() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		service := obj.(*broker.BrokerService)
		serviceAnnotation := fmt.Sprintf("%s:%s", service.Namespace, service.Name)

		// Find all BrokerApps that reference this service via annotation
		appList := &broker.BrokerAppList{}
		if err := r.Client.List(ctx, appList); err != nil {
			r.log.Error(err, "Failed to list BrokerApps for service watch", "service", serviceAnnotation)
			return nil
		}

		var requests []reconcile.Request
		for _, app := range appList.Items {
			if val, ok := app.Annotations[common.AppServiceAnnotation]; ok && val == serviceAnnotation {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: app.Namespace,
						Name:      app.Name,
					},
				})
			}
		}
		return requests
	})
}

func (r *BrokerAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.BrokerApp{}).
		Owns(&corev1.Secret{}).
		Watches(&broker.BrokerService{}, r.enqueueAppsForService()).
		Complete(r)
}
