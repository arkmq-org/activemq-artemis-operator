---
title: "Exchanging messages using port forwarding"  
description: "Steps to get a producer and a consummer exchanging messages over a deployed broker on OpenShift using port forwarding"
draft: false
images: []
menu:
  docs:
    parent: "send-receive"
weight: 110
toc: true
---

### Prerequisite

Before you start, you need to have access to a running Kubernetes cluster
environment. A [Minikube](https://minikube.sigs.k8s.io/docs/start/) instance
running on your laptop will do fine.

#### Start minikube with a parametrized dns domain name

```{"stage":"init", "id":"minikube_start"}
minikube start --profile tutorialtester
minikube profile tutorialtester
```
```shell markdown_runner
* [tutorialtester] minikube v1.36.0 on Fedora 41
  - MINIKUBE_ROOTLESS=true
* Automatically selected the kvm2 driver. Other choices: podman, qemu2, ssh
* Starting "tutorialtester" primary control-plane node in "tutorialtester" cluster
* Creating kvm2 VM (CPUs=2, Memory=6000MB, Disk=20000MB) ...
* Preparing Kubernetes v1.33.1 on Docker 28.0.4 ...
  - Generating certificates and keys ...
  - Booting up control plane ...
  - Configuring RBAC rules ...
* Configuring bridge CNI (Container Networking Interface) ...
* Verifying Kubernetes components...
  - Using image gcr.io/k8s-minikube/storage-provisioner:v5
* Enabled addons: default-storageclass, storage-provisioner
* Done! kubectl is now configured to use "tutorialtester" cluster and "default" namespace by default
! Image was not built for the current minikube version. To resolve this you can delete and recreate your minikube cluster using the latest images. Expected minikube version: v1.35.0 -> Actual minikube version: v1.36.0
* minikube profile was successfully set to tutorialtester
```

### Deploy the operator

#### create the namespace

```{"stage":"init"}
kubectl create namespace send-receive-project
kubectl config set-context --current --namespace=send-receive-project
```
```shell markdown_runner
namespace/send-receive-project created
Context "tutorialtester" modified.
```

Go to the root of the operator repo and install it:

```{"stage":"init", "rootdir":"$initial_dir"}
./deploy/install_opr.sh
```
```shell markdown_runner
Deploying operator to watch single namespace
Client Version: 4.18.5
Kustomize Version: v5.4.2
Kubernetes Version: v1.33.1
customresourcedefinition.apiextensions.k8s.io/activemqartemises.broker.amq.io created
customresourcedefinition.apiextensions.k8s.io/activemqartemisaddresses.broker.amq.io created
customresourcedefinition.apiextensions.k8s.io/activemqartemisscaledowns.broker.amq.io created
customresourcedefinition.apiextensions.k8s.io/activemqartemissecurities.broker.amq.io created
serviceaccount/activemq-artemis-controller-manager created
role.rbac.authorization.k8s.io/activemq-artemis-operator-role created
rolebinding.rbac.authorization.k8s.io/activemq-artemis-operator-rolebinding created
role.rbac.authorization.k8s.io/activemq-artemis-leader-election-role created
rolebinding.rbac.authorization.k8s.io/activemq-artemis-leader-election-rolebinding created
deployment.apps/activemq-artemis-controller-manager created
```

Wait for the Operator to start (status: `running`).

```bash {"stage":"init", "runtime":"bash", "label":"wait for the operator to be running"}
kubectl wait pod --all --for=condition=Ready --namespace=send-receive-project --timeout=600s
```
```shell markdown_runner
pod/activemq-artemis-controller-manager-fdd64476f-kclrm condition met
```

### Deploying the Apache ActiveMQ Artemis Broker

For this tutorial we need to:

* have a broker that is able to listen to any network interface. For that we
  setup an `acceptor` that will be listening on every interfaces on port
  `62626`.
* have queues to exchange messages on. These are configured by the broker
  properties. Two queues are setup, one called `APP.JOBS` that is of type
  `ANYCAST` and one called `APP.COMMANDS` that is of type `MULTICAST`.

```bash {"stage":"deploy", "runtime":"bash", "label":"deploy the broker"}
kubectl apply -f - <<EOF
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemis
metadata:
  name: send-receive
  namespace: send-receive-project
spec:
  acceptors:
    - bindToAllInterfaces: true
      name: acceptall
      port: 62626
  brokerProperties:
    - addressConfigurations."APP.JOBS".routingTypes=ANYCAST
    - addressConfigurations."APP.JOBS".queueConfigs."APP.JOBS".routingType=ANYCAST
    - addressConfigurations."APP.COMMANDS".routingTypes=MULTICAST
EOF
```
```shell markdown_runner
activemqartemis.broker.amq.io/send-receive created
```

Wait for the Broker to be ready:

```{"stage":"deploy"}
kubectl wait ActiveMQArtemis send-receive --for=condition=Ready --namespace=send-receive-project --timeout=240s
```
```shell markdown_runner
activemqartemis.broker.amq.io/send-receive condition met
```


### Exchanging messages between a producer and a consumer

#### Get the actual broker version

```{"stage":"test_setup", "runtime":"bash", "label":"get latest broker version"}
export BROKER_VERSION=$(kubectl get ActiveMQArtemis send-receive --namespace=send-receive-project -o json | jq .status.version.brokerVersion -r)
echo broker version: $BROKER_VERSION
```
```shell markdown_runner
broker version: 2.42.0
```

#### Download the broker

```{"stage":"test_setup", "rootdir":"$tmpdir.1", "runtime":"bash", "label":"download artemis"}
wget --quiet https://repo1.maven.org/maven2/org/apache/activemq/apache-artemis/${BROKER_VERSION}/apache-artemis-${BROKER_VERSION}-bin.tar.gz
tar -zxf apache-artemis-${BROKER_VERSION}-bin.tar.gz apache-artemis-${BROKER_VERSION}/
# make the rest of commands version agnostic
mv apache-artemis-${BROKER_VERSION}/ apache-artemis/
```

#### Produce and consume messages

```bash {"stage":"test", "rootdir":"$tmpdir.1/apache-artemis/bin/", "parallel":true, "runtime":"bash", "label":"anycast: produce and consume 1000 messages"}
# First we need to start port forwarding
kubectl port-forward send-receive-ss-0 62626 -n send-receive-project &

# Then produce and consume some messages
./artemis producer --destination APP.JOBS  --url tcp://localhost:62626
./artemis consumer --destination APP.JOBS  --url tcp://localhost:62626

# Finally we need to kill the port forwarding
pkill kubectl -9
```
```shell markdown_runner
Forwarding from 127.0.0.1:62626 -> 62626
Forwarding from [::1]:62626 -> 62626
Connection brokerURL = tcp://localhost:62626
Handling connection for 62626
Handling connection for 62626
Producer ActiveMQQueue[APP.JOBS], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[APP.JOBS], thread=0 Produced: 1000 messages
Producer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in second : 3 s
Producer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in milli second : 3634 milli seconds
Connection brokerURL = tcp://localhost:62626
Consumer:: filter = null
Handling connection for 62626
Handling connection for 62626
Consumer ActiveMQQueue[APP.JOBS], thread=0 wait 3000ms until 1000 messages are consumed
Received 1000
Consumer ActiveMQQueue[APP.JOBS], thread=0 Consumed: 1000 messages
Consumer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in second : 0 s
Consumer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in milli second : 69 milli seconds
Consumer ActiveMQQueue[APP.JOBS], thread=0 Consumed: 1000 messages
Consumer ActiveMQQueue[APP.JOBS], thread=0 Consumer thread finished
NOTE: Picked up JDK_JAVA_OPTIONS: -Djavax.net.ssl.trustStore=/etc/pki/java/cacerts
NOTE: Picked up JDK_JAVA_OPTIONS: -Djavax.net.ssl.trustStore=/etc/pki/java/cacerts
```

### cleanup

To leave a pristine environment after executing this tutorial you can simply,
delete the minikube cluster.

```{"stage":"teardown", "requires":"init/minikube_start"}
minikube delete --profile tutorialtester
```
```shell markdown_runner
* Deleting "tutorialtester" in kvm2 ...
* Removed all traces of the "tutorialtester" cluster.
```
