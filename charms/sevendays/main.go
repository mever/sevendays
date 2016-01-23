package sevendays

import (
	"github.com/juju/gocharm/charmbits/service"
	"github.com/juju/gocharm/hook"
	"gopkg.in/juju/charm.v5"
)

type Charm struct{
	svc service.Service
	ctx *hook.Context
}

func RegisterHooks(r *hook.Registry) {
	c := Charm{}
	r.RegisterContext(c.setContext, nil)
	r.RegisterHook("install", c.install)
	r.RegisterHook("upgrade-charm", c.install)
	r.RegisterHook("config-changed", c.configChanged)
	c.svc.Register(r.Clone("service"), "", c.startService)

	r.RegisterConfig("message", charm.Option{
		Type:        "string",
		Description: "The massage shown in the deamon",
		Default:     "Hello World",
	})
}

func (c *Charm) configChanged() error {
	reply := ""
	msg, _ := c.ctx.GetConfigString("message")
	err := c.svc.Call("RpcServer.SetMessage", msg, &reply)
	c.ctx.Logf(reply)

	c.svc.Start(msg)

	msg, _ = c.ctx.GetConfigString("message")
	err = c.svc.Call("RpcServer.SetMessage", msg + " new", &reply)
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