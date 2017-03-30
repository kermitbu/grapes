package core

import (
	"encoding/binary"
	"errors"
	log "github.com/kermitbu/grapes/logs"
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

func (c *CoreServer) InitComplete() {
	// 作为客户端，连接服务器，并准备接收数据

	// for i := 0; i < 4; i++ {
	// 	addr, err := net.ResolveTCPAddr("tcp", ":4040")
	// 	checkErr(err)
	// 	conn, err := net.DialTCP("tcp", nil, addr)
	// 	checkErr(err)

	// 	allClientConnects[addr.String()] = conn

	// 	defer conn.Close()
	// 	go handlClientConn(conn)
	// }

	// 作为服务器端监听端口
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4040")
	checkErr(err)
	listen, err := net.ListenTCP("tcp", addr)
	checkErr(err)

	log.Info("服务器正常启动监听")

	complete := make(chan int, 1)

	go func(listen *net.TCPListener) {
		for {
			conn, err := listen.Accept()
			checkErr(err)
			go c.handleServerConn(conn)
		}
	}(listen)

	<-complete
}

const BufLength = 1024

func (c *CoreServer) handleServerConn(conn net.Conn) {
	log.Info("=====>> 处理一个新的连接")

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

		log.Info("接收到数据：%v", unhandledData)

		for nil == head.Unpack(unhandledData) {
			msgLen := head.BodyLen + uint16(head.HeadLen)

			log.Debug("HeadLen = %d, BodyLen = %d, msgLen = %d, head.cmd= %d", head.HeadLen, head.BodyLen, msgLen, head.Cmd)

			msgData := unhandledData[:msgLen]

			unhandledData = unhandledData[msgLen:]

			c.deliverMessage(conn, head, msgData[head.HeadLen:])

			log.Debug("%v", msgData)
		}
	}
	log.Info("*********处理结束********")
}

func (c *CoreServer) Handle(id uint16, f handleFunc) {

	if c.allHandlerFunc == nil {
		c.allHandlerName = make(map[uint16]string)
		c.allHandlerFunc = make(map[uint16]handleFunc)
	}
	if _, ok := c.allHandlerFunc[id]; ok {
		log.Warn("Register called twice for handles ", id)
	}
	c.allHandlerFunc[id] = f

}

func (c *CoreServer) deliverMessage(conn net.Conn, msghead *MessageHead, body []byte) {
	if handler, ok := c.allHandlerFunc[msghead.Cmd]; ok {

		req := &GRequest{connect: &conn, head: msghead, DataLen: msghead.BodyLen, DataBuffer: body}
		rsp := &GResponse{connect: &conn}
		handler(req, rsp)
	} else {
		log.Warn("Never register processing method [%v]", msghead.Cmd)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}
}

type GRequest struct {
	connect    *net.Conn
	head       *MessageHead
	DataLen    uint16
	DataBuffer []byte
}

type GResponse struct {
	connect    *net.Conn
	DataLen    uint16
	DataBuffer []byte
}

func (r *GResponse) Send(data []byte) {
	if len(data) > 0 {

		(*(r.connect)).Write(data)
	} else {
		log.Warn("Send data is empty.")
	}
}

type handleFunc func(request *GRequest, response *GResponse)

type CoreServer struct {
	allHandlerName map[uint16]string
	allHandlerFunc map[uint16]handleFunc

	allClientConnects map[string]*net.TCPConn
}
