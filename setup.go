package alias

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/caddyserver/caddy"
)

func init() {
	caddy.RegisterPlugin("alias", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error("alias", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Alias{Next: next}
	})

	return nil
}
