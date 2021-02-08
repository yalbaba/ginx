package server

import (
	"fmt"
	"log"
	"yalbaba/ginx/iserver"
)

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request iserver.IRequest) {

}

func (br *BaseRouter) Handle(request iserver.IRequest) {
	fmt.Println("BaseRouter")

	//读取数据
	msgId := request.GetMessageId()
	data := request.GetData()
	request.GetDataLen()
	//发送消息到写模块，进行服务端对客户端的响应
	if err := request.GetConn().SendMsg(msgId, data.GetData()); err != nil {
		log.Fatalf("发送消息到写数据模块失败,err:%v", err)
	}

}

func (br *BaseRouter) PostHandle(request iserver.IRequest) {

}
