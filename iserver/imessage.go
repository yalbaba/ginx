package iserver

//封装请求消息的接口
type IMessage interface {
	GetLen() uint32
	GetMessageId() uint32
	GetData() []byte
}
