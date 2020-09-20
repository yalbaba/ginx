package server

import (
	"fmt"
	"yalbaba/ginx/iserver"
)

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request iserver.IRequest) {

}

func (br *BaseRouter) Handle(request iserver.IRequest) {
	fmt.Println("handle")
}

func (br *BaseRouter) PostHandle(request iserver.IRequest) {

}
