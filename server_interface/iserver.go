package server_interface

type IServer interface {
	Start()
	Stop()
	Serve() error
}
