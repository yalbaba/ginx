package server

import "testing"

func TestServer(t *testing.T) {
	s := NewGServer("testServer", "127.0.0.1", 9090)
	s.Serve()
}
