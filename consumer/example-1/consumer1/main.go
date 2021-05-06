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
	greeter.AddStruct("hello", helloRequest{}) //   ...this routing is not yet supported and handlers are not yet called from consumer lib! ....
	greeter.AddFunc("goodbye", goodbye)
	if err := greeter.Run(); err != nil {
		panic(err)
	}
}

type helloRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//consumer return nil to ack the message or error to put it back on the queue
func (req helloRequest) Exec(ctx consumer.IContext) (err error) {
	fmt.Printf("HELLO: %+v\n", req)
	return nil
}

func goodbye(ctx consumer.IContext, req interface{}) (err error) {
	fmt.Printf("GOODBYE: %+v\n", req)
	return nil
}
