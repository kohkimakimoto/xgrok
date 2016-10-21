package command

import (
	"errors"
	"fmt"
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
)

var TunnelCommand = cli.Command{
	Name:      "tunnel",
	Usage:     "Start a tunnel",
	Action:    tunnelAction,
	Flags:     ClientFlags,
	ArgsUsage: "<port>",
	Description: `Start a tunnel.

Example:
  xgrok tunnel 8000                       # forward xgrok server subdomain to local port 8000
  xgrok tunnel --subdomain=bar 8000       # request subdomain name: 'bar.your_server'
  xgrok tunnel --hostname=ex.com 8000     # request tunnel 'ex.com' (DNS CNAME)
`,
}

func tunnelAction(ctx *cli.Context) error {
	opts := LoadClientOptions(ctx)
	opts.Command = "tunnel"
	if ctx.NArg() == 0 {
		return errors.New("Error: Specify a local port to tunnel to, or an xgrok command.\n\nExample: To expose port 80, run 'xgrok tunnel 80'")
	}

	if ctx.NArg() > 1 {
		return fmt.Errorf("You may only specify one port to tunnel to on the command line, got %d: %v",
			ctx.NArg(),
			ctx.Args())
	}
	opts.Args = ctx.Args()

	client.Main(opts)
	return nil
}
