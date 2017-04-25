package core

import (
	"io"
	"net"

	log "github.com/kermitbu/grapes/log"
)

// InitConnectAsClient  a
func (c *CoreServer) InitConnectAsClient() {
	for i := range connectedNodes {
		addr, err := net.ResolveTCPAddr("tcp", connectedNodes[i].GetIp()+":"+connectedNodes[i].GetPort())
		if nil != err {
			log.Error("Resolve %s:%s error:", connectedNodes[i].GetIp(), connectedNodes[i].GetPort())
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		if nil != err {
			log.Error("DialTCP %s:%s error:", connectedNodes[i].GetIp(), connectedNodes[i].GetPort())
		}
		go c.handleClientConn(conn)
	}
}

func (c *CoreServer) handleClientConn(conn net.Conn) {
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
			msgLen := head.BodyLen + uint16(head.HeadLen)
			msgData := unhandledData[:msgLen]
			unhandledData = unhandledData[msgLen:]

			c.deliverMessage(conn, head, msgData[head.HeadLen:])
		}
	}
	log.Info("===>>> Connection closed ===>>>")
}
