package version

import (
	"github.com/kohkimakimoto/xgrok/xgrok"
)

const (
	Proto = "2"
)

func MajorMinor() string {
	return xgrok.Version
}
