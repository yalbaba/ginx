package server

import "yalbaba/ginx/iserver"

type Request struct {
	conn iserver.IConn
	data []byte
}

func (r *Request) GetConn() iserver.IConn {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
