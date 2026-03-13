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
	"testing"

	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestFindServiceWithCapacity(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = brokerv1beta1.AddToScheme(scheme)

	tests := []struct {
		name                string
		app                 *brokerv1beta1.BrokerApp
		services            []brokerv1beta1.BrokerService
		existingApps        []brokerv1beta1.BrokerApp
		expectedServiceName string
		expectError         bool
	}{
		{
			name: "no resource constraints - picks first service",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "new-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					Resources: corev1.ResourceRequirements{}, // No resources specified
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "test"},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service2", Namespace: "test"},
				},
			},
			expectedServiceName: "service1",
			expectError:         false,
		},
		{
			name: "picks service with most available memory",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "new-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service2", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("4Gi"),
							},
						},
					},
				},
			},
			expectedServiceName: "service2", // service2 has more capacity
			expectError:         false,
		},
		{
			name: "considers already provisioned apps",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "new-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service2", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("4Gi"),
							},
						},
					},
				},
			},
			existingApps: []brokerv1beta1.BrokerApp{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "existing-app",
						Namespace: "test",
						Annotations: map[string]string{
							common.AppServiceAnnotation: "test:service2",
						},
					},
					Spec: brokerv1beta1.BrokerAppSpec{
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("3Gi"), // service2 now has less available
							},
						},
					},
				},
			},
			expectedServiceName: "service1", // service1 has more available capacity now
			expectError:         false,
		},
		{
			name: "no service has enough capacity",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "new-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("5Gi"),
						},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
					},
				},
			},
			expectedServiceName: "",
			expectError:         true,
		},
		{
			name: "service with no limit has unlimited capacity",
			app: &brokerv1beta1.BrokerApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "new-app",
					Namespace: "test",
				},
				Spec: brokerv1beta1.BrokerAppSpec{
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("100Gi"),
						},
					},
				},
			},
			services: []brokerv1beta1.BrokerService{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service1", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "service2", Namespace: "test"},
					Spec: brokerv1beta1.BrokerServiceSpec{
						Resources: corev1.ResourceRequirements{
							// No limits specified = unlimited
						},
					},
				},
			},
			expectedServiceName: "service2", // service2 has unlimited capacity
			expectError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fake client with services and existing apps
			objs := make([]runtime.Object, 0, len(tt.services)+len(tt.existingApps)+1)
			objs = append(objs, tt.app)
			for i := range tt.services {
				objs = append(objs, &tt.services[i])
			}
			for i := range tt.existingApps {
				objs = append(objs, &tt.existingApps[i])
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithRuntimeObjects(objs...).
				WithIndex(&brokerv1beta1.BrokerApp{}, common.AppServiceAnnotation, func(obj client.Object) []string {
					app := obj.(*brokerv1beta1.BrokerApp)
					if val, ok := app.Annotations[common.AppServiceAnnotation]; ok {
						return []string{val}
					}
					return nil
				}).
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
			}

			// Create service list
			serviceList := &brokerv1beta1.BrokerServiceList{
				Items: tt.services,
			}

			// Call findServiceWithCapacity
			chosen, err := reconciler.findServiceWithCapacity(serviceList)

			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("findServiceWithCapacity() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Check chosen service
			if tt.expectedServiceName != "" {
				if chosen == nil {
					t.Errorf("expected service %s, got nil", tt.expectedServiceName)
				} else if chosen.Name != tt.expectedServiceName {
					t.Errorf("expected service %s, got %s", tt.expectedServiceName, chosen.Name)
				}
			} else {
				if chosen != nil {
					t.Errorf("expected no service, got %s", chosen.Name)
				}
			}
		})
	}
}
