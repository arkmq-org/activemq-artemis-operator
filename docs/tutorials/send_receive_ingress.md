---
title: "Exchanging messages over an ssl ingress"  
description: "Steps to get a producer and a consummer exchanging messages over a deployed broker on kubernetes using an ingress"
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
```shell tutorial_tester
* [tutorialtester] minikube v1.32.0 on Fedora 40
* Automatically selected the docker driver. Other choices: kvm2, qemu2, ssh
* Using Docker driver with root privileges
* Starting control plane node tutorialtester in cluster tutorialtester
* Pulling base image ...
* Creating docker container (CPUs=2, Memory=15900MB) ...
* Preparing Kubernetes v1.28.3 on Docker 24.0.7 ...
  - Generating certificates and keys ...
  - Booting up control plane ...
  - Configuring RBAC rules ...
* Configuring bridge CNI (Container Networking Interface) ...
  - Using image gcr.io/k8s-minikube/storage-provisioner:v5
* Verifying Kubernetes components...
* Enabled addons: storage-provisioner, default-storageclass
* Done! kubectl is now configured to use "tutorialtester" cluster and "default" namespace by default
* minikube profile was successfully set to tutorialtester
```

#### Enable nginx and ssl passthrough for minikube

```{"stage":"init"}
minikube addons enable ingress
minikube kubectl -- patch deployment -n ingress-nginx ingress-nginx-controller --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'
```
```shell tutorial_tester
* ingress is an addon maintained by Kubernetes. For any concerns contact minikube on GitHub.
You can view the list of minikube maintainers at: https://github.com/kubernetes/minikube/blob/master/OWNERS
  - Using image registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20231011-8b53cabe0
  - Using image registry.k8s.io/ingress-nginx/controller:v1.9.4
  - Using image registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20231011-8b53cabe0
* Verifying ingress addon...
* The 'ingress' addon is enabled
deployment.apps/ingress-nginx-controller patched
```

#### Get minikube's ip

```{"stage":"init", "runtime":"bash", "label":"get the cluster ip"}
export CLUSTER_IP=$(minikube ip --profile tutorialtester)
```

#### Make sure the domain of your cluster is resolvable

If you are running your OpenShift cluster locally, you might not be able to
resolve the urls to IPs out of the blue. Follow [this guide]({{< ref
"/docs/help/hostname_resolution" >}} "set up dnsmasq") to configure your setup.

This tutorial will follow the simple /etc/hosts approach, but feel free to use
the most appropriate one for you.

### Deploy the operator

#### create the namespace

```{"stage":"init"}
kubectl create namespace send-receive-project
kubectl config set-context --current --namespace=send-receive-project
```
```shell tutorial_tester
namespace/send-receive-project created
Context "tutorialtester" modified.
```

Go to the root of the operator repo and install it:

