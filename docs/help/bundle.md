---
title: "Bundle"
description: "Bundle arkmq-org.io"
lead: "Bundle arkmq-org.io"
date: 2020-10-06T08:49:31+00:00
lastmod: 2020-10-06T08:49:31+00:00
draft: false
images: []
menu:
  docs:
    parent: "help"
weight: 630
toc: true
---

# Bundle

## Operator Lifecycle Manager (OLM)
The [Operator Lifecycle Manager](https://olm.operatorframework.io/) can help users to install and manage operators. The arkmq-org operator can be built into a bundle image and installed into OLM.

### Install OLM
Check out the latest [releases on github](https://github.com/operator-framework/operator-lifecycle-manager/releases) for release-specific install instructions.

## Create a repository
Create a repository that Kubernetes will uses to pull your catalog image. You can create a public one for free on quay.io, see [how to create a repo](https://docs.quay.io/guides/create-repo.html).

## Build a catalog image
Set your repository in CATALOG_IMG and execute the following command:
```
make CATALOG_IMG=quay.io/my-org/activemq-artemis-operator-catalog:latest catalog-build
```

## Push a catalog image
Set your repository in CATALOG_IMG and execute the following command:
```
make CATALOG_IMG=quay.io/my-org/activemq-artemis-operator-catalog:latest catalog-push
```

## Create a catalog source (e.g. catalog-source.yaml):
Before creating the catalog source, ensure to update the **image** field within the `spec` section with your own built catalog image specified by the `CATALOG_IMG` environment variable.
For the `CATALOG_IMG`, refer to the [Build a catalog image](#build-a-catalog-image) section.

```
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: activemq-artemis-operator-source
  namespace: operators
spec:
  displayName: ActiveMQ Artemis Operators
  image: quay.io/my-org/activemq-artemis-operator-catalog:latest
  sourceType: grpc
```

and deploy it:

```$xslt
$ kubectl create -f catalog-source.yaml
```
In a moment you will see the catalog image is up and running in namespace **operators**:

```$xslt
$ kubectl get pod -n operators
NAME                                      READY   STATUS    RESTARTS   AGE
activemq-artemis-operator-source-g94fd    1/1     Running   0          42s
```

## Create a subscription (e.g. subscription.yaml)

```
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: activemq-artemis-operator-subscription
  namespace: operators
spec:
  channel: stable
  name: activemq-artemis-operator
  source: activemq-artemis-operator-source
  sourceNamespace: operators
```

and deploy it:
```$xslt
  $ kubectl create -f subscription.yaml
```
An operator will be installed into **operators** namespace.

```$xslt
$ kubectl get pod -n operators
NAME                                                              READY   STATUS      RESTARTS   AGE
069c5d363d51fc04d639086da1c5180883a6cea8ec9d9f9eedde1a55f6v7jsq   0/1     Completed   0          9m55s
activemq-artemis-controller-manager-54c99b9df6-6xdzh              1/1     Running     0          9m28s
activemq-artemis-operator-source-g94fd                            1/1     Running     0          58m
```

## Create a single ActiveMQ Artemis

This step creates a single ActiveMQ Artemis broker instance by applying the custom resource (CR) defined in artemis_single.yaml file.

```$xslt
$ kubectl apply -f examples/artemis/artemis_single.yaml
```

To check the status of the broker, run:

```$xslt 
$ kubectl get ActivemqArtemis 
NAME             READY   AGE
artemis-broker   True    39s
```