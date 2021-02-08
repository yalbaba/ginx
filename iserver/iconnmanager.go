package iserver

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connId uint32) (IConnection, error)
	Clean()
	Len() int
}
