/*
Copyright 2021.

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
package common

import (
	"os"
	"testing"
	"time"

	"github.com/RHsyseng/operator-utils/pkg/olm"
	"github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Common Test", func() {
	It("Default Resync Period", func() {

		currentPeriod := 30 * time.Second // default

		valueFromEnv, defined := os.LookupEnv("RECONCILE_RESYNC_PERIOD")
		if defined {
			currentPeriod, _ = time.ParseDuration(valueFromEnv)
		}

		Expect(GetReconcileResyncPeriod()).To(Equal(currentPeriod))
	})

	It("getDeploymentCondition", func() {

		cr := &v1beta1.ActiveMQArtemis{
			Spec: v1beta1.ActiveMQArtemisSpec{
				DeploymentPlan: v1beta1.DeploymentPlanType{
					Size: Int32ToPtr(2),
				},
			},
		}

		condition := getDeploymentCondition(cr, nil, true, nil)
		Expect(condition.Status).Should(BeEquivalentTo(metav1.ConditionFalse))

		cr.Status.PodStatus = olm.DeploymentStatus{
			Ready: []string{"a", "b"},
		}

		condition = getDeploymentCondition(cr, nil, true, nil)
		Expect(condition.Status).Should(BeEquivalentTo(metav1.ConditionTrue))

		cr.Status.PodStatus = olm.DeploymentStatus{
			Ready: []string{"a", "b", "c"}, // over provisioned still true when scaling down
		}
		cr.Status.Conditions = []metav1.Condition{{Status: metav1.ConditionTrue, Type: v1beta1.ScaleDownPendingConditionType, Reason: v1beta1.ScaleDownPendingConditionPendingEmptyReason}}

		condition = getDeploymentCondition(cr, nil, true, nil)
		Expect(condition.Status).Should(BeEquivalentTo(metav1.ConditionTrue))
		Expect(condition.Reason).ShouldNot(BeEquivalentTo(v1beta1.ScaleDownPendingConditionType))

		cr.Status.PodStatus = olm.DeploymentStatus{
			Ready: []string{"a"},
		}
		condition = getDeploymentCondition(cr, nil, true, nil)
		Expect(condition.Status).Should(BeEquivalentTo(metav1.ConditionFalse))
	})
})
