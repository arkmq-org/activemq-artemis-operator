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
Before you start you need have access to a running Kubernetes cluster environment. A [Minikube](https://minikube.sigs.k8s.io/docs/start/) running on your laptop will just do fine. The arkmq-org operator also runs in a Openshift cluster environment like [CodeReady Container](https://developers.redhat.com/products/codeready-containers/overview). In this blog we assume you have Kubernetes cluster environment. (If you use CodeReady the client tool is **oc** in place of **kubectl**)

### Step 1 - Deploy arkmq-org Operator
In this article we are using the [arkmq-org operator repo](https://github.com/arkmq-org/activemq-artemis-operator). In case you haven't done so, clone it to your local disk:

```shell script
$ git clone https://github.com/arkmq-org/activemq-artemis-operator.git
$ cd activemq-artemis-operator
```

If you are not sure how to deploy the operator take a look at [this tutorial]({{< relref "using_operator.md" >}}).

In this blog post we assume you deployed the operator to a namespace called **myproject**.

After deployment is done, check the operator is up and running. For example run the following command:

```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-m9rfr   1/1     Running   0          2m46s
```

### Step 2 - Deploy ActiveMQ Artemis broker
In this step we'll setup a one-node broker in kubernetes. First we need create a broker custom resource file.

Use your favorite text editor to create a file called **artemis-clustered.yaml** under your repo root directory with the following content:

<a name="broker_clustered_yaml"></a>
```yaml
    apiVersion: broker.amq.io/v2alpha4
    kind: ActiveMQArtemis
    metadata:
      name: ex-aao
    spec:
      deploymentPlan:
        size: 1
        image: quay.io/arkmq-org/activemq-artemis-broker-kubernetes:0.2.1
        persistenceEnabled: true
        messageMigration: true
```

Save it and deploy it using the following command:

```shell script
$ kubectl create -f artemis-clustered.yaml -n myproject
activemqartemis.broker.amq.io/ex-aao created
```

The custom resource file tells the operator to deploy one broker pod from the image **quay.io/arkmq-org/activemq-artemis-broker-kubernetes:0.2.1**,
configured with **persistenceEnabled: true** and **messageMigration: true**.

**persistenceEnabled: true** means the broker persists messages to persistent storage.

**messageMigration: true** means if a broker pod is shut down, its messages will be migrated to another live broker pod so that those messages will be processed.

In a while the broker pod should be up and running:

```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-m9rfr   1/1     Running   0          7m14s
ex-aao-ss-0                                  1/1     Running   0          83s
```

### Step 3 - Scaling up

To inform the operator that we want to scale from one to two broker pods modify artemis-clustered.yaml file to set the **size** to 2

```yaml
    apiVersion: broker.amq.io/v2alpha4
    kind: ActiveMQArtemis
    metadata:
      name: ex-aao
    spec:
      deploymentPlan:
        size: 2
        image: quay.io/arkmq-org/activemq-artemis-broker-kubernetes:0.2.1
        persistenceEnabled: true
        messageMigration: true
```

and apply it:

```shell script
    $ kubectl apply -f artemis-clustered.yaml -n myproject
```

Verify that 2 pods are spawned up in a while:

```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-m9rfr   1/1     Running   0          12m
ex-aao-ss-0                                  1/1     Running   0          6m33s
ex-aao-ss-1                                  1/1     Running   0          75s
```

### Step 4 - Send messages to both brokers

Run the following commands to send 100 messages to each broker:

broker0 at pod ex-aao-ss-0
```shell script
$ kubectl exec ex-aao-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis producer --user admin --password admin --url tcp://ex-aao-ss-0:61616 --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-0:61616
Producer ActiveMQQueue[TEST], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[TEST], thread=0 Produced: 100 messages
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 466 milli seconds
```
broker1 at pod ex-aao-ss-1
```shell script
$ kubectl exec ex-aao-ss-1 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis producer --user admin --password admin --url tcp://ex-aao-ss-1:61616 --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-1:61616
Producer ActiveMQQueue[TEST], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[TEST], thread=0 Produced: 100 messages
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 511 milli seconds
```

The messages are send to a queue **TEST** in each broker. The following commands can show the queue's message count on each broker:

broker0 at pod ex-aao-ss-0
```shell script
$ kubectl exec ex-aao-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://ex-aao-ss-0:61616
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-0:61616
|NAME                     |ADDRESS                  |CONSUMER_COUNT |MESSAGE_COUNT |MESSAGES_ADDED |DELIVERING_COUNT |MESSAGES_ACKED |SCHEDULED_COUNT |ROUTING_TYPE |
|$.artemis.internal.sf.my-cluster.f5e3b3bf-71b2-11eb-ab9a-0242ac110005|$.artemis.internal.sf.my-cluster.f5e3b3bf-71b2-11eb-ab9a-0242ac110005|1              |0             |0              |0                |0              |0               |MULTICAST    |
|DLQ                      |DLQ                      |0              |0             |0              |0                |0              |0               |ANYCAST      |
|ExpiryQueue              |ExpiryQueue              |0              |0             |0              |0                |0              |0               |ANYCAST      |
|TEST                     |TEST                     |0              |100           |100            |0                |0              |0               |ANYCAST      |
|activemq.management.8333dec3-8376-4b50-a9d0-f2197c55e6a8|activemq.management.8333dec3-8376-4b50-a9d0-f2197c55e6a8|1              |0             |0              |0                |0              |0               |MULTICAST    |
|notif.f9cf93f9-71b2-11eb-ab9a-0242ac110005.ActiveMQServerImpl_serverUUID=f5e3b3bf-71b2-11eb-ab9a-0242ac110005|activemq.notifications   |1              |0             |16             |0                |16             |0               |MULTICAST    |
```
broker1 at pod ex-aao-ss-1
```shell script
$ kubectl exec ex-aao-ss-1 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://ex-aao-ss-1:61616
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-1:61616
|NAME                     |ADDRESS                  |CONSUMER_COUNT |MESSAGE_COUNT |MESSAGES_ADDED |DELIVERING_COUNT |MESSAGES_ACKED |SCHEDULED_COUNT |ROUTING_TYPE |
|$.artemis.internal.sf.my-cluster.39079f07-71b2-11eb-88be-0242ac110004|$.artemis.internal.sf.my-cluster.39079f07-71b2-11eb-88be-0242ac110004|1              |0             |0              |0                |0              |0               |MULTICAST    |
|DLQ                      |DLQ                      |0              |0             |0              |0                |0              |0               |ANYCAST      |
|ExpiryQueue              |ExpiryQueue              |0              |0             |0              |0                |0              |0               |ANYCAST      |
|TEST                     |TEST                     |0              |100           |100            |0                |0              |0               |ANYCAST      |
|activemq.management.2a56ebf0-37ad-4cad-9b45-b84cb148f2ff|activemq.management.2a56ebf0-37ad-4cad-9b45-b84cb148f2ff|1              |0             |0              |0                |0              |0               |MULTICAST    |
|notif.fa438b23-71b2-11eb-88be-0242ac110004.ActiveMQServerImpl_serverUUID=39079f07-71b2-11eb-88be-0242ac110004|activemq.notifications   |1              |0             |12             |0                |12             |0               |MULTICAST    |
```

### Step 5 - Scale down with message draining

The operator not only can scale up brokers in a cluster but also can scale them down. Now we have messages on both brokers. If we scale one broker down, will the messsages go offline with it and stuck in its persistence store?

The answer is no. As we set **messageMigration: true** in the [broker cr](#broker_clustered_yaml), the operator will automatically "migrate" the messages on a scaled down broker pod to one of the live broker pod, a process called "message draining".

Internally when the operator detects that a broker pod is down, it starts up a "drainer" broker pod who has access to the down pod's persistence store. It loads the messages in the store and sends those messages to a target live broker pod. When this is done the "drainer" pod shuts down itself.

Now scale down the cluster from 2 pods to one. Edit the [broker cr](#broker_clustered_yaml) file and change the size back to 1:


```yaml
    apiVersion: broker.amq.io/v2alpha4
    kind: ActiveMQArtemis
    metadata:
      name: ex-aao
    spec:
      deploymentPlan:
        size: 1
        image: quay.io/arkmq-org/activemq-artemis-broker-kubernetes:0.2.1
        persistenceEnabled: true
        messageMigration: true
```

Apply it:

```shell script
$ kubectl apply -f artemis-clustered.yaml -n myproject
activemqartemis.broker.amq.io/ex-aao configured
```
In a moment the broker pods will go down to 1. To check, run

```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-m9rfr   1/1     Running   0          26m
ex-aao-ss-0                                  1/1     Running   0          20m
```
Now check the messages in queue TEST at the pod again:
```shell script
$ kubectl exec ex-aao-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://ex-aao-ss-0:61616
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-0:61616
|NAME                     |ADDRESS                  |CONSUMER_COUNT |MESSAGE_COUNT |MESSAGES_ADDED |DELIVERING_COUNT |MESSAGES_ACKED |SCHEDULED_COUNT |ROUTING_TYPE |
|$.artemis.internal.sf.my-cluster.f5e3b3bf-71b2-11eb-ab9a-0242ac110005|$.artemis.internal.sf.my-cluster.f5e3b3bf-71b2-11eb-ab9a-0242ac110005|0              |0             |0              |0                |0              |0               |MULTICAST    |
|DLQ                      |DLQ                      |0              |0             |0              |0                |0              |0               |ANYCAST      |
|ExpiryQueue              |ExpiryQueue              |0              |0             |0              |0                |0              |0               |ANYCAST      |
|TEST                     |TEST                     |0              |200           |200            |0                |0              |0               |ANYCAST      |
|activemq.management.8ee46c3d-c2ba-4e1c-9b6c-f2631e662cb1|activemq.management.8ee46c3d-c2ba-4e1c-9b6c-f2631e662cb1|1              |0             |0              |0                |0              |0               |MULTICAST    |
```
It shows queue TEST's message count is **200** now!

### More information

* Check out [arkmq-org project repo](https://github.com/arkmq-org)
* Reach the [dev team at slack](arkmq-org.slack.com) for questions/issues/help
