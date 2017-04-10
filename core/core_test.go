package core

import "testing"

type MyTestServer1 struct {
	grapes.CoreServer
}

type MyTestServer2 struct {
	grapes.CoreServer
}

type MyTestServer3 struct {
	grapes.CoreServer
}

func TestAsClientConnect(t *testing.T) {
	svr1 := &MyTestServer1{}
	svr1.SetListenPort(10001)
	// 注册事件处理方法
	svr1.Handle(3, func(req *grapes.GRequest, res *grapes.GResponse) {
		test := &kmtt.ResponseTestData{Label: "this is a test string.", Type: 123}
		data, err := proto.Marshal(test)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		res.Send(data)
	})

	append(connectedNodes, NodeInfo{"32", "323"})
	var connectedNodes []NodeInfo = make([]NodeInfo, 0)

	// 读取配置，连接关联的后端服务器。

	// 向Master注册本服务器

	// 服务器初始化完成，开始对外提供服务
	svr.InitComplete()
}
