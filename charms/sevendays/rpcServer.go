package sevendays

import (
	"github.com/juju/gocharm/hook"
	"fmt"
	"github.com/mever/sevendaystodie/charms/sevendays/httpd"
)

type RpcServer struct{
	cmd hook.Command
}

func (rpc *RpcServer) Set(s httpd.State, response *string) error {
	*response = s.User + " pong"
	fmt.Println(s)
	return nil
}