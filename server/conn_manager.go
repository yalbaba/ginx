package server

import (
	"fmt"
	"sync"
	"yalbaba/ginx/iserver"
)

/*
连接管理器
*/

type ConnManager struct {
	connections map[uint32]iserver.IConnection
	sync.Mutex
}

func NewConnManager() iserver.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]iserver.IConnection),
	}
}

func (c *ConnManager) Add(conn iserver.IConnection) {

	c.Lock()
	defer c.Unlock()
	c.connections[conn.GetConnId()] = conn
}

func (c *ConnManager) Remove(conn iserver.IConnection) {

	c.Lock()
	defer c.Unlock()
	delete(c.connections, conn.GetConnId())
}
func (c *ConnManager) Get(connId uint32) (iserver.IConnection, error) {

	c.Lock()
	defer c.Unlock()
	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("connection not found")
}
func (c *ConnManager) Clean() {

	c.Lock()
	defer c.Unlock()
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}
