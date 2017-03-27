# Grapes

协议定义

| MessageHead | MessageBody |
|--------|--------|
|   4字节ID + 2字节body长度|   消息体     |


### 设计想法

### Guide服务器
负责负载均衡

### Connector服务器
负责保持与Client的连接，并把Client的消息转发给Service服务器进行处理。
Service处理完成之后，再通过Connector返回给Clinet

### Service服务器
后端服务器，负责处理业务逻辑。
相同Service服务器的负载分两种方式，等分方式和最大处理数方式
**等分方式** 客户端的请求按照Client的ID被平均路由到指定的服务器
**最大处理数方式** 客户端的请求尽可能的路由到同一台服务器，当Service处理的客户数到达一定程度后，新业务请求被路由到新的一台服务器上。

### Master服务器
用于管理所有的服务器，每个节点启动之后，自动注册到Master服务器，同时Master会把当前所有的服务器信息同步到新启动的节点上。

### 粘包处理