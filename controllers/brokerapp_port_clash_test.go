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
	"testing"

	"github.com/arkmq-org/activemq-artemis-operator/api/v1beta2"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestReconcilePortClashBetweenApps(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta2.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"
	app1Name := "app1"
	app2Name := "app2"
	conflictingPort := int32(61616)

	nsObj := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: ns,
		},
	}

	// Create BrokerService
	svc := &v1beta2.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create first BrokerApp with port 61616
	app1 := &v1beta2.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      app1Name,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: ns + ":" + svcName,
			},
		},
		Spec: v1beta2.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Acceptor: v1beta2.AppAcceptorType{Port: conflictingPort},
		},
	}

	// Create second BrokerApp with the same port 61616 (clash!)
	app2 := &v1beta2.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      app2Name,
			Namespace: ns,
		},
		Spec: v1beta2.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Acceptor: v1beta2.AppAcceptorType{Port: conflictingPort},
		},
	}

	// Setup fake client with both service and app1 already present
	cl := setupBrokerAppIndexer(fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(svc, app1, app2, nsObj).
		WithStatusSubresource(app1, app2, svc)).
		Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile app2 (the one trying to use an already-taken port)
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: app2Name, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)

	// We expect an error because of the port clash (no service capacity)
	assert.Error(t, err, "Expected error due to port clash")

	// Verify BrokerApp status
	updatedApp2 := &v1beta2.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp2)
	assert.NoError(t, err)

	// Check Valid condition - should be True (spec itself is valid)
	validCondition := meta.FindStatusCondition(updatedApp2.Status.Conditions, v1beta2.ValidConditionType)
	assert.NotNil(t, validCondition)
	assert.Equal(t, v1.ConditionTrue, validCondition.Status)
	assert.Equal(t, v1beta2.ValidConditionSuccessReason, validCondition.Reason)

	// Check Deployed condition - should be False due to port clash
	deployedCondition := meta.FindStatusCondition(updatedApp2.Status.Conditions, v1beta2.DeployedConditionType)
	assert.NotNil(t, deployedCondition)
	assert.Equal(t, v1.ConditionFalse, deployedCondition.Status)
	assert.Equal(t, v1beta2.DeployedConditionNoServiceCapacityReason, deployedCondition.Reason)
	assert.Contains(t, deployedCondition.Message, "port")
	assert.Contains(t, deployedCondition.Message, fmt.Sprintf("%d", conflictingPort))

	// Ready should be False
	readyCondition := meta.FindStatusCondition(updatedApp2.Status.Conditions, v1beta2.ReadyConditionType)
	assert.NotNil(t, readyCondition)
	assert.Equal(t, v1.ConditionFalse, readyCondition.Status)
}

func TestReconcileNoPortClashDifferentPorts(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta2.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns := "default"
	svcName := "my-broker-service"
	app1Name := "app1"
	app2Name := "app2"
	port1 := int32(61616)
	port2 := int32(61617) // Different port - no clash

	nsObj := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: ns,
		},
	}

	// Create BrokerService
	svc := &v1beta2.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{"type": "broker"},
		},
	}

	// Create first BrokerApp with port 61616
	app1 := &v1beta2.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      app1Name,
			Namespace: ns,
			Annotations: map[string]string{
				common.AppServiceAnnotation: ns + ":" + svcName,
			},
		},
		Spec: v1beta2.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Acceptor: v1beta2.AppAcceptorType{Port: port1},
		},
	}

	// Create second BrokerApp with different port 61617 (no clash)
	app2 := &v1beta2.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      app2Name,
			Namespace: ns,
		},
		Spec: v1beta2.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Acceptor: v1beta2.AppAcceptorType{Port: port2},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(svc, app1, app2, nsObj).
		WithStatusSubresource(app1, app2, svc).
		WithIndex(&v1beta2.BrokerApp{}, common.AppServiceAnnotation, func(obj client.Object) []string {
			app := obj.(*v1beta2.BrokerApp)
			if val, ok := app.Annotations[common.AppServiceAnnotation]; ok {
				return []string{val}
			}
			return nil
		}).
		Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile app2 (no clash expected)
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: app2Name, Namespace: ns}}
	_, err := r.Reconcile(context.TODO(), req)
	assert.NoError(t, err)

	// Verify BrokerApp status
	updatedApp2 := &v1beta2.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp2)
	assert.NoError(t, err)

	// Check Valid condition - should be True (no port clash)
	validCondition := meta.FindStatusCondition(updatedApp2.Status.Conditions, v1beta2.ValidConditionType)
	assert.NotNil(t, validCondition)
	assert.Equal(t, v1.ConditionTrue, validCondition.Status)
	assert.Equal(t, v1beta2.ValidConditionSuccessReason, validCondition.Reason)
}

