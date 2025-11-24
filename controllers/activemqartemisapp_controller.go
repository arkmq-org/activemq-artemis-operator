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
	"sort"

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
)

type ActiveMQArtemisAppReconciler struct {
	*ReconcilerLoop
}

type ActiveMQArtemisAppInstanceReconciler struct {
	*ActiveMQArtemisAppReconciler
	instance *broker.ActiveMQArtemisApp
}

func NewActiveMQArtemisAppReconciler(client client.Client, scheme *runtime.Scheme, config *rest.Config, logger logr.Logger) *ActiveMQArtemisAppReconciler {
	reconciler := ActiveMQArtemisAppReconciler{ReconcilerLoop: &ReconcilerLoop{KubeBits: &KubeBits{
		Client: client, Scheme: scheme, Config: config, log: logger}}}
	reconciler.ReconcilerLoopType = &reconciler
	return &reconciler
}

//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=activemqartemisapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=activemqartemisapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=activemqartemisapps/finalizers,verbs=update

func (reconciler *ActiveMQArtemisAppReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := reconciler.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "Reconciling", "ActiveMQArtemisApp")

	instance := &broker.ActiveMQArtemisApp{}
	var err = reconciler.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	processor := ActiveMQArtemisAppInstanceReconciler{
		ActiveMQArtemisAppReconciler: &ActiveMQArtemisAppReconciler{ReconcilerLoop: reconciler.ReconcilerLoop},
		instance:                     instance,
	}

	reqLogger.V(2).Info("Reconciler Processing...", "CRD.Name", instance.Name, "CRD ver", instance.ObjectMeta.ResourceVersion, "CRD Gen", instance.ObjectMeta.Generation)

	processor.InitDeployed(instance, processor.getOwned()...)

	// reconcile
	if err = processor.processSpec(instance); err == nil {
		if err = processor.SyncDesiredWithDeployed(instance); err == nil {
			// all good
		} else {
			reqLogger.Error(err, "failed to sync resources")
		}
	} else {
		reqLogger.Error(err, "failed to process spec")
	}

	if err = processor.processStatus(instance, err); err != nil {
		return ctrl.Result{}, err
	}

	reqLogger.V(2).Info("Reconciler Processed...", "CRD.Name", instance.Name, "CRD ver", instance.ObjectMeta.ResourceVersion, "CRD Gen", instance.ObjectMeta.Generation)
	return ctrl.Result{}, nil
}

// instance specifics for a reconciler loop
// TODO: owned could be ordered, to avoid the second operation maybe!
func (r *ActiveMQArtemisAppReconciler) getOwned() []client.ObjectList {
	return []client.ObjectList{&corev1.SecretList{}}
}

func (r *ActiveMQArtemisAppReconciler) getOrderedTypeList() []reflect.Type {
	types := make([]reflect.Type, 1)
	// we want to create/update in this order
	types[0] = reflect.TypeOf(corev1.Secret{})
	return types
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) processSpec(cr *broker.ActiveMQArtemisApp) error {

	// find the matching service to find the matching brokers
	var list = &broker.ActiveMQArtemisServiceList{}
	var opts, err = metav1.LabelSelectorAsSelector(cr.Spec.Selector)
	if err != nil {
		err = fmt.Errorf("failed to evaluate Spec.Selector %v", err)
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
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
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "SpecSelectorListError",
			Message: err.Error(),
		})

		return err
	}

	if len(list.Items) == 0 {
		err = fmt.Errorf("no matchihg services available for selector %v", opts)
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "SpecSelectorNoMatch",
			Message: err.Error(),
		})
		return err
	}

	var service *broker.ActiveMQArtemisService

	// have we got a deployed annotation that matches
	deployedTo, found := cr.Annotations[common.AppServiceAnnotation]

	// find our service
	for index, candidate := range list.Items {
		if found {
			if annotationNameFromService(&candidate) == deployedTo {
				service = &list.Items[index]
				break
			}
		}
	}

	if !found {
		service, err = reconciler.findServiceWithCapacity(list)
		if service != nil {
			// annotate with service identity
			common.ApplyAnnotations(&cr.ObjectMeta, map[string]string{common.AppServiceAnnotation: annotationNameFromService(service)})

			// update immediately such that Auth can find it
			err = resources.Update(reconciler.Client, cr)

		} else {
			err = fmt.Errorf("no service with capacity available for selector %v, %v", opts, err)
			meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{
				Type:    broker.ValidConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "NoServiceCapacity",
				Message: err.Error(),
			})
		}
	}

	if err == nil {
		err = reconciler.processCapabilities(service)
	}

	if err == nil {
		err = reconciler.processAuth(service)
	}

	return err
}

