package main

import (
	"log"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/server"
)

func main() {
	s := server.NewGServer()
	s.AddRouter(1, &MyRouter{})
	s.Serve()
}

type MyRouter struct {
	server.BaseRouter
}

func (mr *MyRouter) Handle(request iserver.IRequest) {
	log.Println("handle....")
	mr.BaseRouter.Handle(request)
}
