package command

import (
	"github.com/urfave/cli"
	"github.com/kohkimakimoto/xgrok/xgrok/client"
	"github.com/kohkimakimoto/xgrok/xgrok/server"
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
		Config:  ctx.String("config"),
		HttpAddr: ctx.String("http-addr"),
		HttpsAddr: ctx.String("https-addr"),
		TunnelAddr: ctx.String("tunnel-addr"),
		Domain: ctx.String("domain"),
		TlsCrt: ctx.String("tls-crt"),
		TlsKey: ctx.String("tls-key"),
		Logto: ctx.String("log"),
		Loglevel: ctx.String("log-level"),
	}

	return opts
}

var ClientFlags = []cli.Flag{
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

	// TODO:
	// authtoken, httpauth is not supported now...

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

	cli.StringFlag{
		Name:  "server-addr",
		EnvVar: "XGROK_SERVER_ADDR",
		Usage: "The xgrok server address to connet with (default: '127.0.0.1:4443').",
	},
}



func LoadClientOptions(ctx *cli.Context) *client.Options {
	opts := &client.Options{
		Logto: ctx.String("log"),
		Loglevel: ctx.String("log-level"),
		Subdomain:  ctx.String("subdomain"),
		Hostname: ctx.String("hostname"),
		Protocol: ctx.String("proto"),
		ServerAddr: ctx.String("server-addr"),
		Args: []string{},
	}

	return opts
}
