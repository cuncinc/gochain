/*
* 网络服务：与其他节点相连
* 实现Net接口
 */

package net

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

/*简单网络实现，使用中心转发的方式，c连接服务器，服务器将所有消息转发到其余客户端*/
type SimpleNetwork struct {
	c       *websocket.Conn
	msgChan chan string //ws收到消息写入管道，主程序使用时读出
}

func (n *SimpleNetwork) Connect() {
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

func (n *SimpleNetwork) Close() {
	close(n.msgChan)
	n.c.Close()
	log.Println("[net] [con] [disconneted]")
}

func (n *SimpleNetwork) BroadcastMsg(m string) {
	err := n.c.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		log.Println("[net] [write]", err)
		return
	}
	log.Printf("[net] [broadcast] %s\n", m)
	// log.Printf("[net] [broadcast] %s\n", m[20:])
}

func (n *SimpleNetwork) MsgChan() <-chan string {
	return n.msgChan
}

func (n *SimpleNetwork) Ping(second int) {
	ticker := time.NewTicker(time.Second * time.Duration(second))
	go func() {
		for {
			select {
			case <-ticker.C:
				m := "ping"
				err := n.c.WriteMessage(websocket.TextMessage, []byte(m))
				if err != nil {
					log.Println("[net] [write]", err)
					return
				}
			}
		}
	}()
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
