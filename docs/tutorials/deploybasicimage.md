---
title: "Deploying the Basic Broker Image"
description: "Deploying the Basic Broker Image."
date: 2020-11-16T13:59:39+01:00
lastmod: 2020-11-16T13:59:39+01:00
draft: false
images: []
menu:
  docs:
    parent: "tutorials"
weight: 110
toc: true
---

The basic Broker Container image is the easiest way to get the broker up and running as a container, we'll explain what it is and how to run it locally.

The Basic Broker Container Image is the simplest of images to get started with, it uses environment variables to configure the broker and then starts it.
You can find the basic Broker Container Image at [quay.io](https://quay.io/repository/arkmq-org/activemq-artemis-broker)

You can use your favourite tool to run the container but for this example we are using docker.

All you need to do is execute the docker run command which will download the basic Broker Image and run it locally for you,
in this instance we are using the latest dev tag but you could choose a released version if needed.  

```shell script
    docker run -e AMQ_USER=admin -e AMQ_PASSWORD=admin --name artemis quay.io/arkmq-org/activemq-artemis-broker:dev.latest
```

This now should download the latest image and run it, you should see:


```shell script
dev.latest: Pulling from arkmq-org/activemq-artemis-broker
eae19a56e9c6: Pull complete
be73321c7956: Pull complete
4b32e1d9d455: Pull complete
Digest: sha256:891dc91d789d93ed474df00355bd173c3980158aa68cba0737a81b920fc0bf2f
Status: Downloaded newer image for quay.io/arkmq-org/activemq-artemis-broker:dev.latest
Creating Broker with args --user XXXXX --password XXXXX --role admin --name broker --allow-anonymous --http-host 172.17.0.2 --host 172.17.0.2   --force
Creating ActiveMQ Artemis instance at: /home/jboss/broker

Auto tuning journal ...
done! Your system can make 0.5 writes per millisecond, your journal-buffer-timeout will be 1988000

You can now start the broker by executing:  

   "/home/jboss/broker/bin/artemis" run

Or you can run the broker in the background using:

   "/home/jboss/broker/bin/artemis-service" start

Running Broker
     _        _               _
    / \  ____| |_  ___ __  __(_) _____
   / _ \|  _ \ __|/ _ \  \/  | |/  __/
  / ___ \ | \/ |_/  __/ |\/| | |\___ \
 /_/   \_\|   \__\____|_|  |_|_|/___ /
 Apache ActiveMQ Artemis 2.16.0


2021-01-29 10:05:07,903 INFO  [org.apache.activemq.artemis.integration.bootstrap] AMQ101000: Starting ActiveMQ Artemis Server
.........
2021-01-29 10:05:09,372 INFO  [org.apache.activemq.artemis] AMQ241004: Artemis Console available at http://172.17.0.2:8161/console
```

Well done you have now deployed your first arkmq-org image. Now we want to expose the broker to the outside world so
stop the broker and remove the image.

```shell script
docker rm artemis
```

Now re run the broker pod and expose the broker by publishing the broker's console port 8161 on the docker hosts machine port 80.

```shell script
docker run -e AMQ_USER=admin -e AMQ_PASSWORD=admin -p80:8161 --name artemis quay.io/arkmq-org/activemq-artemis-broker:dev.latest
```
Now open up a browser and go to http://localhost/console and login using the username and password you provided in the docker command.

Lastly follow the above steps to recreate the broker image but now also pass in the params -p61616:61616 to expose the acceptor
ports for all client protocols. You can now connect an external client to localhost:61616. using the Artemis CLi you can
easily send some messages using..

```shell script
artemis producer
```

For more information on ActiveMQ Artemis please read the [Artemis Documentation](https://activemq.apache.org/components/artemis/documentation/)
and for available environment properties to set you can check the [image.yaml](https://github.com/arkmq-org/activemq-artemis-broker-image/blob/master/image.yaml)
