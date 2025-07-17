---
title: "Service and App CRD Round Trip"
description: "A tutorial on using ActiveMQArtemisService and ActiveMQArtemisApp CRDs, based on the 'round trip simple' e2e test."
draft: false
images: []
menu:
  docs:
    parent: "tutorials"
weight: 121
toc: true
---

This tutorial walks through a complete round trip of sending and receiving messages using the `ActiveMQArtemisService` and `ActiveMQArtemisApp` CRDs. The steps are based on the `round trip simple` e2e test to ensure a working configuration.

### Prerequisites

- A running Kubernetes cluster (this tutorial uses `minikube`).
- `kubectl` configured to interact with your cluster.

### 1. Setup

#### Start Minikube

```bash {"stage":"init", "id":"minikube_start", "runtime":"bash"}
minikube start --profile service-app-tutorial
minikube profile service-app-tutorial
```

#### Create Namespace

```bash {"stage":"init", "runtime":"bash" }
kubectl create namespace service-app-project
kubectl config set-context --current --namespace=service-app-project
```

#### Install Cert-Manager

```bash {"stage":"init", "label":"install cert-manager", "runtime":"bash"}
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.1/cert-manager.yaml
```

Wait for `cert-manager` to be ready.

```bash {"stage":"init", "label":"wait for cert-manager", "runtime":"bash"}
kubectl wait deployment --for=condition=Available -n cert-manager --timeout=600s cert-manager cert-manager-cainjector cert-manager-webhook
```

#### Install Trust Manager

First, add the Jetstack Helm repository.

```bash {"stage":"init", "label":"add jetstack helm repo", "runtime":"bash"}
helm repo add jetstack https://charts.jetstack.io --force-update
```

Now, install `trust-manager`.

```bash {"stage":"init", "label":"install trust-manager", "runtime":"bash"}
helm upgrade trust-manager jetstack/trust-manager --install --namespace cert-manager --set secretTargets.enabled=true --set secretTargets.authorizedSecretsAll=true --wait
```

#### Install the Operator

```bash {"stage":"init", "rootdir":"$operator", "runtime":"bash"}
./deploy/install_opr.sh
```

Wait for the operator pod to become ready.

```bash {"stage":"init", "label":"wait for the operator to be running", "runtime":"bash"}
kubectl wait pod --all --for=condition=Ready --namespace=service-app-project --timeout=600s
```

### 2. Configure Certificates
We'll set up a CA and issue certificates for the operator, the service, and the application.

#### Create Issuers and Root Certificate

First the root issuer.

```bash {"stage":"deploy_certs", "label":"create root issuer", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: root-issuer
spec:
  selfSigned: {}
EOF
```

```bash {"stage":"deploy_certs", "label":"wait for root issuer", "runtime":"bash"}
kubectl wait clusterissuer root-issuer --for=condition=Ready --timeout=300s
```

Then the root certificate.

```bash {"stage":"deploy_certs", "label":"create root cert", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: root-cert
  namespace: cert-manager
spec:
  isCA: true
  commonName: artemis.root.ca
  secretName: artemis-root-cert-secret
  issuerRef:
    name: root-issuer
    kind: ClusterIssuer
EOF
```

```bash {"stage":"deploy_certs", "label":"wait for root cert", "runtime":"bash"}
kubectl wait certificate root-cert --for=condition=Ready -n cert-manager --timeout=300s
```

Then a signing issuer that uses the root certificate.

```bash {"stage":"deploy_certs", "label":"create signing issuer", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: broker-ca-issuer
spec:
  ca:
    secretName: artemis-root-cert-secret
EOF
```

```bash {"stage":"deploy_certs", "label":"wait for signing issuer", "runtime":"bash"}
kubectl wait clusterissuer broker-ca-issuer --for=condition=Ready --timeout=300s
```

#### Create Operator Certificate

##### Install the CA Bundle in the `cert-manager` namespace

