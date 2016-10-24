package client

import (
	"fmt"
	"github.com/kohkimakimoto/xgrok/support/yaml-template"
	"github.com/kohkimakimoto/xgrok/xgrok"
	"github.com/kohkimakimoto/xgrok/xgrok/log"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Configuration struct {
	HttpProxy          string `yaml:"http_proxy,omitempty"`
	ServerAddr         string `yaml:"server_addr,omitempty"`
	InspectAddr        string `yaml:"inspect_addr,omitempty"`
	// I implemented xgrok to be used for self hosting. It should not use embedded crt file as a original 'ngrok'.
	// So it always uses host root cert.
	// TrustHostRootCerts bool                            `yaml:"trust_host_root_certs,omitempty"`
	InsecureSkipVerify bool                            `yaml:"insecure_skip_verify,omitempty"`
	Authtoken          string                          `yaml:"authtoken,omitempty"`
	Tunnels            map[string]*TunnelConfiguration `yaml:"tunnels,omitempty"`
	LogTo              string                          `yaml:"-"`
	Path               string                          `yaml:"-"`
}

type TunnelConfiguration struct {
	Subdomain  string            `yaml:"subdomain,omitempty"`
	Hostname   string            `yaml:"hostname,omitempty"`
	Protocols  map[string]string `yaml:"proto,omitempty"`
	HttpAuth   string            `yaml:"auth,omitempty"`
	RemotePort uint16            `yaml:"remote_port,omitempty"`
}

func LoadConfiguration(opts *Options) (config *Configuration, err error) {
	configPath := opts.Config
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	// deserialize/parse the config
	config = new(Configuration)
	if _, staterr := os.Stat(configPath); staterr == nil {
		log.Info("Reading configuration file %s", configPath)
		configBuf, ioerr := ioutil.ReadFile(configPath)
		if ioerr != nil {
			// failure to read a configuration file is only a fatal error if
			// the user specified one explicitly
			if opts.Config != "" {
				err = fmt.Errorf("Failed to read configuration file %s: %v", configPath, ioerr)
				return
			}
		}

		if ymlerr := template.UnmarshalWithEnv(configBuf, &config); ymlerr != nil {
			err = fmt.Errorf("Error parsing configuration file %s: %v", configPath, ymlerr)
			return
		}

		// xgrok does not need to support BC for ngrok

		//// try to parse the old .ngrok format for backwards compatibility
		//matched := false
		//content := strings.TrimSpace(string(configBuf))
		//if matched, err = regexp.MatchString("^[0-9a-zA-Z_\\-!]+$", content); err != nil {
		//	return
		//} else if matched {
		//	config = &Configuration{AuthToken: content}
		//}
	}

	// set configuration defaults
	if opts.ServerAddr == "" {
		if config.ServerAddr == "" {
			config.ServerAddr = xgrok.DefaultServerAddr
		}
	} else {
		config.ServerAddr = opts.ServerAddr
	}

	if opts.InspectAddr == "" {
		if config.InspectAddr == "" {
			config.InspectAddr = xgrok.DefaultInspectAddr
		}
	} else {
		config.InspectAddr = opts.InspectAddr
	}

	if config.HttpProxy == "" {
		config.HttpProxy = os.Getenv("http_proxy")
	}

	// validate and normalize configuration
	if config.InspectAddr != "disabled" {
		if config.InspectAddr, err = normalizeAddress(config.InspectAddr, "inspect_addr"); err != nil {
			return
		}
	}

	if config.ServerAddr, err = normalizeAddress(config.ServerAddr, "server_addr"); err != nil {
		return
	}

	if config.HttpProxy != "" {
		var proxyUrl *url.URL
		if proxyUrl, err = url.Parse(config.HttpProxy); err != nil {
			return
		} else {
			if proxyUrl.Scheme != "http" && proxyUrl.Scheme != "https" {
				err = fmt.Errorf("Proxy url scheme must be 'http' or 'https', got %v", proxyUrl.Scheme)
				return
			}
		}
	}

	for name, t := range config.Tunnels {
		if t == nil || t.Protocols == nil || len(t.Protocols) == 0 {
			err = fmt.Errorf("Tunnel %s does not specify any protocols to tunnel.", name)
			return
		}

		for k, addr := range t.Protocols {
			tunnelName := fmt.Sprintf("for tunnel %s[%s]", name, k)
			if t.Protocols[k], err = normalizeAddress(addr, tunnelName); err != nil {
				return
			}

			if err = validateProtocol(k, tunnelName); err != nil {
				return
			}
		}

		// use the name of the tunnel as the subdomain if none is specified
		if t.Hostname == "" && t.Subdomain == "" {
			// XXX: a crude heuristic, really we should be checking if the last part
			// is a TLD
			if len(strings.Split(name, ".")) > 1 {
				t.Hostname = name
			} else {
				t.Subdomain = name
			}
		}
	}

	// override configuration with command-line options
	config.LogTo = opts.Logto
	config.Path = configPath
	if opts.InsecureSkipVerify {
		config.InsecureSkipVerify = opts.InsecureSkipVerify
	}
	if opts.Authtoken != "" {
		config.Authtoken = opts.Authtoken
	}

	switch opts.Command {
	// start a single tunnel, the default, simple xgrok behavior
	case "tunnel":
		config.Tunnels = make(map[string]*TunnelConfiguration)
		config.Tunnels["default"] = &TunnelConfiguration{
			Subdomain: opts.Subdomain,
			Hostname:  opts.Hostname,
			HttpAuth:  opts.Httpauth,
			Protocols: make(map[string]string),
		}

		for _, proto := range strings.Split(opts.Protocol, "+") {
			if err = validateProtocol(proto, "default"); err != nil {
				return
			}

			if config.Tunnels["default"].Protocols[proto], err = normalizeAddress(opts.Args[0], ""); err != nil {
				return
			}
		}

	// list tunnels
	case "list":
		for name, _ := range config.Tunnels {
			fmt.Println(name)
		}
		os.Exit(0)

	// start tunnels
	case "start":
		if len(opts.Args) == 0 {
			err = fmt.Errorf("You must specify at least one tunnel to start")
			return
		}

		requestedTunnels := make(map[string]bool)
		for _, arg := range opts.Args {
			requestedTunnels[arg] = true

			if _, ok := config.Tunnels[arg]; !ok {
				err = fmt.Errorf("Requested to start tunnel %s which is not defined in the config file.", arg)
				return
			}
		}

		for name, _ := range config.Tunnels {
			if !requestedTunnels[name] {
				delete(config.Tunnels, name)
			}
		}

	case "start-all":
		return

	default:
		err = fmt.Errorf("Unknown command: %s", opts.Command)
		return
	}

	return
}

func DefaultConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Warn("Failed to get current working directory: %s", err.Error())
		return ""
	}

	return filepath.Join(wd, "xgrok.yml")

	//user, err := user.Current()
	//
	//// user.Current() does not work on linux when cross compiling because
	//// it requires CGO; use os.Getenv("HOME") hack until we compile natively
	//homeDir := os.Getenv("HOME")
	//if err != nil {
	//	log.Warn("Failed to get user's home directory: %s. Using $HOME: %s", err.Error(), homeDir)
	//} else {
	//	homeDir = user.HomeDir
	//}
	//
	//return path.Join(homeDir, ".xgrok")
}

