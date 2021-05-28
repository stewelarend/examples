# DONE
## 2021-05-08
* RPC test send 1 and send 100 works
* User message struct used in server (see rpc-nats_test.go)

# NEXT
* server/consumer test/control concurrency and multiple 
* controlled shutdown in test, server, consumer, ...
* run multiple servers/consumers in one process, with shared channel for processing
* consumer
    * change to use user message structs, move message.go to examples
    * kafka
* stats
* logs
* audits
* partitioning over instances
* example system with docker compose
* templates to create system and services