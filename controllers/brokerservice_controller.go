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
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
//+kubebuilder:rbac:groups=broker.amq.io,namespace=activemq-artemis-operator,resources=brokerapps,verbs=get;list;watch

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

	localLoop := &ReconcilerLoop{
		KubeBits:           reconciler.KubeBits,
		ReconcilerLoopType: reconciler,
	}

	processor := BrokerServiceInstanceReconciler{
		BrokerServiceReconciler: &BrokerServiceReconciler{ReconcilerLoop: localLoop},
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

	if statusErr := processor.processStatus(instance, err); statusErr != nil {
		if err == nil {
			return ctrl.Result{}, statusErr
		}
	}
	return ctrl.Result{}, err
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

	// find all apps that select this service
	apps := &broker.BrokerAppList{}
	serviceKey := fmt.Sprintf("%s:%s", reconciler.instance.Namespace, reconciler.instance.Name)
	if err := reconciler.Client.List(context.TODO(), apps, client.MatchingFields{common.AppServiceAnnotation: serviceKey}); err != nil {
		return err
	}

	// reset data
	desired.Data = make(map[string][]byte)

	for _, app := range apps.Items {
		if app.DeletionTimestamp == nil {
			if err := reconciler.processCapabilities(desired, &app); err != nil {
				reconciler.log.Error(err, "failed to process capabilities for app", "app", app.Name)
			}
			if err := reconciler.processAcceptor(desired, &app); err != nil {
				reconciler.log.Error(err, "failed to process acceptor for app", "app", app.Name)
			}
		}
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
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &broker.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
		app := rawObj.(*broker.BrokerApp)
		val, ok := app.Annotations[common.AppServiceAnnotation]
		if !ok {
			return nil
		}
		return []string{val}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&broker.BrokerService{}).
		Owns(&broker.ActiveMQArtemis{}).
		Watches(&broker.BrokerApp{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, a client.Object) []reconcile.Request {
			app := a.(*broker.BrokerApp)
			if val, ok := app.Annotations[common.AppServiceAnnotation]; ok {
				parts := strings.Split(val, ":")
				if len(parts) == 2 {
					return []reconcile.Request{
						{NamespacedName: types.NamespacedName{
							Namespace: parts[0],
							Name:      parts[1],
						}},
					}
				}
			}
			return []reconcile.Request{}
		})).
		Complete(r)
}

type AddressConfig struct {
	senderRoles   map[string]string
	consumerRoles map[string]string
}

type AddressTracker struct {
	names map[string]AddressConfig
}

func newAddressTracker() *AddressTracker {
	return &AddressTracker{names: map[string]AddressConfig{}}
}

func (t *AddressTracker) newAddressConfig() AddressConfig {
	return AddressConfig{senderRoles: map[string]string{}, consumerRoles: map[string]string{}}
}

func (t *AddressTracker) track(address *broker.AppAddressType) *AddressConfig {

	var present bool
	var entry AddressConfig
	if entry, present = t.names[address.Address]; !present {
		entry = t.newAddressConfig()
		t.names[address.Address] = entry
	}
	return &entry
}

