package iserver

type IMessage interface {
	GetLen() uint32
	GetMessageId() uint32
	GetData() []byte
}
