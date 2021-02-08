package iserver

type IConnManager interface {
	Add(conn IConn)
	Remove(conn IConn)
	Get(connId uint32) (IConn, error)
	Clean()
	Len() int
}
