package controllers

import (
	"encoding/json"

	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	v1beta2 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConvertArtemisToBroker converts a v1beta1.ActiveMQArtemis to a v1beta2.Broker
// so the reconciler can work with Broker as the canonical type.
// The ActiveMQArtemis GVK is preserved on the converted object so that owner
// references on child resources (StatefulSets, Services, etc.) correctly
// point back to the ActiveMQArtemis CR.
func ConvertArtemisToBroker(artemis *brokerv1beta1.ActiveMQArtemis) (*v1beta2.Broker, error) {
	data, err := json.Marshal(artemis)
	if err != nil {
		return nil, err
	}
	broker := &v1beta2.Broker{}
	if err = json.Unmarshal(data, broker); err != nil {
		return nil, err
	}
	broker.TypeMeta = metav1.TypeMeta{
		APIVersion: brokerv1beta1.GroupVersion.String(),
		Kind:       "ActiveMQArtemis",
	}
	return broker, nil
}

// ConvertBrokerStatusToArtemis copies the status from a reconciled
// Broker object back to the original ActiveMQArtemis CR so that
// Kubernetes can persist the updated status for the ActiveMQArtemis resource.
func ConvertBrokerStatusToArtemis(broker *v1beta2.Broker, artemis *brokerv1beta1.ActiveMQArtemis) error {
	data, err := json.Marshal(broker.Status)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &artemis.Status)
}
