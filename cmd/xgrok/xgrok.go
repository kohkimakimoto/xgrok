package main

import (
	"fmt"
	"github.com/kohkimakimoto/xgrok/command"
	"github.com/kohkimakimoto/xgrok/xgrok"
	"github.com/urfave/cli"
	"os"
)

func main() {
	os.Exit(realMain())
}

func realMain() (status int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			status = 1
		}
	}()

	CLI := cli.NewApp()
	CLI.Name = xgrok.Name
	CLI.Usage = xgrok.Usage
	CLI.Version = xgrok.Version
	CLI.Commands = []cli.Command{
		command.TunnelCommand,
		command.ServeCommand,
	}

	if err := CLI.Run(os.Args); err != nil {
		status = 1
	}

	return status
}

func init() {
	cli.AppHelpTemplate = `Usage: {{.Name}}{{if .VisibleFlags}} [<options...>]{{end}} <command>

  {{.Name}}{{if .Usage}} -- {{.Usage}}{{end}}{{if .Version}}
  version {{.Version}}{{end}}{{if .Flags}}

Options:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}{{if .VisibleCommands}}
Commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{end}}
`
	cli.CommandHelpTemplate = `Usage: {{.Name}}{{if .VisibleFlags}} [<options...>]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[<arguments...>]{{end}}
{{if .Usage}}
  {{.Usage}}{{end}}
{{if .VisibleFlags}}
Options:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}{{if .Description}}
Description:
  {{.Description}}
{{end}}
`
}
