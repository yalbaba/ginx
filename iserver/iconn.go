package iserver

import "net"

//封装连接的接口
type IConn interface {
	Start()

	Stop()

	GetConnId() uint32

	GetTcpConn() *net.TCPConn

	GetRemoteAddr() net.Addr

	Send(msgId uint32, data []byte) error
}
