package client

import (
	"flag"
	"fmt"
	"github.com/kohkimakimoto/xgrok/xgrok/version"
	"os"
)

const usage1 string = `Usage: %s [OPTIONS] <local port or address>
Options:
`

const usage2 string = `
Examples:
	xgrok 80
	xgrok -subdomain=example 8080
	xgrok -proto=tcp 22
	xgrok -hostname="example.com" -httpauth="user:password" 10.0.0.1


Advanced usage: xgrok [OPTIONS] <command> [command args] [...]
Commands:
	xgrok start [tunnel] [...]    Start tunnels by name from config file
	ngork start-all               Start all tunnels defined in config file
	xgrok list                    List tunnel names from config file
	xgrok help                    Print help
	xgrok version                 Print xgrok version

Examples:
	xgrok start www api blog pubsub
	xgrok -log=stdout -config=xgrok.yml start ssh
	xgrok start-all
	xgrok version

`

type Options struct {
	Config    string
	Logto     string
	Loglevel  string
	Authtoken string
	Httpauth  string
	Hostname  string
	Protocol  string
	ServerAddr string
	Subdomain string
	Command   string
	Args      []string
}







func ParseArgs() (opts *Options, err error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage2)
	}

	config := flag.String(
		"config",
		"",
		"Path to xgrok configuration file. (default: $HOME/.xgrok)")

	logto := flag.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings")

	loglevel := flag.String(
		"log-level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")

	authtoken := flag.String(
		"authtoken",
		"",
		"Authentication token for identifying an xgrok.com account")

	httpauth := flag.String(
		"httpauth",
		"",
		"username:password HTTP basic auth creds protecting the public tunnel endpoint")

	subdomain := flag.String(
		"subdomain",
		"",
		"Request a custom subdomain from the xgrok server. (HTTP only)")

	hostname := flag.String(
		"hostname",
		"",
		"Request a custom hostname from the xgrok server. (HTTP only) (requires CNAME of your DNS)")

	protocol := flag.String(
		"proto",
		"http+https",
		"The protocol of the traffic over the tunnel {'http', 'https', 'tcp'} (default: 'http+https')")

	flag.Parse()

	opts = &Options{
		Config:    *config,
		Logto:     *logto,
		Loglevel:  *loglevel,
		Httpauth:  *httpauth,
		Subdomain: *subdomain,
		Protocol:  *protocol,
		Authtoken: *authtoken,
		Hostname:  *hostname,
		Command:   flag.Arg(0),
	}

	switch opts.Command {
	case "list":
		opts.Args = flag.Args()[1:]
	case "start":
		opts.Args = flag.Args()[1:]
	case "start-all":
		opts.Args = flag.Args()[1:]
	case "version":
		fmt.Println(version.MajorMinor())
		os.Exit(0)
	case "help":
		flag.Usage()
		os.Exit(0)
	case "":
		err = fmt.Errorf("Error: Specify a local port to tunnel to, or " +
			"an xgrok command.\n\nExample: To expose port 80, run " +
			"'xgrok 80'")
		return

	default:
		if len(flag.Args()) > 1 {
			err = fmt.Errorf("You may only specify one port to tunnel to on the command line, got %d: %v",
				len(flag.Args()),
				flag.Args())
			return
		}

		opts.Command = "default"
		opts.Args = flag.Args()
	}

	return
}