func annotationNameFromService(service *broker.ActiveMQArtemisService) string {
	return fmt.Sprintf("%s:%s", service.Namespace, service.Name)
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) findServiceWithCapacity(list *broker.ActiveMQArtemisServiceList) (chosen *broker.ActiveMQArtemisService, err error) {
	// no notion of resource constraints yet
	first := list.Items[0]
	return &first, reconciler.checkAuthMatch(&first)
}

type AddressConfig struct {
	senderRoles   map[string]string
	consumerRoles map[string]string
}

type AddressTracker struct {
	multiCast map[string]AddressConfig
	anyCast   map[string]AddressConfig
}

func newAddressTracker() *AddressTracker {
	return &AddressTracker{multiCast: map[string]AddressConfig{}, anyCast: map[string]AddressConfig{}}
}

func (t *AddressTracker) newAddressConfig() AddressConfig {
	return AddressConfig{senderRoles: map[string]string{}, consumerRoles: map[string]string{}}
}

func (t *AddressTracker) track(address *broker.AppAddressType) *AddressConfig {

	var trackerMap map[string]AddressConfig
	var name string
	if address.QueueName != "" {
		trackerMap = t.anyCast
		name = address.QueueName
	} else {
		trackerMap = t.multiCast
		name = address.Name
	}

	var present bool
	var entry AddressConfig
	if entry, present = trackerMap[name]; !present {
		entry = t.newAddressConfig()
		trackerMap[name] = entry
	}
	return &entry
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) processCapabilities(service *broker.ActiveMQArtemisService) (err error) {

	err = reconciler.verifyCapabilityAddressType()
	if err != nil {
		return err
	}

	key := types.NamespacedName{
		Name:      AppPropertiesSecretName(service.Name),
		Namespace: service.Namespace,
	}

	secret, err := secrets.RetriveSecret(key, nil, reconciler.Client)
	if err == nil {

		buf := NewPropsWithHeader()

		addressTracker := newAddressTracker()

		for _, capability := range reconciler.instance.Spec.Capabilities {

			var role = capability.Role
			if role == "" {
				role = reconciler.AppIdentity()
			}
			var entry *AddressConfig

			for _, address := range capability.ProducerOf {
				entry = addressTracker.track(&address)
				entry.senderRoles[role] = role
			}

			for _, address := range capability.ConsumerOf {
				entry = addressTracker.track(&address)
				entry.consumerRoles[role] = role
			}

			for _, address := range capability.ProducerAndConsumerOf {
				entry = addressTracker.track(&address)
				entry.consumerRoles[role] = role
				entry.senderRoles[role] = role
			}
		}

		for name, addr := range addressTracker.anyCast {
			fmt.Fprintf(buf, "addressConfigurations.\"%s\".routingTypes=ANYCAST\n", name)
			fmt.Fprintf(buf, "addressConfigurations.\"%s\".queueConfigs.\"%s\".routingType=ANYCAST\n", name, name)

			for _, role := range addr.senderRoles {
				fmt.Fprintf(buf, "securityRoles.\"%s\".\"%s\".send=true\n", name, producerRole(role))
			}
			for _, role := range addr.consumerRoles {
				fmt.Fprintf(buf, "securityRoles.\"%s\".\"%s\".consume=true\n", name, consumerRole(role))
			}
		}

		for name, addr := range addressTracker.multiCast {
			fmt.Fprintf(buf, "addressConfigurations.\"%s\".routingTypes=MULTICAST\n", name)

			for _, role := range addr.senderRoles {
				fmt.Fprintf(buf, "securityRoles.\"%s\".\"%s\".send=true\n", name, producerRole(role))
			}
			for _, role := range addr.consumerRoles {
				fmt.Fprintf(buf, "securityRoles.\"%s\".\"%s\".consume=true\n", name, consumerRole(role))
			}
		}

		secret.Data = map[string][]byte{
			reconciler.AppIdentity() + "_capabilities.properties": buf.Bytes(),
		}
		err = resources.Update(reconciler.Client, secret)
	}
	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) AppIdentity() string {
	return NamespacedAppId(reconciler.instance)
}

