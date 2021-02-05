package main

import (
	"fmt"
	"io"
	"net"
	"testing"
	"yalbaba/ginx/server"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9099")
	if err != nil {
		fmt.Println(err.Error())
	}

	//构建两个包的内容
	dp := &server.Package{}
	msg1 := &server.Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	send, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack1 err:", err)
		return
	}
	conn.Write(send)

	//获取每个包的头
	dataHead := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(conn, dataHead); err != nil {
		fmt.Println("err1::::", err)
		return
	}
	//获取每个包体
	msgHead, err := dp.UnPack(dataHead)
	if err != nil {
		fmt.Println("err::::", err)
		return
	}
	//根据头信息获取包体
	msg := msgHead.(*server.Message)
	if msgHead.GetLen() > 0 {
		//表示该包有数据
		msg.Data = make([]byte, msg.GetLen())
		_, err := io.ReadFull(conn, msg.Data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("id:", msg.Id, "len:", msg.DataLen, "data:", string(msg.Data))

}