func (reconciler *BrokerServiceInstanceReconciler) processCapabilities(secret *corev1.Secret, app *broker.BrokerApp) (err error) {

	err = reconciler.verifyCapabilityAddressType(app)
	if err != nil {
		return err
	}

	addressTracker := newAddressTracker()

	for _, capability := range app.Spec.Capabilities {

		var role = capability.Role
		if role == "" {
			role = AppIdentity(app)
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

		for _, address := range capability.SubscriberOf {
			entry = addressTracker.track(&address)
			entry.consumerRoles[role] = role
		}
	}

	props := map[string]string{} // need to dedup
	for addressName, addr := range addressTracker.names {
		fqqn := strings.SplitN(addressName, "::", 2)
		if len(fqqn) > 1 {
			address := escapeForProperties(fqqn[0])
			queueName := escapeForProperties(fqqn[1])
			props[fmt.Sprintf("addressConfigurations.\"%s\".routingTypes=ANYCAST,MULTICAST\n", address)] = ""
			props[fmt.Sprintf("addressConfigurations.\"%s\".queueConfigs.\"%s\".routingType=MULTICAST\n", address, queueName)] = ""
			props[fmt.Sprintf("addressConfigurations.\"%s\".queueConfigs.\"%s\".address=%s\n", address, queueName, address)] = ""
		} else {
			props[fmt.Sprintf("addressConfigurations.\"%s\".routingTypes=ANYCAST,MULTICAST\n", addressName)] = ""
			props[fmt.Sprintf("addressConfigurations.\"%s\".queueConfigs.\"%s\".routingType=ANYCAST\n", addressName, addressName)] = ""
		}

		// use fqqn as is for RBAC
		addressName = escapeForProperties(addressName)

		for _, role := range addr.senderRoles {
			props[fmt.Sprintf("securityRoles.\"%s\".\"%s\".send=true\n", addressName, producerRole(role))] = ""

		}
		for _, role := range addr.consumerRoles {
			props[fmt.Sprintf("securityRoles.\"%s\".\"%s\".consume=true\n", addressName, consumerRole(role))] = ""
		}
	}

	buf := NewPropsWithHeader()
	for _, k := range sortedKeys(props) {
		fmt.Fprint(buf, k)
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[AppIdentityPrefixed(app, "capabilities.properties")] = buf.Bytes()

	return err
}

func (reconciler *BrokerServiceInstanceReconciler) verifyCapabilityAddressType(app *broker.BrokerApp) (err error) {

	for _, capability := range app.Spec.Capabilities {
		for index, address := range capability.SubscriberOf {
			if !strings.Contains(address.Address, "::") {
				err = fmt.Errorf("Spec.Capability.SubscriptionOf[%d] address must specify a FQQN, %v", index, err)
				break
			}
		}
	}
	if err != nil {
		meta.SetStatusCondition(&app.Status.Conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  "AddressTypeError",
			Message: err.Error(),
		})
	}
	return err
}

func (reconciler *BrokerServiceInstanceReconciler) processAcceptor(serverConfigPropertiesSecret *corev1.Secret, app *broker.BrokerApp) (err error) {

	// TODO: need data plane trust store, this is access to the control plane trust store
	// they could be the same but we need two accessors
	trustStorePath, err := reconciler.getTrustStorePath(reconciler.instance)
	if err != nil {
		return err
	}

	namespacedName := AppIdentity(app)

	pemCfgkey := UnderscoreAppIdentityPrefixed(app, "tls.pemcfg")
	serverConfigPropertiesSecret.Data[pemCfgkey] = reconciler.makePemCfgProps(reconciler.instance)

	realmName := jaasConfigRealmName(app)

	// process authN cert login module params

	/* TODO: pull down full DN from app cert

	var appCert *tls.Certificate
	if appCert, err = common.ExtractCertFromSecret(...); err != nil {
		return nil, err
	}

	var appCertSubject *pkix.Name
	if operatorCertSubject, err = common.ExtractCertSubject(appCert); err != nil {
		return nil, err
	}
	*/
	usersBuf := NewPropsWithHeader()
	fmt.Fprintf(usersBuf, "%s=/.*%s.*/\n", namespacedName, app.Name)

	certUsersCfgKey := UnderscoreAppIdentityPrefixed(app, common.GetCertUsersKey(realmName))
	serverConfigPropertiesSecret.Data[certUsersCfgKey] = usersBuf.Bytes()

	dedupMap := map[string]string{}
	for _, capability := range app.Spec.Capabilities {

		roleName := capability.Role
		if roleName == "" {
			roleName = namespacedName
		}

		if len(capability.ConsumerOf) > 0 || len(capability.SubscriberOf) > 0 {
			dedupMap[fmt.Sprintf("%s=%s\n", consumerRole(roleName), namespacedName)] = ""
		}
		if len(capability.ProducerOf) > 0 {
			dedupMap[fmt.Sprintf("%s=%s\n", producerRole(roleName), namespacedName)] = ""
		}
	}

	rolesBuf := NewPropsWithHeader()
	for _, k := range sortedKeys(dedupMap) {
		fmt.Fprint(rolesBuf, k)
	}

	certRolesCfgKey := UnderscoreAppIdentityPrefixed(app, common.GetCertRolesKey(realmName))
	serverConfigPropertiesSecret.Data[certRolesCfgKey] = rolesBuf.Bytes()

	acceptorCfgKey := AppIdentityPrefixed(app, "acceptor.properties")

	buf := NewPropsWithHeader()

	name := fmt.Sprintf("%d", app.Spec.Acceptor.Port)
	fmt.Fprintln(buf, "# tls acceptor")

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".factoryClassName=org.apache.activemq.artemis.core.remoting.impl.netty.NettyAcceptorFactory\n", name)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.securityDomain=%s\n", name, realmName)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.host=${HOSTNAME}\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.port=%d\n", name, app.Spec.Acceptor.Port)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.sslEnabled=true\n", name)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.needClientAuth=true\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.saslMechanisms=EXTERNAL\n", name)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.keyStoreType=PEMCFG\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.keyStorePath=/amq/extra/secrets/%s/%s\n", name, AppPropertiesSecretName(reconciler.instance.Name), pemCfgkey)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.trustStoreType=PEMCA\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.trustStorePath=%s\n", name, trustStorePath)

	// need a matching realm
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.loginModuleClass=org.apache.activemq.artemis.spi.core.security.jaas.TextFileCertificateLoginModule\n", realmName)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.controlFlag=required\n", realmName)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.debug=true\n", realmName)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.\"org.apache.activemq.jaas.textfiledn.role\"=%s\n", realmName, certRolesCfgKey)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.\"org.apache.activemq.jaas.textfiledn.user\"=%s\n", realmName, certUsersCfgKey)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.baseDir=%s%s\n", realmName, secretPathBase, AppPropertiesSecretName(reconciler.instance.Name))

	serverConfigPropertiesSecret.Data[acceptorCfgKey] = buf.Bytes()

	return err
}

