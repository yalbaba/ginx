package server

import (
	"fmt"
	"log"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util/global_conf"
)

type MsgHandler struct {
	ApisHandler     map[uint32]iserver.IRouter
	WorkerPoolSize  uint32                  //原本是来一个请求分配一个协程，为了控制协程数量而使用工作池原理
	TaskWorkerQueue []chan iserver.IRequest //队列个数（与worker数量保持一致）
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		ApisHandler:     make(map[uint32]iserver.IRouter),
		WorkerPoolSize:  global_conf.GlobalConfObj.WorkerPoolSize,
		TaskWorkerQueue: make([]chan iserver.IRequest, global_conf.GlobalConfObj.WorkerPoolSize),
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

//开启工作池
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//初始化队列
		m.TaskWorkerQueue[i] = make(chan iserver.IRequest, global_conf.GlobalConfObj.MaxQueueTaskSize)
		//开启worker，等待消息
		go m.StartOneWorker(uint32(i), m.TaskWorkerQueue[i])
	}
}

func (m *MsgHandler) StartOneWorker(workerId uint32, queue chan iserver.IRequest) {
	fmt.Printf("id为:%d 的worker启动了\n", workerId)
	for {
		select {
		case request := <-queue:
			m.DoHandle(request)
		}
	}
}

//提交请求到对应工作池中
func (m *MsgHandler) SendMessageToPool(connId uint32, request iserver.IRequest) {
	workerId := connId % m.WorkerPoolSize
	m.TaskWorkerQueue[workerId] <- request
}
