
[TOC]


# Grapes

### 设计想法
最简单的方法来创建一个节点，每个节点定义一个继承于CoreServer的结构即可。
```go
package main

import (
	"flag"
	grapes "github.com/kermitbu/grapes/core"
)

type GrapesMQ struct {
	grapes.CoreServer
}

func main() {
	flag.Parse()
	svr := new(GrapesMQ)
	// 注册事件处理方法
	svr.Handle(1, func(req *grapes.GRequest, res *grapes.GResponse) {
		// 处理业务逻辑
	})

	// 服务器初始化完成，开始对外提供服务
	svr.InitComplete()
}
```


协议定义
采用消息头+PB序列化数据的方式
```proto
type MessageHead struct {
	Cmd     uint16
	Version byte
	HeadLen byte
	BodyLen uint16
}
```
