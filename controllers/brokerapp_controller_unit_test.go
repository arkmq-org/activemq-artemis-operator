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
// +kubebuilder:docs-gen:collapse=Apache License
package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func TestSimpleReconcile(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"
	appName := "my-app"

	// Create BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create BrokerApp matching the service
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, app).WithStatusSubresource(app, svc).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify BrokerApp has annotation
	updatedApp := &v1beta1.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	expectedAnnotation := ns + ":" + svcName
	assert.Equal(t, expectedAnnotation, updatedApp.Annotations[common.AppServiceAnnotation])

	// Verify Status
	assert.False(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, v1beta1.DeployedConditionType))
	assert.False(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, v1beta1.ReadyConditionType))
	assert.NotNil(t, updatedApp.Status.Binding)

	bindingSecret := &corev1.Secret{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: updatedApp.Status.Binding.Name, Namespace: ns}, bindingSecret)
	assert.NoError(t, err)

	assert.Equal(t, fmt.Sprintf("%s.%s.svc.%s", svcName, ns, common.GetClusterDomain()), string(bindingSecret.Data["host"]))
	assert.Equal(t, fmt.Sprintf("%d", app.Spec.Acceptor.Port), string(bindingSecret.Data["port"]))
	assert.Equal(t, fmt.Sprintf("amqps://%s.%s.svc.%s:%d", svcName, ns, common.GetClusterDomain(), app.Spec.Acceptor.Port), string(bindingSecret.Data["uri"]))

	// update broker service status to reflect ready with deployed app
	svc.Status.ProvisionedApps = []string{AppIdentity(app)}
	err = cl.Status().Update(context.TODO(), svc)
	assert.NoError(t, err)

	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Status
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)
	assert.True(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, v1beta1.DeployedConditionType))
	assert.True(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, v1beta1.ReadyConditionType))
	assert.NotNil(t, updatedApp.Status.Binding)

}

func TestReconcileNoMatchingService(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	appName := "my-app"

	// Create BrokerApp with a selector that won't match anything
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "non-existent"},
			},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(app).WithStatusSubresource(app).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.Error(t, err)

	// Verify BrokerApp status
	updatedApp := &v1beta1.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	// Check Valid condition
	validCondition := meta.FindStatusCondition(updatedApp.Status.Conditions, v1beta1.ValidConditionType)
	assert.NotNil(t, validCondition)
	assert.Equal(t, v1.ConditionFalse, validCondition.Status)
	assert.Equal(t, "SpecSelectorNoMatch", validCondition.Reason)
}

func TestReconcileValidConditionTransition(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"
	appName := "my-app"

	// Create BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create BrokerApp with non-matching selector
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "non-existent"},
			},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, app).WithStatusSubresource(app).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// 1. Reconcile with non-matching selector
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.Error(t, err) // Expect error because no service found

	// Verify Valid condition is False
	updatedApp := &v1beta1.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	validCond := meta.FindStatusCondition(updatedApp.Status.Conditions, v1beta1.ValidConditionType)
	assert.NotNil(t, validCond)
	assert.Equal(t, v1.ConditionFalse, validCond.Status)
	creationTime := validCond.LastTransitionTime

	// Wait a bit to ensure time difference
	time.Sleep(1 * time.Second)

	// 2. Update App to match service
	updatedApp.Spec.ServiceSelector.MatchLabels["type"] = "broker"
	err = cl.Update(context.TODO(), updatedApp)
	assert.NoError(t, err)

	// Reconcile again
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Valid condition is True and time updated
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	validCond = meta.FindStatusCondition(updatedApp.Status.Conditions, v1beta1.ValidConditionType)
	assert.NotNil(t, validCond)
	assert.Equal(t, v1.ConditionTrue, validCond.Status)
	assert.True(t, validCond.LastTransitionTime.After(creationTime.Time))
}

func TestReconcileStatusUpdateFailure(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	appName := "my-app"
	svcName := "my-broker-service"

	// Create BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create BrokerApp matching the service
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
		},
	}

	// Setup fake client with interceptor to fail Status Update
	interceptorFuncs := interceptor.Funcs{
		SubResourceUpdate: func(ctx context.Context, client client.Client, subResourceName string, obj client.Object, opts ...client.SubResourceUpdateOption) error {
			return fmt.Errorf("simulated status update error")
		},
	}

	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, app).WithStatusSubresource(app).WithInterceptorFuncs(interceptorFuncs).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	result, err := r.Reconcile(context.TODO(), req)

	// Verify error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated status update error")
	assert.Equal(t, time.Duration(0), result.RequeueAfter)
}

