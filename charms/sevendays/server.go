package sevendays

import (
	"github.com/juju/gocharm/charmbits/service"
	"github.com/juju/gocharm/hook"
	"gopkg.in/errgo.v1"
	"gopkg.in/tomb.v2"
	"net/http"
	"strconv"
	"sync"
)

type server struct {
	mu   sync.Mutex
	tomb tomb.Tomb
	rpc  RpcServer
}

func startServer(ctx *service.Context, args []string) (hook.Command, error) {
	var err error
	srv := &server{}
	srv.rpc.cmd, err = ctx.ServeLocalRPC(&srv.rpc)
	if err != nil {
		return nil, errgo.Notef(err, "serve local RPC failed")
	}

	srv.bindTerm(srv.rpc.cmd)
	srv.run()
	return srv, nil
}

var counter = 0

func (srv *server) run() {
	srv.tomb.Go(func() error {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World "))
			w.Write([]byte(strconv.Itoa(counter)))
			counter++
		})

		return http.ListenAndServe(":80", nil)
	})
}

// bindTerm ensures mutual termination of this server
// and the given command.
func (srv *server) bindTerm(cmd hook.Command) {
	srv.tomb.Go(func() error {
		srv.tomb.Kill(cmd.Wait())
		return nil
	})
	srv.tomb.Go(func() error {
		<-srv.tomb.Dying()
		// The command has been killed. Stop the RPC listener
		// and close any current listeners to cause the handlers
		// to terminate.
		cmd.Kill()
		err := cmd.Wait()
		if err != nil {
			err = errgo.Notef(err, "local RPC server")
		}

		//		srv.mu.Lock()
		//		defer srv.mu.Unlock()
		//		srv.closeResources()
		return err
	})
}

func (s *server) Kill() {
	s.tomb.Kill(nil)
}

func (s *server) Wait() error {
	return s.tomb.Wait()
}
