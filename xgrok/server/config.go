package server

import (
	"errors"
	"github.com/kohkimakimoto/xgrok/xgrok"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

type Configuration struct {
	TunnelAddr      string
	HttpAddr        string
	HttpsAddr       string
	Domain          string
	TlsCrt          string
	TlsKey          string
	DisableTcp      bool
	DisableHostname bool
	Logto           string `gluamapper:"-"`
	Loglevel        string `gluamapper:"-"`

	UserAuth UserAuthConfiguration `gluamapper:"-"`
	Hooks    HooksConfiguration    `gluamapper:"-"`
}

type UserAuthConfiguration struct {
	Enable    bool
	Tokens    []string
	TokensMap map[string]bool `gluamapper:"-"`
}

type HooksConfiguration struct {
	PreRegisterTunnel  *lua.LFunction
	PreOutputNewTunnel *lua.LFunction
	PostRegisterTunnel *lua.LFunction
	AuthResponseFilter *lua.LFunction
	PreShutdownTunnel  *lua.LFunction
	PostShutdownTunnel *lua.LFunction
}

func LoadConfiguration(opts *Options, L *lua.LState) (*Configuration, error) {
	config := &Configuration{
		TunnelAddr: xgrok.DefaultTunnelAddr,
		HttpAddr:   xgrok.DefaultHttpAddr,
		UserAuth: UserAuthConfiguration{
			Enable:    false,
			Tokens:    []string{},
			TokensMap: map[string]bool{},
		},
		Hooks: HooksConfiguration{},
	}

	configPath := opts.Config
	if configPath != "" {
		if err := L.DoFile(configPath); err != nil {
			return nil, err
		}

		if err := gluamapper.Map(L.GetGlobal("server").(*lua.LTable), config); err != nil {
			return nil, err
		}

		userAuthConfig := config.UserAuth

		if err := gluamapper.Map(L.GetGlobal("user_auth").(*lua.LTable), &userAuthConfig); err != nil {
			return nil, err
		}

		if userAuthConfig.Tokens != nil && len(userAuthConfig.Tokens) > 0 {
			for _, t := range userAuthConfig.Tokens {
				userAuthConfig.TokensMap[t] = true
			}
		}

		config.UserAuth = userAuthConfig

		hooksConfig := config.Hooks

		if err := gluamapper.Map(L.GetGlobal("hooks").(*lua.LTable), &hooksConfig); err != nil {
			return nil, err
		}

		config.Hooks = hooksConfig
	}

	// override configuration with command-line options
	config.Logto = opts.Logto
	config.Loglevel = opts.Loglevel

	if opts.HttpAddr != "" {
		config.HttpAddr = opts.HttpAddr
	}

	if opts.HttpsAddr != "" {
		config.HttpAddr = opts.HttpAddr
	}

	if opts.TunnelAddr != "" {
		config.TunnelAddr = opts.TunnelAddr
	}

	if opts.Domain != "" {
		config.Domain = opts.Domain
	}

	if config.Domain == "" {
		return nil, errors.New("xgrok server requires 'domain' config. try to set a value to the '--domain' option.")
	}

	if opts.TlsCrt != "" {
		config.TlsCrt = opts.TlsCrt
	}

	if opts.TlsKey != "" {
		config.TlsKey = opts.TlsKey
	}

	if opts.DisableTCP {
		config.DisableTcp = true
	}

	return config, nil
}
