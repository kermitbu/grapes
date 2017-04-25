package core

import (
	"io"
	"net"

	log "github.com/kermitbu/grapes/log"
)

// InitConnectAsServer  a
func (c *CoreServer) InitConnectAsServer() {

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+*port)
	if err != nil {
		log.Fatal(err.Error())
	}
	listen, err := net.ListenTCP("tcp", addr)
	defer listen.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("服务器正常启动,开始监听%v端口", *port)

	// 监听
	go func(listen *net.TCPListener) {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err.Error())
			}
			go c.handleServerConn(conn)
		}
	}(listen)

	// 处理
	go func() {
		log.Debug("Channel开始等待数据")
		buffer := <-BufferChan
		log.Debug("从Channel中读取出数据：%v", buffer)
	}()

	complete := make(chan int, 1)
	<-complete
}

func (c *CoreServer) handleServerConn(conn net.Conn) {
	log.Info("===>>> New Connection ===>>>")

	head := new(MessageHead)
	unhandledData := make([]byte, 0)

DISCONNECT:
	for {
		buf := make([]byte, BufLength)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				log.Error("读取缓冲区出错，有可能是连接断开了: %v", err.Error())
				break DISCONNECT
			}

			unhandledData = append(unhandledData, buf[:n]...)

			if n != BufLength {
				break
			}
		}

		if len(unhandledData) == 0 {
			log.Error("读取到的数据长度为0，有可能是连接断开了")
			break
		}
		log.Debug("接收到数据：%v", unhandledData)

		for nil == head.Unpack(unhandledData) {
			log.Debug("解析出消息：%v", unhandledData)

			msgLen := head.BodyLen + uint16(head.HeadLen)
			msgData := unhandledData[:msgLen]
			unhandledData = unhandledData[msgLen:]

			BufferChan <- msgData

			c.deliverMessage(conn, head, msgData[head.HeadLen:])
		}
	}
	log.Info("===>>> Connection closed ===>>>")
}
