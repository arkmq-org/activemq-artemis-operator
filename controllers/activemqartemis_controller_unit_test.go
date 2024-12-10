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
	"strings"
	"testing"

	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/arkmq-org/activemq-artemis-operator/pkg/utils/selectors"
	"k8s.io/apimachinery/pkg/api/meta"
)

func TestValidate(t *testing.T) {

	cr := &brokerv1beta1.ActiveMQArtemis{
		Spec: brokerv1beta1.ActiveMQArtemisSpec{
			ResourceTemplates: []brokerv1beta1.ResourceTemplate{
				{
					// reserved key
					Labels: map[string]string{selectors.LabelAppKey: "myAppKey"},
				},
			},
		},
	}

	namer := MakeNamers(cr)

	r := NewActiveMQArtemisReconciler(&NillCluster{}, ctrl.Log, isOpenshift)
	ri := NewActiveMQArtemisReconcilerImpl(cr, r)

	valid, retry := ri.validate(cr, k8sClient, *namer)

	assert.False(t, valid)
	assert.False(t, retry)

	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, brokerv1beta1.ValidConditionType))

	condition := meta.FindStatusCondition(cr.Status.Conditions, brokerv1beta1.ValidConditionType)
	assert.Equal(t, condition.Reason, brokerv1beta1.ValidConditionFailedReservedLabelReason)
	assert.True(t, strings.Contains(condition.Message, "Templates[0]"))
}

func TestValidateBrokerPropsDuplicate(t *testing.T) {

	cr := &brokerv1beta1.ActiveMQArtemis{
		Spec: brokerv1beta1.ActiveMQArtemisSpec{
			BrokerProperties: []string{
				"min=X",
				"min=y",
			},
		},
	}

	namer := MakeNamers(cr)

	r := NewActiveMQArtemisReconciler(&NillCluster{}, ctrl.Log, isOpenshift)
	ri := NewActiveMQArtemisReconcilerImpl(cr, r)

	valid, retry := ri.validate(cr, k8sClient, *namer)

	assert.False(t, valid)
	assert.False(t, retry)

	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, brokerv1beta1.ValidConditionType))

	condition := meta.FindStatusCondition(cr.Status.Conditions, brokerv1beta1.ValidConditionType)
	assert.Equal(t, condition.Reason, brokerv1beta1.ValidConditionFailedDuplicateBrokerPropertiesKey)
	assert.True(t, strings.Contains(condition.Message, "min"))
}

func TestValidateBrokerPropsDuplicateOnFirstEquals(t *testing.T) {

	cr := &brokerv1beta1.ActiveMQArtemis{
		Spec: brokerv1beta1.ActiveMQArtemisSpec{
			BrokerProperties: []string{
				"nameWith\\=equals_not_matched=X",
				"nameWith\\=equals_not_matched=Y",
			},
		},
	}

	namer := MakeNamers(cr)

	r := NewActiveMQArtemisReconciler(&NillCluster{}, ctrl.Log, isOpenshift)
	ri := NewActiveMQArtemisReconcilerImpl(cr, r)

	valid, retry := ri.validate(cr, k8sClient, *namer)

	assert.False(t, valid)
	assert.False(t, retry)

	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, brokerv1beta1.ValidConditionType))

	condition := meta.FindStatusCondition(cr.Status.Conditions, brokerv1beta1.ValidConditionType)
	assert.Equal(t, condition.Reason, brokerv1beta1.ValidConditionFailedDuplicateBrokerPropertiesKey)
	assert.True(t, strings.Contains(condition.Message, "nameWith"))
}

func TestValidateBrokerPropsDuplicateOnFirstEqualsIncorrectButUnrealisticForOurBrokerConfigUsecase(t *testing.T) {

	cr := &brokerv1beta1.ActiveMQArtemis{
		Spec: brokerv1beta1.ActiveMQArtemisSpec{
			BrokerProperties: []string{
				"nameWith\\=equals_A_not_matched=X",
				"nameWith\\=equals_B_not_matched=Y",
			},
		},
	}

	namer := MakeNamers(cr)

	r := NewActiveMQArtemisReconciler(&NillCluster{}, ctrl.Log, isOpenshift)
	ri := NewActiveMQArtemisReconcilerImpl(cr, r)

	valid, retry := ri.validate(cr, k8sClient, *namer)

	assert.False(t, valid)
	assert.False(t, retry)

	assert.True(t, meta.IsStatusConditionFalse(cr.Status.Conditions, brokerv1beta1.ValidConditionType))

	condition := meta.FindStatusCondition(cr.Status.Conditions, brokerv1beta1.ValidConditionType)
	assert.Equal(t, condition.Reason, brokerv1beta1.ValidConditionFailedDuplicateBrokerPropertiesKey)
	assert.True(t, strings.Contains(condition.Message, "nameWith"))
}
