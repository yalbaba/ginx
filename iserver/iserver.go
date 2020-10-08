package iserver

type IServer interface {
	Start()
	Stop()
	Serve() error
	AddRouter(msgId uint32, router IRouter) //根据msgId不同处理的逻辑不同
}
