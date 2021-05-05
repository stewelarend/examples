# Consumer Example 1

This consumer runs with a events stream from NATS
(not nats streaming, just getting a stream of events from a NATS topic)

To build and run, do:
```
cd consumer/exampl1
docker-compose -f build.yml build
docker-compose -f build.yml up -d
```
The config.json indicates the topic. If you push events to that topic, the consumer will be triggered.

If you have NATS running locally, you can also go build and run the consumer in your console.
If not, the docker-compose file (build.yml) includes a NATS server.

When started with docker, the follow test shows that it works:
```
example-1$ curl -XGET http://localhost:8080/
Basic NATS based microservice example v0.0.1
example-1$ curl -XGET http://localhost:8080/request
NATS request success. Duration(486.042µs) Response: Done!
example-1$ curl -XGET http://localhost:8080/publish
NATS publish success. Duration(1.431µs)
```
The response "Done!" comes from the consumer.
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
