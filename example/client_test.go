package main

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"yalbaba/ginx/server"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9099")
	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		dp := &server.Package{}
		send, _ := dp.Pack(server.NewMessage(1, []byte("ginx test")))

		fmt.Println("11111111")
		_, err := conn.Write(send)
		if err != nil {
			fmt.Println("Write err:", err)
			return
		}

		fmt.Println("2222222")
		//获取每个包的头
		dataHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, dataHead); err != nil && err != io.EOF {
			fmt.Println("dataHead::::", err)
			return
		}

		fmt.Println("3333333")
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
		time.Sleep(1 * time.Second)
	}

}
