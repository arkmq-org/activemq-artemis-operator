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
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(svc, app).WithStatusSubresource(app).Build()

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
	assert.True(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, "Reconciled"))
	assert.True(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, v1beta1.DeployedConditionType))
	assert.True(t, meta.IsStatusConditionTrue(updatedApp.Status.Conditions, v1beta1.ReadyConditionType))
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
