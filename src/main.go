package main

import (
	"fmt"
	. "gochain/net"
	. "gochain/tx"
	"log"
)

func main() {
	tran := NewTx("alice", "bob", 12, "0xfff")
	fmt.Println(tran)
	fmt.Println(tran.VerifyTx())

	tran = NewTxWithOptions(
		WithFrom("bob"),
		WithTo("alice"),
		WithAmount(34),
		WithSignature("bob's sig"),
	)

	fmt.Println(tran)
	fmt.Println(tran.VerifyTx())

	msgChan := make(chan string, 10)
	defer close(msgChan)

	var net Net = &Network{}
	net.Connect()
	defer net.Close()

	net.BroadcastMsg("hello, world")
	net.RegistMsgHandler(handler, msgChan)

	go func() {
		for {
			msg := <-msgChan // 从管道接收消息（阻塞等待消息）
			log.Println("[msg]", msg)
		}
	}()

	select {}

}

func handler(m string, msgChan ReceivedMsgChan) {
	// log.Printf("[recv]: [%s]", m)
	msgChan <- m
}
