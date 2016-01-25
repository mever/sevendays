package sevendays

import (
	"github.com/juju/gocharm/charmbits/service"
	"github.com/juju/gocharm/hook"
	"github.com/mever/sevendaystodie/charms/sevendays/httpd"
	"gopkg.in/juju/charm.v5"
)

type Charm struct {
	svc service.Service
	ctx *hook.Context
}

func RegisterHooks(r *hook.Registry) {
	c := Charm{}
	r.RegisterContext(c.setContext, nil)
	r.RegisterHook("start", c.start)
	r.RegisterHook("upgrade-charm", c.start)
	r.RegisterHook("config-changed", c.configChanged)
	c.svc.Register(r.Clone("service"), "", c.startService)

	r.RegisterConfig("user", charm.Option{
		Type:        "string",
		Description: "Provide Steam user name",
		Default:     "",
	})
}

func (c *Charm) start() error {
	return c.svc.Start(c.ctx.CharmDir)
}

func (c *Charm) configChanged() error {
	user, _ := c.ctx.GetConfigString("user")
	ss := httpd.State{
		User: user,
	}

	var reply string
	err := c.svc.Call("RpcServer.Set", ss, &reply)
	c.ctx.Logf(reply)
	return err
}

func (c *Charm) startService(ctx *service.Context, args []string) (hook.Command, error) {
	return startServer(ctx, args)
}

func (c *Charm) setContext(ctx *hook.Context) error {
	c.ctx = ctx
	return nil
}
