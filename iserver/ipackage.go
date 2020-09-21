package iserver

type IPackage interface {
	GetHeadLen() uint32
	Pack(m IMessage) ([]byte, error)
	UnPack(data []byte) (IMessage, error)
}
