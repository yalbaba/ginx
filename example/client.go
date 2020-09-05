package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println(err.Error())
	}
	conn.Write([]byte("我请求了"))

	b := make([]byte, 512)
	conn.Read(b)
	fmt.Println("服务器回答内容：", string(b))
}
