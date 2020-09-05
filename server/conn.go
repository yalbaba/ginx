package server

import (
	"fmt"
	"net"
	"yalbaba/ginx/iserver"
)

type GConn struct {
	ConnId uint32

	TcpConn *net.TCPConn

	isClosed bool

	HandleFuc iserver.Handler

	CloseCh chan struct{}
}

func NewGConn(conn *net.TCPConn, connId uint32, handler iserver.Handler) *GConn {
	return &GConn{
		ConnId:    connId,
		TcpConn:   conn,
		HandleFuc: handler,
		isClosed:  false,
		CloseCh:   make(chan struct{}),
	}
}

// 这是服务器内部的读取数据的方法，不包括业务，具体业务是handlerfunc
func (c *GConn) StartRead() {

	//断开连接后要关闭连接
	defer c.Stop()

	for {

		buf := make([]byte, 512)
		//todo 解决eof错误
		if _, err := c.TcpConn.Read(buf); err != nil {
			fmt.Println("read data err:", err.Error())
			continue
		}

		fmt.Println("accept data:", string(buf))

		// 执行自定义的业务
		if err := c.HandleFuc(c.TcpConn, buf, 512); err != nil {
			fmt.Println("handleFuc err:", err.Error())
			break
		}

	}
}

func (c *GConn) Start() {

	// 服务器内部读取数据后执行的流程
	go c.StartRead()
	// todo 写数据的方法
}

func (c *GConn) Stop() {
	if err := c.TcpConn.Close(); err != nil {
		fmt.Println("stop err", err.Error())
		return
	}
	if c.isClosed {
		fmt.Println("conn is closed")
		return
	}

	c.isClosed = true
	close(c.CloseCh)
}

func (c *GConn) GetConnId() uint32 {
	return c.ConnId
}

func (c *GConn) GetTcpConn() *net.TCPConn {
	return c.TcpConn
}

func (c *GConn) GetRemoteAddr() net.Addr {
	return c.TcpConn.RemoteAddr()
}

func (c *GConn) Send(data []byte) error {
	if _, err := c.TcpConn.Write(data); err != nil {
		fmt.Println("send err:", err.Error())
		return err
	}
	return nil
}
