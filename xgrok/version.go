package xgrok

var (
	Name       = "xgrok"
	Version    = "unknown"
	CommitHash = "unknown"
	Usage      = "Introspected tunnels to localhost."
)

var (
	DefaultTunnelAddr  = ":9690"
	DefaultHttpAddr    = ":9680"
	DefaultServerAddr  = "127.0.0.1" + DefaultTunnelAddr
	DefaultInspectAddr = "127.0.0.1:4040"
)
