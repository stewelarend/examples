package main

import (
	"github.com/stewelarend/rpc"

	_ "github.com/stewelarend/rpc/http"
)

func main() {
	greeter := rpc.New("greeter")
	greeter.AddFunc("hello", hello)
	greeter.AddFunc("goodbye", goodbye)
	if err := greeter.Run(); err != nil {
		panic(err)
	}
}

func hello(ctx rpc.IContext, req interface{}) (res interface{}, err error) {
	return "Hi!", nil
}

func goodbye(ctx rpc.IContext, req interface{}) (res interface{}, err error) {
	return "Cheers!", nil
}