```bash {"stage":"deploy_certs", "label":"create ca bundle", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: trust.cert-manager.io/v1alpha1
kind: Bundle
metadata:
  name: activemq-artemis-manager-ca
  namespace: cert-manager
spec:
  sources:
  - secret:
      name: artemis-root-cert-secret
      key: "tls.crt"
  target:
    secret:
      key: "ca.pem"
EOF
```

```bash {"stage":"deploy_certs", "label":"wait for ca bundle", "runtime":"bash"}
kubectl wait bundle activemq-artemis-manager-ca -n cert-manager --for=condition=Synced --timeout=300s
```

##### Create the certificate for the operator

```bash {"stage":"deploy_certs", "label":"create operator cert", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: activemq-artemis-manager-cert
  namespace: service-app-project
spec:
  secretName: activemq-artemis-manager-cert
  commonName: activemq-artemis-operator
  issuerRef:
    name: broker-ca-issuer
    kind: ClusterIssuer
EOF
```

```bash {"stage":"deploy_certs", "label":"wait for operator cert", "runtime":"bash"}
kubectl wait certificate activemq-artemis-manager-cert -n service-app-project --for=condition=Ready --timeout=300s
```

### 3. Deploy the Messaging Service and Application

#### Create Service Certificate

The service needs a certificate customized with a matching common name to enable
mTLS communication.

```bash {"stage":"deploy_service", "label":"create broker cert", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: messaging-service-broker-cert
  namespace: service-app-project
spec:
  secretName: messaging-service-broker-cert
  commonName: messaging-service
  dnsNames:
  - messaging-service
  - '*.messaging-service-hdls-svc.service-app-project.svc.cluster.local'
  issuerRef:
    name: broker-ca-issuer
    kind: ClusterIssuer
EOF
```

```bash {"stage":"deploy_certs", "label":"wait for operator cert", "runtime":"bash"}
kubectl wait certificate messaging-service-broker-cert -n service-app-project --for=condition=Ready --timeout=300s
```

#### Deploy `ActiveMQArtemisService`

```bash {"stage":"deploy_service", "label":"deploy service crd", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemisService
metadata:
  name: messaging-service
  namespace: service-app-project
  labels:
    forWorkQueue: "true"
spec:
  resources:
    limits:
      memory: "1Gi"
  env:
    - name: JAVA_ARGS_APPEND
      value: "-Dlog4j2.level=INFO"
  auth:
  - mtls
  acceptors:
  - name: amqp
EOF
```

Wait for the resource to be ready.

```bash {"stage":"deploy_service", "label":"wait for service"}
kubectl wait ActiveMQArtemisService messaging-service -n service-app-project --for=condition=Ready --timeout=300s
```

#### Create Service Certificate

```bash {"stage":"deploy_app", "label":"create app cert", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: first-app-app-cert
  namespace: service-app-project
spec:
  secretName: first-app-app-cert
  commonName: first-app
  issuerRef:
    name: broker-ca-issuer
    kind: ClusterIssuer
EOF
```

#### Deploy `ActiveMQArtemisApp`

```bash {"stage":"deploy_app", "label":"deploy app crd", "HereTag":"EOF", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: broker.amq.io/v1beta1
kind: ActiveMQArtemisApp
metadata:
  name: first-app
  namespace: service-app-project
spec:
  selector:
    matchLabels:
      forWorkQueue: "true"
  auth:
  - mtls
  capabilities:
  - producerandconsumerof:
    - queuename: "APP.JOBS"
EOF
```

Wait for the resource to be ready.

```bash {"stage":"deploy_app", "label":"wait for app", "runtime":"bash"}
kubectl wait ActiveMQArtemisApp first-app -n service-app-project --for=condition=Ready --timeout=300s
```

### 4. Test Messaging

#### Create Client Configuration

```bash {"stage":"test_messaging", "label":"create pemcfg secret", "runtime":"bash"}
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: cert-pemcfg
  namespace: service-app-project
type: Opaque
stringData:
  tls.pemcfg: |
    source.key=/app/tls/client/tls.key
    source.cert=/app/tls/client/tls.crt
  java.security: security.provider.6=de.dentrassi.crypto.pem.PemKeyStoreProvider
EOF
```

