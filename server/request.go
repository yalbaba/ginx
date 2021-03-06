package server

import "yalbaba/ginx/iserver"

type Request struct {
	conn iserver.IConnection
	data iserver.IMessage
}

func (r *Request) GetConn() iserver.IConnection {
	return r.conn
}

func (r *Request) GetData() iserver.IMessage {
	return r.data
}

func (r *Request) GetDataLen() uint32 {
	return r.data.GetLen()
}

func (r *Request) GetMessageId() uint32 {
	return r.data.GetMessageId()
}