func (reconciler *BrokerServiceInstanceReconciler) getTrustStorePath(_ *broker.BrokerService) (string, error) {

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

func (reconciler *BrokerServiceInstanceReconciler) makePemCfgProps(service *broker.BrokerService) []byte {

	buf := NewPropsWithHeader()

	certSecretName := certSecretName(service)

	fmt.Fprintf(buf, "source.key=/amq/extra/secrets/%s/tls.key\n", certSecretName)
	fmt.Fprintf(buf, "source.cert=/amq/extra/secrets/%s/tls.crt\n", certSecretName)

	return buf.Bytes()
}

func jaasConfigRealmName(app *broker.BrokerApp) string {
	realmName := fmt.Sprintf("port-%d", app.Spec.Acceptor.Port)
	return realmName
}

func escapeForProperties(s string) string {
	s = strings.Replace(s, "::", "\\:\\:", 1)
	s = strings.Replace(s, "=", "\\=", -1)
	s = strings.Replace(s, " ", "\\ ", -1)
	return s
}

func producerRole(prefix string) string {
	return fmt.Sprintf("%s-producer", prefix)
}

func consumerRole(prefix string) string {
	return fmt.Sprintf("%s-consumer", prefix)
}

func AppIdentity(app *broker.BrokerApp) string {
	return NameSpacedValue(app, app.Name)
}

func AppIdentityPrefixed(app *broker.BrokerApp, v string) string {
	return DashPrefixValue(AppIdentity(app), v)
}

func UnderscoreAppIdentityPrefixed(app *broker.BrokerApp, v string) string {
	return fmt.Sprintf("_%s", AppIdentityPrefixed(app, v))
}

func NameSpacedValue(app *broker.BrokerApp, v string) string {
	return DashPrefixValue(app.Namespace, v)
}

func DashPrefixValue(prefix, value string) string {
	return fmt.Sprintf("%s-%s", prefix, value)
}
