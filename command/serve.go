package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/server"
	"github.com/urfave/cli"
)

var ServeCommand = cli.Command{
	Name:   "serve",
	Usage:  "Run xgrok server.",
	Action: serveAction,
}

func serveAction(ctx *cli.Context) error {
	server.Main()

	return nil
}
