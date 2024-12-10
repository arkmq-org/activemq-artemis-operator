---
title: "Setting up SSL with cert-manager and trust-manager"  
description: "An example for setting up ssl connections using cert-manager and trust-manager projects"
draft: false
images: []
menu:
  docs:
    parent: "tutorials"
weight: 110
toc: true
---

A lot of Kubernetes clusters already use cert-manager and trust-manager to handle certificates management.

The goal of this tutorial is to show how to configure ActiveMQ Artemis Operator resources to utilize both projects mentioned above for ssl communication.

## Prerequisites

- running Kubernetes cluster
- cert-manager should be installed in the cluster in the "cert-manager" namespace
- trust-manager should be installed in the cluster in the "cert-manager" namespace

Installation guides for cert-manager and trust-manager you can find on the [cert-manager project website](https://cert-manager.io/docs/).

## SSL

There are various scenarios how SSL can be achieved for your ActiveMQ Artemis brokers, to list some:

- ActiveMQ Artemis brokers handling SSL termination
  - you need to provide keystore with certificate and truststore
  - you need to configure acceptor to use ssl
- ActiveMQ Artemis brokers exposed over plain tcp (no ssl) and Istio handling SSL termination on ingress level, with mutual TLS between pods in cluster
  - you need to have Istio installed in cluster
  - you need to configure acceptor with no ssl
  - you need to setup Istio resources accordingly (Istio knowledge required)

In this tutorial we will cover the first scenario "ActiveMQ Artemis brokers handling SSL termination" with no assumption it is the best for your cluster, please consider the best setup yourself.

Let's start!

## Prepare keystore and truststore

Here we are gonna use cert-manager and trust-manager to create keystore and truststore.

Our assumption is that you are gonna use local **Issuer** or local **ClusterIssuer** in cert-manager.

Another assumption is your cert-manager and trust-manager are installed in the "cert-manager" namespace.

### Create local CA (Certificate authority)

:information_source: If you have already some local CA configured please go directly to the next step.

To create local CA, first we need to have base certificate for that issuer to allow it issue certificates.
This base certificate for our custom local CA should be saved into the secret in cert-manager namespace with the secret type kubernetes.io/tls.

One way of creating it is to use **selfsigned-cluster-issuer**:

Use kubectl apply command with file created based on the yaml below.

```shell
kubectl apply -f <file_name>
```

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-cluster-issuer
  namespace: cert-manager
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-selfsigned-ca
  namespace: cert-manager
spec:
  isCA: true
  commonName: my-selfsigned-ca
  secretName: root-secret
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: selfsigned-cluster-issuer
    kind: ClusterIssuer
    group: cert-manager.io
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: my-ca-issuer
  namespace: cert-manager
spec:
  ca:
    secretName: root-secret
```

Now you have your own CA which can sign new certificates for you.

### Create secret for keystore password

Replace **myproject** with actual namespace your brokers are installed in.

Replace **dummy_password** with actual password you want to use.

You can create a secret directly using command line:

```shell
kubectl create secret generic -n myproject jks-password-secret --from-literal=password=dummy_password
```

### Create certificate for ActiveMQ Artemis brokers

Create certificate resource using kubectl apply.

Replace **all occuriences of myproject** with actual namespace your brokers are installed in.

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: amq-tls-acceptor-cert
  namespace: myproject
spec:
  secretName: amq-ssl-secret
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  commonName: artemis-broker-ssl-0.myproject.svc.cluster.local
  dnsNames:
  - artemis-broker-ssl-0-svc.myproject.svc
  - artemis-broker-ssl-0-svc.myproject
  - artemis-broker-ssl-0-svc
  - artemis-broker-ss-0
  issuerRef:
    name: my-ca-issuer
    kind: ClusterIssuer
    group: cert-manager.io
  keystores:
    jks:
      create: true
      passwordSecretRef: # Password used to encrypt the keystore and truststore
        key: password
        name: jks-password-secret
```

### Create truststore using trust-manager

Here we will utlize trust-manager to create truststore in config map.

:information_source: Truststore does not contain any secrets, that's why it is stored by trust-manager in a config map, additionally password for truststore is set by design to "changeit", because of the same reason (no secrets).

It is advised to create a copy of your base CA certificate to not have an issue during CA certificate rotation event, so we will create a copy of base certificate of CA. In the phase of rotation of that certificate you would need to support old and new base CA certificate until all services use already certificates issued based on the new one. You can read more in trust-manager documentation about that topic.

Copy base ca secret:

```shell
kubectl get secret root-secret -n=cert-manager -o yaml | sed 's/name: .+/name: local-ca-cert-copy-trust-manager/' | kubectl apply -f -
```

Create bundle resource which trust-manager will use to create ca budles in **all** namespaces of your cluster. Use ```kubectl apply``` with content below:

```yaml
apiVersion: trust.cert-manager.io/v1alpha1
kind: Bundle
metadata:
  name: ca-bundle
  namespace: cert-manager
spec:
  sources:
  # all default CAs
  - useDefaultCAs: true
  # plus our custom local CA
  - secret:
      name: "local-ca-cert-copy-trust-manager"
      key: "tls.crt"
  target:
    configMap:
      key: "trust-bundle.pem"
    additionalFormats:
      jks:
        key: "truststore.jks"
```

This will produce configmap in all namespaces with truststore.jks and trust-bundle.pem which you can use internally for any services inside your cluster (java based or others).

### Create a secret with the details of your ssl setup for the ActiveMQ Artemis broker

Replace **dummy_password** with the password you have used in section **Create secret for keystore password**.

**Do not** replace **changeit** password, it is like that by design, for more informations please check above in trust-manager section.

```shell
kubectl create secret generic ssl-acceptor-ssl-secret -n myproject \
--from-literal=keyStorePath=/amq/extra/secrets/amq-ssl-secret/keystore.jks \
--from-literal=trustStorePath=/amq/extra/configmaps/ca-bundle/truststore.jks \
--from-literal=keyStorePassword=dummy_password \
--from-literal=trustStorePassword=changeit
```

## Deploy Arkmq-org operator

Operator will create ActiveMQ resources based on custom resources definitions (CRD).

If you are not sure how to deploy the operator take a look at [here]({{< relref "using_operator.md" >}}).

In this tutorial we assume you deployed the operator to a namespace called **activemq-artemis-operator**.

Make sure the operator is in "Runing" status before going to the next step.
You can run this command and observe the output:

```shell script
$ kubectl get pod -n activemq-artemis-operator
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-zcqmw   1/1     Running   0          7m32s
```

## Deploy a broker

Example of the broker CRD is located in the local repository cloned during activities from section **Deploy Arkmq-org operator**.

Replace **myproject** with actual namespace your brokers are installed in.

Deploy a broker with (beeing on the root folder of this repository):

```shell
kubectl apply -f examples/artemis/artemis_ssl_acceptor_cert_and_trust_managers.yaml -n myproject
```

## Deploy a queue

Example of the queue CRD is located in the local repository cloned during activities from section **Deploy Arkmq-org operator**.

Deploy a queue with (beeing on the root folder of this repository):

```shell
kubectl apply -f examples/address/address_queue.yaml -n myproject
```

The CR tells the Operator to create a queue named myQueue0 on address myAddress0 on each broker that it manages.

After the CR is deployed, you can observe the queue on the broker:

```shell
kubectl exec artemis-broker-ss-0 -n myproject -- /bin/bash /home/jboss/amq-broker/bin/artemis queue stat --user admin --password admin --url tcp://artemis-broker-ss-0:61616
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://artemis-broker-ss-0:61616
|NAME                     |ADDRESS                  |CONSUMER_COUNT |MESSAGE_COUNT |MESSAGES_ADDED |DELIVERING_COUNT |MESSAGES_ACKED |SCHEDULED_COUNT |ROUTING_TYPE |
|DLQ                      |DLQ                      |0              |0             |0              |0                |0              |0               |ANYCAST      |
|ExpiryQueue              |ExpiryQueue              |0              |0             |0              |0                |0              |0               |ANYCAST      |
|myQueue0                 |myAddress0               |0              |0             |0              |0                |0              |0               |ANYCAST      |
```

## Test messaging over a SSL connection

Log into the broker pod first to get a shell command environment:

```shell
kubectl exec --stdin --tty artemis-broker-ss-0 -n myproject -- /bin/bash
[jboss@artemis-broker-ss-0 ~]$
```

Then send 100 messages through port 61618, but before replace **dummy_password** with the password you have used in section **Create secret for keystore password**.

```shell
cd amq-broker/bin
[jboss@artemis-broker-ss-0 bin]$ ./artemis producer --user admin --password admin --url tcp://artemis-broker-ss-0:61618?sslEnabled=true\&trustStorePath=/amq/extra/configmaps/ca-bundle/truststore.jks\&trustStorePassword=changeit --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://artemis-broker-ss-0:61618?sslEnabled=true\&trustStorePath=/amq/extra/configmaps/ca-bundle/truststore.jks\&trustStorePassword=changeit
Producer ActiveMQQueue[TEST], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[TEST], thread=0 Produced: 100 messages
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 724 milli seconds
```

Finally you can receive those 100 messages (again replace **dummy_password** with the password you have used in section **Create secret for keystore password**):

```shell
[jboss@artemis-broker-ss-0 bin]$ ./artemis consumer --user admin --password admin --url tcp://artemis-broker-ss-0:61618?sslEnabled=true\&trustStorePath=/amq/extra/configmaps/ca-bundle/truststore.jks\&trustStorePassword=changeit --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://artemis-broker-ss-0:61618?sslEnabled=true\&trustStorePath=/amq/extra/configmaps/ca-bundle/truststore.jks\&trustStorePassword=changeit
Consumer:: filter = null
Consumer ActiveMQQueue[TEST], thread=0 wait until 100 messages are consumed
Consumer ActiveMQQueue[TEST], thread=0 Consumed: 100 messages
Consumer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Consumer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 160 milli seconds
Consumer ActiveMQQueue[TEST], thread=0 Consumed: 100 messages
Consumer ActiveMQQueue[TEST], thread=0 Consumer thread finished
```

This confirms everything works as expected!

Have fun!