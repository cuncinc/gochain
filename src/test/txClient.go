//go:build client
// +build client

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	. "gochain/tx"
	. "gochain/wallet"
	"net"
	"time"

	"math/rand"
)

func getWallet(filename string) *Wallet {
	var wallet *Wallet
	if filename == "" {
		wallet = NewWallet()
	} else {
		wallet, _ = LoadWalletFromFile(filename)
	}
	return wallet
}

func main() {
	// port := flag.Int("port", 59100, "connect tx server")
	// host := flag.String("host", "127.0.0.1", "host")
	flag.Parse()

	addr := "127.0.0.1:59100"
	fmt.Println(addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server", addr)

	alice := getWallet("alice.key")
	bob := getWallet("bob.key")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				var sender *Wallet
				var reciever *Wallet
				if rand.Intn(2) == 0 {
					sender = alice
					reciever = bob
				} else {
					sender = bob
					reciever = alice
				}

				t := NewTx(sender.GetAddress(), reciever.GetAddress(), rand.Intn(2000)+10, sender.GetPublicKeyAsString())
				t.Signature, _ = sender.Sign(t.Data)

				message, _ := json.Marshal(t)

				_, err = conn.Write([]byte(message))
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "send tx")
			}

		}
	}()

	select {}
}
