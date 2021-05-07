package main

import (
	"fmt"

	"github.com/stewelarend/examples/smsc/submit/msg"
	"github.com/stewelarend/rpc"
)

func main() {
	submit := rpc.New("submit")
	submit.AddStruct("submit", SubmitRequest{})
	submit.Run()
}

type SubmitRequest msg.Message

func (req SubmitRequest) Exec(rpc.IContext) (interface{}, error) {
	return nil, fmt.Errorf("NYI")
}
