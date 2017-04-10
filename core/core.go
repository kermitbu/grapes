package core

import (
	"flag"
	log "github.com/kermitbu/grapes/log"
	proto "github.com/kermitbu/grapes/proto"
	"io"
	"net"
	"strconv"
)

// 外部用于注册事件处理方法的方法
func (c *CoreServer) Handle(id uint16, f handleFunc) {

	if c.allHandlerFunc == nil {
		c.allHandlerFunc = make(map[uint16]handleFunc)
	}
	if _, ok := c.allHandlerFunc[id]; ok {
		log.Warn("Register called twice for handles ", id)
	}
	c.allHandlerFunc[id] = f
}

// 派发事件
func (c *CoreServer) deliverMessage(conn net.Conn, msghead *MessageHead, body []byte) {
	if handler, ok := c.allHandlerFunc[msghead.Cmd]; ok {

		req := &GRequest{connect: &conn, head: msghead, DataLen: msghead.BodyLen, DataBuffer: body}
		rsp := &GResponse{connect: &conn}
		handler(req, rsp)
	} else {
		log.Warn("Never register processing method [%v]", msghead.Cmd)
	}
}

var port = flag.String("port", "10000", "指定服务器监听的端口号")
var conf = flag.String("conf", "", "指定服务器的配置文件")

func (c *CoreServer) SetListenPort(p uint16) {
	*port = strconv.Itoa(int(p))
}

func (c *CoreServer) initHandleJoinRequest() {

	if c.allHandlerFunc == nil {
		c.allHandlerFunc = make(map[uint16]handleFunc)
	}

	c.allHandlerFunc[uint16(SystemEvent_REQUEST_JOIN_CLUSTER)] = func(req *GRequest, res *GResponse) {
		// 1. 从 req中解析出来报过来的IP和端口号
		// 进行解码
		nodeInfo := &NodeInfo{}
		err := proto.Unmarshal(req.DataBuffer, nodeInfo)
		if err != nil {
			log.Error("解析收到的节点信息出错: ", err)
		}
		log.Debug("%v:%v  %v", nodeInfo.GetIp(), nodeInfo.GetPort(), nodeInfo.GetType())

		// 2. 把集群的信息发回去
		data, err := proto.Marshal(clusterNodes)
		if err != nil {
			log.Error("合成集群信息出错: ", err)
		}
		log.Debug("集群信息：%v", data)
		res.Send(data)
		res.Close()
	}

	c.allHandlerFunc[uint16(SystemEvent_NOTIFY_SYNC_CLUSTER)] = func(req *GRequest, res *GResponse) {
		// 同步集群的状态
	}
}

var connectedNodes []NodeInfo = make([]NodeInfo, 0)

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

func (c *CoreServer) InitComplete() {
	c.initHandleJoinRequest()

	// 连接相关的服务器
	c.InitConnectAsClient()

	// 作为服务器端监听端口, 正常传输数据使用
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

	complete := make(chan int, 1)

	go func(listen *net.TCPListener) {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err.Error())
			}
			go c.handleServerConn(conn)
		}
	}(listen)

	<-complete
}

const BufLength = 1024

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
			msgLen := head.BodyLen + uint16(head.HeadLen)
			msgData := unhandledData[:msgLen]
			unhandledData = unhandledData[msgLen:]

			c.deliverMessage(conn, head, msgData[head.HeadLen:])
		}
	}
	log.Info("===>>> Connection closed ===>>>")
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

/////////////////////////////////////////////////////////////////
// 存储集群中所有节点的信息
var clusterNodes *ClusterInfos = &ClusterInfos{}

type handleFunc func(request *GRequest, response *GResponse)

type CoreServer struct {
	allHandlerFunc    map[uint16]handleFunc
	allClientConnects map[string]*net.TCPConn
	groupName         string // 服务组名
}

func (c *CoreServer) SetGroupName(n string) {
	c.groupName = n
}

func (c *CoreServer) GetGroupName() string {
	return c.groupName
}

// 向所有与自己有关系的节点发送数据
func (c *CoreServer) NotifyConnectedNodes(b []byte) error {
	return nil
}

// 向指定节点发送数据
func (c *CoreServer) RequestSpecifiedNode(addr string, data []byte, f handleFunc) error {
	return nil
}

// 向指定节点组发送数据
func (c *CoreServer) RequestSpecifiedGroup(grpname string, data []byte) (ret []byte) {
	// 先找到组内的所有节点的conn
	// addr, err := net.ResolveTCPAddr("tcp", ":10000")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// conn, err := net.DialTCP("tcp", nil, addr)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// 向规则路由的节点发信息
	// head

	// data = BytesCombine(rpchead.Pack(), data)

	// size, _ := conn.Write(data)

	return nil
}
