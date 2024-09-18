//go:build client
// +build client

/*
* block chain node client
* 区块链网络节点测试
* 在net包中实现
 */

package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type ReceivedMsgChan chan<- string
type ReceivedMsgHandler func(string, ReceivedMsgChan)

type Net interface {
	Connect()
	Close()
	BroadcastMsg(string)
	RegistMsgHandler(ReceivedMsgHandler, ReceivedMsgChan)
}

type Network struct {
	c *websocket.Conn
}

func (n *Network) Connect() {
	addr := "localhost:58080"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	var e error
	n.c, _, e = websocket.DefaultDialer.Dial(u.String(), nil)
	if e != nil {
		log.Fatal("[net] [con] [dial]", e)
	} else {
		log.Printf("[net] [con] [connected] to %s \n", u.String())
	}
}
func (n *Network) Close() {
	// // Cleanly close the connection by sending a close message and then
	// // waiting (with timeout) for the server to close the connection.
	// err := n.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// if err != nil {
	// 	log.Println("[net] [con] [close] write close:", err)
	// 	return
	// } else {
	// 	log.Println("[net] [con] [close] [send ok]")
	// }
	n.c.Close()
	log.Println("[net] [con] [disconneted]")
}
func (n *Network) BroadcastMsg(m string) {
	err := n.c.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		log.Println("write:", err)
		return
	}
	log.Printf("[net] [broadcast] [%s]\n", m)
	// log.Printf("[net] [broadcast] [%s]\n", m[20:])
}

func (n *Network) RegistMsgHandler(handler ReceivedMsgHandler, msgChan ReceivedMsgChan) {
	go func() { //在协程里处理收到的消息
		for {
			_, message, err := n.c.ReadMessage()
			if err != nil {
				log.Println("[net] [read_error]", err)
				return
			}
			handler(string(message), msgChan)
		}
	}()
}

/////此处的设计有点麻烦且绕，chan在使用处注册
////////////usage:

func handler(m string, msgChan ReceivedMsgChan) {
	// log.Printf("[recv]: [%s]", m)
	msgChan <- m
}

func main() {

	msgChan := make(chan string, 10)
	defer close(msgChan)

	var net Network = Network{}
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
