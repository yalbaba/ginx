package iserver

// 多路由调度和注册
type IMsgHandler interface {
	DoHandle(request IRequest)
	AddRouter(msgId uint32, router IRouter)
	StartWorkerPool()
	StartOneWorker(workerId uint32, queue chan IRequest)
	SendMessageToPool(connId uint32, request IRequest)
}
