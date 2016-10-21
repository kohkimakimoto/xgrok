package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
)

var StartCommand = cli.Command{
	Name:   "start",
	Usage:  "Start tunnels by name from config file",
	Action: startAction,
	Flags: ClientFlags,
	ArgsUsage: "[<tunnel...>]",
	Description: `Start tunnels by name from config file.

Example:
  xgrok start www api blog pubsub
  xgrok --log=stdout --config=ngrok.yml start ssh
`,
}

func startAction(ctx *cli.Context) error {
	opts := LoadClientOptions(ctx)
	opts.Command = "start"
	opts.Args = ctx.Args()

	client.Main(opts)
	return nil
}