```bash {"stage":"test_messaging", "label":"wait for pemcfg secret", "runtime":"bash"}
until kubectl get secret cert-pemcfg -n service-app-project &> /dev/null; do echo "Waiting for secret..." && sleep 2; done
```

#### Run Producer Job

```bash {"stage":"test_messaging", "label":"run producer", "HereTag":"EOF", "runtime":"bash"}
cat <<'EOT' | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  name: producer
  namespace: service-app-project
spec:
  template:
    spec:
      containers:
      - name: producer
        image: quay.io/arkmq-org/activemq-artemis-broker-kubernetes:artemis.2.40.0
        command:
        - "/bin/sh"
        - "-c"
        - exec java -classpath /opt/amq/lib/*:/opt/amq/lib/extra/* org.apache.activemq.artemis.cli.Artemis producer --protocol=AMQP --user p --password passwd --url 'amqps://messaging-service:61616?transport.trustStoreType=PEMCA&transport.trustStoreLocation=/app/tls/ca/ca.pem&transport.keyStoreType=PEMCFG&transport.keyStoreLocation=/app/tls/pem/tls.pemcfg' --message-count 1 --destination queue://APP.JOBS;
        env:
        - name: JDK_JAVA_OPTIONS
          value: "-Djava.security.properties=/app/tls/pem/java.security"
        volumeMounts:
        - name: trust
          mountPath: /app/tls/ca
        - name: cert
          mountPath: /app/tls/client
        - name: pem
          mountPath: /app/tls/pem
      volumes:
      - name: trust
        secret:
          secretName: activemq-artemis-manager-ca
      - name: cert
        secret:
          secretName: first-app-app-cert
      - name: pem
        secret:
          secretName: cert-pemcfg
      restartPolicy: OnFailure
EOT
```

#### Run Consumer Job

```bash {"stage":"test_messaging", "label":"run consumer", "HereTag":"EOF", "runtime":"bash"}
cat <<'EOT' | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  name: consumer
  namespace: service-app-project
spec:
  template:
    spec:
      containers:
      - name: consumer
        image: quay.io/arkmq-org/activemq-artemis-broker-kubernetes:artemis.2.40.0
        command:
        - "/bin/sh"
        - "-c"
        - exec java -classpath /opt/amq/lib/*:/opt/amq/lib/extra/* org.apache.activemq.artemis.cli.Artemis consumer --protocol=AMQP --user p --password passwd --url 'amqps://messaging-service:61616?transport.trustStoreType=PEMCA&transport.trustStoreLocation=/app/tls/ca/ca.pem&transport.keyStoreType=PEMCFG&transport.keyStoreLocation=/app/tls/pem/tls.pemcfg' --message-count 1 --destination queue://APP.JOBS --receive-timeout 10000;
        env:
        - name: JDK_JAVA_OPTIONS
          value: "-Djava.security.properties=/app/tls/pem/java.security"
        volumeMounts:
        - name: trust
          mountPath: /app/tls/ca
        - name: cert
          mountPath: /app/tls/client
        - name: pem
          mountPath: /app/tls/pem
      volumes:
      - name: trust
        secret:
          secretName: activemq-artemis-manager-ca
      - name: cert
        secret:
          secretName: first-app-app-cert
      - name: pem
        secret:
          secretName: cert-pemcfg
      restartPolicy: OnFailure
EOT
```

Wait for jobs to complete.

```bash {"stage":"test_messaging", "label":"wait for jobs", "runtime":"bash"}
kubectl wait job producer -n service-app-project --for=condition=Complete --timeout=120s
kubectl wait job consumer -n service-app-project --for=condition=Complete --timeout=120s
```

### 5. Cleanup

Finally, delete the minikube cluster.

```bash {"stage":"teardown", "requires":"init/minikube_start", "runtime":"bash"}
minikube delete --profile service-app-tutorial
```

