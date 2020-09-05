package main

import "yalbaba/ginx/server"

func main() {
	s := server.NewGServer("testServer", "127.0.0.1", 9090)
	s.Serve()
}