func TestReconcilePortClashCrossNamespace(t *testing.T) {
	// Setup scheme
	scheme := runtime.NewScheme()
	_ = v1beta2.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	// Data
	ns1 := "default"
	ns2 := "other"
	svcName := "my-broker-service"
	app1Name := "app1"
	app2Name := "app2"
	conflictingPort := int32(61616)

	ns1Obj := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: ns1,
		},
	}

	ns2Obj := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: ns2,
		},
	}

	// Create BrokerService that allows both namespaces
	svc := &v1beta2.BrokerService{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns1,
			Labels:    map[string]string{"type": "broker"},
		},
		Spec: v1beta2.BrokerServiceSpec{
			AppSelectorExpression: fmt.Sprintf(`app.metadata.namespace in ["%s", "%s"]`, ns1, ns2), // Allow both namespaces
		},
	}

	// Create first BrokerApp in default namespace with port 61616
	app1 := &v1beta2.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      app1Name,
			Namespace: ns1,
			Annotations: map[string]string{
				common.AppServiceAnnotation: ns1 + ":" + svcName,
			},
		},
		Spec: v1beta2.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Acceptor: v1beta2.AppAcceptorType{Port: conflictingPort},
		},
	}

	// Create second BrokerApp in other namespace with same port (clash!)
	app2 := &v1beta2.BrokerApp{
		ObjectMeta: v1.ObjectMeta{
			Name:      app2Name,
			Namespace: ns2,
		},
		Spec: v1beta2.BrokerAppSpec{
			ServiceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{"type": "broker"},
			},
			Acceptor: v1beta2.AppAcceptorType{Port: conflictingPort},
		},
	}

	// Setup fake client
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(svc, app1, app2, ns1Obj, ns2Obj).
		WithStatusSubresource(app1, app2, svc).
		WithIndex(&v1beta2.BrokerApp{}, common.AppServiceAnnotation, func(obj client.Object) []string {
			app := obj.(*v1beta2.BrokerApp)
			if val, ok := app.Annotations[common.AppServiceAnnotation]; ok {
				return []string{val}
			}
			return nil
		}).
		Build()

	// Create Reconciler
	r := NewBrokerAppReconciler(cl, scheme, nil, logr.New(log.NullLogSink{}))

	// Reconcile app2 in other namespace (should detect clash even cross-namespace)
	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: app2Name, Namespace: ns2}}
	_, err := r.Reconcile(context.TODO(), req)

	// We expect an error because of the port clash
	assert.Error(t, err, "Expected error due to cross-namespace port clash")

	// Verify BrokerApp status
	updatedApp2 := &v1beta2.BrokerApp{}
	err = cl.Get(context.TODO(), req.NamespacedName, updatedApp2)
	assert.NoError(t, err)

	// Check Valid condition - should be True (spec itself is valid)
	validCondition := meta.FindStatusCondition(updatedApp2.Status.Conditions, v1beta2.ValidConditionType)
	assert.NotNil(t, validCondition)
	assert.Equal(t, v1.ConditionTrue, validCondition.Status)
	assert.Equal(t, v1beta2.ValidConditionSuccessReason, validCondition.Reason)

	// Check Deployed condition - should be False due to port clash
	deployedCondition := meta.FindStatusCondition(updatedApp2.Status.Conditions, v1beta2.DeployedConditionType)
	assert.NotNil(t, deployedCondition)
	assert.Equal(t, v1.ConditionFalse, deployedCondition.Status)
	assert.Equal(t, v1beta2.DeployedConditionNoServiceCapacityReason, deployedCondition.Reason)
	assert.Contains(t, deployedCondition.Message, "port")
	assert.Contains(t, deployedCondition.Message, fmt.Sprintf("%d", conflictingPort))
}
