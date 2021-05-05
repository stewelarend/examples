package main

import (
	"fmt"

	"github.com/stewelarend/config/source/configfile"
	"github.com/stewelarend/consumer"
	_ "github.com/stewelarend/consumer/stream/nats"
)

func main() {
	if err := configfile.Add("./conf/config.json"); err != nil {
		panic(err)
	}
	greeter := consumer.New("greeter")
	greeter.AddFunc("hello", hello)
	greeter.AddFunc("goodbye", goodbye)
	if err := greeter.Run(); err != nil {
		panic(err)
	}
}

//consumer return nil to ack the message or error to put it back on the queue
func hello(ctx consumer.IContext, req interface{}) (err error) {
	fmt.Printf("HELLO: %+v\n", req)
	return nil
}

func goodbye(ctx consumer.IContext, req interface{}) (err error) {
	fmt.Printf("GOODBYE: %+v\n", req)
	return nil
}
