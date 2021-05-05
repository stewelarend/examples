package main

import (
	"github.com/stewelarend/consumer"
)

func main() {
	greeter := consumer.New("greeter")
	greeter.AddFunc("hello", hello)
	greeter.AddFunc("goodbye", goodbye)
	if err := greeter.Run(); err != nil {
		panic(err)
	}
}

//consumer return nil to ack the message or error to put it back on the queue
func hello(ctx rpc.IContext, req interface{}) (err error) {
	return "Hi!", nil
}

func goodbye(ctx rpc.IContext, req interface{}) (res interface{}, err error) {
	return "Cheers!", nil
}
