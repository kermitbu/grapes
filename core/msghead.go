package core

import (
	"encoding/binary"
	"errors"
	log "github.com/kermitbu/grapes/log"
	utils "github.com/kermitbu/grapes/utils"
	"net"
)

type MessageHead struct {
	Cmd     uint16
	Version byte
	HeadLen byte
	BodyLen uint16
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

func (r *GResponse) Close() {
	(*(r.connect)).Close()
}

// 服务器地址信息
type NodeAddr struct {
	Ip   string
	Port string
}
