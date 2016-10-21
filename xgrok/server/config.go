package server

import (
	"github.com/kohkimakimoto/xgrok/support/yaml-template"
	"io/ioutil"
)

type Configuration struct {
	HttpAddr   string `yaml:"http_addr,omitempty"`
	HttpsAddr  string `yaml:"https_addr,omitempty"`
	TunnelAddr string `yaml:"tunnel_addr,omitempty"`
	Domain     string `yaml:"domain,omitempty"`
	TlsCrt     string `yaml:"tls_crt,omitempty"`
	TlsKey     string `yaml:"tls_key,omitempty"`
	Logto      string `yaml:"-"`
	Loglevel   string `yaml:"-"`
}

func LoadConfiguration(opts *Options) (*Configuration, error) {
	config := &Configuration{
		HttpAddr:   ":80",
		TunnelAddr: ":4443",
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

	if opts.TlsCrt != "" {
		config.TlsCrt = opts.TlsCrt
	}

	if opts.TlsKey != "" {
		config.TlsKey = opts.TlsKey
	}

	return config, nil
}
