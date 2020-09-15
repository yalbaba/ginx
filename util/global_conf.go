package util

import (
	"encoding/json"
	"io/ioutil"
)

// 配置文件对象
type GlobalConf struct {
	Name           string `json:"name"`
	Port           int    `json:"port"`
	Host           string `json:"host"`
	MaxConn        int    `json:"max_conn"`
	MaxPackageSize int    `json:"max_package_size"`
}

var ConfObj *GlobalConf

func init() {
	ConfObj = &GlobalConf{
		Name:           "default_server",
		Port:           9090,
		Host:           "127.0.0.1",
		MaxConn:        100,
		MaxPackageSize: 60,
	}
	ConfObj.Reload()
}

func (g *GlobalConf) Reload() {
	data, _ := ioutil.ReadFile("../conf/ginx.json")
	json.Unmarshal(data, g)
}
