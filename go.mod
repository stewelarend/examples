module github.com/stewelarend/examples

go 1.16

replace github.com/stewelarend/rpc => ../rpc

replace github.com/stewelarend/consumer => ../consumer

require (
	github.com/stewelarend/config v0.0.2
	github.com/stewelarend/rpc v0.0.1
	github.com/stewelarend/util v0.0.1 // indirect
)
