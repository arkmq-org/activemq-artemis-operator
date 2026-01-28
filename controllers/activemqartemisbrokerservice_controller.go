/*
Copyright 2026.

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
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultServicePort int32 = 61616
	EmptyBrokerXml           = "empty-broker-xml"
)

type BrokerServiceReconciler struct {
	*ReconcilerLoop
}

type BrokerServiceInstanceReconciler struct {
	*BrokerServiceReconciler
	instance *broker.BrokerService
}

func NewBrokerServiceReconciler(client client.Client, scheme *runtime.Scheme, config *rest.Config, logger logr.Logger) *BrokerServiceReconciler {
	reconciler := BrokerServiceReconciler{
		ReconcilerLoop: &ReconcilerLoop{KubeBits: &KubeBits{client, scheme, config, logger}},
	}
	reconciler.ReconcilerLoopType = &reconciler
	return &reconciler
}

//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerservices/finalizers,verbs=update

func (reconciler *BrokerServiceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := reconciler.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "Reconciling", "BrokerService")

	instance := &broker.BrokerService{}
	var err = reconciler.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	processor := BrokerServiceInstanceReconciler{
		BrokerServiceReconciler: reconciler,
		instance:                instance,
	}

	reqLogger.V(2).Info("Reconciler Processing...", "CRD.Name", instance.Name, "CRD ver", instance.ObjectMeta.ResourceVersion, "CRD Gen", instance.ObjectMeta.Generation)

	err = processor.InitDeployed(instance, processor.getOwned()...)
	if err != nil {
		reqLogger.Error(err, "failed to initialised owned resources")
		return ctrl.Result{}, err
	}

	// reconcile
	if err = processor.processSpec(); err == nil {
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

	reqLogger.V(1).Info("OK, return result")
	return ctrl.Result{}, nil
}

// instance specifics for a reconciler loop
func (r *BrokerServiceReconciler) getOwned() []client.ObjectList {
	return []client.ObjectList{
		&corev1.SecretList{},
		&broker.ActiveMQArtemisList{},
		&corev1.ServiceList{}}
}

func (r *BrokerServiceReconciler) getOrderedTypeList() []reflect.Type {
	// we want to create/update in this order
	return []reflect.Type{
		reflect.TypeOf(corev1.Secret{}),
		reflect.TypeOf(broker.ActiveMQArtemis{}),
		reflect.TypeOf(corev1.Service{})}
}

func (reconciler *BrokerServiceInstanceReconciler) processSpec() (err error) {
	if err = reconciler.processBroker(); err == nil {
		err = reconciler.processService()
	}
	return err
}

func (reconciler *BrokerServiceInstanceReconciler) processBroker() error {

	var desired *broker.ActiveMQArtemis
	obj := reconciler.CloneOfDeployed(reflect.TypeOf(broker.ActiveMQArtemis{}), reconciler.instance.Name)
	if obj != nil {
		desired = obj.(*broker.ActiveMQArtemis)
	} else {
		desired = common.GenerateArtemis(reconciler.instance.Name, reconciler.instance.Namespace)
	}
	desired.Spec.Restricted = common.NewTrue()
	desired.Spec.DeploymentPlan.PersistenceEnabled = false
	desired.Spec.DeploymentPlan.Clustered = common.NewFalse()
	desired.Spec.DeploymentPlan.Labels = map[string]string{
		fmt.Sprintf("%s-peer-index", reconciler.instance.Name): fmt.Sprintf("%v", 0),
		getPeerLabelKey(reconciler.instance):                   reconciler.instance.Name,
	}
	desired.Spec.Env = reconciler.instance.Spec.Env
	desired.Spec.DeploymentPlan.Resources = reconciler.instance.Spec.Resources

	if reconciler.instance.Spec.Image != nil {
		desired.Spec.DeploymentPlan.Image = *reconciler.instance.Spec.Image
	}

	desired.Spec.DeploymentPlan.ExtraMounts.Secrets = []string{
		reconciler.appPropertiesSecretName(),
	}

	// a place the app controller can modify
	reconciler.processAppSecrets()

	reconciler.TrackDesired(desired)

	return nil
}

func (reconciler *BrokerServiceInstanceReconciler) processAppSecrets() error {
	// avoid restart for app onboarding with existing mount points
	// TODO potentially N app-secrets to overcome 1Mb size limit
	resourceName := types.NamespacedName{
		Namespace: reconciler.instance.Namespace,
		Name:      reconciler.appPropertiesSecretName(),
	}

	var desired *corev1.Secret

	obj := reconciler.CloneOfDeployed(reflect.TypeOf(corev1.Secret{}), resourceName.Name)
	if obj != nil {
		desired = obj.(*corev1.Secret)
	} else {
		desired = secrets.NewSecret(resourceName, nil, nil)
	}

	reconciler.TrackDesired(desired)
	return nil
}

func (reconciler *BrokerServiceInstanceReconciler) appPropertiesSecretName() string {
	return AppPropertiesSecretName(reconciler.instance.Name)
}

func AppPropertiesSecretName(name string) string {
	return fmt.Sprintf("%s-app%s", name, BrokerPropsSuffix)
}

func PropertiesSecretName(name string) string {
	return fmt.Sprintf("%s%s", name, BrokerPropsSuffix)
}

func certSecretName(cr *broker.BrokerService) string {
	return fmt.Sprintf("%s-%s", cr.Name, common.DefaultOperandCertSecretName)
}

func (reconciler *BrokerServiceInstanceReconciler) processStatus(cr *broker.BrokerService, reconcilerError error) (err error) {

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

	var deployedCondition metav1.Condition = metav1.Condition{
		Type: broker.DeployedConditionType,
	}

	var ready bool = false
	obj := reconciler.CloneOfDeployed(reflect.TypeOf(broker.ActiveMQArtemis{}), cr.Name)
	if obj != nil {
		deployed := obj.(*broker.ActiveMQArtemis)
		brokerReady := meta.FindStatusCondition(deployed.Status.Conditions, broker.ReadyConditionType)

		if brokerReady != nil && brokerReady.Status == metav1.ConditionTrue {
			ready = true
		} else {
			deployedCondition.Message = fmt.Sprintf("not ready status %v", deployed.Status)
		}
	}

	if ready {
		deployedCondition.Status = metav1.ConditionTrue
		deployedCondition.Reason = broker.ReadyConditionReason
	} else {
		deployedCondition.Status = metav1.ConditionFalse
		deployedCondition.Reason = broker.DeployedConditionNotReadyReason
	}
	meta.SetStatusCondition(&conditions, deployedCondition)

	common.SetReadyCondition(&conditions)

	if !reflect.DeepEqual(cr.Status.Conditions, conditions) {
		cr.Status.Conditions = conditions
		err = resources.UpdateStatus(reconciler.Client, cr)
	}
	return err
}

func getPeerLabelKey(cr *broker.BrokerService) string {
	return fmt.Sprintf("%s-peers", cr.Name)
}

func (reconciler *BrokerServiceInstanceReconciler) processService() error {

	var desired *corev1.Service

	obj := reconciler.CloneOfDeployed(reflect.TypeOf(corev1.Service{}), reconciler.instance.Name)
	if obj != nil {
		desired = obj.(*corev1.Service)
	} else {
		// TODO labels ?
		desired = &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      reconciler.instance.Name,
				Namespace: reconciler.instance.Namespace,
			},
			Spec: corev1.ServiceSpec{},
		}
	}

	desired.Spec.Selector = map[string]string{
		getPeerLabelKey(reconciler.instance): reconciler.instance.Name,
	}
	desired.Spec.Ports = []corev1.ServicePort{
		{
			Port:       DefaultServicePort,
			TargetPort: intstr.IntOrString{IntVal: DefaultServicePort},
		},
	}
	reconciler.TrackDesired(desired)
	return nil
}

func (r *BrokerServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.BrokerService{}).
		Owns(&broker.ActiveMQArtemis{}).
		Complete(r)
}
