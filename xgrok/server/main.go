package server

import (
	"crypto/tls"
	"fmt"
	"github.com/kohkimakimoto/xgrok/xgrok/conn"
	log "github.com/kohkimakimoto/xgrok/xgrok/log"
	"github.com/kohkimakimoto/xgrok/xgrok/msg"
	"github.com/kohkimakimoto/xgrok/xgrok/util"
	"github.com/yuin/gopher-lua"
	"math/rand"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
	"sync"
)

const (
	registryCacheSize uint64        = 1024 * 1024 // 1 MB
	connReadTimeout   time.Duration = 10 * time.Second
)

// GLOBALS
var (
	LState          *lua.LState
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	config    *Configuration
	listeners map[string]*conn.Listener
	hooksMutex *sync.Mutex
)

func NewProxy(pxyConn conn.Conn, regPxy *msg.RegProxy) {
	// fail gracefully if the proxy connection fails to register
	defer func() {
		if r := recover(); r != nil {
			pxyConn.Warn("Failed with error: %v", r)
			pxyConn.Close()
		}
	}()

	// set logging prefix
	pxyConn.SetType("pxy")

	// look up the control connection for this proxy
	pxyConn.Info("Registering new proxy for %s", regPxy.ClientId)
	ctl := controlRegistry.Get(regPxy.ClientId)

	if ctl == nil {
		panic("No client found for identifier: " + regPxy.ClientId)
	}

	ctl.RegisterProxy(pxyConn)
}

func startTunnelListener(addr string, tlsConfig *tls.Config) {
	go tunnelListener(addr, tlsConfig)
}

// Listen for incoming control and proxy connections
// We listen for incoming control and proxy connections on the same port
// for ease of deployment. The hope is that by running on port 443, using
// TLS and running all connections over the same port, we can bust through
// restrictive firewalls.
func tunnelListener(addr string, tlsConfig *tls.Config) {
	// listen for incoming connections
	listener, err := conn.Listen(addr, "tun", tlsConfig)
	if err != nil {
		panic(err)
	}

	log.Info("Listening for control and proxy connections on %s", listener.Addr.String())
	for c := range listener.Conns {
		go func(tunnelConn conn.Conn) {
			// don't crash on panics
			defer func() {
				if r := recover(); r != nil {
					tunnelConn.Info("tunnelListener failed with error %v: %s", r, debug.Stack())
				}
			}()

			tunnelConn.SetReadDeadline(time.Now().Add(connReadTimeout))
			var rawMsg msg.Message
			if rawMsg, err = msg.ReadMsg(tunnelConn); err != nil {
				tunnelConn.Warn("Failed to read message: %v", err)
				tunnelConn.Close()
				return
			}

			// don't timeout after the initial read, tunnel heartbeating will kill
			// dead connections
			tunnelConn.SetReadDeadline(time.Time{})

			switch m := rawMsg.(type) {
			case *msg.Auth:
				NewControl(tunnelConn, m)
			case *msg.RegProxy:
				NewProxy(tunnelConn, m)
			default:
				tunnelConn.Close()
			}
		}(c)
	}
}

func waitHandlingSignals() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			log.Info("Received SIGINT")
			for _, t := range tunnelRegistry.tunnels {
				t.Shutdown()
			}

			return
		case syscall.SIGTERM:
			log.Info("Received SIGTERM")
			for _, t := range tunnelRegistry.tunnels {
				t.Shutdown()
			}

			return
		}
	}
}

func wait() {
	waitHandlingSignals()

	// workaround: wait to output log
	// Sometimes the logger does not output some log messages before the process exits.
	time.Sleep(1 * time.Millisecond)
}

func Main(opts *Options) {
	log.LogTo(opts.Logto, opts.Loglevel)

	// init lua state
	LState = lua.NewState()
	defer LState.Close()
	initLuaState(LState)

	hooksMutex = new(sync.Mutex)

	// read configuration file
	c, err := LoadConfiguration(opts, LState)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config = c

	log.Debug("server config: %v", config)

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		panic(err)
	}
	rand.Seed(seed)

	// init tunnel/control registry
	registryCacheFile := os.Getenv("REGISTRY_CACHE_FILE")
	tunnelRegistry = NewTunnelRegistry(registryCacheSize, registryCacheFile)
	controlRegistry = NewControlRegistry()

	// start listeners
	listeners = make(map[string]*conn.Listener)

	// load tls configuration
	if config.TlsCrt == "" {
		log.Warn("Using bundled snakeoil certificate, so you have a risk about secure networking of the tunnel. Check configuration 'tls_crt' and 'tls_key'.")
	}

	tlsConfig, err := LoadTLSConfig(config.TlsCrt, config.TlsKey)
	if err != nil {
		panic(err)
	}

	// listen for http
	if config.HttpAddr != "" {
		listeners["http"] = startHttpListener(config.HttpAddr, nil)
	}

	// listen for https
	if config.HttpsAddr != "" {
		listeners["https"] = startHttpListener(config.HttpsAddr, tlsConfig)
	}

	startTunnelListener(config.TunnelAddr, tlsConfig)
	wait()
}
