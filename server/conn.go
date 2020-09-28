package server

import (
	"fmt"
	"io"
	"net"
	"yalbaba/ginx/iserver"
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
	dp := NewPackage()
	dataHead := make([]byte, dp.GetHeadLen())
	for {
		//获取每个包的头
		if _, err := io.ReadFull(c.Conn, dataHead); err != nil {
			fmt.Println(err)
			break
		}
		//获取每个包体
		msgHead, err := dp.UnPack(dataHead)
		if err != nil {
			fmt.Println(err)
			break
		}
		//根据头信息获取包体
		msg := msgHead.(*Message)
		if msgHead.GetLen() > 0 {
			//表示该包有数据
			msg.Data = make([]byte, msg.GetLen())
			_, err := io.ReadFull(c.Conn, msg.Data)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("这是本包的内容:", msg)
		}
		// 获取请求对象
		request := &Request{
			conn: c,
			data: msg,
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

func (c *GConn) Send(msgId uint32, data []byte) error {
	dp := NewPackage()
	//对消息进行打包
	msg := NewMessage(msgId, data)
	dataByte, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	//进行回写
	c.Conn.Write(dataByte)
	return nil
}