func normalizeAddress(addr string, propName string) (string, error) {
	// normalize port to address
	if _, err := strconv.Atoi(addr); err == nil {
		addr = ":" + addr
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("Invalid address %s '%s': %s", propName, addr, err.Error())
	}

	if host == "" {
		host = "127.0.0.1"
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

func validateProtocol(proto, propName string) (err error) {
	switch proto {
	case "http", "https", "http+https", "tcp":
	default:
		err = fmt.Errorf("Invalid protocol for %s: %s", propName, proto)
	}

	return
}

//func SaveAuthToken(configPath, authtoken string) (err error) {
//	// empty configuration by default for the case that we can't read it
//	c := new(Configuration)
//
//	// read the configuration
//	oldConfigBytes, err := ioutil.ReadFile(configPath)
//	if err == nil {
//		// unmarshal if we successfully read the configuration file
//		if err = yaml.Unmarshal(oldConfigBytes, c); err != nil {
//			return
//		}
//	}
//
//	// no need to save, the authtoken is already the correct value
//	if c.AuthToken == authtoken {
//		return
//	}
//
//	// update auth token
//	c.AuthToken = authtoken
//
//	// rewrite configuration
//	newConfigBytes, err := yaml.Marshal(c)
//	if err != nil {
//		return
//	}
//
//	err = ioutil.WriteFile(configPath, newConfigBytes, 0600)
//	return
//}