func NamespacedAppId(app *broker.ActiveMQArtemisApp) string {
	return fmt.Sprintf("%s-%s", app.Name, app.Namespace)
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) checkAuthMatch(service *broker.ActiveMQArtemisService) (err error) {

	for _, requested := range reconciler.instance.Spec.Auth {
		if !contains(service.Spec.Auth, requested) {
			err = fmt.Errorf("no matchihg auth capability available for %v", requested)
			meta.SetStatusCondition(&reconciler.instance.Status.Conditions, metav1.Condition{
				Type:    broker.ValidConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "SpecAuthNoMatch",
				Message: err.Error(),
			})
			break
		}
	}
	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) verifyCapabilityAddressType() (err error) {

	for _, capability := range reconciler.instance.Spec.Capabilities {

		for index, address := range capability.ProducerOf {

			if address.Name == "" && address.QueueName == "" {
				err = fmt.Errorf("Spec.Capability.ProducerOf[%d] address must specify name or queue name %v", index, err)
				break
			}
		}
		if err == nil {
			for index, address := range capability.ConsumerOf {

				if address.Name == "" && address.QueueName == "" {
					err = fmt.Errorf("Spec.Capability.ConsumerOf[%d] address must specify name or queue name %v", index, err)
					break
				}
			}
		}
		if err == nil {
			for index, address := range capability.ProducerAndConsumerOf {

				if address.Name == "" && address.QueueName == "" {
					err = fmt.Errorf("Spec.Capability.ProducerAndConsumerOf[%d] address must specify name or queue name %v", index, err)
					break
				}
			}
		}
	}
	if err != nil {
		meta.SetStatusCondition(&reconciler.instance.Status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "AddressTypeError",
			Message: err.Error(),
		})
	}
	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) verifyNoRoleInCapabilitiesForMtls() (err error) {

	for index, requested := range reconciler.instance.Spec.Capabilities {
		if requested.Role != "" {
			err = fmt.Errorf("invalid non empty Spec.Capabilities[%d].Role %s, mtls Auth does not support role assigment, the single role is derived from the app certificate", index, requested.Role)
			meta.SetStatusCondition(&reconciler.instance.Status.Conditions, metav1.Condition{
				Type:    broker.ValidConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "NonNilRoleWithMtlsAuth",
				Message: err.Error(),
			})
			break
		}
	}
	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) processAuth(service *broker.ActiveMQArtemisService) (err error) {

	if contains(service.Spec.Auth, broker.MTLS) {
		err = reconciler.verifyNoRoleInCapabilitiesForMtls()
		if err == nil {
			err = reconciler.doMtlsAuthN(service)
		}
	}
	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) doMtlsAuthN(service *broker.ActiveMQArtemisService) (err error) {

	key := types.NamespacedName{
		Name:      JaasConfigSecretName(service.Name),
		Namespace: service.Namespace,
	}

	secret, err := secrets.RetriveSecret(key, nil, reconciler.Client)
	if err == nil {

		// shared secret, reconcile from first principals for *all* apps assigned to this service

		var list = &broker.ActiveMQArtemisAppList{}
		err = reconciler.Client.List(context.TODO(), list, &client.ListOptions{})
		if err != nil {
			err = fmt.Errorf("failed to list app to find auth match, err: %v", err)
			meta.SetStatusCondition(&reconciler.instance.Status.Conditions, metav1.Condition{
				Type:    broker.DeployedConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "AppListError",
				Message: err.Error(),
			})
			return err
		}

		if len(list.Items) == 0 {
			err = fmt.Errorf("no matchihg apps for service auth, at least one is expected")
			meta.SetStatusCondition(&reconciler.instance.Status.Conditions, metav1.Condition{
				Type:    broker.DeployedConditionType,
				Status:  metav1.ConditionFalse,
				Reason:  "AppListEmptyError",
				Message: err.Error(),
			})
			return err
		}

		// have we got a deployed annotation that matches
		deployedApps := make([]broker.ActiveMQArtemisApp, 0)
		matchedServiceAnnotation := annotationNameFromService(service)

		for _, candidate := range list.Items {

			deployedTo, found := candidate.Annotations[common.AppServiceAnnotation]
			if found && deployedTo == matchedServiceAnnotation {
				deployedApps = append(deployedApps, candidate)
			}
		}

		// order is important for reconcile consistency
		sort.SliceStable(deployedApps, func(i, j int) bool {
			return deployedApps[i].Namespace < deployedApps[j].Namespace &&
				deployedApps[i].Name < deployedApps[j].Name
		})

		usersBuf := NewPropsWithHeader()
		rolesBuf := NewPropsWithHeader()

		for _, candidate := range deployedApps {

			userName := NamespacedAppId(&candidate)
			// TODO: pull down full DN from app cert
			fmt.Fprintf(usersBuf, "%s=/.*%s.*/\n", userName, candidate.Name)

			for _, capability := range candidate.Spec.Capabilities {

				// single identity role (userName) for mtls
				if len(capability.ConsumerOf) > 0 {
					fmt.Fprintf(rolesBuf, "%s=%s\n", consumerRole(userName), userName)
				}
				if len(capability.ProducerOf) > 0 {
					fmt.Fprintf(rolesBuf, "%s=%s\n", producerRole(userName), userName)
				}
				if len(capability.ProducerAndConsumerOf) > 0 {
					fmt.Fprintf(rolesBuf, "%s=%s\n", consumerRole(userName), userName)
					fmt.Fprintf(rolesBuf, "%s=%s\n", producerRole(userName), userName)
				}
			}
		}

		secret.Data[common.GetCertUsersKey(common.JaasRealm)] = usersBuf.Bytes()
		secret.Data[common.GetCertRolesKey(common.JaasRealm)] = rolesBuf.Bytes()

		err = resources.Update(reconciler.Client, secret)
	}
	return err
}

