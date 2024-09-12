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

	var net Net = &Network{}
	net.Connect()
	defer net.Close()

	net.BroadcastMsg("hello, world")

	for {
		select {
		case msg := <-net.MsgChan():
			log.Println("[msg]", msg)
		}
	}
}
