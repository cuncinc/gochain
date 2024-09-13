//go:build client
// +build client

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	. "gochain/tx"
	"net"
)

func main() {
	// port := flag.Int("port", 59100, "connect tx server")
	// host := flag.String("host", "127.0.0.1", "host")
	flag.Parse()
	// addr := fmt.Sprintf("%s:%d", host, port)
	addr := "127.0.0.1:59100"
	fmt.Println(addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server", addr)

	var slice []Tx
	for i := 0; i < 5; i++ {
		tran := NewTx("alice", "bob", i, "0xfff")
		slice = append(slice, *tran)
	}
	message, _ := json.Marshal(slice)

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("send ok")

}
