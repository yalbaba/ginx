package server

import (
	"fmt"
	"log"
	"net"
)

type GServer struct {
	Name      string `json:"name"`
	IpVersion string `json:"ip_version"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	close     chan interface{}
}

func NewGServer(name, addr string, port int) *GServer {
	return &GServer{
		Name:      name,
		IpVersion: "tcp4",
		Addr:      addr,
		Port:      port,
		close:     make(chan interface{}),
	}
}

func (s *GServer) Start() {

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
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalf("accept err:%v", err)
				continue
			}

			// 获取请求内容，执行操作
			go func() {
				for {
					buf := make([]byte, 512)
					_, err := conn.Read(buf)
					if err != nil {
						log.Fatalf("get client msg err:%v", err)
						continue
					}

					log.Printf("client msg:%s\n", string(buf))
					// 应答客户端
					if _, err = conn.Write([]byte("收到请求")); err != nil {
						log.Fatalf("reply err:%v", err)
					}
				}
			}()
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
	select {
	case <-s.close:
		// 关闭服务器
		s.Stop()
	}
	return nil
}
