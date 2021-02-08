package iserver

import "net"

//封装连接的接口
type IConn interface {
	Start()

	Stop()

	GetConnId() uint32

	GetTcpConn() *net.TCPConn

	GetRemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error

	SendBuffMsg(msgId uint32, data []byte) error //发送缓冲消息
}