func producerRole(prefix string) string {
	return fmt.Sprintf("%s-producer", prefix)
}

func consumerRole(prefix string) string {
	return fmt.Sprintf("%s-consumer", prefix)
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) processStatus(cr *broker.ActiveMQArtemisApp, reconcilerError error) (err error) {

	var conditions []metav1.Condition = nil

	var reconciledCondition metav1.Condition = metav1.Condition{
		Type: "Reconciled",
	}
	if reconcilerError != nil {
		reconciledCondition.Status = metav1.ConditionFalse
		reconciledCondition.Reason = broker.DeployedConditionCrudKindErrorReason
		reconciledCondition.Message = fmt.Sprintf("error on resource crud %v", reconcilerError)
	} else {
		reconciledCondition.Status = metav1.ConditionTrue
		reconciledCondition.Reason = "ok"
	}

	meta.SetStatusCondition(&conditions, reconciledCondition)

	/* need some knowledge, typically only one broker or ss, read status in leader follower is
	type dependent */

	// check for deployed annotation set by service

	var deployedCondition metav1.Condition = metav1.Condition{
		Type: broker.DeployedConditionType,
	}

	deployedCondition.Status = metav1.ConditionTrue
	deployedCondition.Reason = broker.ReadyConditionReason
	meta.SetStatusCondition(&conditions, deployedCondition)

	common.SetReadyCondition(&conditions)

	if !reflect.DeepEqual(cr.Status.Conditions, conditions) {
		cr.Status.Conditions = conditions
		err = resources.UpdateStatus(reconciler.Client, cr)
	}
	return err
}

func (r *ActiveMQArtemisAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.ActiveMQArtemisApp{}).Complete(r)
}
