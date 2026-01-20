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

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ActiveMQArtemisAppReconciler struct {
	*ReconcilerLoop
}

type ActiveMQArtemisAppInstanceReconciler struct {
	*ActiveMQArtemisAppReconciler
	instance *broker.ActiveMQArtemisApp
}

func (reconciler ActiveMQArtemisAppInstanceReconciler) processDelete(service *broker.ActiveMQArtemisService, serverConfigPropertiesSecret *corev1.Secret) (err error) {

	namespacePrefix := reconciler.AppIdentity()
	underStorePrefix := fmt.Sprintf("_%s", namespacePrefix)
	for k := range serverConfigPropertiesSecret.Data {
		if strings.HasPrefix(k, underStorePrefix) || strings.HasPrefix(k, namespacePrefix) {
			delete(serverConfigPropertiesSecret.Data, k)
		}
	}
	serverConfigPropertiesSecret.Data[reconciler.UnderscoreAppIdentityPrefixed("rm.properties")] = []byte(fmt.Sprintf("jaasConfigs.%s=-\n", reconciler.jaasConfigRealmName()))
	if err = resources.Update(reconciler.Client, serverConfigPropertiesSecret); err == nil {

		controllerutil.RemoveFinalizer(reconciler.instance, common.AppServiceAnnotation)

		err = resources.Update(reconciler.Client, reconciler.instance)

		// TODO condition for deletion pending resource version config applied to brokerservice
		// TODO: any other update would be a problem with matching, as if we can only deal with one finalizer at a time, that is not good.
	}
	return err
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
	var opts, err = metav1.LabelSelectorAsSelector(cr.Spec.ServiceSelector)
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
			controllerutil.AddFinalizer(cr, common.AppServiceAnnotation)
			// update immediately such that AuthN can find it
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
		var serverConfigPropertiesSecret *corev1.Secret
		serverConfigPropertiesSecret, err = reconciler.getServiceAppProperiesSecret(service)

		if err == nil {

			if reconciler.instance.ObjectMeta.DeletionTimestamp == nil {

				err = reconciler.processCapabilities(service, serverConfigPropertiesSecret)
				if err == nil {
					err = reconciler.processAcceptor(service, serverConfigPropertiesSecret)
				}
			} else {

				if controllerutil.ContainsFinalizer(reconciler.instance, common.AppServiceAnnotation) {
					if err = reconciler.processDelete(service, serverConfigPropertiesSecret); err == nil {

					}
				}
			}
		}
	}
	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) getServiceAppProperiesSecret(service *broker.ActiveMQArtemisService) (*corev1.Secret, error) {
	key := types.NamespacedName{
		Name:      AppPropertiesSecretName(service.Name),
		Namespace: service.Namespace,
	}
	return secrets.RetriveSecret(key, nil, reconciler.Client)
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

func (reconciler *ActiveMQArtemisAppInstanceReconciler) processAcceptor(service *broker.ActiveMQArtemisService, serverConfigPropertiesSecret *corev1.Secret) (err error) {

	// TODO: need data plane trust store, this is access to the control plane trust store
	// they could be the same but we need two accessors
	trustStorePath, err := reconciler.getTrustStorePath(service)
	if err != nil {
		return err
	}

	namespacedName := reconciler.AppIdentity()

	pemCfgkey := reconciler.UnderscoreAppIdentityPrefixed("tls.pemcfg")
	serverConfigPropertiesSecret.Data[pemCfgkey] = reconciler.makePemCfgProps(service)

	realmName := reconciler.jaasConfigRealmName()

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
	fmt.Fprintf(usersBuf, "%s=/.*%s.*/\n", namespacedName, reconciler.instance.Name)

	certUsersCfgKey := reconciler.UnderscoreAppIdentityPrefixed(common.GetCertUsersKey(realmName))
	serverConfigPropertiesSecret.Data[certUsersCfgKey] = usersBuf.Bytes()

	dedupMap := map[string]string{}
	for _, capability := range reconciler.instance.Spec.Capabilities {

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

	certRolesCfgKey := reconciler.UnderscoreAppIdentityPrefixed(common.GetCertRolesKey(realmName))
	serverConfigPropertiesSecret.Data[certRolesCfgKey] = rolesBuf.Bytes()

	acceptorCfgKey := reconciler.AppIdentityPrefixed("acceptor.properties")

	buf := NewPropsWithHeader()

	name := fmt.Sprintf("%d", reconciler.instance.Spec.Acceptor.Port)
	fmt.Fprintln(buf, "# tls acceptor")

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".factoryClassName=org.apache.activemq.artemis.core.remoting.impl.netty.NettyAcceptorFactory\n", name)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.securityDomain=%s\n", name, realmName)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.host=${HOSTNAME}\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.port=%d\n", name, reconciler.instance.Spec.Acceptor.Port)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.sslEnabled=true\n", name)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.needClientAuth=true\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.saslMechanisms=EXTERNAL\n", name)

	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.keyStoreType=PEMCFG\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.keyStorePath=/amq/extra/secrets/%s/%s\n", name, AppPropertiesSecretName(service.Name), pemCfgkey)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.trustStoreType=PEMCA\n", name)
	fmt.Fprintf(buf, "acceptorConfigurations.\"%s\".params.trustStorePath=%s\n", name, trustStorePath)

	// need a matching realm
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.loginModuleClass=org.apache.activemq.artemis.spi.core.security.jaas.TextFileCertificateLoginModule\n", realmName)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.controlFlag=required\n", realmName)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.debug=true\n", realmName)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.\"org.apache.activemq.jaas.textfiledn.role\"=%s\n", realmName, certRolesCfgKey)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.\"org.apache.activemq.jaas.textfiledn.user\"=%s\n", realmName, certUsersCfgKey)
	fmt.Fprintf(buf, "jaasConfigs.\"%s\".modules.cert.params.baseDir=%s%s\n", realmName, secretPathBase, AppPropertiesSecretName(service.Name))

	serverConfigPropertiesSecret.Data[acceptorCfgKey] = buf.Bytes()

	err = resources.Update(reconciler.Client, serverConfigPropertiesSecret)
	if err != nil {
		return err
	}

	return err
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) jaasConfigRealmName() string {
	realmName := fmt.Sprintf("port-%d", reconciler.instance.Spec.Acceptor.Port)
	return realmName
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) getTrustStorePath(service *broker.ActiveMQArtemisService) (string, error) {

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

func (reconciler *ActiveMQArtemisAppInstanceReconciler) makePemCfgProps(service *broker.ActiveMQArtemisService) []byte {

	buf := NewPropsWithHeader()

	certSecretName := certSecretName(service)

	fmt.Fprintf(buf, "source.key=/amq/extra/secrets/%s/tls.key\n", certSecretName)
	fmt.Fprintf(buf, "source.cert=/amq/extra/secrets/%s/tls.crt\n", certSecretName)

	return buf.Bytes()
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) processCapabilities(service *broker.ActiveMQArtemisService, secret *corev1.Secret) (err error) {

	err = reconciler.verifyCapabilityAddressType()
	if err != nil {
		return err
	}

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

	secret.Data = map[string][]byte{
		reconciler.AppIdentityPrefixed("capabilities.properties"): buf.Bytes(),
	}
	err = resources.Update(reconciler.Client, secret)
	return err
}

func escapeForProperties(s string) string {
	s = strings.Replace(s, "::", "\\:\\:", 1)
	s = strings.Replace(s, "=", "\\=", -1)
	s = strings.Replace(s, " ", "\\ ", -1)
	return s
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) AppIdentity() string {
	return reconciler.NameSpacedValue(reconciler.instance.Name)
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) AppIdentityPrefixed(v string) string {
	return DashPrefixValue(reconciler.AppIdentity(), v)
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) UnderscoreAppIdentityPrefixed(v string) string {
	return fmt.Sprintf("_%s", reconciler.AppIdentityPrefixed(v))
}

func (reconciler *ActiveMQArtemisAppInstanceReconciler) NameSpacedValue(v string) string {
	return DashPrefixValue(reconciler.instance.Namespace, v)
}

func DashPrefixValue(prefix, value string) string {
	return fmt.Sprintf("%s-%s", prefix, value)
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
		for index, address := range capability.SubscriberOf {
			if !strings.Contains(address.Address, "::") {
				err = fmt.Errorf("Spec.Capability.SubscriptionOf[%d] address must specify a FQQN, %v", index, err)
				break
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

func (reconciler *ActiveMQArtemisAppInstanceReconciler) doMtlsAuthN(service *broker.ActiveMQArtemisService) (err error) {

	key := types.NamespacedName{
		Name:      service.Name,
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

		dedupMap := map[string]string{}

		for _, candidate := range deployedApps {

			userName := DashPrefixValue(candidate.Name, candidate.Namespace)
			// TODO: pull down full DN from app cert
			fmt.Fprintf(usersBuf, "%s=/.*%s.*/\n", userName, candidate.Name)

			for _, capability := range candidate.Spec.Capabilities {

				roleName := capability.Role
				if roleName == "" {
					roleName = userName
				}

				if len(capability.ConsumerOf) > 0 || len(capability.SubscriberOf) > 0 {
					dedupMap[fmt.Sprintf("%s=%s\n", consumerRole(roleName), userName)] = ""
				}
				if len(capability.ProducerOf) > 0 {
					dedupMap[fmt.Sprintf("%s=%s\n", producerRole(roleName), userName)] = ""
				}
			}
		}

		secret.Data[common.GetCertUsersKey(common.JaasRealm)] = usersBuf.Bytes()

		rolesBuf := NewPropsWithHeader()
		for _, k := range sortedKeys(dedupMap) {
			fmt.Fprint(rolesBuf, k)
		}
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

	var conditions []metav1.Condition = cr.Status.Conditions

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
