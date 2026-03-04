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

type BrokerAppReconciler struct {
	*ReconcilerLoop
}

type BrokerAppInstanceReconciler struct {
	*BrokerAppReconciler
	instance *broker.BrokerApp
	service  *broker.BrokerService
	status   *broker.BrokerAppStatus
}

func (reconciler BrokerAppInstanceReconciler) processBindingSecret() error {

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

	if reconciler.service != nil {
		desired.Data = map[string][]byte{
			// host as FQQN to work everywhere in the cluster
			"host": []byte(fmt.Sprintf("%s.%s.svc.%s", reconciler.service.Name, reconciler.service.Namespace, common.GetClusterDomain())),
			"port": []byte(fmt.Sprintf("%d", reconciler.instance.Spec.Acceptor.Port)),
			"uri":  []byte(fmt.Sprintf("amqps://%s.%s.svc:%d", reconciler.service.Name, reconciler.service.Namespace, reconciler.instance.Spec.Acceptor.Port)),
		}
	} else {
		desired.Data = nil
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
	if err = processor.verifyCapabilityAddressType(); err == nil {
		if err = processor.resolveBrokerService(); err == nil {
			if err = processor.InitDeployed(instance, processor.getOwned()...); err == nil {
				if err = processor.processBindingSecret(); err == nil {
					err = processor.SyncDesiredWithDeployed(processor.instance)
				}
			}
		}
	}

	if statusErr := processor.processStatus(err); statusErr != nil {
		if err == nil {
			return ctrl.Result{}, statusErr
		}
	}
	reqLogger.V(2).Info("Reconciler Processed...", "CRD.Name", instance.Name, "CRD ver", instance.ObjectMeta.ResourceVersion, "CRD Gen", instance.ObjectMeta.Generation, "error", err)
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

	var service *broker.BrokerService

	// have we got a deployed annotation that matches
	deployedTo, found := reconciler.instance.Annotations[common.AppServiceAnnotation]

	// find our service
	for index, candidate := range list.Items {
		if found {
			if annotationNameFromService(&candidate) == deployedTo {
				service = &list.Items[index]
				break
			}
		}
	}

	if found && service == nil {
		err = fmt.Errorf("deployed to service from annotation %s not found", deployedTo)
		meta.SetStatusCondition(&reconciler.status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "DeployedToNotFound",
			Message: err.Error(),
		})
		return err
	}

	if !found {
		service, err = reconciler.findServiceWithCapacity(list)
		if service != nil {
			// annotate with service identity, the service operator will reconcile
			common.ApplyAnnotations(&reconciler.instance.ObjectMeta, map[string]string{common.AppServiceAnnotation: annotationNameFromService(service)})
			err = resources.Update(reconciler.Client, reconciler.instance)
		} else {
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

func (reconciler *BrokerAppInstanceReconciler) findServiceWithCapacity(list *broker.BrokerServiceList) (chosen *broker.BrokerService, err error) {
	// no notion of resource constraints yet
	first := list.Items[0]
	return &first, nil
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

func (reconciler *BrokerAppInstanceReconciler) processStatus(reconcilerError error) (err error) {

	var deployedCondition metav1.Condition = metav1.Condition{
		Type: "Deployed",
	}
	if reconcilerError != nil {
		deployedCondition.Status = metav1.ConditionUnknown
		deployedCondition.Reason = broker.DeployedConditionCrudKindErrorReason
		deployedCondition.Message = fmt.Sprintf("error on resource crud %v", reconcilerError)
	} else if _, found := reconciler.instance.Annotations[common.AppServiceAnnotation]; found {

		// conditional on the relevant properties in the broker status or the service status?
		deployedCondition.Status = metav1.ConditionTrue
		deployedCondition.Reason = broker.ReadyConditionReason

	} else {
		deployedCondition.Status = metav1.ConditionFalse
		deployedCondition.Reason = "NoMatchingService"
	}
	meta.SetStatusCondition(&reconciler.status.Conditions, deployedCondition)
	common.SetReadyCondition(&reconciler.status.Conditions)

	if !reflect.DeepEqual(reconciler.instance.Status, reconciler.status) {
		reconciler.instance.Status = *reconciler.status
		err = resources.UpdateStatus(reconciler.Client, reconciler.instance)
	}
	return err
}

func (r *BrokerAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.BrokerApp{}).Complete(r)
}
