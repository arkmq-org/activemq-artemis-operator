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

type ActiveMQArtemisServiceReconciler struct {
	*ReconcilerLoop
}

type ActiveMQArtemisServiceInstanceReconciler struct {
	*ActiveMQArtemisServiceReconciler
	instance *broker.ActiveMQArtemisService
}

func NewActiveMQArtemisServiceReconciler(client client.Client, scheme *runtime.Scheme, config *rest.Config, logger logr.Logger) *ActiveMQArtemisServiceReconciler {
	reconciler := ActiveMQArtemisServiceReconciler{
		ReconcilerLoop: &ReconcilerLoop{KubeBits: &KubeBits{client, scheme, config, logger}},
	}
	reconciler.ReconcilerLoopType = &reconciler
	return &reconciler
}

//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=activemqartemisservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=activemqartemisservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=activemqartemisservices/finalizers,verbs=update

func (reconciler *ActiveMQArtemisServiceReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := reconciler.log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "Reconciling", "ActiveMQArtemisService")

	instance := &broker.ActiveMQArtemisService{}
	var err = reconciler.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	processor := ActiveMQArtemisServiceInstanceReconciler{
		ActiveMQArtemisServiceReconciler: reconciler,
		instance:                         instance,
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
func (r *ActiveMQArtemisServiceReconciler) getOwned() []client.ObjectList {
	return []client.ObjectList{
		&corev1.SecretList{},
		&broker.ActiveMQArtemisList{},
		&corev1.ServiceList{}}
}

func (r *ActiveMQArtemisServiceReconciler) getOrderedTypeList() []reflect.Type {
	// we want to create/update in this order
	return []reflect.Type{
		reflect.TypeOf(corev1.Secret{}),
		reflect.TypeOf(broker.ActiveMQArtemis{}),
		reflect.TypeOf(corev1.Service{})}
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processSpec() (err error) {
	if err = reconciler.processBroker(); err == nil {
		err = reconciler.processService()
	}
	return err
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processBroker() error {

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

	desired.Spec.DeploymentPlan.ExtraMounts.Secrets = []string{
		reconciler.appPropertiesSecretName(),
		reconciler.propertiesSecretName(),
		reconciler.jaasConfigSecretName(),
		EmptyBrokerXml,
	}

	// the broker needs am xml to do config reload - fix and revisit
	reconciler.processXmlSecret()

	desired.Spec.DeploymentPlan.ExtraVolumeMounts = []corev1.VolumeMount{
		{
			Name:      EmptyBrokerXml,
			MountPath: "/app/etc",
		},
	}
	desired.Spec.DeploymentPlan.ExtraVolumes = []corev1.Volume{
		{
			Name: EmptyBrokerXml,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: EmptyBrokerXml,
				},
			},
		},
	}
	reconciler.processJaas()

	// a place the app controller can modify
	reconciler.processAppSecrets()

	propertiesSecret := reconciler.processPropertiesSecrets()

	trustStorePath, err := reconciler.getTrustStorePath()
	if err != nil {
		return err
	}

	// a key in the properties secret
	reconciler.processAcceptors(propertiesSecret, trustStorePath)

	reconciler.TrackDesired(desired)

	return nil
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) getTrustStorePath() (string, error) {

	var caCertSecret *corev1.Secret
	var caSecretKey string
	var err error
	if caCertSecret, err = common.GetOperatorCASecret(reconciler.Client); err == nil {
		if caSecretKey, err = common.GetOperatorCASecretKey(reconciler.Client, caCertSecret); err == nil {
			return fmt.Sprintf("/amq/extra/secrets/%s/%s", caCertSecret.Name, caSecretKey), nil
		}
	}
	return "", err
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processAppSecrets() error {
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

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) appPropertiesSecretName() string {
	return AppPropertiesSecretName(reconciler.instance.Name)
}

func AppPropertiesSecretName(name string) string {
	return fmt.Sprintf("%s-app%s", name, BrokerPropsSuffix)
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) propertiesSecretName() string {
	return PropertiesSecretName(reconciler.instance.Name)
}

func PropertiesSecretName(name string) string {
	return fmt.Sprintf("%s%s", name, BrokerPropsSuffix)
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) jaasConfigSecretName() string {
	return JaasConfigSecretName(reconciler.instance.Name)
}

func JaasConfigSecretName(name string) string {
	return fmt.Sprintf("%s%s", name, jaasConfigSuffix)
}

func certSecretName(cr *broker.ActiveMQArtemisService) string {
	return fmt.Sprintf("%s-%s", cr.Name, common.DefaultOperandCertSecretName)
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processXmlSecret() *corev1.Secret {
	resourceName := types.NamespacedName{
		Namespace: reconciler.instance.Namespace,
		Name:      EmptyBrokerXml,
	}

	var desired *corev1.Secret

	obj := reconciler.CloneOfDeployed(reflect.TypeOf(corev1.Secret{}), resourceName.Name)
	if obj != nil {
		desired = obj.(*corev1.Secret)
	} else {
		desired = secrets.NewSecret(resourceName, nil, nil)
		desired.Data = map[string][]byte{}

	}
	minimalConfig := `<configuration xmlns="urn:activemq"
               xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
               xmlns:xi="http://www.w3.org/2001/XInclude"
               xsi:schemaLocation="urn:activemq /schema/artemis-configuration.xsd">

   <core xmlns="urn:activemq:core" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="urn:activemq:core ">
   </core>
</configuration>
`
	desired.Data["broker.xml"] = []byte(minimalConfig)
	reconciler.TrackDesired(desired)
	return desired
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processPropertiesSecrets() *corev1.Secret {
	resourceName := types.NamespacedName{
		Namespace: reconciler.instance.Namespace,
		Name:      reconciler.propertiesSecretName(),
	}

	var desired *corev1.Secret

	obj := reconciler.CloneOfDeployed(reflect.TypeOf(corev1.Secret{}), resourceName.Name)
	if obj != nil {
		desired = obj.(*corev1.Secret)
	} else {
		// TODO labels ?
		desired = secrets.NewSecret(resourceName, nil, nil)
		desired.Data = map[string][]byte{}
	}
	reconciler.TrackDesired(desired)
	return desired
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processJaas() error {
	resourceName := types.NamespacedName{
		Namespace: reconciler.instance.Namespace,
		Name:      reconciler.jaasConfigSecretName(),
	}

	var desired *corev1.Secret

	obj := reconciler.CloneOfDeployed(reflect.TypeOf(corev1.Secret{}), resourceName.Name)
	if obj != nil {
		desired = obj.(*corev1.Secret)
	} else {

		// TODO
		// if it exists we just use it!
		// may need to validate it too!

		desired = secrets.NewSecret(resourceName, nil, nil)
		desired.Data = map[string][]byte{}
		desired.Data[common.GetCertUsersKey(common.JaasRealm)] = NewPropsWithHeader().Bytes()
		desired.Data[common.GetCertRolesKey(common.JaasRealm)] = NewPropsWithHeader().Bytes()
	}

	// we care about the login modules, not the user/role mappings

	login_config := newBufferWithHeader("//")
	fmt.Fprintf(login_config, "%s {\n", common.JaasRealm)
	if contains(reconciler.instance.Spec.Auth, broker.MTLS) {
		fmt.Fprintln(login_config, "  org.apache.activemq.artemis.spi.core.security.jaas.TextFileCertificateLoginModule required")
		fmt.Fprintln(login_config, "   reload=true")
		fmt.Fprintln(login_config, "   debug=true")
		fmt.Fprintf(login_config, "   org.apache.activemq.jaas.textfiledn.user=%s\n", common.GetCertUsersKey(common.JaasRealm))
		fmt.Fprintf(login_config, "   org.apache.activemq.jaas.textfiledn.role=%s\n", common.GetCertRolesKey(common.JaasRealm))
		fmt.Fprintf(login_config, "   baseDir=\"%s%s\"\n", secretPathBase, resourceName.Name)
		fmt.Fprintln(login_config, "  ;")
	}
	fmt.Fprintln(login_config, "};")
	desired.Data["login.config"] = login_config.Bytes()

	reconciler.TrackDesired(desired)
	return nil
}

func contains(appAuthTypes []broker.AppAuthType, value broker.AppAuthType) bool {
	for _, authType := range appAuthTypes {
		if authType == value {
			return true
		}
	}
	return false
}

/*
## reference user role files with login.config, -Djava.security.auth.login.config=jaas/login.config
## the broker JAAS domain in login.config
acceptorConfigurations.tcp.params.securityDomain=broker

## lock down broker acceptor
## to SCRAM AMQP
acceptorConfigurations.tcp.params.saslMechanisms=SCRAM-SHA-512
acceptorConfigurations.tcp.params.protocols=AMQP
acceptorConfigurations.tcp.params.saslLoginConfigScope=amqp-sasl-scram

*/

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processAcceptors(propertiesSecret *corev1.Secret, trustStorePath string) {

	if len(reconciler.instance.Spec.Acceptors) == 0 {
		return
	}

	pemCfgkey := "_acceptor.pemcfg"
	propertiesSecret.Data[pemCfgkey] = reconciler.makePemCfgProps()
	propertiesSecret.Data["acceptor.properties"] = reconciler.makeAcceptorProps(trustStorePath, pemCfgkey)

}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) makePemCfgProps() []byte {

	buf := NewPropsWithHeader()

	certSecretName := certSecretName(reconciler.instance)

	fmt.Fprintf(buf, "source.key=/amq/extra/secrets/%s/tls.key\n", certSecretName)
	fmt.Fprintf(buf, "source.cert=/amq/extra/secrets/%s/tls.crt\n", certSecretName)

	return buf.Bytes()
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) makeAcceptorProps(trustStorePath, pemCfgkey string) []byte {

	buf := NewPropsWithHeader()

	fmt.Fprintln(buf, "# tls acceptors")
	for _, acceptor := range reconciler.instance.Spec.Acceptors {

		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".factoryClassName=org.apache.activemq.artemis.core.remoting.impl.netty.NettyAcceptorFactory\n", acceptor.Name)

		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.securityDomain=activemq\n", acceptor.Name)

		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.host=${HOSTNAME}\n", acceptor.Name)
		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.port=%d\n", acceptor.Name, DefaultServicePort)

		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.sslEnabled=true\n", acceptor.Name)

		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.needClientAuth=true\n", acceptor.Name)
		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.saslMechanisms=EXTERNAL\n", acceptor.Name)

		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.keyStoreType=PEMCFG\n", acceptor.Name)
		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.keyStorePath=/amq/extra/secrets/%s/%s\n", acceptor.Name, reconciler.propertiesSecretName(), pemCfgkey)
		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.trustStoreType=PEMCA\n", acceptor.Name)
		fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.trustStorePath=%s\n", acceptor.Name, trustStorePath)
	}
	return buf.Bytes()
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processStatus(cr *broker.ActiveMQArtemisService, reconcilerError error) (err error) {

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

func getPeerLabelKey(cr *broker.ActiveMQArtemisService) string {
	return fmt.Sprintf("%s-peers", cr.Name)
}

func (reconciler *ActiveMQArtemisServiceInstanceReconciler) processService() error {

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

func (r *ActiveMQArtemisServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.ActiveMQArtemisService{}).
		Owns(&broker.ActiveMQArtemis{}).
		Complete(r)
}
