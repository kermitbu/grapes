package main

import (
	grapes "github.com/kermitbu/grapes/core"
	log "github.com/kermitbu/grapes/log"
	proto "github.com/kermitbu/grapes/proto"
	"io"
	"kmtt"
	"net"
)

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

	// 创建一个服务器实例
	svr := new(GrapesMQ)

	// 注册事件处理方法
	svr.Handle(1, func(req *grapes.GRequest, res *grapes.GResponse) {

		// 构造返回的数据
		test := &kmtt.ResponseTestData{Label: "this is a test string.", Type: 123}
		data, err := proto.Marshal(test)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}

		res.Send(data)
	})

	// 读取配置，连接关联的后端服务器。

	// 向Master注册本服务器

	// 服务器初始化完成，开始对外提供服务
	svr.InitComplete()
}

const (
	BufLength = 1024
)

func Handle(conn net.Conn) {
	for {
		data := make([]byte, 0)
		buf := make([]byte, BufLength)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err.Error())
			}
			data = append(data, buf[:n]...)
			if n != BufLength {
				break
			}
		}
		log.Debug("Receive message")
	}
}
