package iserver

//包装和拆包消息的接口
type IPackage interface {
	GetHeadLen() uint32
	Pack(m IMessage) ([]byte, error)
	UnPack(data []byte) (IMessage, error)
}
