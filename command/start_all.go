package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
)

var StartAllCommand = cli.Command{
	Name:   "start-all",
	Usage:  "Start all tunnels defined in config file",
	Action: startAllAction,
	Flags: ClientFlags,
	ArgsUsage: " ",
	Description: `Start all tunnels defined in config file.

Example:
  xgrok start-all
`,
}

func startAllAction(ctx *cli.Context) error {
	opts := LoadClientOptions(ctx)
	opts.Command = "start-all"
	opts.Args = ctx.Args()

	client.Main(opts)
	return nil
}
