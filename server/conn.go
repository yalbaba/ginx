package server

import (
	"fmt"
	"net"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util/global_conf"
)

type GConn struct {
	ConnId uint32

	Conn *net.TCPConn

	isClosed bool

	Router iserver.IRouter

	CloseCh chan struct{}
}

func NewGConn(conn *net.TCPConn, connId uint32, router iserver.IRouter) *GConn {
	return &GConn{
		ConnId:   connId,
		Conn:     conn,
		Router:   router,
		isClosed: false,
		CloseCh:  make(chan struct{}),
	}
}

// 这是服务器内部的读取数据的方法，不包括业务，具体业务是handlerfunc
func (c *GConn) StartRead() {

	//断开连接后要关闭连接
	defer c.Stop()

	for {

		buf := make([]byte, global_conf.GlobalConfObj.MaxPackageSize)
		//todo 解决eof错误
		if _, err := c.Conn.Read(buf); err != nil {
			fmt.Println("read data err:", err.Error())
			break
		}

		fmt.Println("accept data:", string(buf))
		// 获取请求对象
		request := &Request{
			conn: c,
			data: buf,
		}

		// 执行用户添加的的业务
		//c.Router.PreHandle(request)
		c.Router.Handle(request)
		//c.Router.PostHandle(request)

	}
}

func (c *GConn) Start() {

	// 服务器内部读取数据后执行的流程
	go c.StartRead()
	// todo 写数据的方法
}

func (c *GConn) Stop() {
	if err := c.Conn.Close(); err != nil {
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
	return c.Conn
}

func (c *GConn) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *GConn) Send(data []byte) error {
	if _, err := c.Conn.Write(data); err != nil {
		fmt.Println("send err:", err.Error())
		return err
	}
	return nil
}
