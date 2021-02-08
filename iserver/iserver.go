package iserver

type IServer interface {
	Start()
	Stop()
	Serve() error
	AddRouter(msgId uint32, router IRouter) //根据msgId不同处理的逻辑不同
	GetConnManager() IConnManager
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
	GetOnConnStart() func(IConnection)
	GetOnConnStop() func(IConnection)
}
