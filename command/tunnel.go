package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
	"fmt"
	"errors"
)

var TunnelCommand = cli.Command{
	Name:   "tunnel",
	Usage:  "Start a tunnel",
	Action: tunnelAction,
	Flags: ClientFlags,
	Description: `Start a tunnel.

Example:
  xgrok tunnel 8080
  xgrok tunnel --subdomain=example 8080
`,
}

func tunnelAction(ctx *cli.Context) error {
	opts := LoadClientOptions(ctx)
	opts.Command = "tunnel"
	if ctx.NArg() == 0{
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
