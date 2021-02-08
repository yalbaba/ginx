package iserver

//封装一次请求的结构，包含本次请求的连接和数据
type IRequest interface {
	// 获取连接
	GetConn() IConnection

	// 获取数据
	GetData() IMessage

	GetDataLen() uint32

	GetMessageId() uint32
}
