package server

import (
	"errors"
	"github.com/kohkimakimoto/xgrok/support/yaml-template"
	"github.com/kohkimakimoto/xgrok/xgrok"
	"io/ioutil"
)

type Configuration struct {
	TunnelAddr    string            `yaml:"tunnel_addr,omitempty"`
	HttpAddr      string            `yaml:"http_addr,omitempty"`
	HttpsAddr     string            `yaml:"https_addr,omitempty"`
	Domain        string            `yaml:"domain,omitempty"`
	TlsCrt        string            `yaml:"tls_crt,omitempty"`
	TlsKey        string            `yaml:"tls_key,omitempty"`
	DisableTCP    bool              `yaml:"disable_tcp,omitempty"`
	UseAuth       bool              `yaml:"use_auth,omitempty"`
	Authtokens    []string          `yaml:"authtokens,omitempty"`
	AuthtokensMap map[string]bool    `yaml:"-"`
	PreRegister []string           `yaml:"pre_register,omitempty"`
	PostRegister []string           `yaml:"post_register,omitempty"`
	PreShutdown []string           `yaml:"pre_shutdown,omitempty"`
	PostShutdown []string           `yaml:"post_shutdown,omitempty"`
	Logto         string            `yaml:"-"`
	Loglevel      string            `yaml:"-"`
}

func LoadConfiguration(opts *Options) (*Configuration, error) {
	config := &Configuration{
		TunnelAddr: xgrok.DefaultTunnelAddr,
		HttpAddr:   xgrok.DefaultHttpAddr,
		AuthtokensMap: map[string]bool{},
	}

	configPath := opts.Config
	if configPath != "" {
		configBuf, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := template.UnmarshalWithEnv(configBuf, &config); err != nil {
			return nil, err
		}
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
		config.DisableTCP = true
	}

	if config.Authtokens != nil && len(config.Authtokens) > 0{
		for _, t := range config.Authtokens {
			config.AuthtokensMap[t] = true
		}
	}

	return config, nil
}
