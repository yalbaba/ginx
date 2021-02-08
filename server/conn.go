package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util/global_conf"
)

type GConn struct {
	TcpServer iserver.IServer

	ConnId uint32

	Conn *net.TCPConn

	isClosed bool

	MsgHandler iserver.IMsgHandler

	//用于写模块接收消息的通道
	msgChannel chan []byte

	//用于写模块接收消息的通道（带缓冲）
	msgBuffChannel chan []byte

	//客户端断开连接，关闭当前连接的通道
	closeCh chan struct{}
}

func NewGConn(s iserver.IServer, conn *net.TCPConn, connId uint32, msgHandler iserver.IMsgHandler) *GConn {
	return &GConn{
		TcpServer:      s,
		ConnId:         connId,
		Conn:           conn,
		MsgHandler:     msgHandler,
		isClosed:       false,
		msgChannel:     make(chan []byte),
		msgBuffChannel: make(chan []byte, global_conf.GlobalConfObj.MaxMsgBuff),
		closeCh:        make(chan struct{}),
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
			fmt.Println("id:", msg.Id, "len:", msg.DataLen, "data:", string(msg.Data))
		}

		// 构建请求对象
		request := &Request{
			conn: c,
			data: msg,
		}

		// 执行用户添加的的业务 不能来一个请求就开一个协程，要控制数量
		//c.Router.PreHandle(request)
		//go c.MsgHandler.DoHandle(request)
		//c.Router.PostHandle(request)

		// 之前的开启协程来读写改为提交到工作池中
		c.MsgHandler.SendMessageToPool(c.GetConnId(), request)
	}
}

func (c *GConn) StartWrite() {
	log.Println("开始进行回写消息给客户端")
	defer log.Println("客户端已断开连接")
	select {
	case data := <-c.msgChannel:
		//服务端给客户端写数据
		_, err := c.Conn.Write(data)
		if err != nil {
			log.Fatalf("服务端给客户端写数据失败,err:%v", err)
			return
		}
	case data := <-c.msgBuffChannel:
		//服务端给客户端写数据
		_, err := c.Conn.Write(data)
		if err != nil {
			log.Fatalf("服务端给客户端写数据失败,err:%v", err)
			return
		}
	case <-c.closeCh:
		return
	}
}

func (c *GConn) Start() {

	// 服务器内部读取数据后执行的流程
	go c.StartRead()
	// 服务端回写数据流程
	go c.StartWrite()
}

func (c *GConn) Stop() {

	if c.isClosed {
		fmt.Println("conn is closed")
		return
	}

	if err := c.Conn.Close(); err != nil {
		fmt.Println("stop err", err.Error())
		return
	}

	c.isClosed = true
	c.closeCh <- struct{}{}

	c.TcpServer.GetConnManager().Remove(c)

	close(c.closeCh)
	close(c.msgChannel)

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

func (c *GConn) SendMsg(msgId uint32, data []byte) error {

	if c.isClosed {
		return fmt.Errorf("连接已关闭")
	}

	dp := NewPackage()
	//对消息进行打包
	msg := NewMessage(msgId, data)
	dataByte, err := dp.Pack(msg)
	if err != nil {
		return err
	}

	//发送消息到通道给写数据协程
	c.msgChannel <- dataByte
	return nil
}

func (c *GConn) SendBuffMsg(msgId uint32, data []byte) error {

	if c.isClosed {
		return fmt.Errorf("连接已关闭")
	}

	dp := NewPackage()
	//对消息进行打包
	msg := NewMessage(msgId, data)
	dataByte, err := dp.Pack(msg)
	if err != nil {
		return err
	}

	//发送消息到通道给写数据协程
	c.msgChannel <- dataByte
	return nil
}
