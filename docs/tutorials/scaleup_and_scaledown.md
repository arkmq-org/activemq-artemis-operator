---
title: "Scaling Up and Down Brokers with arkmq-org Operator"  
description: "How to use operator to scale up and down broker pods"
draft: false
images: []
menu:
  docs:
    parent: "tutorials"
weight: 110
toc: true
---

With arkmq-org operator one can easily manage the broker clusters.
Either scaling up number of nodes(pods) when workload is high, or scaling down when some is not needed -- without messages being lost or stuck.

### Prerequisite
Before you start you need have access to a running Kubernetes cluster environment. A [Minikube](https://minikube.sigs.k8s.io/docs/start/) running on your laptop will just do fine. The arkmq-org operator also runs in a Openshift cluster environment like [CodeReady Container](https://developers.redhat.com/products/openshift-local/overview). In this blog we assume you have Kubernetes cluster environment. (If you use CodeReady the client tool is **oc** in place of **kubectl**)

### Step 1 - Deploy arkmq-org Operator
In this article we are using the [arkmq-org operator repo](https://github.com/arkmq-org/activemq-artemis-operator). In case you haven't done so, clone it to your local disk:

```shell
git clone https://github.com/arkmq-org/activemq-artemis-operator.git
cd activemq-artemis-operator
```

If you are not sure how to deploy the operator take a look at [this tutorial](using_operator.md).

In this tutorial we assume you deployed the operator to a namespace called **myproject**.

After deployment is done, check the operator is up and running. For example run the following command:

```{"stage":"init","id":"check_operator"}
kubectl get pod -n myproject
```
```shell markdown_runner
NAME                                                   READY   STATUS    RESTARTS   AGE
activemq-artemis-controller-manager-6b78d758bf-hgh6z   1/1     Running   0          106m
```

### Step 2 - Deploy Apache ActiveMQ Artemis broker
In this step we'll setup a one-node broker in kubernetes. First we need create a broker custom resource file.

Use your favorite text editor to create a file called **artemis-clustered.yaml** under your repo root directory with the following content:

<a name="broker_clustered_yaml"></a>
```yaml
    apiVersion: broker.amq.io/v1beta1
    kind: ActiveMQArtemis
    metadata:
      name: ex-aao
    spec:
      deploymentPlan:
        size: 1
        clustered: true
        persistenceEnabled: true
        messageMigration: true
      brokerProperties:
        - "HAPolicyConfiguration=PRIMARY_ONLY"
        - "HAPolicyConfiguration.scaleDownConfiguration.discoveryGroup=my-discovery-group"
        - "HAPolicyConfiguration.scaleDownConfiguration.enabled=false"
```

Or create it using kubectl:

```{"stage":"deploy","id":"create_broker","runtime":"bash"}
kubectl apply -f - -n myproject <<EOF
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemis
metadata:
  name: ex-aao
spec:
  deploymentPlan:
    size: 1
    clustered: true
    persistenceEnabled: true
    messageMigration: true
  brokerProperties:
    - "HAPolicyConfiguration=PRIMARY_ONLY"
    - "HAPolicyConfiguration.scaleDownConfiguration.discoveryGroup=my-discovery-group"
    - "HAPolicyConfiguration.scaleDownConfiguration.enabled=false"
EOF
```
```shell markdown_runner
activemqartemis.broker.amq.io/ex-aao created
```

The custom resource tells the operator to deploy one broker pod with **persistenceEnabled: true** and **messageMigration: true**.

**persistenceEnabled: true** means the broker persists messages to persistent storage.

**messageMigration: true** means if a broker pod is shut down, its messages will be migrated to another live broker pod so that those messages will be processed.

Wait for the broker to be ready:

```{"stage":"deploy","id":"wait_broker"}
kubectl wait ActiveMQArtemis ex-aao --for=condition=Ready --namespace=myproject --timeout=300s
```
```shell markdown_runner
activemqartemis.broker.amq.io/ex-aao condition met
```

Check the pods:

```{"stage":"deploy","id":"check_pods"}
kubectl get pod -n myproject
```
```shell markdown_runner
NAME                                                   READY   STATUS    RESTARTS   AGE
activemq-artemis-controller-manager-6b78d758bf-hgh6z   1/1     Running   0          106m
ex-aao-ss-0                                            1/1     Running   0          14s
```

### Step 3 - Scaling up

To inform the operator that we want to scale from one to two broker pods modify artemis-clustered.yaml file to set the **size** to 2

```yaml
    apiVersion: broker.amq.io/v1beta1
    kind: ActiveMQArtemis
    metadata:
      name: ex-aao
    spec:
      deploymentPlan:
        size: 2
        clustered: true
        persistenceEnabled: true
        messageMigration: true
      brokerProperties:
        - "HAPolicyConfiguration=PRIMARY_ONLY"
        - "HAPolicyConfiguration.scaleDownConfiguration.discoveryGroup=my-discovery-group"
        - "HAPolicyConfiguration.scaleDownConfiguration.enabled=false"
```

and apply it:

```{"stage":"scale_up","id":"scale_to_2"}
kubectl patch activemqartemis ex-aao -n myproject --type='json' -p='[{"op": "replace", "path": "/spec/deploymentPlan/size", "value": 2}]'
```
```shell markdown_runner
activemqartemis.broker.amq.io/ex-aao patched
```

Wait for the scale up to complete:

```{"stage":"scale_up","id":"wait_scale_up"}
kubectl wait ActiveMQArtemis ex-aao --for=condition=Ready --namespace=myproject --timeout=300s
```
```shell markdown_runner
activemqartemis.broker.amq.io/ex-aao condition met
```

Check the pods:

```{"stage":"scale_up","id":"check_pods_scaled"}
kubectl get pod -n myproject
```
```shell markdown_runner
NAME                                                   READY   STATUS     RESTARTS   AGE
activemq-artemis-controller-manager-6b78d758bf-hgh6z   1/1     Running    0          106m
ex-aao-ss-0                                            1/1     Running    0          15s
ex-aao-ss-1                                            0/1     Init:0/1   0          0s
```

### Step 4 - Send messages

Send 100 messages to broker0 (pod `ex-aao-ss-0`):

```{"stage":"send","id":"produce_100_broker0"}
kubectl exec ex-aao-ss-0 -n myproject -c ex-aao-container -- /bin/bash /home/jboss/amq-broker/bin/artemis producer --user admin --password admin --url tcp://ex-aao-ss-0:61616 --message-count 100
```
```shell markdown_runner
Connection brokerURL = tcp://ex-aao-ss-0:61616
Producer ActiveMQQueue[TEST], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[TEST], thread=0 Produced: 100 messages
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 547 milli seconds
NOTE: Picked up JDK_JAVA_OPTIONS: -Dbroker.properties=/amq/extra/secrets/ex-aao-props/,/amq/extra/secrets/ex-aao-props/broker-${STATEFUL_SET_ORDINAL}/,/amq/extra/secrets/ex-aao-props/?filter=.*\.for_ordinal_${STATEFUL_SET_ORDINAL}_only
```

Check the queue's message count on broker0:

```{"stage":"send","id":"queue_stat_broker0"}
kubectl exec ex-aao-ss-0 -n myproject -c ex-aao-container -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://ex-aao-ss-0:61616
```
```shell markdown_runner
Connection brokerURL = tcp://ex-aao-ss-0:61616
|NAME              |ADDRESS           |CONSUMER|MESSAGE|MESSAGES|DELIVERING|MESSAGES|SCHEDULED|ROUTING|INTERNAL|
|                  |                  | COUNT  | COUNT | ADDED  |  COUNT   | ACKED  |  COUNT  | TYPE  |        |
|$sys.mqtt.sessions|$sys.mqtt.sessions|   0    |   0   |   0    |    0     |   0    |    0    |ANYCAST|  true  |
|DLQ               |DLQ               |   0    |   0   |   0    |    0     |   0    |    0    |ANYCAST| false  |
|ExpiryQueue       |ExpiryQueue       |   0    |   0   |   0    |    0     |   0    |    0    |ANYCAST| false  |
|TEST              |TEST              |   0    |  300  |  300   |    0     |   0    |    0    |ANYCAST| false  |
NOTE: Picked up JDK_JAVA_OPTIONS: -Dbroker.properties=/amq/extra/secrets/ex-aao-props/,/amq/extra/secrets/ex-aao-props/broker-${STATEFUL_SET_ORDINAL}/,/amq/extra/secrets/ex-aao-props/?filter=.*\.for_ordinal_${STATEFUL_SET_ORDINAL}_only
```

### Step 5 - Scale down with message draining

The operator not only can scale up brokers in a cluster but also can scale them down. As we set **messageMigration: true** in the [broker cr](#broker_clustered_yaml), the operator will automatically migrate messages when scaling down.

When a broker pod is scaled down, the operator waits for it to forward all messages to live brokers before removing it from the cluster.

Now scale down the cluster from 2 pods to one. Edit the [broker cr](#broker_clustered_yaml) file and change the size back to 1:


```yaml
    apiVersion: broker.amq.io/v1beta1
    kind: ActiveMQArtemis
    metadata:
      name: ex-aao
    spec:
      deploymentPlan:
        size: 1
        clustered: true
        persistenceEnabled: true
        messageMigration: true
      brokerProperties:
        - "HAPolicyConfiguration=PRIMARY_ONLY"
        - "HAPolicyConfiguration.scaleDownConfiguration.discoveryGroup=my-discovery-group"
        - "HAPolicyConfiguration.scaleDownConfiguration.enabled=false"
```

and apply it:

```{"stage":"scale_down","id":"scale_to_1"}
kubectl patch activemqartemis ex-aao -n myproject --type='json' -p='[{"op": "replace", "path": "/spec/deploymentPlan/size", "value": 1}]'
```
```shell markdown_runner
activemqartemis.broker.amq.io/ex-aao patched
```

Wait for the scale down and message migration to complete:

```{"stage":"scale_down","id":"wait_scale_down"}
kubectl wait ActiveMQArtemis ex-aao --for=condition=Ready --namespace=myproject --timeout=300s
```
```shell markdown_runner
activemqartemis.broker.amq.io/ex-aao condition met
```

Check the pods:

```{"stage":"scale_down","id":"check_pods_final"}
kubectl get pod -n myproject
```
```shell markdown_runner
NAME                                                   READY   STATUS    RESTARTS   AGE
activemq-artemis-controller-manager-6b78d758bf-hgh6z   1/1     Running   0          106m
ex-aao-ss-0                                            1/1     Running   0          20s
ex-aao-ss-1                                            0/1     Running   0          5s
```

Now check the messages in queue TEST at the pod:

```{"stage":"scale_down","id":"final_queue_stat"}
kubectl exec ex-aao-ss-0 -n myproject -c ex-aao-container -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://ex-aao-ss-0:61616
```
```shell markdown_runner
Connection brokerURL = tcp://ex-aao-ss-0:61616
|NAME                     |ADDRESS                  |CONSUMER|MESSAGE|MESSAGES|DELIVERING|MESSAGES|SCHEDULED| ROUTING |INTERNAL|
|                         |                         | COUNT  | COUNT | ADDED  |  COUNT   | ACKED  |  COUNT  |  TYPE   |        |
|$.artemis.internal.sf.my-|$.artemis.internal.sf.my-|   1    |   0   |   0    |    0     |   0    |    0    |MULTICAST|  true  |
|  cluster.bc7d2415-fc56-1|  cluster.bc7d2415-fc56-1|        |       |        |          |        |         |         |        |
|  1f0-90d8-1a8a18bf2694  |  1f0-90d8-1a8a18bf2694  |        |       |        |          |        |         |         |        |
|$sys.mqtt.sessions       |$sys.mqtt.sessions       |   0    |   0   |   0    |    0     |   0    |    0    | ANYCAST |  true  |
|DLQ                      |DLQ                      |   0    |   0   |   0    |    0     |   0    |    0    | ANYCAST | false  |
|ExpiryQueue              |ExpiryQueue              |   0    |   0   |   0    |    0     |   0    |    0    | ANYCAST | false  |
|TEST                     |TEST                     |   0    |  300  |  300   |    0     |   0    |    0    | ANYCAST | false  |
|notif.cd313013-fc63-11f0-|activemq.notifications   |   1    |   0   |   10   |    0     |   10   |    0    |MULTICAST|  true  |
|  83f1-fa5c918fa7de.Activ|                         |        |       |        |          |        |         |         |        |
|  eMQServerImpl_name=amq-|                         |        |       |        |          |        |         |         |        |
|  broker                 |                         |        |       |        |          |        |         |         |        |


NOTE: Picked up JDK_JAVA_OPTIONS: -Dbroker.properties=/amq/extra/secrets/ex-aao-props/,/amq/extra/secrets/ex-aao-props/broker-${STATEFUL_SET_ORDINAL}/,/amq/extra/secrets/ex-aao-props/?filter=.*\.for_ordinal_${STATEFUL_SET_ORDINAL}_only
```

### More information

* Check out [arkmq-org project repo](https://github.com/arkmq-org)
