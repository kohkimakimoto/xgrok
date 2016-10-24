package command

import (
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/kohkimakimoto/xgrok/xgrok/server"
	"github.com/urfave/cli"
)

var ServeFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config",
		Usage: "Path to xgrok server configuration `file`.",
	},
	cli.StringFlag{
		Name:  "http-addr",
		Usage: "Public address for HTTP connections, empty string to disable",
	},
	cli.StringFlag{
		Name:  "https-addr",
		Usage: "Public address listening for HTTPS connections, emptry string to disable",
	},
	cli.StringFlag{
		Name:  "tunnel-addr",
		Usage: "Public address listening for xgrok client",
	},
	cli.StringFlag{
		Name:  "domain",
		Usage: "Domain where the tunnels are hosted",
	},
	cli.StringFlag{
		Name:  "tls-crt",
		Usage: "Path to a TLS certificate `file`",
	},
	cli.StringFlag{
		Name:  "tls-key",
		Usage: "Path to a TLS key `file`",
	},
	cli.BoolFlag{
		Name:  "disable-tcp",
		Usage: "disable TCP protocol proxy.",
	},
	cli.StringFlag{
		Name:  "log",
		Value: "stdout",
		Usage: "Write log messages to this file. 'stdout' and 'none' have special meanings",
	},
	cli.StringFlag{
		Name:  "log-level",
		Value: "INFO",
		Usage: "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR",
	},
}

func LoadServerOptions(ctx *cli.Context) *server.Options {
	opts := &server.Options{
		Config:     ctx.String("config"),
		HttpAddr:   ctx.String("http-addr"),
		HttpsAddr:  ctx.String("https-addr"),
		TunnelAddr: ctx.String("tunnel-addr"),
		Domain:     ctx.String("domain"),
		TlsCrt:     ctx.String("tls-crt"),
		TlsKey:     ctx.String("tls-key"),
		DisableTCP: ctx.Bool("disable-tcp"),
		Logto:      ctx.String("log"),
		Loglevel:   ctx.String("log-level"),
	}

	return opts
}

var ClientStartFlags = []cli.Flag{
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
	cli.StringFlag{
		Name:   "server-addr, s",
		EnvVar: "XGROK_SERVER_ADDR",
		Usage:  "The xgrok server address to connet with.",
	},
	cli.StringFlag{
		Name:   "inspect-addr",
		EnvVar: "XGROK_INSPECT_ADDR",
		Usage:  "The client inspect address.",
	},
	cli.BoolFlag{
		Name:   "insecure-skip-verify, i",
		EnvVar: "XGROK_INSECURE_SKIP_VERIFY",
		Usage:  "TLS accepts any certificate. This should be used only for testing.",
	},
	cli.StringFlag{
		Name:  "authtoken",
		Usage: "Authentication token for identifying an xgrok server account",
	},
}

var ClientTunnelFlags = append(ClientStartFlags, []cli.Flag{
	cli.StringFlag{
		Name:  "subdomain",
		Usage: "Request a custom subdomain from the xgrok server. (HTTP only)",
	},
	cli.StringFlag{
		Name:  "hostname",
		Usage: "Request a custom hostname from the xgrok server. (HTTP only) (requires CNAME of your DNS)",
	},
	cli.StringFlag{
		Name:  "proto",
		Value: "http",
		Usage: "The protocol of the traffic over the tunnel ('http', 'https', 'tcp', 'http+https')",
	},
}...)

func LoadClientOptions(ctx *cli.Context) *client.Options {
	opts := &client.Options{
		Logto:              ctx.String("log"),
		Loglevel:           ctx.String("log-level"),
		Subdomain:          ctx.String("subdomain"),
		Hostname:           ctx.String("hostname"),
		Protocol:           ctx.String("proto"),
		ServerAddr:         ctx.String("server-addr"),
		InspectAddr:        ctx.String("inspect-addr"),
		InsecureSkipVerify: ctx.Bool("insecure-skip-verify"),
		Authtoken:          ctx.String("authtoken"),
		Args:               []string{},
	}

	return opts
}