func TestReconcileAddressTypeError(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	appName := "my-app"

	// Create BrokerApp with invalid subscriber address (simple address instead of FQQN)
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Capabilities: []v1beta1.AppCapabilityType{
				{
					SubscriberOf: []v1beta1.AppAddressType{
						{Address: "simple-address"}, // Missing "::"
					},
				},
			},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(app).WithStatusSubresource(app).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.Error(t, err)

	// Verify BrokerApp status
	updatedApp := &v1beta1.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	// Check Valid condition
	validCondition := meta.FindStatusCondition(updatedApp.Status.Conditions, v1beta1.ValidConditionType)
	assert.NotNil(t, validCondition)
	assert.Equal(t, v1.ConditionFalse, validCondition.Status)
	assert.Equal(t, "AddressTypeError", validCondition.Reason)
	assert.Contains(t, validCondition.Message, "must specify a FQQN")
}

func TestReconcileDeployedConditionFromBrokerServiceStatus(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"
	appName := "my-app"

	// Create BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create BrokerApp matching the service
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, app).WithStatusSubresource(app, svc).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// 1. Reconcile with BrokerService status not having the app.
	// This first reconcile will annotate the app. The Deployed condition will be False.
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Deployed condition is False
	updatedApp := &v1beta1.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	deployedCond := meta.FindStatusCondition(updatedApp.Status.Conditions, v1beta1.DeployedConditionType)
	assert.NotNil(t, deployedCond)
	assert.Equal(t, v1.ConditionFalse, deployedCond.Status)
	assert.Equal(t, "ProvisioningPending", deployedCond.Reason)

	// 2. Update BrokerService status to include the app
	updatedSvc := &v1beta1.BrokerService{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: svcName, Namespace: ns}, updatedSvc)
	assert.NoError(t, err)

	appIdentity := AppIdentity(app)
	updatedSvc.Status.ProvisionedApps = []string{appIdentity}
	err = cl.Status().Update(context.TODO(), updatedSvc)
	assert.NoError(t, err)

	// Reconcile again to pick up the status change
	_, err = r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify Deployed condition is True
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	deployedCond = meta.FindStatusCondition(updatedApp.Status.Conditions, v1beta1.DeployedConditionType)
	assert.NotNil(t, deployedCond)
	assert.Equal(t, v1.ConditionTrue, deployedCond.Status)
	assert.Equal(t, "Provisioned", deployedCond.Reason)
}

func TestReconcileIdempotentStatus(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"
	appName := "my-app"

	// Create BrokerService
	svc := &v1beta1.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create BrokerApp matching the service
	app := &v1beta1.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: ns,
		},
		Spec: v1beta1.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
		},
	}

	// Setup fake client for first reconcile
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, app).WithStatusSubresource(app, svc).Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// 1. First Reconcile to establish a status
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: appName, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Get the updated app from the fake client
	updatedApp := &v1beta1.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp)
	assert.NoError(t, err)

	// 2. Second Reconcile with the already updated app
	// We need a new client with the updated app and an interceptor to track status updates.
	updateCalled := false
	interceptorFuncs := interceptor.Funcs{
		SubResourceUpdate: func(ctx context.Context, client client.Client, subResourceName string, obj client.Object, opts ...client.SubResourceUpdateOption) error {
			if _, ok := obj.(*v1beta1.BrokerApp); ok {
				updateCalled = true
			}
			return client.SubResource(subResourceName).Update(ctx, obj, opts...)
		},
	}

	cl2 := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, updatedApp).WithStatusSubresource(updatedApp, svc).WithInterceptorFuncs(interceptorFuncs).Build()
	r2 := NewBrokerAppReconciler(cl2, scheme, nil, logr.New(log.NullLogSink{}))
	_, err = r2.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Assert that status update was not called
	assert.False(t, updateCalled, "Status update should not be called on second reconcile if status is unchanged")
}
