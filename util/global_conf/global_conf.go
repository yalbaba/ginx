package global_conf

import (
	"encoding/json"
	"io/ioutil"
)

// 配置文件对象
type GlobalConf struct {
	Name             string `json:"name"`
	Port             int    `json:"port"`
	Host             string `json:"host"`
	IpVersion        string `json:"ip_version"`
	MaxConn          int    `json:"max_conn"`            //最大连接数
	MaxPackageSize   uint32 `json:"max_package_size"`    //最大传输包大小
	WorkerPoolSize   uint32 `json:"worker_pool_size"`    //消息队列的工作池数量
	MaxQueueTaskSize uint32 `json:"max_queue_task_size"` //每个消息队列的最大任务数
}

var GlobalConfObj *GlobalConf

func init() {
	GlobalConfObj = &GlobalConf{
		Name:             "default_server",
		Port:             9090,
		Host:             "127.0.0.1",
		IpVersion:        "tcp4",
		MaxConn:          100,
		MaxPackageSize:   512,
		WorkerPoolSize:   10,
		MaxQueueTaskSize: 1024,
	}
	GlobalConfObj.Reload()
}

func (g *GlobalConf) Reload() {
	data, _ := ioutil.ReadFile("/Users/dr.yang/go/src/yalbaba/ginx/config/ginx.json")
	json.Unmarshal(data, g)
}
