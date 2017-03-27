package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	utils "github.com/kermitbu/grapes/utils"
	"io"
	"net"
	"os"
)

type MessageHead struct {
	Cmd     uint16
	Version byte
	HeadLen byte
	BodyLen uint16
	adfe    uint16
}

func (m *MessageHead) Unpack(buf []byte) error {
	if len(buf) > utils.SizeStruct(m) {
		m.Cmd = binary.BigEndian.Uint16(buf[:2])
		m.Version = buf[2]
		m.HeadLen = buf[3]
		m.BodyLen = binary.BigEndian.Uint16(buf[4:6])
		return nil
	}
	return errors.New("数据长度小于最小协议的长度")
}

func (m *MessageHead) Pack() (buf []byte) {

	size := utils.SizeStruct(m)
	buf = make([]byte, size)

	binary.BigEndian.PutUint16(buf[:2], m.Cmd)
	buf[2] = m.Version
	buf[3] = byte(size)
	binary.BigEndian.PutUint16(buf[4:6], m.BodyLen)
	return buf
}

type handleFunc func(uint16, []byte)
type CoreServer struct {
	allHandlerName map[uint16]string
	allHandlerFunc map[uint16]handleFunc

	allConnects map[string]*net.TCPConn
}

func (c *CoreServer) InitComplete() {
	addr, err := net.ResolveTCPAddr("tcp", ":4040")
	checkErr(err)
	listen, err := net.ListenTCP("tcp", addr)
	checkErr(err)

	fmt.Println("服务器正常启动监听")

	complete := make(chan int, 1)

	go func(listen *net.TCPListener) {
		for {
			conn, err := listen.Accept()
			checkErr(err)
			go c.handleConn(conn)
		}
	}(listen)

	<-complete
}

const BufLength = 1024

func (c *CoreServer) handleConn(conn net.Conn) {
	fmt.Println("=====>> 处理一个新的连接")

	head := new(MessageHead)

	hasError := false
	unhandledData := make([]byte, 0)

	for false == hasError {
		buf := make([]byte, BufLength)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				hasError = true
			}

			unhandledData = append(unhandledData, buf[:n]...)

			if n != BufLength {
				break
			}
		}

		fmt.Println("接收到数据：", unhandledData)

		for nil == head.Unpack(unhandledData) {
			msgLen := head.BodyLen + uint16(head.HeadLen)

			fmt.Printf("HeadLen = %d, BodyLen = %d, msgLen = %d, head.cmd= %d\n", head.HeadLen, head.BodyLen, msgLen, head.Cmd)

			msgData := unhandledData[:msgLen]

			unhandledData = unhandledData[msgLen:]

			c.deliverMessage(head, msgData[head.HeadLen:])

			fmt.Println(msgData)
		}
	}
	fmt.Println("*********处理结束********")
}

func (c *CoreServer) Register(id uint16, f handleFunc) {
	if c.allHandlerFunc == nil {
		fmt.Println("c.allHandlerFunc")
		c.allHandlerName = make(map[uint16]string)
		c.allHandlerFunc = make(map[uint16]handleFunc)
	}

	if c.allHandlerFunc == nil {
		fmt.Println("core: Register handles is nil")
	}
	if _, ok := c.allHandlerFunc[id]; ok {
		fmt.Println("core: Register called twice for handles ", id)
	}
	c.allHandlerFunc[id] = f

}

func (c *CoreServer) deliverMessage(head *MessageHead, body []byte) {

	fmt.Println(c.allHandlerFunc)
	if handler, ok := c.allHandlerFunc[head.Cmd]; ok {
		handler(head.BodyLen, body)
	} else {
		fmt.Println("从未注册过", head.Cmd, "的处理方法")
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
