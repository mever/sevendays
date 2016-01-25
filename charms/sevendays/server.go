package sevendays

import (
	"github.com/juju/gocharm/charmbits/service"
	"github.com/juju/gocharm/hook"
	"github.com/mever/sevendaystodie/charms/sevendays/httpd"
	"golang.org/x/net/context"
	"gopkg.in/errgo.v1"
	"sync"
)

type server struct {
	mu   sync.Mutex
	rpc  RpcServer
	ctx  context.Context
	stop context.CancelFunc
}

func startServer(ctx *service.Context, args []string) (hook.Command, error) {
	var (
		err error
		cmd hook.Command
		srv = &server{}
	)

	cmd, err = ctx.ServeLocalRPC(&srv.rpc)
	if err != nil {
		return nil, errgo.Notef(err, "serve local RPC failed")
	}

	srv.ctx, srv.stop = context.WithCancel(context.Background())
	go func() {
		<-srv.ctx.Done()
		cmd.Kill()
		cmd.Wait()
	}()

	if len(args) > 0 {
		httpd.AssetsDir = args[0] + "/assets"
	}

	return srv, httpd.Serve(srv.ctx)
}

func (s *server) Kill() {
	s.stop()
}

func (s *server) Wait() error {
	<-s.ctx.Done()
	return s.ctx.Err()
}
