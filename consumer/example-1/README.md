# Consumer Example 1

This consumer runs with a events stream from NATS
(not nats streaming, just getting a stream of events from a NATS topic)

Also busy adding option to run with kafka...

To build and run, do:
```
cd consumer/example-1
docker-compose -f build.yml build
docker-compose -f build.yml up -d
```
The config.json indicates the topic. If you push events to that topic, the consumer will be triggered.

If you have NATS running locally, you can also go build and run the consumer in your console.
If not, the docker-compose file (build.yml) includes a NATS server.

When started with docker, the follow test shows that it works:
```
% curl -XGET http://localhost:8080/
Basic NATS based microservice example v0.0.1
% curl -XPOST http://localhost:8080/nats/request/123 -d'{}'
NATS request to topic(test) failed: nats: no responders available for request
% curl -XPOST http://localhost:8080/nats/publish/123 -d'{}'
NATS publish success. Duration(79.115µs)
% curl -XPOST http://localhost:8080/kafka/produce/123 -d'{}'
Kafka produce success. Duration(100.048µs)
```
The .../123 at the end of each URL is used for event or request type name. Any value will do for this example.
The /nats/request/123 is failing when consumer is not running.
When consumer is running, it will show response "Done!" comes from the consumer when it is running.

You can see consumer logs with:
```
$ docker logs example-1_consumer1_1
Registered constructor[nats] = &{URI: Topic:}
2021-05-05 21:07:12.742 DEBUG          source.go(   87): Get(consumer.stream)...
2021-05-05 21:07:12.743 DEBUG     nats-stream.go(   59): Connected to NATS(nats://nats:4222)
2021-05-05 21:07:12.743 DEBUG     nats-stream.go(   71): Subscribed to 'example-1' for processing requests...
2021-05-05 21:07:28.904 DEBUG     nats-stream.go(   62): Got task request on:%!(EXTRA string=example-1)
2021-05-05 21:07:28.904 DEBUG     nats-stream.go(   64): Sending reply to "_INBOX.LMD4InC5FiK6eTQrPt5YkO.6JWD3JQ6"
2021-05-05 21:07:33.576 DEBUG     nats-stream.go(   62): Got task request on:%!(EXTRA string=example-1)
2021-05-05 21:07:33.576 DEBUG     nats-stream.go(   67): NOT Sending reply
```

You can change the topic from `example-1` to another value in config and the build.yml file (use the same value in both).

## Kafka or NATS
The consumer can consumer either using config.
TBD...

## Faster Development Cycles
To run NATS and API in docker, comment out the consumer part in build.yml
Then do docker-compose up to start only nats and the API that allow you to puplish or produce events and send requests...
If you are improving the api, it too can be commented out of build.yml then run it on the command line with only NATS in docker.
Then build consumer locally and change as you like and run locally from the console
It makes the dev cycle a lot faster and allow you to import other libraries with replace statements in go.mod
which you cannot do in the docker builds.


## Testing with API
Publish an event with this on the console:
```
curl -XPOST 'http://localhost:8080/publish/hello' -d'{"name":"Janneman"}'
NATS publish success. Duration(121.558µs)
```
The event type is on the URL after /publish/, i.e. "hello" in this example.
The event data is in the POST body.
