package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
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
	client.Main(opts)
	return nil
}
