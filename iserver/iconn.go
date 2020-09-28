package iserver

import "net"

type IConn interface {
	Start()

	Stop()

	GetConnId() uint32

	GetTcpConn() *net.TCPConn

	GetRemoteAddr() net.Addr

	Send(msgId uint32, data []byte) error
}

// 收到请求后的回调方法
type Handler func(conn *net.TCPConn, bytes []byte, cnt int) error
