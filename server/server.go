package server

import (
	"fmt"
	"log"
	"net"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util"
)

type GServer struct {
	Name      string          `json:"name"`
	IpVersion string          `json:"ip_version"`
	Addr      string          `json:"addr"`
	Port      int             `json:"port"`
	Router    iserver.IRouter //自定义的处理业务
}

func NewGServer(name, addr string, port int) *GServer {
	return &GServer{
		Name:      name,
		IpVersion: "tcp4",
		Addr:      addr,
		Port:      port,
		Router:    nil,
	}
}

func (s *GServer) Start() {
	if s.Router == nil {
		log.Fatalf("路由方法为空")
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
			dealConn := NewGConn(conn, util.GetConnId(), s.Router)
			go dealConn.Start()
		}
	}()
}

func (s *GServer) Stop() {
	return
}

func (s *GServer) Serve() error {
	s.Start()

	// todo 做点其他初始化服务器业务

	// 阻塞
	select {}
	return nil
}

func (s *GServer) AddRouter(router iserver.IRouter) {
	s.Router = router
}
