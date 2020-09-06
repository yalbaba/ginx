package main

import (
	"log"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/server"
)

func main() {
	s := server.NewGServer("testServer", "127.0.0.1", 9090)
	s.AddRouter(&MyRouter{})
	s.Serve()
}

type MyRouter struct {
	server.BaseRouter
}

func (mr *MyRouter) Handle(request iserver.IRequest) {
	log.Println("handle....")
}
