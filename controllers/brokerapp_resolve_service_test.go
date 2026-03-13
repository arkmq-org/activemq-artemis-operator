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
	"testing"

	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestResolveBrokerService(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = brokerv1beta1.AddToScheme(scheme)

	tests := []struct {
		name                   string
		app                    *brokerv1beta1.BrokerApp
		services               []brokerv1beta1.BrokerService
		expectedServiceName    string
		expectedAnnotation     string
		expectedError          bool
		expectedValidCondition metav1.ConditionStatus
	}{
		{
			name: "initial assignment - finds matching service",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					ServiceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"env": "dev"},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "test",
						Labels:    map[string]string{"env": "dev"},
					},
				},
			},
			expectedServiceName:    "service1",
			expectedAnnotation:     "test:service1",
			expectedError:          false,
			expectedValidCondition: metav1.ConditionTrue,
		},
		{
			name: "initial assignment - no matching services",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					ServiceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"env": "prod"},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "test",
						Labels:    map[string]string{"env": "dev"},
					},
				},
			},
			expectedServiceName:    "",
			expectedAnnotation:     "",
			expectedError:          true,
			expectedValidCondition: metav1.ConditionFalse,
		},
		{
			name: "existing annotation - service still matches",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app",
					Namespace: "test",
					Annotations: map[string]string{
						common.AppServiceAnnotation: "test:service1",
					},
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					ServiceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"env": "dev"},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "test",
						Labels:    map[string]string{"env": "dev"},
					},
				},
			},
			expectedServiceName:    "service1",
			expectedAnnotation:     "test:service1",
			expectedError:          false,
			expectedValidCondition: metav1.ConditionTrue,
		},
		{
			name: "selector changed - reassign to new service",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app",
					Namespace: "test",
					Annotations: map[string]string{
						common.AppServiceAnnotation: "test:service1", // Previously bound to service1
					},
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					ServiceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"env": "prod"}, // Now selecting prod
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service1",
						Namespace: "test",
						Labels:    map[string]string{"env": "dev"}, // service1 is dev
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service2",
						Namespace: "test",
						Labels:    map[string]string{"env": "prod"}, // service2 is prod
					},
				},
			},
			expectedServiceName:    "service2",
			expectedAnnotation:     "test:service2", // Should reassign to service2
			expectedError:          false,
			expectedValidCondition: metav1.ConditionTrue,
		},
		{
			name: "service deleted - annotation exists but no matching services",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app",
					Namespace: "test",
					Annotations: map[string]string{
						common.AppServiceAnnotation: "test:service1",
					},
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					ServiceSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"env": "dev"},
					},
				},
			},
			services:               []brokerv1beta1.BrokerService{}, // Service deleted
			expectedServiceName:    "",
			expectedAnnotation:     "test:service1", // Annotation preserved
			expectedError:          false,           // No error, just no service available
			expectedValidCondition: metav1.ConditionTrue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake client with app and services
			objs := make([]runtime.Object, 0, len(tt.services)+1)
			objs = append(objs, tt.app)
			for i := range tt.services {
				objs = append(objs, &tt.services[i])
			}
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(objs...).
				Build()

			// Create reconciler
			reconciler := &BrokerAppInstanceReconciler{
				BrokerAppReconciler: &BrokerAppReconciler{
					ReconcilerLoop: &ReconcilerLoop{
						KubeBits: &KubeBits{
							Client: fakeClient,
							Scheme: scheme,
						},
					},
				},
				instance: tt.app,
				status:   &brokerv1beta1.BrokerAppStatus{},
			}

			// Call resolveBrokerService
			err := reconciler.resolveBrokerService()

			// Check error expectation
			if (err != nil) != tt.expectedError {
				t.Errorf("resolveBrokerService() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			// Check service assignment
			if tt.expectedServiceName != "" {
				if reconciler.service == nil {
					t.Errorf("expected service to be assigned to %s, got nil", tt.expectedServiceName)
				} else if reconciler.service.Name != tt.expectedServiceName {
					t.Errorf("expected service name %s, got %s", tt.expectedServiceName, reconciler.service.Name)
				}
			} else {
				if reconciler.service != nil {
					t.Errorf("expected no service assignment, got %s", reconciler.service.Name)
				}
			}

			// Check annotation was updated correctly
			if tt.expectedAnnotation != "" {
				// Fetch the updated app from the fake client
				updatedApp := &brokerv1beta1.BrokerApp{}
				key := types.NamespacedName{Name: tt.app.Name, Namespace: tt.app.Namespace}
				if err := fakeClient.Get(context.TODO(), key, updatedApp); err != nil {
					t.Fatalf("failed to get updated app: %v", err)
				}
				actualAnnotation := updatedApp.Annotations[common.AppServiceAnnotation]
				if actualAnnotation != tt.expectedAnnotation {
					t.Errorf("expected annotation %s, got %s", tt.expectedAnnotation, actualAnnotation)
				}
			}

			// Check Valid condition
			validCond := meta.FindStatusCondition(reconciler.status.Conditions, brokerv1beta1.ValidConditionType)
			if validCond != nil && validCond.Status != tt.expectedValidCondition {
				t.Errorf("expected Valid condition status %s, got %s", tt.expectedValidCondition, validCond.Status)
			}
		})
	}
}