```{"stage":"init", "rootdir":"$operator"}
./deploy/install_opr.sh
```
```shell tutorial_tester
Deploying operator to watch single namespace
Client Version: 4.15.0-0.okd-2024-01-27-070424
Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
Kubernetes Version: v1.28.3
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

```{"stage":"init", "runtime":"bash", "label":"wait for the operator to be running"}
kubectl wait pod --all --for=condition=Ready --namespace=send-receive-project --timeout=600s
```
```shell tutorial_tester
pod/activemq-artemis-controller-manager-55b8c479df-rpzg4 condition met
```
### Deploy the ActiveMQ Artemis Broker

For this tutorial we need to:

* have a broker that is able to listen to any network interface. For that we
  setup an `acceptor` that will be listening on every interfaces on port
  `62626`.
* have the ssl protocol configured for the `acceptor`
* have queues to exchange messages on. These are configured by the broker
  properties. Two queues are setup, one called `APP.JOBS` that is of type
  `ANYCAST` and one called `APP.COMMANDS` that is of type `MULTICAST`.

#### Create the certs

We'll take some inspiration from the [ssl broker
setup](https://github.com/arkmq-org/send-receive-project/blob/main/docs/tutorials/ssl_broker_setup.md)
to configure the certificates.

> [!NOTE]
> In this tutorial:
> * The password used for the certificates is `000000`.
> * The secret name is `send-receive-sslacceptor-secret` composed from the broker
>   name `send-receive` and the acceptor name `sselacceptor`

```{"stage":"etc", "runtime":"bash", "label":"get the cluster ip"}
export CLUSTER_IP=$(minikube ip --profile tutorialtester)
```

```{"stage":"cert-creation", "rootdir":"$tmpdir.1", "runtime":"bash", "label":"generate cert"}
printf "000000\n000000\n${CLUSTER_IP}.nip.io\narkmq-org\nRed Hat\nGrenoble\nAuvergne Rhône Alpes\nFR\nyes\n" | keytool -genkeypair -alias artemis -keyalg RSA -keysize 2048 -storetype PKCS12 -keystore broker.ks -validity 3000
printf '000000\n' | keytool -export -alias artemis -file broker.cert -keystore broker.ks
printf '000000\n000000\nyes\n' | keytool -import -v -trustcacerts -alias artemis -file broker.cert -keystore client.ts
```
```shell tutorial_tester
Owner: CN=192.168.49.2.nip.io, OU=arkmq-org, O=Red Hat, L=Grenoble, ST=Auvergne Rhône Alpes, C=FR
Issuer: CN=192.168.49.2.nip.io, OU=arkmq-org, O=Red Hat, L=Grenoble, ST=Auvergne Rhône Alpes, C=FR
Serial number: 2ff151d6
Valid from: Wed Sep 04 14:33:24 CEST 2024 until: Sun Nov 21 13:33:24 CET 2032
Certificate fingerprints:
	 SHA1: 2A:E3:2B:FF:21:FB:5D:93:7E:33:F3:C6:5D:3A:7F:09:53:54:77:E3
	 SHA256: 04:B1:FB:77:98:35:BF:5D:48:89:8B:9B:F7:8F:1B:9D:56:1F:32:C5:77:0B:A3:92:F8:8E:39:2A:67:49:76:D6
Signature algorithm name: SHA256withRSA
Subject Public Key Algorithm: 2048-bit RSA key
Version: 3

Extensions: 

#1: ObjectId: 2.5.29.14 Criticality=false
SubjectKeyIdentifier [
KeyIdentifier [
0000: 26 E7 BD 4F 95 22 A5 72   12 31 CC 0A 6A D1 0A 5D  &..O.".r.1..j..]
0010: 01 7A 25 A1                                        .z%.
]
]

```

Create the secret in kubernetes

```{"stage":"cert-creation", "rootdir":"$tmpdir.1"}
kubectl create secret generic send-receive-sslacceptor-secret --from-file=broker.ks --from-file=client.ts --from-literal=keyStorePassword='000000' --from-literal=trustStorePassword='000000' -n send-receive-project
```
```shell tutorial_tester
secret/send-receive-sslacceptor-secret created
```

Get the path of the cert folder for later

```{"stage":"cert-creation", "rootdir":"$tmpdir.1", "runtime":"bash", "label":"get cert folder"}
export CERT_FOLDER=$(pwd)
```

#### Start the broker


```{"stage":"deploy", "HereTag":"EOF", "runtime":"bash", "label":"deploy the broker", "env":["CLUSTER_IP"]}
kubectl apply -f - << EOF
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemis
metadata:
  name: send-receive
  namespace: send-receive-project
spec:
  ingressDomain: ${CLUSTER_IP}.nip.io
  acceptors:
    - name: sslacceptor
      port: 62626
      expose: true
      sslEnabled: true
      sslSecret: send-receive-sslacceptor-secret
  brokerProperties:
    - addressConfigurations."APP.JOBS".routingTypes=ANYCAST
    - addressConfigurations."APP.JOBS".queueConfigs."APP.JOBS".routingType=ANYCAST
    - addressConfigurations."APP.COMMANDS".routingTypes=MULTICAST
