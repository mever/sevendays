package sevendays

import (
	"github.com/juju/gocharm/hook"
	"fmt"
)

type RpcServer struct{
	cmd hook.Command
}

//type ServerState struct{}
//type Feedback struct{}

func (rpc *RpcServer) SetMessage(p string, response *string) error {
	*response = p + " pong"
	fmt.Println("Rpc server received: " + p)
	return nil
}