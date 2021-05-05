package main

import (
	"fmt"

	"github.com/stewelarend/config/source/configfile"
	"github.com/stewelarend/rpc"

	_ "github.com/stewelarend/rpc/server/http"
)

func main() {
	if err := configfile.Add("./conf/config.json"); err != nil {
		panic(err)
	}
	greeter := rpc.New("greeter")
	greeter.AddStruct("hello", helloRequest{})
	greeter.AddFunc("goodbye", goodbye)
	if err := greeter.Run(); err != nil {
		panic(err)
	}
}

type helloRequest struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Age8   int8   `json:"age8,omitempty"`
	Age16  int16  `json:"age16,omitempty"`
	Age32  int32  `json:"age32,omitempty"`
	Age64  int64  `json:"age64,omitempty"`
	Ageu8  uint8  `json:"ageu8,omitempty"`
	Ageu16 uint16 `json:"ageu16,omitempty"`
	Ageu32 uint32 `json:"ageu32,omitempty"`
	Ageu64 uint64 `json:"ageu64,omitempty"`
}

func (req helloRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("missing name")
	}
	return nil
}

func (req helloRequest) Exec(ctx rpc.IContext) (res interface{}, err error) {
	return fmt.Sprintf("Hi %+v!", req), nil
}

func goodbye(ctx rpc.IContext, req interface{}) (res interface{}, err error) {
	return "Cheers!", nil
}
