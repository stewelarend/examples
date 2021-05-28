package main

import (
	"fmt"

	"github.com/stewelarend/config"
	"github.com/stewelarend/config/source/configfile"
	"github.com/stewelarend/consumer"
	_ "github.com/stewelarend/consumer/stream/nats"
	"github.com/stewelarend/controller"
	"github.com/stewelarend/logger"
)

var log = logger.New()

func main() {
	log.SetLevel(logger.LevelDebug)
	logger.SetGlobalLevel(logger.LevelDebug)

	//load config from file
	if err := configfile.Add("./conf/config.json"); err != nil {
		panic(fmt.Errorf("failed to load config file: %v", err))
	}

	//define the consumer events and handlers
	greeter := consumer.New("greeter")
	greeter.AddStruct("hello", helloRequest{}) //   ...this routing is not yet supported and handlers are not yet called from consumer lib! ....
	greeter.AddFunc("goodbye", goodbye)

	//define the handler that will decode the stream messages
	//this should be an agreement between the producer and consumer
	//so ideally this is a library shared by both
	//or based on a shared protocol specification if developed independantly...
	//when it decoded a message it will send it to the consumer for processing
	handler := handler{
		consumer: greeter,
	}

	//get stream config
	streamName, streamConfig, ok := config.GetNamed("consumer.stream")
	if !ok {
		panic(fmt.Errorf("missing consumer.stream config"))
	}

	log.Debugf("consumer.stream: %s: %+v", streamName, streamConfig)

	//connect to the event stream using config
	//this allows the consumer to run with kafka or NATS or other streams you add...
	stream, err := consumer.NewStream(streamName, streamConfig.(map[string]interface{}))
	if err != nil {
		panic(fmt.Errorf("failed to create configured stream: %v", err))
	}

	//run the controller
	if err := controller.Run(
		controller.Config{},
		stream,
		handler,
	); err != nil {
		panic(err)
	}
}

type helloRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//consumer return nil to ack the message or error to put it back on the queue
func (req helloRequest) Exec(ctx controller.Context) (err error) {
	fmt.Printf("HELLO: %+v\n", req)
	return nil
}

func goodbye(ctx controller.Context, req interface{}) (err error) {
	fmt.Printf("GOODBYE: %+v\n", req)
	return nil
}

type handler struct {
	consumer consumer.IConsumer
}

func (h handler) Handle(ctx controller.Context, event []byte) error {
	log.Debugf("Not yet handling message: %v", string(event))
	return fmt.Errorf("NYI handler...")
}
