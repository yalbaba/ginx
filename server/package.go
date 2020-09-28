package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"yalbaba/ginx/iserver"
	"yalbaba/ginx/util/global_conf"
)

/*
消息格式:  head(长度|id)+data
*/

type Package struct {
}

func NewPackage() *Package {
	return &Package{}
}

//获取头部的长度
func (p *Package) GetHeadLen() uint32 {
	//这里定义数据包协议的头部长度
	return 8
}

//将消息以一定协议来进行封装
func (p *Package) Pack(m iserver.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	//写头部的消息长度
	if err := binary.Write(buf, binary.LittleEndian, m.GetLen()); err != nil {
		return nil, err
	}

	//写消息id
	if err := binary.Write(buf, binary.LittleEndian, m.GetMessageId()); err != nil {
		return nil, err
	}

	//写消息体
	if err := binary.Write(buf, binary.LittleEndian, m.GetData()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//将消息拆包，得到消息对象
func (p *Package) UnPack(data []byte) (iserver.IMessage, error) {
	buf := bytes.NewReader(data)

	msg := &Message{}

	//读长度
	binary.Read(buf, binary.LittleEndian, &msg.DataLen)

	if global_conf.GlobalConfObj.MaxPackageSize > 0 && global_conf.GlobalConfObj.MaxPackageSize < msg.DataLen {
		return nil, fmt.Errorf("data is too big")
	}

	binary.Read(buf, binary.LittleEndian, &msg.Id)
	return msg, nil
}
