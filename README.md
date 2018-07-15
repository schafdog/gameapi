
# Introduction

This project implements the specified API for SYBO Games test in Go and uses Cassandra DB for storage.

The Database is specified through a Interface and thus can be replace by other implementations.
The REST API is implemented in the user package (should have been renamed).

The service is unit tested using httptest and a small shell script implements a integration test.

# Installation

Both Cassandra and the API can be deployed using a existing kubernetes cluster:
## Cassandra

> kubectl create -f local-storages.yaml 
> kubectl create -f cassandra-statefulset.yaml
> kubectl create -f cassamdra-service.yaml

When the cassandra cluster is running we need to initialize the keyspace and setup the table

> kubectl cp keyspace.cql cassandra-0:keyspace.cql
> kubectl cp keyspace.cql cassandra-0:user.cql

> kubectl exec -it cassandra-0 /bin/bash
> cat keyspace.cql |cqlsh
> cat user.cql |cqlsh
> exit

## GAMEAPI 
It can be deployed using the following commands:

> kubectl create -f gameapi-deployment.yaml
> kubectl expose gameapi --type LoadBalancer --external-ip=<your-external-ip>
