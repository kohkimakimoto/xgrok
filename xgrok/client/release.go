// +build release

package client

var (
	rootCrtPaths = []string{"assets/client/tls/xgrokroot.crt"}
)

func useInsecureSkipVerify() bool {
	return false
}
