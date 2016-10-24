package command

import (
	"fmt"
	"github.com/kohkimakimoto/xgrok/xgrok"
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

var InitCommand = cli.Command{
	Name:   "init",
	Usage:  "Create client config template.",
	Action: initAction,
}

func initAction(ctx *cli.Context) error {
	p := client.DefaultConfigPath()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := ioutil.WriteFile(p, []byte(configContent), 0644); err != nil {
			return err
		}

		fmt.Printf("Created '%s'\n", p)
	} else {
		return fmt.Errorf("%s is already exists", p)
	}

	return nil
}

var configContent = `
server_addr:           "` + xgrok.DefaultServerAddr + `"
inspect_addr:          "` + xgrok.DefaultInspectAddr + `"
insecure_skip_verify:  true
tunnels:
  foo:
    subdomain:  "foo"
    proto:      {"http": "8000"}
  # bar:
    #   subcomain:
    #   hostname:
    #   remote_port:
    #   proto:
`
