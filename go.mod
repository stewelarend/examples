module github.com/stewelarend/examples

go 1.16

replace github.com/stewelarend/rpc => ../rpc

replace github.com/stewelarend/config => ../config

replace github.com/stewelarend/logger => ../logger

require (
	github.com/stewelarend/config v0.0.0-00010101000000-000000000000 // indirect
	github.com/stewelarend/rpc v0.0.0-00010101000000-000000000000
)
