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
	"errors"
	"fmt"
	"testing"

	broker "github.com/arkmq-org/activemq-artemis-operator/api/v1beta2"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewConditionError(t *testing.T) {
	t.Run("creates error with reason and message", func(t *testing.T) {
		reason := broker.DeployedConditionNoMatchingServiceReason

		err := NewConditionError(reason, "no matching services available")

		assert.NotNil(t, err)
		assert.Equal(t, reason, err.Reason)
		assert.Equal(t, "no matching services available", err.Message)
		assert.Equal(t, "no matching services available", err.Error())
	})

	t.Run("supports formatted message", func(t *testing.T) {
		reason := broker.DeployedConditionNoServiceCapacityReason

		err := NewConditionError(reason, "no service with capacity: required %s, available %s", "1Gi", "512Mi")

		assert.NotNil(t, err)
		assert.Equal(t, reason, err.Reason)
		assert.Equal(t, "no service with capacity: required 1Gi, available 512Mi", err.Message)
	})
}

func TestAsConditionError(t *testing.T) {
	t.Run("returns ConditionError when error is ConditionError", func(t *testing.T) {
		original := NewConditionError(broker.DeployedConditionNoMatchingServiceReason, "test message")

		result, ok := AsConditionError(original)

		assert.True(t, ok)
		assert.NotNil(t, result)
		assert.Equal(t, original.Reason, result.Reason)
		assert.Equal(t, original.Message, result.Message)
	})

	t.Run("returns nil and false for non-ConditionError", func(t *testing.T) {
		normalError := errors.New("regular error")

		result, ok := AsConditionError(normalError)

		assert.False(t, ok)
		assert.Nil(t, result)
	})

	t.Run("returns nil and false for nil error", func(t *testing.T) {
		result, ok := AsConditionError(nil)

		assert.False(t, ok)
		assert.Nil(t, result)
	})

	t.Run("returns nil and false for wrapped error", func(t *testing.T) {
		wrappedError := fmt.Errorf("wrapped: %w", errors.New("inner error"))

		result, ok := AsConditionError(wrappedError)

		assert.False(t, ok)
		assert.Nil(t, result)
	})
}

func TestValidateResourceNameAndSetCondition(t *testing.T) {
	t.Run("sets Valid=False for invalid resource name", func(t *testing.T) {
		conditions := []metav1.Condition{}
		invalidName := "invalid/name"

		err := ValidateResourceNameAndSetCondition(invalidName, &conditions)

		assert.Error(t, err)
		assert.Len(t, conditions, 1)

		validCond := meta.FindStatusCondition(conditions, broker.ValidConditionType)
		assert.NotNil(t, validCond)
		assert.Equal(t, metav1.ConditionFalse, validCond.Status)
		assert.Equal(t, broker.ValidConditionInvalidResourceName, validCond.Reason)
		assert.Contains(t, validCond.Message, "invalid")
	})

	t.Run("does not set condition for valid resource name", func(t *testing.T) {
		conditions := []metav1.Condition{}
		validName := "valid-name-123"

		err := ValidateResourceNameAndSetCondition(validName, &conditions)

		assert.NoError(t, err)
		assert.Len(t, conditions, 0)
	})

	t.Run("updates existing condition", func(t *testing.T) {
		conditions := []metav1.Condition{
			{
				Type:   broker.ValidConditionType,
				Status: metav1.ConditionTrue,
				Reason: broker.ValidConditionSuccessReason,
			},
		}
		invalidName := "invalid/name"

		err := ValidateResourceNameAndSetCondition(invalidName, &conditions)

		assert.Error(t, err)
		assert.Len(t, conditions, 1)

		validCond := meta.FindStatusCondition(conditions, broker.ValidConditionType)
		assert.NotNil(t, validCond)
		assert.Equal(t, metav1.ConditionFalse, validCond.Status)
		assert.Equal(t, broker.ValidConditionInvalidResourceName, validCond.Reason)
	})

	t.Run("validates common invalid patterns", func(t *testing.T) {
		invalidNames := []string{
			"name/with/slashes",
			"name/../with-parent-ref",
			".starts-with-dot",
		}

		for _, name := range invalidNames {
			t.Run(name, func(t *testing.T) {
				conditions := []metav1.Condition{}
				err := ValidateResourceNameAndSetCondition(name, &conditions)
				assert.Error(t, err, "expected error for name: %s", name)

				validCond := meta.FindStatusCondition(conditions, broker.ValidConditionType)
				assert.NotNil(t, validCond)
				assert.Equal(t, metav1.ConditionFalse, validCond.Status)
			})
		}
	})
}
