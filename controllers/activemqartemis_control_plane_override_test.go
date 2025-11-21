package controllers

import (
	"context"
	"testing"

	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// TestApplyControlPlaneOverrides_CRSpecificSecret tests that CR-specific override secret is used when it exists
func TestApplyControlPlaneOverrides_CRSpecificSecret(t *testing.T) {
	// Setup
	crName := "test-broker"
	namespace := "test-ns"

	customResource := &brokerv1beta1.ActiveMQArtemis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
	}

	// CR-specific override secret
	crSpecificSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName + "-control-plane-override",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"_cert-users": []byte("operator=/.*operator.*/\nnew-user=/.*new-user.*/\n"),
			"_cert-roles": []byte("metrics=operator,new-user\n"),
		},
	}

	// Create fake client with the secret
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = brokerv1beta1.AddToScheme(scheme)
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crSpecificSecret).Build()

	// Initial broker properties
	brokerPropertiesMapData := map[string]string{
		"_cert-users":  "operator=/.*operator.*/\n",
		"_cert-roles":  "metrics=operator\n",
		"other-config": "value",
	}

	// Execute
	err := applyControlPlaneOverrides(customResource, fakeClient, brokerPropertiesMapData)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "operator=/.*operator.*/\nnew-user=/.*new-user.*/\n", brokerPropertiesMapData["_cert-users"], "CR-specific secret should override _cert-users")
	assert.Equal(t, "metrics=operator,new-user\n", brokerPropertiesMapData["_cert-roles"], "CR-specific secret should override _cert-roles")
	assert.Equal(t, "value", brokerPropertiesMapData["other-config"], "Unspecified keys should remain unchanged")
}

// TestApplyControlPlaneOverrides_SharedSecretFallback tests that shared override secret is used when CR-specific doesn't exist
func TestApplyControlPlaneOverrides_SharedSecretFallback(t *testing.T) {
	// Setup
	crName := "test-broker"
	namespace := "test-ns"

	customResource := &brokerv1beta1.ActiveMQArtemis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
	}

	// Shared override secret (no CR-specific secret exists)
	sharedSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "control-plane-override",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"_cert-users": []byte("operator=/.*operator.*/\nshared-user=/.*shared-user.*/\n"),
		},
	}

	// Create fake client with only the shared secret
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = brokerv1beta1.AddToScheme(scheme)
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(sharedSecret).Build()

	// Initial broker properties
	brokerPropertiesMapData := map[string]string{
		"_cert-users": "operator=/.*operator.*/\n",
	}

	// Execute
	err := applyControlPlaneOverrides(customResource, fakeClient, brokerPropertiesMapData)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "operator=/.*operator.*/\nshared-user=/.*shared-user.*/\n", brokerPropertiesMapData["_cert-users"], "Shared secret should be used as fallback")
}

// TestApplyControlPlaneOverrides_NoSecretFound tests that no error occurs when neither secret exists
func TestApplyControlPlaneOverrides_NoSecretFound(t *testing.T) {
	// Setup
	crName := "test-broker"
	namespace := "test-ns"

	customResource := &brokerv1beta1.ActiveMQArtemis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
	}

	// Create fake client with no secrets
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = brokerv1beta1.AddToScheme(scheme)
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	// Initial broker properties
	brokerPropertiesMapData := map[string]string{
		"_cert-users": "operator=/.*operator.*/\n",
		"_cert-roles": "metrics=operator\n",
	}

	originalUsers := brokerPropertiesMapData["_cert-users"]
	originalRoles := brokerPropertiesMapData["_cert-roles"]

	// Execute
	err := applyControlPlaneOverrides(customResource, fakeClient, brokerPropertiesMapData)

	// Verify
	assert.NoError(t, err, "No error should occur when no override secret exists")
	assert.Equal(t, originalUsers, brokerPropertiesMapData["_cert-users"], "Original values should be unchanged")
	assert.Equal(t, originalRoles, brokerPropertiesMapData["_cert-roles"], "Original values should be unchanged")
}

// TestApplyControlPlaneOverrides_CRSpecificTakesPrecedence tests that CR-specific secret takes precedence over shared
func TestApplyControlPlaneOverrides_CRSpecificTakesPrecedence(t *testing.T) {
	// Setup
	crName := "test-broker"
	namespace := "test-ns"

	customResource := &brokerv1beta1.ActiveMQArtemis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
	}

	// Both CR-specific and shared secrets exist
	crSpecificSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName + "-control-plane-override",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"_cert-users": []byte("from-cr-specific\n"),
		},
	}

	sharedSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "control-plane-override",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"_cert-users": []byte("from-shared\n"),
		},
	}

	// Create fake client with both secrets
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = brokerv1beta1.AddToScheme(scheme)
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crSpecificSecret, sharedSecret).Build()

	// Initial broker properties
	brokerPropertiesMapData := map[string]string{
		"_cert-users": "original\n",
	}

	// Execute
	err := applyControlPlaneOverrides(customResource, fakeClient, brokerPropertiesMapData)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "from-cr-specific\n", brokerPropertiesMapData["_cert-users"], "CR-specific secret should take precedence")
}

// TestApplyControlPlaneOverrides_CompleteReplacement tests that override completely replaces the key value
func TestApplyControlPlaneOverrides_CompleteReplacement(t *testing.T) {
	// Setup
	crName := "test-broker"
	namespace := "test-ns"

	customResource := &brokerv1beta1.ActiveMQArtemis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
	}

	// Override secret with completely new content
	overrideSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName + "-control-plane-override",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"_cert-users": []byte("completely-new-content\n"),
		},
	}

	// Create fake client
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = brokerv1beta1.AddToScheme(scheme)
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(overrideSecret).Build()

	// Initial broker properties with existing content
	brokerPropertiesMapData := map[string]string{
		"_cert-users": "operator=/.*operator.*/\nprobe=/.*probe.*/\n",
	}

	// Execute
	err := applyControlPlaneOverrides(customResource, fakeClient, brokerPropertiesMapData)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "completely-new-content\n", brokerPropertiesMapData["_cert-users"], "Override should completely replace original value")
	assert.NotContains(t, brokerPropertiesMapData["_cert-users"], "operator", "Original content should be completely replaced")
	assert.NotContains(t, brokerPropertiesMapData["_cert-users"], "probe", "Original content should be completely replaced")
}

// TestApplyControlPlaneOverrides_GetError tests error handling when getting secret fails for non-NotFound errors
func TestApplyControlPlaneOverrides_GetError(t *testing.T) {
	// Setup
	crName := "test-broker"
	namespace := "test-ns"

	customResource := &brokerv1beta1.ActiveMQArtemis{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: namespace,
		},
	}

	// Create a fake client that will return an error
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = brokerv1beta1.AddToScheme(scheme)
	fakeClient := &errorClient{
		Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
	}

	brokerPropertiesMapData := map[string]string{
		"_cert-users": "original\n",
	}

	// Execute
	err := applyControlPlaneOverrides(customResource, fakeClient, brokerPropertiesMapData)

	// Verify
	assert.Error(t, err, "Should return error when Get fails with non-NotFound error")
}

// errorClient is a fake client that returns errors for Get operations
type errorClient struct {
	client.Client
}

func (e *errorClient) Get(ctx context.Context, key types.NamespacedName, obj client.Object, opts ...client.GetOption) error {
	return assert.AnError
}
