package controllers

import (
	"encoding/json"

	v1 "github.com/arkmq-org/activemq-artemis-operator/api/v1"
	brokerv1beta1 "github.com/arkmq-org/activemq-artemis-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConvertBrokerToArtemis converts a v1.Broker to a v1beta1.ActiveMQArtemis
// so the existing reconciler logic can be reused without duplication.
// The Broker's GVK is preserved on the converted object so that owner
// references on child resources (StatefulSets, Services, etc.) correctly
// point back to the Broker CR.
func ConvertBrokerToArtemis(broker *v1.Broker) (*brokerv1beta1.ActiveMQArtemis, error) {
	data, err := json.Marshal(broker)
	if err != nil {
		return nil, err
	}
	artemis := &brokerv1beta1.ActiveMQArtemis{}
	if err = json.Unmarshal(data, artemis); err != nil {
		return nil, err
	}
	artemis.TypeMeta = metav1.TypeMeta{
		APIVersion: v1.GroupVersion.String(),
		Kind:       "Broker",
	}
	return artemis, nil
}

// ConvertArtemisStatusToBroker copies the status from a reconciled
// ActiveMQArtemis object back to the original Broker CR so that
// Kubernetes can persist the updated status for the Broker resource.
func ConvertArtemisStatusToBroker(artemis *brokerv1beta1.ActiveMQArtemis, broker *v1.Broker) error {
	data, err := json.Marshal(artemis.Status)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &broker.Status)
}
