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
	"maps"
	"strings"
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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func TestSimpleReconcile(t *testing.T) {

	fakes := map[string]client.Object{}
	var brokerListEntry *v1beta1.BrokerService
	var secretListEntry []corev1.Secret

	var updateCalled = 0

	interceptorFuncs := interceptor.Funcs{
		Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			if o, found := fakes[key.Name]; found {
				obj.SetName(o.GetName())
				obj.SetAnnotations(o.GetAnnotations())
				obj.SetDeletionTimestamp(o.GetDeletionTimestamp())
				obj.SetFinalizers(o.GetFinalizers())

				switch v := obj.(type) {
				case *corev1.Secret:
					if v.Data == nil {
						v.Data = map[string][]byte{}
					}
					fake := o.(*corev1.Secret)
					maps.Copy(v.Data, fake.Data)
				}
				return nil
			}
			return apierrors.NewNotFound(schema.GroupResource{}, key.Name)
		},
		Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
			updateCalled++
			if o, found := fakes[obj.GetName()]; found {
				o.SetAnnotations(obj.GetAnnotations())
				o.SetFinalizers(obj.GetFinalizers())
				o.SetDeletionTimestamp(obj.GetDeletionTimestamp())

				switch v := obj.(type) {
				case *corev1.Secret:
					fake := o.(*corev1.Secret)
					if fake.Data == nil {
						fake.Data = map[string][]byte{}
					}
					maps.Copy(fake.Data, v.Data)
				}

			}
			return nil
		},
		SubResourceUpdate: func(ctx context.Context, client client.Client, subResourceName string, obj client.Object, opts ...client.SubResourceUpdateOption) error {
			if o, found := fakes[obj.GetName()]; found {
				// accept status so we can verify
				o.(*v1beta1.BrokerApp).Status = obj.(*v1beta1.BrokerApp).Status
				return nil
			}
			return nil
		},
		List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {

			switch v := list.(type) {
			case *v1beta1.BrokerServiceList:
				if brokerListEntry != nil {
					v.Items = append(v.Items, *brokerListEntry)
				}
			case *corev1.SecretList:
				if secretListEntry != nil {
					v.Items = append(v.Items, secretListEntry...)
				}
			}
			return nil
		},
	}

	fakeClient := fake.NewClientBuilder().WithInterceptorFuncs(interceptorFuncs).Build()
	v1beta1.AddToScheme(fakeClient.Scheme())

	r := NewBrokerAppReconciler(fakeClient, nil, nil, logr.New(log.NullLogSink{}))

	result, err := r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "not-found", Name: "not-found"}})

	assert.Nil(t, err)
	assert.False(t, result.Requeue)

	cr := v1beta1.BrokerApp{ObjectMeta: v1.ObjectMeta{Name: "one"}}
	fakes[cr.Name] = &cr

	nsName := types.NamespacedName{Namespace: "ns", Name: cr.Name}
	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})

	assert.Equal(t, 3, len(cr.Status.Conditions))
	assert.NotNil(t, err)
	assert.False(t, result.Requeue)
	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, "Ready"))
	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, "Valid"))
	assert.Equal(t, meta.FindStatusCondition(cr.Status.Conditions, "Valid").Reason, "SpecSelectorNoMatch")

	brokerListEntry = &v1beta1.BrokerService{ObjectMeta: v1.ObjectMeta{Name: "B"}}

	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})

	assert.Equal(t, 3, len(cr.Status.Conditions))
	assert.NotNil(t, err)
	assert.False(t, result.Requeue)

	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, "Ready"))
	assert.Equal(t, 1, len(cr.GetFinalizers()))

	// provide service owned secret
	secretListEntry = append(secretListEntry, corev1.Secret{ObjectMeta: v1.ObjectMeta{Name: AppPropertiesSecretName(brokerListEntry.Name)}})
	gvk := brokerListEntry.GetObjectKind().GroupVersionKind()
	isController := true
	ref := v1.OwnerReference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		UID:        brokerListEntry.GetUID(),
		Name:       brokerListEntry.GetName(),
		Controller: &isController,
	}
	secretListEntry[0].SetOwnerReferences([]v1.OwnerReference{ref})
	// track so update can find it
	fakes[secretListEntry[0].Name] = &secretListEntry[0]

	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})

	assert.Equal(t, 3, len(cr.Status.Conditions))
	assert.NotNil(t, err)
	assert.False(t, result.Requeue)
	assert.True(t, strings.Contains(err.Error(), "ca"))

	// provide tls creds secrets
	common.SetOperatorNameSpace("test")
	fakes[common.DefaultOperatorCASecretName] = &corev1.Secret{ObjectMeta: v1.ObjectMeta{Name: common.DefaultOperatorCASecretName}, Data: map[string][]byte{"cert.pem": {}}}

	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})

	assert.Equal(t, 3, len(cr.Status.Conditions))
	assert.Nil(t, err)
	assert.Equal(t, 2, updateCalled)

	// repeat to see no further update
	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})

	assert.Equal(t, 3, len(cr.Status.Conditions))
	assert.Nil(t, err)
	assert.Equal(t, 2, updateCalled)

	// delete
	cr.SetDeletionTimestamp(&v1.Time{})
	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})
	assert.Equal(t, 4, len(cr.Status.Conditions))
	assert.Nil(t, err)
	assert.Equal(t, 3, updateCalled)
	assert.Equal(t, 1, len(cr.GetFinalizers()))
	assert.False(t, result.Requeue)
	assert.Greater(t, result.RequeueAfter, time.Duration(0))
	assert.True(t, meta.IsStatusConditionTrue(cr.Status.Conditions, TerminatingConditionType))

	// final stage
	result, err = r.Reconcile(context.TODO(), ctrl.Request{NamespacedName: nsName})

	assert.Equal(t, 3, len(cr.Status.Conditions))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(cr.GetFinalizers()))
	assert.False(t, meta.IsStatusConditionTrue(cr.Status.Conditions, TerminatingConditionType))

}
