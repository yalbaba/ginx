package server

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestPackage(t *testing.T) {
	//模拟服务端
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("监听 err:", err)
		return
	}

	go func(listener net.Listener) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Accept err:", err)
				break
			}

			go func(conn net.Conn) {
				// 对数据进行拆包
				dp := &Package{}
				dataHead := make([]byte, dp.GetHeadLen())
				for {
					count, err := io.ReadFull(conn, dataHead)
					if err != nil {
						fmt.Println("ReadFull err:", err, "count:", count)
						break
					}
					msgHead, err := dp.UnPack(dataHead)
					if err != nil {
						fmt.Println("UnPack err:", err)
						return
					}
					msg := msgHead.(*Message)
					if msg.GetLen() > 0 {
						//有数据，继续读
						msg.Data = make([]byte, msg.GetLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("ReadFulldata err:", err)
							return
						}
						//读完了头和消息体
						fmt.Println("id:", msg.Id, "len:", msg.DataLen, "data:", string(msg.Data))
					}
				}
			}(conn)
		}

	}(listener)

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
	//构建两个包的内容
	dp := &Package{}
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	send1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack1 err:", err)
		return
	}
	fmt.Println("send1::", send1)
	msg2 := &Message{
		Id:      2,
		DataLen: 4,
		Data:    []byte{'w', 'o', 'r', 'd'},
	}
	send2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack2 err:", err)
		return
	}
	fmt.Println("send2::", send2)
	send1 = append(send1, send2...)
	conn.Write(send1)

	select {}

}
