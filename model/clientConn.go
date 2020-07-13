package model

import "net"

// ClientConn 客户端连接
type ClientConn struct {
	CliConn  net.Conn
	Username string
}
