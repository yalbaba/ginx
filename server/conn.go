package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util/global_conf"
)

type GConnection struct {
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

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	sync.RWMutex
}

func NewGConn(s iserver.IServer, conn *net.TCPConn, connId uint32, msgHandler iserver.IMsgHandler) *GConnection {
	return &GConnection{
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
func (c *GConnection) StartRead() {

	//断开连接后要关闭连接
	defer c.Stop()

	dp := NewPackage()
	dataHead := make([]byte, dp.GetHeadLen())
	for {
		//获取每个包的头
		if _, err := io.ReadFull(c.Conn, dataHead); err != nil && err != io.EOF {
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

func (c *GConnection) StartWrite() {
	log.Println("开始进行回写消息给客户端...")
	defer c.Stop()

	for {
		select {
		case data := <-c.msgChannel:
			//服务端给客户端写数据
			_, err := c.Conn.Write(data)
			if err != nil {
				log.Fatalf("服务端给客户端写数据失败,err:%v", err)
				return
			}
			log.Println("回写消息结束...")
		case data := <-c.msgBuffChannel:
			//服务端给客户端写数据
			_, err := c.Conn.Write(data)
			if err != nil {
				log.Fatalf("服务端给客户端写数据失败,err:%v", err)
				return
			}
			log.Println("回写缓冲消息结束...")
		case <-c.closeCh:
			log.Println("连接关闭...")
			return
		}
	}

}

func (c *GConnection) Start() {

	// 服务器内部读取数据后执行的流程
	go c.StartRead()
	// 服务端回写数据流程
	go c.StartWrite()

	//开启链接回调函数
	if c.TcpServer.GetOnConnStart() != nil {
		c.TcpServer.CallOnConnStart(c)
	}
}

func (c *GConnection) Stop() {

	if c.isClosed {
		fmt.Println("conn is closed")
		return
	}

	if err := c.Conn.Close(); err != nil {
		fmt.Println("stop err", err.Error())
		return
	}

	c.isClosed = true

	if c.TcpServer.GetOnConnStop() != nil {
		c.TcpServer.CallOnConnStop(c)
	}

	c.closeCh <- struct{}{}
	c.TcpServer.GetConnManager().Remove(c)
	close(c.closeCh)
	close(c.msgChannel)

}

func (c *GConnection) GetConnId() uint32 {
	return c.ConnId
}

func (c *GConnection) GetTcpConn() *net.TCPConn {
	return c.Conn
}

func (c *GConnection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *GConnection) SendMsg(msgId uint32, data []byte) error {
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

func (c *GConnection) SendBuffMsg(msgId uint32, data []byte) error {

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

//设置链接属性
func (c *GConnection) SetProperty(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.property[key] = value
}

//获取链接属性
func (c *GConnection) GetProperty(key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("该属性不存在")
}

//移除链接属性
func (c *GConnection) RemoveProperty(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.property, key)
}