EOF
```
```shell tutorial_tester
activemqartemis.broker.amq.io/send-receive created
```

Wait for the Broker to be ready:

```{"stage":"deploy"}
kubectl wait ActiveMQArtemis send-receive --for=condition=Ready --namespace=send-receive-project --timeout=240s
```
```shell tutorial_tester
activemqartemis.broker.amq.io/send-receive condition met
```

#### Create a route to access the ingress:

Check for the ingress availability:

```{"stage":"deploy"}
kubectl get ingress --show-labels
```
```shell tutorial_tester
NAME                                 CLASS   HOSTS                                                                         ADDRESS        PORTS     AGE   LABELS
send-receive-sslacceptor-0-svc-ing   nginx   send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io   192.168.49.2   80, 443   41s   ActiveMQArtemis=send-receive,application=send-receive-app,statefulset.kubernetes.io/pod-name=send-receive-ss-0
```

### Exchanging messages between a producer and a consumer

Download the [latest
release](https://activemq.apache.org/components/artemis/download/) of ActiveMQ
Artemis, decompress the tarball and locate the artemis executable.

```{"stage":"test_setup", "rootdir":"$tmpdir.1", "runtime":"bash", "label":"download artemis"}
wget --quiet https://archive.apache.org/dist/activemq/activemq-artemis/2.36.0/apache-artemis-2.36.0-bin.tar.gz
tar -zxf apache-artemis-2.36.0-bin.tar.gz apache-artemis-2.36.0/
```

#### Figure out the broker endpoint

The `artemis` will need to point to the https endpoint generated in earlier with
a couple of parameters set:
* `sslEnabled` = `true`
* `verifyHost` = `false`
* `trustStorePath` = `/some/path/broker.ks`
* `trustStorePassword` = `000000`

To use the consumer and the producer you'll need to give the path to the
`broker.ks` file you've created earlier. In the following commands the file is
located to `${CERT_FOLDER}/broker.ks`.


```{"stage":"test_setup", "runtime":"bash", "label":"get the ingress host"}
export INGRESS_URL=$(kubectl get ingress send-receive-sslacceptor-0-svc-ing -o json | jq -r '.spec.rules[] | .host')
```

Craft the broker url for artemis

```{"stage":"test_setup", "runtime":"bash", "label":"compute the broker url"}
export BROKER_URL="tcp://${INGRESS_URL}:443?sslEnabled=true&verifyHost=false&trustStorePath=${CERT_FOLDER}/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false"
```

##### Test the connection

```{"stage":"test0", "rootdir":"$tmpdir.1/apache-artemis-2.36.0/bin/", "runtime":"bash", "label":"test connection"}
./artemis check queue --name TEST --produce 10 --browse 10 --consume 10 --url ${BROKER_URL} --verbose
```
```shell tutorial_tester
Executing org.apache.activemq.artemis.cli.commands.check.QueueCheck check queue --name TEST --produce 10 --browse 10 --consume 10 --url tcp://send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io:443?sslEnabled=true&verifyHost=false&trustStorePath=/tmp/4159481985/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false --verbose 
Home::/tmp/4159481985/apache-artemis-2.36.0, Instance::null
Connection brokerURL = tcp://send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io:443?sslEnabled=true&verifyHost=false&trustStorePath=/tmp/4159481985/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false
Running QueueCheck
Checking that a producer can send 10 messages to the queue TEST ... success
Checking that a consumer can browse 10 messages from the queue TEST ... success
Checking that a consumer can consume 10 messages from the queue TEST ... success
Checks run: 3, Failures: 0, Errors: 0, Skipped: 0, Time elapsed: 0.415 sec - QueueCheck
```

#### ANYCAST

For this use case, run first the producer, then the consumer.

```{"stage":"test1", "rootdir":"$tmpdir.1/apache-artemis-2.36.0/bin/", "runtime":"bash", "label":"anycast: produce 1000 messages"}
./artemis producer --destination queue://APP.JOBS --url ${BROKER_URL}
```
```shell tutorial_tester
Connection brokerURL = tcp://send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io:443?sslEnabled=true&verifyHost=false&trustStorePath=/tmp/4159481985/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false
Producer ActiveMQQueue[APP.JOBS], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[APP.JOBS], thread=0 Produced: 1000 messages
Producer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in second : 8 s
Producer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in milli second : 8735 milli seconds
```

```{"stage":"test1", "rootdir":"$tmpdir.1/apache-artemis-2.36.0/bin/", "runtime":"bash", "label":"anycast: consume 1000 messages"}
./artemis consumer --destination queue://APP.JOBS --url ${BROKER_URL}
```
```shell tutorial_tester
Connection brokerURL = tcp://send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io:443?sslEnabled=true&verifyHost=false&trustStorePath=/tmp/4159481985/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false
Consumer:: filter = null
Consumer ActiveMQQueue[APP.JOBS], thread=0 wait until 1000 messages are consumed
Received 1000
Consumer ActiveMQQueue[APP.JOBS], thread=0 Consumed: 1000 messages
Consumer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in second : 0 s
Consumer ActiveMQQueue[APP.JOBS], thread=0 Elapsed time in milli second : 99 milli seconds
Consumer ActiveMQQueue[APP.JOBS], thread=0 Consumed: 1000 messages
Consumer ActiveMQQueue[APP.JOBS], thread=0 Consumer thread finished
```

#### MULTICAST

For this use case, run first the consumer(s), then the producer.
[More details there](https://activemq.apache.org/components/artemis/documentation/2.0.0/address-model.html).

1. in `n` other terminal(s) connect `n` consumer(s):

```{"stage":"test2", "rootdir":"$tmpdir.1/apache-artemis-2.36.0/bin/", "parallel": true, "runtime":"bash", "env":["BROKER_URL"], "label":"multicast: consume 1000 messages"}
./artemis consumer --destination topic://APP.COMMANDS --url ${BROKER_URL}
```
```shell tutorial_tester
Connection brokerURL = tcp://send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io:443?sslEnabled=true&verifyHost=false&trustStorePath=/tmp/4159481985/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false
Consumer:: filter = null
Consumer ActiveMQTopic[APP.COMMANDS], thread=0 wait until 1000 messages are consumed
Received 1000
Consumer ActiveMQTopic[APP.COMMANDS], thread=0 Consumed: 1000 messages
Consumer ActiveMQTopic[APP.COMMANDS], thread=0 Elapsed time in second : 5 s
Consumer ActiveMQTopic[APP.COMMANDS], thread=0 Elapsed time in milli second : 5719 milli seconds
Consumer ActiveMQTopic[APP.COMMANDS], thread=0 Consumed: 1000 messages
Consumer ActiveMQTopic[APP.COMMANDS], thread=0 Consumer thread finished
```

2. connect the producer to start broadcasting messages.

```{"stage":"test2", "rootdir":"$tmpdir.1/apache-artemis-2.36.0/bin/", "parallel": true, "runtime":"bash", "env":["BROKER_URL"], "label":"multicast: produce 1000 messages"}
sleep 5s
./artemis producer --destination topic://APP.COMMANDS --url ${BROKER_URL}
```
```shell tutorial_tester
Connection brokerURL = tcp://send-receive-sslacceptor-0-svc-ing-send-receive-project.192.168.49.2.nip.io:443?sslEnabled=true&verifyHost=false&trustStorePath=/tmp/4159481985/broker.ks&trustStorePassword=000000&useTopologyForLoadBalancing=false
Producer ActiveMQTopic[APP.COMMANDS], thread=0 Started to calculate elapsed time ...

Producer ActiveMQTopic[APP.COMMANDS], thread=0 Produced: 1000 messages
Producer ActiveMQTopic[APP.COMMANDS], thread=0 Elapsed time in second : 0 s
Producer ActiveMQTopic[APP.COMMANDS], thread=0 Elapsed time in milli second : 652 milli seconds
```

### cleanup

To leave a pristine environment after executing this tutorial you can simply,
delete the minikube cluster and clean the `/etc/hosts` file.

```{"stage":"teardown", "requires":"init/minikube_start"}
minikube delete --profile tutorialtester
```
```shell tutorial_tester
* Deleting "tutorialtester" in docker ...
* Deleting container "tutorialtester" ...
* Removing /home/tlavocat/.minikube/machines/tutorialtester ...
* Removed all traces of the "tutorialtester" cluster.
```
