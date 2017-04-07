package core

import (
	"flag"
	log "github.com/kermitbu/grapes/log"
	proto "github.com/kermitbu/grapes/proto"
	"io"
	"net"
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

func (c *CoreServer) initHandleJoinRequest() {

	if c.allHandlerFunc == nil {
		c.allHandlerFunc = make(map[uint16]handleFunc)
	}

	c.allHandlerFunc[1] = func(req *GRequest, res *GResponse) {
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
}

func (c *CoreServer) SetNodeType(t NodeType) {
	c.t = t
}

func (c *CoreServer) GetNodeType() NodeType {
	return c.t
}

func (c *CoreServer) InitComplete() {
	c.initHandleJoinRequest()

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
	allHandlerName    map[uint16]string
	allHandlerFunc    map[uint16]handleFunc
	allClientConnects map[string]*net.TCPConn
	t                 NodeType
}

type ServiceCollection map[string][]ServiceNode
type ServiceNode struct {
	name    string
	connect *net.Conn
}

var services = make(ServiceCollection)

func (c *CoreServer) GetServiceNodes() ServiceCollection {
	return services
}

func (c *CoreServer) GetServiceNodeByName(name string) []ServiceNode {
	nodes := make([]ServiceNode, 0)

	return nodes
}
