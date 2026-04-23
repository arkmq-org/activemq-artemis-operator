/*
Copyright 2026.

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
	"fmt"

	broker "github.com/arkmq-org/arkmq-org-broker-operator/api/v1beta2"
	"github.com/arkmq-org/arkmq-org-broker-operator/pkg/utils/common"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionError wraps an error with a condition reason string
// This allows errors to carry structured metadata about which condition
// reason should be set when the error is processed
type ConditionError struct {
	Reason  string
	Message string
}

func (e *ConditionError) Error() string {
	return e.Message
}

// NewConditionError creates a new ConditionError with the given reason and message
func NewConditionError(reason string, format string, args ...interface{}) *ConditionError {
	return &ConditionError{
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// AsConditionError attempts to extract a ConditionError from an error
// Returns the ConditionError and true if the error is a ConditionError, nil and false otherwise
func AsConditionError(err error) (*ConditionError, bool) {
	if err == nil {
		return nil, false
	}
	if ce, ok := err.(*ConditionError); ok {
		return ce, true
	}
	return nil, false
}

// ValidateResourceNameAndSetCondition validates a resource name and sets the Valid condition accordingly
// Returns an error if the name is invalid, nil otherwise
func ValidateResourceNameAndSetCondition(resourceName string, conditions *[]metav1.Condition) error {
	err := common.ValidateResourceName(resourceName)
	if err != nil {
		meta.SetStatusCondition(conditions, metav1.Condition{
			Type:    broker.ValidConditionType,
			Status:  metav1.ConditionFalse,
			Reason:  broker.ValidConditionInvalidResourceName,
			Message: err.Error(),
		})
	}
	return err
}
