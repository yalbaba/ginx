package server

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
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
