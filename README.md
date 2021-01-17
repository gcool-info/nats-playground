# NATS JetStream Playground

Playground for a Secure, Highly Available NATS Cluster with message persistence (using JetStream).

This repo contains a cluster with 3 nodes and a Go client sending/receiving messages to it. 

## Pre-requisites
To use this repo, please install:
- [Docker](https://docs.docker.com/get-docker/)
- [mkcert](https://github.com/FiloSottile/mkcert) (zero-config tool for locally-trusted development certificates)


## Installation
Authentication for client-server and server-server (for the cluster) use X.509 certificatesn. To install them locally:

```shell
make certificates
```

## Run it
Spin up the cluster with:

```shell
docker-compose up
```

## Clean up
Once you are done testing, remove the CA from your local system trust store(s):

```shell
make cleanup
```

===========


```shell
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```


To set up the nats client context with the appropriate certificates:
```shell
nats context save local --server nats://localhost:4222 --description 'Local client' --tlscert /Users/gkoul/Desktop/blog/NATS/nats-playground/config/client-cert.pem --tlskey /Users/gkoul/Desktop/blog/NATS/nats-playground/config/client-key.pem --tlsca /Users/gkoul/Desktop/blog/NATS/nats-playground/config/rootCA.pem --select
```

To create an `orders` stream:
```shell
nats str add ORDERS --config config/orders-stream.json 
```

To check status of `orders` stream:
```shell
nats str info ORDERS
```

To publish into `orders.processed` with storage acknowledgment:
```shell
nats req ORDERS.processed hello
```

To create an `orders processed` pull-based consumer and receive the message:
```shell
nats con add ORDERS PROCESSED --config config/orders-processed-pull-consumer.json
nats con next ORDERS PROCESSED
```