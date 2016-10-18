// +build !release

package client

var (
	rootCrtPaths = []string{"assets/client/tls/xgrokroot.crt", "assets/client/tls/snakeoilca.crt"}
)

func useInsecureSkipVerify() bool {
	return true
}
