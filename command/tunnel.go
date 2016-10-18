package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
)

var TunnelCommand = cli.Command{
	Name:   "tunnel",
	Usage:  "start a tunnel",
	Action: tunnelAction,
}

func tunnelAction(ctx *cli.Context) error {
	client.Main()

	return nil
}
