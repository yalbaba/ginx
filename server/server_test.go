package server

import "testing"

func TestServer(t *testing.T) {
	s := NewGServer()
	s.AddRouter(&myrouter{})
	s.Serve()
}

type myrouter struct {
	*BaseRouter
}
