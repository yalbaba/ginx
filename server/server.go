package server

import (
	"fmt"
	"log"
	"net"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util"
	"yalbaba/ginx/util/global_conf"
)

type GServer struct {
	Name      string              `json:"name"`
	IpVersion string              `json:"ip_version"`
	Addr      string              `json:"addr"`
	Port      int                 `json:"port"`
	MaxConns  int                 `json:"max_conns"`
	Handler   iserver.IMsgHandler //自定义的处理业务
}

func NewGServer() *GServer {
	return &GServer{
		Name:      global_conf.GlobalConfObj.Name,
		IpVersion: global_conf.GlobalConfObj.IpVersion,
		Addr:      global_conf.GlobalConfObj.Host,
		Port:      global_conf.GlobalConfObj.Port,
		MaxConns:  global_conf.GlobalConfObj.MaxConn,
		Handler:   NewMsgHandler(),
	}
}

func (s *GServer) Start() {
	if len(s.Handler.(*MsgHandler).ApisHandler) == 0 {
		log.Fatalf("未注册router")
		return
	}

	go func() {
		// 获取addr
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Addr, s.Port))
		if err != nil {
			log.Fatalf("get tcp addr err:%v", err)
			return
		}
		// 监听端口
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			log.Fatalf("listen addr err:%v", err)
			return
		}

		log.Printf("server start success addr:%s:%d\n", s.Addr, s.Port)

		// 阻塞获取客户端请求
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Fatalf("accept err:%v", err)
				continue
			}

			// 获取请求内容，执行操作
			dealConn := NewGConn(conn, util.GetConnId(), s.Handler)
			go dealConn.Start()
		}
	}()
}

func (s *GServer) Stop() {
	return
}

func (s *GServer) Serve() error {
	s.Start()

	// todo 做点其他初始化服务器业务(比如初始化中间件和数据库)

	// 阻塞
	select {}
	return nil
}

func (s *GServer) AddRouter(msgId uint32, router iserver.IRouter) {
	s.Handler.AddRouter(msgId, router)
}
