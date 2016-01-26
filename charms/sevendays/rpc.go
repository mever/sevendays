package sevendays

import (
	"fmt"
	"github.com/juju/gocharm/hook"
	"github.com/mever/sevendaystodie/httpd"
)

type RpcServer struct {
	cmd hook.Command
}

func (rpc *RpcServer) Set(s httpd.State, response *string) error {
	*response = s.User + " pong"
	fmt.Println(s)
	return nil
}
