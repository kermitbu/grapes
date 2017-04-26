/*
 * @Author: kermit.bu
 * @Date: 2017-04-24 15:40:36
 * @Last Modified by: kermit.bu
 * @Last Modified time: 2017-04-24 16:15:47
 */
package main

import (
	"flag"

	grapes "github.com/kermitbu/grapes/core"
)

// GrapesMQ asd
type GrapesMQ struct {
	grapes.CoreServer
}

type ServerNodeInfo struct {
	Name  string `json:"name"`
	Host  string `json:"host"`
	Port  uint16 `json:"port"`
	Group string `json:"group"`
}

type ServerGroup struct {
	DevItems []ServerNodeInfo `json:"development"`
	ProItems []ServerNodeInfo `json:"production"`
}

func main() {
	flag.Parse()
	// 创建一个服务器实例
	svr := new(GrapesMQ)

	// 注册事件处理方法
	svr.Handle(1, func(req *grapes.GRequest, res *grapes.GResponse) {
		// req.Head.Cmd
	})

	// 读取配置，连接关联的后端服务器。

	// 向Master注册本服务器

	// 服务器初始化完成，开始对外提供服务
	svr.InitComplete()
}
