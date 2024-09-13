package main

import (
	"flag"
	. "gochain/net"
	"log"
)

func main() {
	// tran := NewTx("alice", "bob", 12, "0xfff")
	// fmt.Println(tran)
	// fmt.Println(tran.VerifyTx())

	port := flag.Int("port", 59100, "The port number of server to receive tx")
	serverMode := flag.Bool("server", false, "Run with server start")
	pingMode := flag.Bool("ping", false, "Run with ping")
	pingDuration := flag.Int("duration", 5, "ping duration (seconds)")
	flag.Parse()

	// var net Net = &SimpleNetwork{}
	net := SimpleNetwork{}
	net.Connect()
	defer net.Close()

	if *pingMode {
		net.Ping(*pingDuration)
	}

	server := ServerNet{}
	if *serverMode {
		server.Listen(*port)
		defer server.Close()
	}

	// net.BroadcastMsg("hello, world")

	for {
		select {
		case msg := <-net.MsgChan():
			log.Println("[msg]", msg)
		case serverMsg := <-server.MsgChan():
			log.Println("[server_msg]", serverMsg)
			net.BroadcastMsg(serverMsg)
		}
	}
}
