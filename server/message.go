package server

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMessage(msgId uint32, data []byte) *Message {
	return &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMessageId() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}
