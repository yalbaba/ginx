package iserver

import "net"

//封装连接的接口
type IConnection interface {
	Start()

	Stop()

	GetConnId() uint32

	GetTcpConn() *net.TCPConn

	GetRemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error

	SendBuffMsg(msgId uint32, data []byte) error //发送缓冲消息

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}
