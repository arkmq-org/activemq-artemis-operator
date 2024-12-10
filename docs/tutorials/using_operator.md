---
title: "Using the Arkmq-org Operator"  
description: "Steps to get operator up and running and basic broker operations"
draft: false
images: []
menu:
  docs:
    parent: "tutorials"
weight: 110
toc: true
---

The [arkmq-org](https://github.com/arkmq-org) Operator is a powerful tool that allows you to configure and
manage ActiveMQ Artemis broker resources in a cloud environment. You can get the Operator running in just a few steps.

### Prerequisite
Before you start, you need to have access to a running Kubernetes cluster environment. A [Minikube](https://minikube.sigs.k8s.io/docs/start/)
instance running on your laptop will do fine. The arkmq-org Operator can also run in an Openshift cluster environment such as [CodeReady Containers](https://developers.redhat.com/products/codeready-containers/overview).

In this blog post, we assume that you have a Kubernetes cluster environment.

**_NOTE:_**  If you use CodeReady Containers, the client tool is **oc** rather than **kubectl**

### Step 1 - Preparing for deployment
Clone the arkmq-org Operator repo:
```shell script
      $ git clone https://github.com/arkmq-org/activemq-artemis-operator.git
```
We will use a namespace called **myproject** to deploy the operator and other resources.
If you don't specify a namespace the **default** namespace will be used.

The following command will create the namespace:

```shell script
$ kubectl create namespace myproject
namespace/myproject created
```

Go to the root of the local repo and install the operator:

```shell script      
      $ cd activemq-artemis-operator
      $ ./deploy/install_opr.sh
```

> **_NOTE:_**    You might see some warning messages while deploying the CRDs. For example:
    _"Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use
    apiextensions.k8s.io/v1 CustomResourceDefinition customresourcedefinition.apiextensions.k8s.io/activemqartemises.broker.amq.io created"_.
    You can safely ignore these warnings.

### Step 2 - Deploying the Operator
Deploy the Operator:

```shell script
$ kubectl create -f deploy/operator.yaml --namespace myproject
deployment.apps/activemq-artemis-operator created
```
You might need to wait a few moments for the Operator to fully start. You can verify the Operator status by running the command and looking at the output:
```shell script
$ kubectl get pod --namespace myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-gthwb   1/1     Running   0          12s

```
Make sure that the **STATUS** is **Running**.

By default the operator watches the namespace where it is deployed (i.e. **myproject**) for any custome resources it supports.

### Step 3 - Deploying ActiveMQ Artemis Broker in the cloud
Now, with a running Operator, it's time to deploy the broker via a Custom Resource (CR) instance:
```shell script
kubectl create -f examples/artemis-basic-deployment.yaml -n myproject
```
Watch the broker Pod start up:
```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-gthwb   1/1     Running   0          14m
artemis-broker-ss-0                          1/1     Running   0          61s

```
Behind the scenes, the Operator watches CR deployments in the target namespace. When the broker CR is deployed, the Operator configures and deploys the broker Pod into the cluster.

To see details for startup of the broker Pod you can get the console log from the Pod:
```shell script
$ kubectl logs artemis-broker-ss-0 -n myproject
-XX:+UseParallelOldGC -XX:MinHeapFreeRatio=10 -XX:MaxHeapFreeRatio=20 -XX:GCTimeRatio=4 -XX:AdaptiveSizePolicyWeight=90 -XX:MaxMetaspaceSize=100m -XX:+ExitOnOutOfMemoryError
Removing provided -XX:+UseParallelOldGC in favour of artemis.profile provided option
Running server env: home: /home/jboss AMQ_HOME /opt/amq CONFIG_BROKER false RUN_BROKER
NO RUN_BROKER defined
Using custom configuration. Copy from /amq/init/config to /home/jboss/amq-broker
bin
data
etc
lib
log
tmp
Running Broker in /home/jboss/amq-broker
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
     _        _               _
    / \  ____| |_  ___ __  __(_) _____
   / _ \|  _ \ __|/ _ \  \/  | |/  __/
  / ___ \ | \/ |_/  __/ |\/| | |\___ \
 /_/   \_\|   \__\____|_|  |_|_|/___ /
 Apache ActiveMQ Artemis 2.16.0


2021-02-18 06:00:06,958 INFO  [org.apache.activemq.artemis.integration.bootstrap] AMQ101000: Starting ActiveMQ Artemis Server
2021-02-18 06:00:07,033 INFO  [org.apache.activemq.artemis.core.server] AMQ221000: live Message Broker is starting with configuration Broker Configuration (clustered=true,journalDirectory=data/journal,bindingsDirectory=data/bindings,largeMessagesDirectory=data/large-messages,pagingDirectory=data/paging)
2021-02-18 06:00:07,186 INFO  [org.apache.activemq.artemis.core.server] AMQ221013: Using NIO Journal
2021-02-18 06:00:07,349 INFO  [org.apache.activemq.artemis.core.server] AMQ221057: Global Max Size is being adjusted to 1/2 of the JVM max size (-Xmx). being defined as 1,045,430,272
2021-02-18 06:00:07,781 WARNING [org.jgroups.stack.Configurator] JGRP000014: BasicTCP.use_send_queues has been deprecated: will be removed in 4.0
2021-02-18 06:00:07,814 WARNING [org.jgroups.stack.Configurator] JGRP000014: Discovery.timeout has been deprecated: GMS.join_timeout should be used instead
2021-02-18 06:00:07,943 INFO  [org.jgroups.protocols.openshift.DNS_PING] serviceName [artemis-broker-ping-svc] set; clustering enabled
2021-02-18 06:00:11,047 INFO  [org.openshift.ping.common.Utils] 3 attempt(s) with a 1000ms sleep to execute [GetServicePort] failed. Last failure was [javax.naming.NameNotFoundException: DNS name not found [response code 3]]
2021-02-18 06:00:11,048 WARNING [org.jgroups.protocols.openshift.DNS_PING] No DNS SRV record found for service [artemis-broker-ping-svc]

-------------------------------------------------------------------
GMS: address=artemis-broker-ss-0-6103, cluster=activemq_broadcast_channel, physical address=172.17.0.4:7800
-------------------------------------------------------------------
2021-02-18 06:00:14,217 INFO  [org.apache.activemq.artemis.core.server] AMQ221043: Protocol module found: [artemis-server]. Adding protocol support for: CORE
2021-02-18 06:00:14,219 INFO  [org.apache.activemq.artemis.core.server] AMQ221043: Protocol module found: [artemis-amqp-protocol]. Adding protocol support for: AMQP
2021-02-18 06:00:14,221 INFO  [org.apache.activemq.artemis.core.server] AMQ221043: Protocol module found: [artemis-hornetq-protocol]. Adding protocol support for: HORNETQ
2021-02-18 06:00:14,226 INFO  [org.apache.activemq.artemis.core.server] AMQ221043: Protocol module found: [artemis-mqtt-protocol]. Adding protocol support for: MQTT
2021-02-18 06:00:14,227 INFO  [org.apache.activemq.artemis.core.server] AMQ221043: Protocol module found: [artemis-openwire-protocol]. Adding protocol support for: OPENWIRE
2021-02-18 06:00:14,228 INFO  [org.apache.activemq.artemis.core.server] AMQ221043: Protocol module found: [artemis-stomp-protocol]. Adding protocol support for: STOMP
2021-02-18 06:00:14,327 INFO  [org.apache.activemq.artemis.core.server] AMQ221034: Waiting indefinitely to obtain live lock
2021-02-18 06:00:14,328 INFO  [org.apache.activemq.artemis.core.server] AMQ221035: Live Server Obtained live lock
2021-02-18 06:00:14,510 INFO  [org.apache.activemq.artemis.core.server] AMQ221080: Deploying address DLQ supporting [ANYCAST]
2021-02-18 06:00:14,539 INFO  [org.apache.activemq.artemis.core.server] AMQ221003: Deploying ANYCAST queue DLQ on address DLQ
2021-02-18 06:00:14,658 INFO  [org.apache.activemq.artemis.core.server] AMQ221080: Deploying address ExpiryQueue supporting [ANYCAST]
2021-02-18 06:00:14,660 INFO  [org.apache.activemq.artemis.core.server] AMQ221003: Deploying ANYCAST queue ExpiryQueue on address ExpiryQueue
2021-02-18 06:00:15,019 INFO  [org.apache.activemq.artemis.core.server] AMQ221020: Started EPOLL Acceptor at artemis-broker-ss-0.artemis-broker-hdls-svc.myproject.svc.cluster.local:61616 for protocols [CORE]
2021-02-18 06:00:15,022 INFO  [org.apache.activemq.artemis.core.server] AMQ221007: Server is now live
2021-02-18 06:00:15,023 INFO  [org.apache.activemq.artemis.core.server] AMQ221001: Apache ActiveMQ Artemis Message Broker version 2.16.0 [amq-broker, nodeID=8d13dfd2-71ae-11eb-b7ec-0242ac110004]
2021-02-18 06:00:15,981 INFO  [org.apache.activemq.hawtio.branding.PluginContextListener] Initialized activemq-branding plugin
2021-02-18 06:00:16,075 INFO  [org.apache.activemq.hawtio.plugin.PluginContextListener] Initialized artemis-plugin plugin
2021-02-18 06:00:16,724 INFO  [io.hawt.HawtioContextListener] Initialising hawtio services
2021-02-18 06:00:16,763 INFO  [io.hawt.system.ConfigManager] Configuration will be discovered via system properties
2021-02-18 06:00:16,782 INFO  [io.hawt.jmx.JmxTreeWatcher] Welcome to Hawtio 2.11.0
2021-02-18 06:00:16,816 INFO  [io.hawt.web.auth.AuthenticationConfiguration] Starting hawtio authentication filter, JAAS realm: "activemq" authorized role(s): "admin" role principal classes: "org.apache.activemq.artemis.spi.core.security.jaas.RolePrincipal"
2021-02-18 06:00:16,855 INFO  [io.hawt.web.proxy.ProxyServlet] Proxy servlet is disabled
2021-02-18 06:00:16,866 INFO  [io.hawt.web.servlets.JolokiaConfiguredAgentServlet] Jolokia overridden property: [key=policyLocation, value=file:/home/jboss/amq-broker/etc/jolokia-access.xml]
2021-02-18 06:00:17,038 INFO  [org.apache.activemq.artemis] AMQ241001: HTTP Server started at http://artemis-broker-ss-0.artemis-broker-hdls-svc.myproject.svc.cluster.local:8161
2021-02-18 06:00:17,038 INFO  [org.apache.activemq.artemis] AMQ241002: Artemis Jolokia REST API available at http://artemis-broker-ss-0.artemis-broker-hdls-svc.myproject.svc.cluster.local:8161/console/jolokia
2021-02-18 06:00:17,039 INFO  [org.apache.activemq.artemis] AMQ241004: Artemis Console available at http://artemis-broker-ss-0.artemis-broker-hdls-svc.myproject.svc.cluster.local:8161/console
```
### Step 4 - Create a queue using the Operator
Now, let's create a message queue in the broker:
```shell script
$ kubectl create -f examples/address/address_queue.yaml -n myproject
activemqartemisaddress.broker.amq.io/artemis-address-queue created
```
The _address_queue.yaml_ is another CR supported by the arkmq-org Operator. Its content is shown below:
```yaml
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemisAddress
metadata:
  name: artemis-address-queue
spec:
  addressName: myAddress0
  queueName: myQueue0
  routingType: anycast
  removeFromBrokerOnDelete: true
```
The CR tells the Operator to create a queue named **myQueue0** on address **myAddress0** on each broker that it manages.

After the CR is deployed, you can observe the queue on the broker:

<a name="queuestat"></a>
```shell script
$ kubectl exec artemis-address-queue-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://artemis-address-queue-ss-0:61616
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://artemis-address-queue-ss-0:61616
|NAME                     |ADDRESS                  |CONSUMER_COUNT |MESSAGE_COUNT |MESSAGES_ADDED |DELIVERING_COUNT |MESSAGES_ACKED |SCHEDULED_COUNT |ROUTING_TYPE |
|DLQ                      |DLQ                      |0              |0             |0              |0                |0              |0               |ANYCAST      |
|ExpiryQueue              |ExpiryQueue              |0              |0             |0              |0                |0              |0               |ANYCAST      |
|activemq.management.a5037b80-fbb7-48b8-93d2-f505d7a39aae|activemq.management.a5037b80-fbb7-48b8-93d2-f505d7a39aae|1              |0             |0              |0                |0              |0               |MULTICAST    |
|myQueue0                 |myAddress0               |0              |0             |0              |0                |0              |0               |ANYCAST      |
```

### Step 5 - Sending and Receiving messages
Finally, you can send some messages to the broker and receive them. Here, we use the `artemis` CLI tool that comes with the deployed broker instance for the test.

Send 100 messages:
```shell script
$ kubectl exec artemis-address-queue-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis producer --user admin --password admin --url tcp://artemis-address-queue-ss-0:61616 --destination myQueue0::myAddress0 --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://artemis-address-queue-ss-0:61616
Producer ActiveMQQueue[myQueue0::myAddress0], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[myQueue0::myAddress0], thread=0 Produced: 100 messages
Producer ActiveMQQueue[myQueue0::myAddress0], thread=0 Elapsed time in second : 0 s
Producer ActiveMQQueue[myQueue0::myAddress0], thread=0 Elapsed time in milli second : 505 milli seconds
````
Now, if you check the queue statistics using the command mentioned in [Step 4](#queuestat), you see that the message count is 100:
```shell script
$ kubectl exec artemis-address-queue-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://artemis-address-queue-ss-0:61616
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://artemis-address-queue-ss-0:61616
|NAME                     |ADDRESS                  |CONSUMER_COUNT |MESSAGE_COUNT |MESSAGES_ADDED |DELIVERING_COUNT |MESSAGES_ACKED |SCHEDULED_COUNT |ROUTING_TYPE |
|DLQ                      |DLQ                      |0              |0             |0              |0                |0              |0               |ANYCAST      |
|ExpiryQueue              |ExpiryQueue              |0              |0             |0              |0                |0              |0               |ANYCAST      |
|activemq.management.67ad22bd-f726-4b56-8949-fa37e214842c|activemq.management.67ad22bd-f726-4b56-8949-fa37e214842c|1              |0             |0              |0                |0              |0               |MULTICAST    |
|myAddress0               |myQueue0                 |0              |100           |100            |0                |0              |0               |ANYCAST      |
|myQueue0                 |myAddress0               |0              |0             |0              |0                |0              |0               |ANYCAST      |
```
### Further information

* [arkmq-org Github Repo](https://github.com/arkmq-org)
