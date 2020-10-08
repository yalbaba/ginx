package iserver

//路由方法接口
type IRouter interface {
	//PreHandle(request IRequest)
	Handle(request IRequest)
	//PostHandle(request IRequest)
}
