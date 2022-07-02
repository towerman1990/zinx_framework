package conf

import (
	"encoding/json"
	"io/ioutil"

	"towerman1990.cn/zinx_framework/ziface"
)

type Configure struct {
	TcpServer ziface.IServer

	Host string

	TcpPort int

	Name string

	Version string

	MaxConn int

	MaxPackageSize uint32

	WorkerPoolSize uint32

	MaxWorkerTaskLen uint32
}

var Config *Configure

func (c *Configure) Reload() {
	data, err := ioutil.ReadFile("conf/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}
}

func init() {
	Config = &Configure{
		Name:             "ZinxServerApp",
		Version:          "v0.8",
		TcpPort:          8088,
		Host:             "0.0.0.1",
		MaxConn:          100,
		MaxPackageSize:   4096,
		WorkerPoolSize:   8,
		MaxWorkerTaskLen: 1024,
	}

	Config.Reload()
}
