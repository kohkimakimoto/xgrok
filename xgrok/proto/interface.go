package proto

import (
	"github.com/kohkimakimoto/xgrok/xgrok/conn"
)

type Protocol interface {
	GetName() string
	WrapConn(conn.Conn, interface{}) conn.Conn
}
