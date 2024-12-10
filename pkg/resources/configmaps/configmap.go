package configmaps

import (
	"github.com/arkmq-org/activemq-artemis-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RetrieveConfigMap(namespace string, configMapName string, client client.Client) (*corev1.ConfigMap, error) {
	data := make(map[string]string)
	namespacedName := types.NamespacedName{
		Name:      configMapName,
		Namespace: namespace,
	}
	configMapDefinition := MakeConfigMap(namespace, configMapName, data)
	if err := resources.Retrieve(namespacedName, client, configMapDefinition); err != nil {
		return nil, err
	}
	return configMapDefinition, nil
}

func MakeConfigMap(namespace string, configMapName string, stringData map[string]string) *corev1.ConfigMap {

	configMapDefinition := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: stringData,
	}

	return &configMapDefinition
}
