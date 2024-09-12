package net

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Net interface {
	Connect()
	Close()
	BroadcastMsg(string)
	MsgChan() <-chan string // 网络消息管道的只读接口
}

/*简单网络实现，使用中心转发的方式，c连接服务器，服务器将所有消息转发到其余客户端*/
type Network struct {
	c       *websocket.Conn
	msgChan chan string //ws收到消息写入管道，主程序使用时读出
}

func (n *Network) Connect() {
	addr := "localhost:58080"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	var e error
	n.c, _, e = websocket.DefaultDialer.Dial(u.String(), nil) //升级到ws网络
	if e != nil {
		log.Fatal("[net] [con] [dial]", e)
	} else {
		log.Printf("[net] [con] [connected] to %s \n", u.String())
	}

	n.msgChan = make(chan string, 10) //初始化管道
	go func() {                       //在协程将收到的网络消息写入管道
		for {
			_, message, err := n.c.ReadMessage()
			if err != nil {
				log.Println("[net] [read_error]", err)
				// n.Close()
				return
			}
			n.msgChan <- string(message)
		}
	}()
}

func (n *Network) Close() {
	close(n.msgChan)
	n.c.Close()
	log.Println("[net] [con] [disconneted]")
}

func (n *Network) BroadcastMsg(m string) {
	err := n.c.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		log.Println("[net] [write]", err)
		return
	}
	log.Printf("[net] [broadcast] [%s]\n", m)
	// log.Printf("[net] [broadcast] [%s]\n", m[20:])
}

func (n *Network) MsgChan() <-chan string {
	return n.msgChan
}

////////////usage:

// func main() {
// 	var net Net = &Network{}
// 	net.Connect()
// 	defer net.Close()

// 	net.BroadcastMsg("hello, world")

// 	for {
// 		select {
// 		case msg := <-net.MsgChan():
// 			log.Println("[msg]", msg)
// 		}
// 	}
// }
