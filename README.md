# Micro-Service Examples
In this repo you will find a few examples of micro-services developed with the libraries of github.com/stewelarend

# Overview of Examples
## rpc/example-1
A very basic service with two operations running on the default HTTP server port 8000

## rpc/example-2
Same as above, but uses a request structure with validation and different types of fields.
Also demonstrates how server config is loaded from file.

## rpc/example-batch
TODO
Replace the HTTP server with another server implementation, e.g. here requests are read from a batch file and responses written to stdout.

## example-logs
TODO
Set stats, write audits and write logs from the service

## example-call
TODO
Call another service

## example-timer
TODO
With either in memory or persistent db

## example-docker
TODO
Build the service into a docker container and run with docker-compose

## example-kubernetes
TODO
Deploy the docker image as a kubernetes pod

## example CRUD
TODO
An RPC service for queries and a consumer for updates, both deployed in kubernetes

## example-rpc-scaling
    - Register
    - Partition
    - Load-balance

## example-discover
Discover another service from config
    - see what is running
    - see partitioning


## example-consumer-simple
Run a simple consumer with producer utility
    - Use kafka
    - Use rabbitmq

## example-consumer-scaling
Scale the consumer up/down at run-time with partitioning
    - show metrics to see how busy and lag
    - use metrics to scale up/down
    - show partitioning
    - show serialization in workers

## example-db-mysql

## example-db-mongo

## example-db-postgresql

## example-schema
Show schema registration and retrieval
    - multiple versions
    - persistent storage
    - admin to cleanup
