package main

import (
	"fmt"
	"time"

	"github.com/stewelarend/config/source/configfile"
	"github.com/stewelarend/logger"
	"github.com/stewelarend/rpc"

	_ "github.com/stewelarend/rpc/server/nats"
)

var log = logger.New("nats-server")

func main() {
	if err := configfile.Add("./conf/config.json"); err != nil {
		panic(err)
	}
	greeter := rpc.New("test")
	greeter.AddStruct("echo", echoRequest{})
	if err := greeter.Run(); err != nil {
		panic(err)
	}
}

type echoRequest struct {
	Message      string `json:"message" doc:"This message is echoed in the response"`
	DelaySeconds int    `json:"delay_seconds" doc:"Wait this nr of seconds before responding"`
}

type echoResponse struct {
	Message string `json:"message" doc:"This message is copied from the request"`
}

func (req echoRequest) Validate() error {
	if req.DelaySeconds < 0 {
		return fmt.Errorf("negative delay")
	}
	return nil
}

func (req echoRequest) Exec(ctx rpc.IContext) (res interface{}, err error) {
	time.Sleep(time.Duration(req.DelaySeconds) * time.Second)
	return echoResponse{Message: req.Message}, nil
}
