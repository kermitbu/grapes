package core

import (
	"fmt"
)

type MessageHead struct {
	Id      int
	BodyLen int
}

type GMessage struct {
	MessageHead
	MessageBody []byte
}

type handlefunc func(int, []byte)

var allhandelrName map[int]string
var allHandler map[int]handlefunc = make(map[int]handlefunc)

type GrapesCore struct {
}

func (g *GrapesCore) Register(id int, f handlefunc) {
	if _, ok := allHandler[id]; ok {
		fmt.Println("已注册过了。")
	} else {
		allHandler[id] = f
	}
}

func (g *GrapesCore) HandleMessage(msg *GMessage) {

	if msg == nil {
		fmt.Println("Param error. msg is nil")
		return
	}

	if handler, ok := allHandler[msg.Id]; ok {
		handler(msg.BodyLen, msg.MessageBody)
	} else {
		fmt.Println(msg.Id)
	}

	fmt.Println("This is Core.HandleMessage")
}
