package slack

import (
	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("slack", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new slack plugin instance.
func setup(c *caddy.Controller) error {
	var rules []Rule

	for c.Next() {
		var rule Rule
		args := c.RemainingArgs()
		switch len(args) {
		case 1:
			rule.Token = args[0]
		default:
			return c.ArgErr()
		}
		rules = append(rules, rule)
	}

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Slack{Rules: rules, Next: next}
	}
	cfg.AddMiddleware(mid)

	return nil
}
