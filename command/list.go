package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
)

var ListCommand = cli.Command{
	Name:      "list",
	Usage:     " List tunnel names from config file",
	Action:    listAction,
	Flags:     []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "Path to xgrok client configuration `file` (default: $(pwd)/.xgrok.yml).",
		},
		cli.StringFlag{
			Name:  "log",
			Value: "none",
			Usage: "Write log messages to this file. 'stdout' and 'none' have special meanings",
		},
		cli.StringFlag{
			Name:  "log-level",
			Value: "DEBUG",
			Usage: "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR",
		},
	},
}

func listAction(ctx *cli.Context) error {
	opts := LoadClientOptions(ctx)
	opts.Command = "list"
	opts.Args = ctx.Args()

	client.Main(opts)
	return nil
}
