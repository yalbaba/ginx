package iserver

type IServer interface {
	Start()
	Stop()
	Serve() error
	AddRouter(router IRouter)
}
