module github.com/stewelarend/examples

go 1.16

replace github.com/stewelarend/rpc => ../rpc

replace github.com/stewelarend/consumer => ../consumer

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/nats-io/nats-server/v2 v2.2.2 // indirect
	github.com/stewelarend/config v0.0.2
	github.com/stewelarend/consumer v0.0.0-00010101000000-000000000000 // indirect
	github.com/stewelarend/rpc v0.0.1
)
