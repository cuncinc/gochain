package main

import (
	"flag"
	. "gochain/net"
	"log"
)

func main() {
	port := flag.Int("port", 59100, "The port number of server to receive tx")
	serverMode := flag.Bool("server", false, "Run with server start")
	pingMode := flag.Bool("ping", false, "Run with ping")
	pingDuration := flag.Int("duration", 5, "ping duration (seconds)")
	flag.Parse()

	// 区块链节点连接中转服务器
	// var net Net = &SimpleNetwork{}
	net := SimpleNetwork{}
	net.Connect()
	defer net.Close()

	if *pingMode {
		net.Ping(*pingDuration)
	}

	// 交易信息接收服务器
	server := ServerNet{}
	if *serverMode {
		server.Listen(*port)
		defer server.Close()
	}

	for {
		select {
		case msg := <-net.MsgChan(): //其他节点发来的消息
			log.Println("[msg]", msg[:60])
		case serverMsg := <-server.MsgChan(): //交易客户端发来的消息
			log.Println("[server_msg]", serverMsg[:60])
			net.BroadcastMsg(serverMsg)
		}
	}
}
