---
title: "Setting up SSL connections with Arkmq-org Operator"  
description: "An example for setting up ssl connections for broker in kubernetes with operator"
draft: false
images: []
menu:
  docs:
    parent: "tutorials"
weight: 110
toc: true
---

Security is always a concern in a production environment. With arkmq-org Operator
You can easily configure and set up a broker with ssl-enabled acceptors. The blog explains how to do it.

The [ActiveMQ Artemis](https://activemq.apache.org/components/artemis/) broker supports a variety of network protocols(tcp, http, etc) including [SSL(TLS)](https://en.wikipedia.org/wiki/Transport_Layer_Security) secure connections. Underneath it uses [Netty](https://netty.io/) as the base transport layer.

This article guides you through the steps to set up a broker to run in kubernetes (Minikube). The broker will listen on a secure port 61617 (ssl over tcp). It also demonstrates sending and receiving messages over secure connections using one-way authentication.

### Prerequisite
Before you start you need have access to a running Kubernetes cluster environment. A [Minikube](https://minikube.sigs.k8s.io/docs/start/) running on your laptop will just do fine. The arkmq-org operator also runs in a Openshift cluster environment like [CodeReady Container](https://developers.redhat.com/products/codeready-containers/overview). In this blog we assume you have Kubernetes cluster environment. (If you use CodeReady the client tool is **oc** in place of **kubectl**)

### Deploy Arkmq-org operator
First you need to deploy the arkmq-org operator.
If you are not sure how to deploy the operator take a look at [this blog]({{< relref "using_operator.md" >}}).

In this blog post we assume you deployed the operator to a namespace called **myproject**.

Make sure the operator is in "Runing" status before going to the next step.
You can run this command and observe the output:

```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-zcqmw   1/1     Running   0          7m32s
```

### Prepare keystore and truststore
To establish a SSL connection you need certificates. Here for demonstration purpose we prepare a self-signed certificate.

We'll use the "keytool" utility that comes with JDK:

```shell script
$ keytool -genkeypair -alias artemis -keyalg RSA -keysize 2048 -storetype PKCS12 -keystore broker.ks -validity 3000
Enter keystore password:  
Re-enter new password:
What is your first and last name?
  [Unknown]:  Howard Gao
What is the name of your organizational unit?
  [Unknown]:  JBoss
What is the name of your organization?
  [Unknown]:  Red Hat
What is the name of your City or Locality?
  [Unknown]:  Beijing
What is the name of your State or Province?
  [Unknown]:  Beijing
What is the two-letter country code for this unit?
  [Unknown]:  CN
Is CN=Howard Gao, OU=JBoss, O=Red Hat, L=Beijing, ST=Beijing, C=CN correct?
  [no]:  yes
```
It creates a keystore file named **broker.ks** under the current directory.
Let's give the password as **password** when prompted above.

Next make a truststore using the same cert in the keystore.

```shell script
$ keytool -export -alias artemis -file broker.cert -keystore broker.ks
Enter keystore password:  
Certificate stored in file <broker.cert>
```
```shell script
$ keytool -import -v -trustcacerts -alias artemis -file broker.cert -keystore client.ts
Enter keystore password:  
Re-enter new password:
Owner: CN=Howard Gao, OU=JBoss, O=Red Hat, L=Beijing, ST=Beijing, C=CN
Issuer: CN=Howard Gao, OU=JBoss, O=Red Hat, L=Beijing, ST=Beijing, C=CN
Serial number: 582b8fd4
Valid from: Mon Feb 08 20:17:45 CST 2021 until: Fri Apr 27 20:17:45 CST 2029
Certificate fingerprints:
	 MD5:  01:89:A9:B0:07:A1:2F:19:FC:43:5C:27:2E:E8:D7:C3
	 SHA1: D4:25:61:9F:AA:B6:05:1F:CC:F0:CD:65:A8:BC:B0:E1:70:49:1B:81
	 SHA256: 64:63:4E:68:2E:98:59:DA:A4:6B:FF:8E:E7:8C:AC:65:A2:F2:37:CB:12:BC:96:3C:AE:70:44:63:BD:0D:41:AE
Signature algorithm name: SHA256withRSA
Subject Public Key Algorithm: 2048-bit RSA key
Version: 3

Extensions:

#1: ObjectId: 2.5.29.14 Criticality=false
SubjectKeyIdentifier [
KeyIdentifier [
0000: C6 35 D2 14 85 C3 A2 68   E5 A3 78 D3 9F 3F D2 C7  .5.....h..x..?..
0010: 8F 9D B6 A9                                        ....
]
]

Trust this certificate? [no]:  yes
Certificate was added to keystore
[Storing client.ts]
```
Make sure the password for your truststore **client.ts** is also **password**.

By default the operator fetches the truststore and keystore from a secret in kubernetes in order to configure SSL acceptors for a broker. The secret name is deducted from broker CR's name combined with the acceptor's name.

Here we'll use "ex-aao" for CR's name and "sslacceptor" for the acceptor's name. So the truststore and keystore should be stored in a secret named **ex-aao-sslacceptor-secret**.

Run the following command to create the secret we need:
```shell script
$ kubectl create secret generic ex-aao-sslacceptor-secret --from-file=broker.ks --from-file=client.ts --from-literal=keyStorePassword='password' --from-literal=trustStorePassword='password' -n myproject
secret/ex-aao-sslacceptor-secret created
```

### Prepare the broker CR with SSL enabled
Now create a file named "broker_ssl_enabled.yaml" with the following contents:

```yaml
apiVersion: broker.amq.io/v2alpha4
kind: ActiveMQArtemis
metadata:
  name: ex-aao
spec:
  deploymentPlan:
    size: 1
    image: quay.io/arkmq-org/activemq-artemis-broker-kubernetes:0.2.1
  acceptors:
  - name: sslacceptor
    protocols: all
    port: 61617
    sslEnabled: true
```
In this broker CR we configure an acceptor named "sslacceptor" that listens on tcp port 61617. The **sslEnabled: true** tells the operator to make this acceptor to use SSL transport.

### Deploy the broker
Deploy the above **broker_ssl_enabled.yaml** to the cluster:
```shell script
$ kubectl create -f broker_ssl_enabled.yaml -n myproject
activemqartemis.broker.amq.io/ex-aao created
```
In a moment the broker should be up and running. Run the command to check it out:
```shell script
$ kubectl get pod -n myproject
NAME                                         READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-58bb658f4c-zcqmw   1/1     Running   0          18m
ex-aao-ss-0                                  1/1     Running   0          71s
```
Using "kubectl logs ex-aao-ss-0 -n myproject" command you can checkout the console log of the broker. You'll be seeing a line like this in the log:

```
2021-02-08 12:54:12,837 INFO  [org.apache.activemq.artemis.core.server] AMQ221020: Started EPOLL Acceptor at ex-aao-ss-0.ex-aao-hdls-svc.default.svc.cluster.local:61617 for protocols [CORE,MQTT,AMQP,HORNETQ,STOMP,OPENWIRE]
```
which means the acceptor is now listening on port 61617. Although it doesn't give us whether it's SSL or plain tcp we can check out in the following steps that it accepts SSL connections only.

One way to check out that this acceptor is indeed SSL enabled is to log in to the broker pod and take a look at it's configure file in /home/jboss/amq-broker/etc/broker.xml. In it there should be an element like this:

```xml
<acceptor name="sslacceptor">tcp://ex-aao-ss-0.ex-aao-hdls-svc.default.svc.cluster.local:61617?protocols=AMQP,CORE,HORNETQ,MQTT,OPENWIRE,STOMP;sslEnabled=true;keyStorePath=/etc/ex-aao-sslacceptor-secret-volume/broker.ks;keyStorePassword=password;trustStorePath=/etc/ex-aao-sslacceptor-secret-volume/client.ts;trustStorePassword=password;tcpSendBufferSize=1048576;tcpReceiveBufferSize=1048576;useEpoll=true;amqpCredits=1000;amqpMinCredits=300</acceptor>
```

### Test messaging over a SSL connection
With the broker pod in running status we can proceed to make some connections against it and do some simple messaging. We'll use Artemis broker's built in CLI commands to do this.

Log into the broker pod first to get a shell command environment:

```shell script
$ kubectl exec --stdin --tty ex-aao-ss-0 -- /bin/bash
[jboss@ex-aao-ss-0 ~]$
```
Then send 100 messages through port 61617:

```shell script
$ cd amq-broker/bin
[jboss@ex-aao-ss-0 bin]$ ./artemis producer --user admin --password admin --url tcp://ex-aao-ss-0:61617?sslEnabled=true\&trustStorePath=/etc/ex-aao-sslacceptor-secret-volume/client.ts\&trustStorePassword=password --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-0:61617?sslEnabled=true&trustStorePath=/etc/ex-aao-sslacceptor-secret-volume/client.ts&trustStorePassword=password
Producer ActiveMQQueue[TEST], thread=0 Started to calculate elapsed time ...

Producer ActiveMQQueue[TEST], thread=0 Produced: 100 messages
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Producer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 724 milli seconds
```
Pay attention to the **--url** option that is required to make an SSL connection to the broker.

You may also wonder how it gets the **trustStorePath** for the connection.

This is because the truststore and keystore are mounted automatically by the operator when it processes the broker CR. The mount path follows the pattern derived from CR's name (ex-aao) and the acceptor's name (sslacceptor, thus **/etc/ex-aao-sslacceptor-secret-volume**).

Now receive the messages we just sent -- also using SSL over the same port (61617):

```shell script
[jboss@ex-aao-ss-0 bin]$ ./artemis consumer --user admin --password admin --url tcp://ex-aao-ss-0:61617?sslEnabled=true\&trustStorePath=/etc/ex-aao-sslacceptor-secret-volume/client.ts\&trustStorePassword=password --message-count 100
OpenJDK 64-Bit Server VM warning: If the number of processors is expected to increase from one, then you should configure the number of parallel GC threads appropriately using -XX:ParallelGCThreads=N
Connection brokerURL = tcp://ex-aao-ss-0:61617?sslEnabled=true&trustStorePath=/etc/ex-aao-sslacceptor-secret-volume/client.ts&trustStorePassword=password
Consumer:: filter = null
Consumer ActiveMQQueue[TEST], thread=0 wait until 100 messages are consumed
Consumer ActiveMQQueue[TEST], thread=0 Consumed: 100 messages
Consumer ActiveMQQueue[TEST], thread=0 Elapsed time in second : 0 s
Consumer ActiveMQQueue[TEST], thread=0 Elapsed time in milli second : 160 milli seconds
Consumer ActiveMQQueue[TEST], thread=0 Consumed: 100 messages
Consumer ActiveMQQueue[TEST], thread=0 Consumer thread finished
```
Now you get an idea how an SSL acceptor is configured and processed by the operator and see it in action!

### Secure cluster connections
The internal cluster connections rely on the internal acceptor listening on the port `61616` and the internal connector with the name `artemis`. They can be secured with the following steps, create a secret with the secure stores, enable ssl in the internal acceptor by using the acceptor fields `sslEnabled` and `sslSecret`, and enable ssl in the internal connector by using broker properties

The server certificate included in the secure stores must include a wildcard DNS name for the internal broker instances in the `Subject Alternative Name`, i.e. for an ActiveMQArtemis CR with name `ex-aao` deployed in the namespace `test` a key and trust store with a self -signed certificate could be generated with the following commands:
```
keytool -storetype jks -keystore server-keystore.jks -storepass artemis -keypass artemis -alias server -genkey -keyalg "RSA" -keysize 2048 -dname "CN=ActiveMQ Artemis Server, OU=Artemis, O=ActiveMQ, L=AMQ, S=AMQ, C=AMQ" -validity 365 -ext bc=ca:false -ext eku=sA -ext san=dns:*.ex-aao-hdls-svc.test.svc.cluster.local

keytool -storetype jks -keystore server-keystore.jks -storepass artemis -alias server -exportcert -rfc > server.crt

keytool -storetype jks -keystore server-truststore.jks -storepass artemis -keypass artemis -importcert -alias server -file server.crt -noprompt
```
If the wildcard DNS names are not supported a DNS name for every internal broker instance could be included in the `Subject Alternative Name`, i.e.
```
keytool -storetype jks -keystore server-keystore.jks -storepass artemis -keypass artemis -alias server -genkey -keyalg "RSA" -keysize 2048 -dname "CN=ActiveMQ Artemis Server, OU=Artemis, O=ActiveMQ, L=AMQ, S=AMQ, C=AMQ" -validity 365 -ext bc=ca:false -ext eku=sA -ext san=dns:ex-aao-ss-0.ex-aao-hdls-svc.test.svc.cluster.local,dns:ex-aao-ss-1.ex-aao-hdls-svc.test.svc.cluster.local
```
If DNS names are not supported the host verification must be disabled setting the connector parameter `verifyHost` to false by using the broker properties, i.e.
```
  brokerProperties:
    - 'connectorConfigurations.artemis.params.verifyHost=false'
```

The secret with the secure stores can be created by using the following command:
```
kubectl create secret generic artemis-ssl-secret --namespace test \
--from-file=broker.ks=server-keystore.jks \
--from-file=client.ts=server-truststore.jks \
--from-literal=keyStorePassword=artemis \
--from-literal=trustStorePassword=artemis
```

The ActiveMQ Artemis with the secured internal acceptor and connector can be created by using the following command:
```
kubectl apply -f - <<EOF
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemis
metadata:
  name: ex-aao
  namespace: test
spec:
  deploymentPlan:
    size: 2
  acceptors:
    - name: artemis
      port: 61616
      sslEnabled: true
      sslSecret: artemis-ssl-secret
  brokerProperties:
    - 'connectorConfigurations.artemis.params.sslEnabled=true'
    - 'connectorConfigurations.artemis.params.trustStorePath=/etc/artemis-ssl-secret-volume/broker.ks'
    - 'connectorConfigurations.artemis.params.trustStorePassword=artemis'
EOF
```

### More SSL options
We have just demonstrated a simplified SSL configuration. In fact the operator supports quite a few more SSL options through the CRD definitions.
You can checkout those options in broker CRD [down here](https://github.com/arkmq-org/activemq-artemis-operator/blob/5183ddc4c2f66e0d270233a3f37340b14e225d80/deploy/crds/broker_activemqartemis_crd.yaml#L45)
and also read the [Artemis Doc on configuring transports](https://activemq.apache.org/components/artemis/documentation/latest/configuring-transports.html) for more information.
