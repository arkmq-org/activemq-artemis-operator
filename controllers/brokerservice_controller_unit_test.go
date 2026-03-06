/*
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
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestBrokerServiceReconcileWithAppMove(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	s1Name := "service1"
	s2Name := "service2"
	appName := "my-app"

	common.SetOperatorCASecretName("op_ca")
	t.Cleanup(common.UnsetOperatorCASecretName)

	common.SetOperatorNameSpace(ns)
	t.Cleanup(common.UnsetOperatorNameSpace)

	oc := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "op_ca",
			Namespace: ns,
		},
		Data: map[string][]byte{"ca.pem": []byte("bla")},
	}
	s1 := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s1Name,
			Namespace: ns,
			UID:       types.UID("uid-s1"),
		},
	}
	s2 := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s2Name,
			Namespace: ns,
			UID:       types.UID("uid-s2"),
		},
	}
	app := &v1beta1.BrokerApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: ns + ":" + s1Name,
			},
			UID: types.UID("uid-app"),
		},
		Spec: v1beta1.BrokerAppSpec{
			Acceptor: v1beta1.AppAcceptorType{Port: 61616},
		},
	}

	// Setup fake client with indexer
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(oc, s1, s2, app).WithStatusSubresource(s1, s2, app)
	builder.WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
		app := rawObj.(*v1beta1.BrokerApp)
		val, ok := app.Annotations[common.AppServiceAnnotation]
		if !ok {
			return nil
		}
		return []string{val}
	})

	cl := builder.Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile S1
	reqS1 := ctrl.Request{NamespacedName: types.NamespacedName{Name: s1Name, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), reqS1)
	assert.NoError(t, err)

	// Check S1 secret has app config
	secretS1 := &corev1.Secret{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: AppPropertiesSecretName(s1Name), Namespace: ns}, secretS1)
	assert.NoError(t, err)
	// Check for some key related to the app
	assert.True(t, hasKeyContaining(secretS1.Data, appName), "S1 secret should contain app config")

	// Move App to S2
	err = cl.Get(context.TODO(), types.NamespacedName{Name: appName, Namespace: ns}, app)
	assert.NoError(t, err)
	app.Annotations[common.AppServiceAnnotation] = ns + ":" + s2Name
	assert.NoError(t, cl.Update(context.TODO(), app))

	// Reconcile S1 (should remove app)
	_, err = r.Reconcile(context.TODO(), reqS1)
	assert.NoError(t, err)

	err = cl.Get(context.TODO(), types.NamespacedName{Name: AppPropertiesSecretName(s1Name), Namespace: ns}, secretS1)
	assert.NoError(t, err)
	assert.False(t, hasKeyContaining(secretS1.Data, appName), "S1 secret should NOT contain app config after move")

	// Reconcile S2 (should add app)
	reqS2 := ctrl.Request{NamespacedName: types.NamespacedName{Name: s2Name, Namespace: ns}}
	_, err = r.Reconcile(context.TODO(), reqS2)
	assert.NoError(t, err)

	secretS2 := &corev1.Secret{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: AppPropertiesSecretName(s2Name), Namespace: ns}, secretS2)
	assert.NoError(t, err)
	assert.True(t, hasKeyContaining(secretS2.Data, appName), "S2 secret should contain app config after move")
}

func hasKeyContaining(data map[string][]byte, substring string) bool {
	for k := range data {
		if strings.Contains(k, substring) {
			return true
		}
	}
	return false
}

func TestBrokerServiceReconcileErrorPropagation(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	s1Name := "service1"
	ns := "default"
	s1 := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s1Name,
			Namespace: ns,
			UID:       types.UID("uid-s1"),
		},
	}

	// Setup fake client with interceptor to fail List
	interceptorFuncs := interceptor.Funcs{
		List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			if _, ok := list.(*corev1.SecretList); ok {
				return fmt.Errorf("simulated list error")
			}
			return client.List(ctx, list, opts...)
		},
	}

	// Setup fake client with indexer
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(s1).WithStatusSubresource(s1).WithInterceptorFuncs(interceptorFuncs)
	builder.WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
		app := rawObj.(*v1beta1.BrokerApp)
		val, ok := app.Annotations[common.AppServiceAnnotation]
		if !ok {
			return nil
		}
		return []string{val}
	})

	cl := builder.Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile S1
	reqS1 := ctrl.Request{NamespacedName: types.NamespacedName{Name: s1Name, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), reqS1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated list error")

	err = cl.Get(context.TODO(), reqS1.NamespacedName, s1)
	assert.Nil(t, err)

	assert.True(t, meta.IsStatusConditionPresentAndEqual(s1.Status.Conditions, v1beta1.DeployedConditionType, metav1.ConditionUnknown))
	assert.True(t, meta.IsStatusConditionFalse(s1.Status.Conditions, v1beta1.ReadyConditionType))
}

func TestBrokerServiceReconcileStatusUpdateFailure(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	s1Name := "service1"
	ns := "default"
	s1 := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s1Name,
			Namespace: ns,
			UID:       types.UID("uid-s1"),
		},
	}

	// Setup fake client with interceptor to fail Status Update
	interceptorFuncs := interceptor.Funcs{
		SubResourceUpdate: func(ctx context.Context, client client.Client, subResourceName string, obj client.Object, opts ...client.SubResourceUpdateOption) error {
			return fmt.Errorf("simulated status update error")
		},
	}

	// Setup fake client with indexer
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(s1).WithStatusSubresource(s1).WithInterceptorFuncs(interceptorFuncs)
	builder.WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
		return nil
	})

	cl := builder.Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile S1
	reqS1 := ctrl.Request{NamespacedName: types.NamespacedName{Name: s1Name, Namespace: ns}}
	result, err := r.Reconcile(context.TODO(), reqS1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated status update error")
	assert.Equal(t, time.Duration(0), result.RequeueAfter)
}

func TestBrokerServiceReconcileRequiresIndex(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	s1Name := "service1"
	appName := "my-app"

	s1 := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s1Name,
			Namespace: ns,
			UID:       types.UID("uid-s1"),
		},
	}
	app := &v1beta1.BrokerApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: ns + ":" + s1Name,
			},
			UID: types.UID("uid-app"),
		},
		Spec: v1beta1.BrokerAppSpec{
			Acceptor: v1beta1.AppAcceptorType{Port: 61616},
		},
	}

	// Setup fake client WITHOUT indexer
	// This simulates what happens if SetupWithManager doesn't register the indexer
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(s1, app).WithStatusSubresource(s1, app).Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile S1
	reqS1 := ctrl.Request{NamespacedName: types.NamespacedName{Name: s1Name, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), reqS1)

	// Expect error because List with MatchingFields requires an index in the fake client (and real client)
	assert.Error(t, err)
	// The error message from controller-runtime fake client when index is missing
	assert.Contains(t, err.Error(), "index")
}

func TestReconcileDeployedConditionTransition(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"

	// Create BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
		},
	}

	// Setup fake client with indexer required by controller
	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc).WithStatusSubresource(svc, &v1beta1.ActiveMQArtemis{})
	builder.WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
		app := rawObj.(*v1beta1.BrokerApp)
		val, ok := app.Annotations[common.AppServiceAnnotation]
		if !ok {
			return nil
		}
		return []string{val}
	})
	cl := builder.Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// 1. Reconcile - should create broker but it won't be ready
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: svcName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Deployed condition is False (NotReady)
	updatedSvc := &v1beta1.BrokerService{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)

	deployedCond := meta.FindStatusCondition(updatedSvc.Status.Conditions, v1beta1.DeployedConditionType)
	assert.NotNil(t, deployedCond)
	assert.Equal(t, metav1.ConditionFalse, deployedCond.Status)
	assert.Equal(t, "NotReady", deployedCond.Reason)
	creationTime := deployedCond.LastTransitionTime

	// Wait a bit to ensure time difference
	time.Sleep(1 * time.Second)

	// 2. Update underlying ActiveMQArtemis to be Ready
	brokerCR := &v1beta1.ActiveMQArtemis{}
	err = cl.Get(context.TODO(), req.NamespacedName, brokerCR)
	assert.NoError(t, err)

	// Update status of brokerCR
	meta.SetStatusCondition(&brokerCR.Status.Conditions, metav1.Condition{
		Type:   v1beta1.ReadyConditionType,
		Status: metav1.ConditionTrue,
	})
	meta.SetStatusCondition(&brokerCR.Status.Conditions, metav1.Condition{
		Type:   v1beta1.DeployedConditionType,
		Status: metav1.ConditionTrue,
	})
	err = cl.Status().Update(context.TODO(), brokerCR)
	assert.NoError(t, err)

	// Reconcile again
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Deployed condition is True and time updated
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)

	deployedCond = meta.FindStatusCondition(updatedSvc.Status.Conditions, v1beta1.DeployedConditionType)
	assert.NotNil(t, deployedCond)
	assert.Equal(t, metav1.ConditionTrue, deployedCond.Status)
	assert.Equal(t, v1beta1.ReadyConditionReason, deployedCond.Reason)
	assert.True(t, deployedCond.LastTransitionTime.After(creationTime.Time))
}

func TestBrokerServiceReconcileStatusAppliedApps(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-service"
	appName := "my-app"

	common.SetOperatorCASecretName("op_ca")
	t.Cleanup(common.UnsetOperatorCASecretName)

	common.SetOperatorNameSpace(ns)
	t.Cleanup(common.UnsetOperatorNameSpace)

	oc := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "op_ca",
			Namespace: ns,
		},
		Data: map[string][]byte{"ca.pem": []byte("bla")},
	}

	// BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerServiceSpec{
			Image: StringToPtr("placeholder"),
		},
	}

	// BrokerApp
	app := &v1beta1.BrokerApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: fmt.Sprintf("%s:%s", ns, svcName),
			},
		},
		Spec: v1beta1.BrokerAppSpec{
			Acceptor: v1beta1.AppAcceptorType{Port: 61616},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(oc, svc, app).
		WithStatusSubresource(svc, &v1beta1.ActiveMQArtemis{}).
		WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
			app := rawObj.(*v1beta1.BrokerApp)
			val, ok := app.Annotations[common.AppServiceAnnotation]
			if !ok {
				return nil
			}
			return []string{val}
		}).Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: svcName, Namespace: ns}}

	// 1. First Reconcile: Creates resources
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify BrokerService status is not yet populated with AppliedApps
	updatedSvc := &v1beta1.BrokerService{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)
	assert.Empty(t, updatedSvc.Status.ProvisionedApps)

	// 2. Get generated Secret and its ResourceVersion
	secretName := AppPropertiesSecretName(svcName)
	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: ns}, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, secret.ResourceVersion)
	// Verify annotation is present on the secret
	assert.Equal(t, fmt.Sprintf("%s-%s", ns, appName), secret.Annotations[common.ProvisionedAppsAnnotation])

	// 3. Update ActiveMQArtemis status to simulate broker picking up the config
	brokerCR := &v1beta1.ActiveMQArtemis{}
	err = cl.Get(context.TODO(), req.NamespacedName, brokerCR)
	assert.NoError(t, err)

	brokerCR.Status.Conditions = []metav1.Condition{
		{
			Type:   v1beta1.ReadyConditionType,
			Status: metav1.ConditionTrue,
			Reason: "Ready",
		},
	}
	brokerCR.Status.ExternalConfigs = []v1beta1.ExternalConfigStatus{
		{
			Name:            secretName,
			ResourceVersion: secret.ResourceVersion,
		},
	}
	err = cl.Status().Update(context.TODO(), brokerCR)
	assert.NoError(t, err)

	// 4. Second Reconcile: Should update AppliedApps
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify BrokerService status
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)
	assert.Equal(t, []string{fmt.Sprintf("%s-%s", ns, appName)}, updatedSvc.Status.ProvisionedApps)
}

func TestBrokerServiceReconcileStatusAppliedAppsIncremental(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-service"
	app1Name := "my-app-1"
	app2Name := "my-app-2"

	common.SetOperatorCASecretName("op_ca")
	t.Cleanup(common.UnsetOperatorCASecretName)

	common.SetOperatorNameSpace(ns)
	t.Cleanup(common.UnsetOperatorNameSpace)

	oc := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "op_ca",
			Namespace: ns,
		},
		Data: map[string][]byte{"ca.pem": []byte("bla")},
	}

	// BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerServiceSpec{
			Image: StringToPtr("placeholder"),
		},
	}

	// BrokerApp 1
	app1 := &v1beta1.BrokerApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app1Name,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: fmt.Sprintf("%s:%s", ns, svcName),
			},
		},
		Spec: v1beta1.BrokerAppSpec{
			Acceptor: v1beta1.AppAcceptorType{Port: 61616},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(oc, svc, app1).
		WithStatusSubresource(svc, &v1beta1.ActiveMQArtemis{}).
		WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
			app := rawObj.(*v1beta1.BrokerApp)
			val, ok := app.Annotations[common.AppServiceAnnotation]
			if !ok {
				return nil
			}
			return []string{val}
		}).Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: svcName, Namespace: ns}}

	// 1. Reconcile with App1
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Get generated Secret (v1)
	secretName := AppPropertiesSecretName(svcName)
	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: ns}, secret)
	assert.NoError(t, err)
	secretV1 := secret.ResourceVersion

	// Update Broker Status to point to Secret v1
	brokerCR := &v1beta1.ActiveMQArtemis{}
	err = cl.Get(context.TODO(), req.NamespacedName, brokerCR)
	assert.NoError(t, err)
	brokerCR.Status.Conditions = []metav1.Condition{
		{
			Type:   v1beta1.ReadyConditionType,
			Status: metav1.ConditionTrue,
			Reason: "Ready",
		},
	}
	brokerCR.Status.ExternalConfigs = []v1beta1.ExternalConfigStatus{
		{
			Name:            secretName,
			ResourceVersion: secretV1,
		},
	}
	err = cl.Status().Update(context.TODO(), brokerCR)
	assert.NoError(t, err)

	// Reconcile again to update AppliedApps
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify AppliedApps has App1
	updatedSvc := &v1beta1.BrokerService{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)
	assert.Equal(t, []string{fmt.Sprintf("%s-%s", ns, app1Name)}, updatedSvc.Status.ProvisionedApps)

	// 2. Add App2
	app2 := &v1beta1.BrokerApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app2Name,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: fmt.Sprintf("%s:%s", ns, svcName),
			},
		},
		Spec: v1beta1.BrokerAppSpec{
			Acceptor: v1beta1.AppAcceptorType{Port: 61617},
		},
	}
	err = cl.Create(context.TODO(), app2)
	assert.NoError(t, err)

	// Reconcile to pick up App2. This updates the Secret to v2.
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Secret is updated
	err = cl.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: ns}, secret)
	assert.NoError(t, err)
	assert.NotEqual(t, secretV1, secret.ResourceVersion)
	secretV2 := secret.ResourceVersion

	// Verify AppliedApps STILL has App1 (and not App2 yet, because Broker Status still points to v1)
	// IMPORTANT: It should NOT be empty.
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)
	assert.Equal(t, []string{fmt.Sprintf("%s-%s", ns, app1Name)}, updatedSvc.Status.ProvisionedApps)

	// 3. Update Broker Status to point to Secret v2
	err = cl.Get(context.TODO(), req.NamespacedName, brokerCR)
	assert.NoError(t, err)
	brokerCR.Status.ExternalConfigs = []v1beta1.ExternalConfigStatus{
		{
			Name:            secretName,
			ResourceVersion: secretV2,
		},
	}
	err = cl.Status().Update(context.TODO(), brokerCR)
	assert.NoError(t, err)

	// Reconcile again
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify AppliedApps has App1 and App2
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)
	expectedApps := []string{fmt.Sprintf("%s-%s", ns, app1Name), fmt.Sprintf("%s-%s", ns, app2Name)}
	sort.Strings(expectedApps)
	sort.Strings(updatedSvc.Status.ProvisionedApps)
	assert.Equal(t, expectedApps, updatedSvc.Status.ProvisionedApps)
}

func TestBrokerServiceReconcileAppsProvisionedCondition(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-service"

	// BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerServiceSpec{
			Image: StringToPtr("placeholder"),
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(svc).
		WithStatusSubresource(svc, &v1beta1.ActiveMQArtemis{}).
		WithIndex(&v1beta1.BrokerApp{}, common.AppServiceAnnotation, func(rawObj client.Object) []string {
			app := rawObj.(*v1beta1.BrokerApp)
			val, ok := app.Annotations[common.AppServiceAnnotation]
			if !ok {
				return nil
			}
			return []string{val}
		}).
		Build()

	// Create Reconciler
	r := NewBrokerServiceReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: svcName, Namespace: ns}}

	// 1. First Reconcile: Creates resources
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify AppsProvisioned condition is False (WaitingForBroker)
	updatedSvc := &v1beta1.BrokerService{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)

	cond := meta.FindStatusCondition(updatedSvc.Status.Conditions, "AppsProvisioned")
	assert.NotNil(t, cond)
	assert.Equal(t, metav1.ConditionFalse, cond.Status)
	assert.Equal(t, "WaitingForBroker", cond.Reason)

	// 2. Get generated Secret and its ResourceVersion
	secretName := AppPropertiesSecretName(svcName)
	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: ns}, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, secret.ResourceVersion)

	// 3. Update ActiveMQArtemis status to simulate broker picking up the config
	brokerCR := &v1beta1.ActiveMQArtemis{}
	err = cl.Get(context.TODO(), req.NamespacedName, brokerCR)
	assert.NoError(t, err)

	brokerCR.Status.Conditions = []metav1.Condition{
		{
			Type:   v1beta1.ReadyConditionType,
			Status: metav1.ConditionTrue,
			Reason: "Ready",
		},
	}
	brokerCR.Status.ExternalConfigs = []v1beta1.ExternalConfigStatus{
		{
			Name:            secretName,
			ResourceVersion: secret.ResourceVersion,
		},
	}
	err = cl.Status().Update(context.TODO(), brokerCR)
	assert.NoError(t, err)

	currentSecretResourceVersion := secret.ResourceVersion

	// 4. Second Reconcile: Should update AppsProvisioned
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// verify no resource version change, still no apps
	err = cl.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: ns}, secret)
	assert.NoError(t, err)
	assert.Equal(t, currentSecretResourceVersion, secret.ResourceVersion)

	// Verify AppsProvisioned condition is True
	err = cl.Get(context.TODO(), req.NamespacedName, updatedSvc)
	assert.NoError(t, err)

	cond = meta.FindStatusCondition(updatedSvc.Status.Conditions, "AppsProvisioned")
	assert.NotNil(t, cond)
	assert.Equal(t, metav1.ConditionTrue, cond.Status)
	assert.Equal(t, "Synced", cond.Reason)
}
