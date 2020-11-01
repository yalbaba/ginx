package server

import (
	"log"
	"yalbaba/ginx/iserver"
)

type MsgHandler struct {
	ApisHandler map[uint32]iserver.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		ApisHandler: make(map[uint32]iserver.IRouter),
	}
}

func (m *MsgHandler) DoHandle(request iserver.IRequest) {
	m.ApisHandler[request.GetMessageId()].Handle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router iserver.IRouter) {
	if _, ok := m.ApisHandler[msgId]; ok {
		log.Fatalf("该路由已注册过")
		return
	}
	m.ApisHandler[msgId] = router
}
