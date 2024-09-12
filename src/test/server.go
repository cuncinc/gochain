//go:build server
// +build server

package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("[con_info] [connetd] from %s \n", r.RemoteAddr)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("[con_info] [upgrade]:", err)
		return
	}
	defer c.Close()

	// 将新连接添加到客户端列表中
	mutex.Lock()
	clients[c] = true
	mutex.Unlock()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("[read err]:", err)
			// 移除断开连接的客户端
			mutex.Lock()
			delete(clients, c)
			mutex.Unlock()
			break
		}

		// log.Printf("[recv] [%s] [%s]", r.RemoteAddr, message)
		// rsp := "got it"
		// err = c.WriteMessage(mt, []byte(rsp))
		// if err != nil {
		// 	log.Println("[write err]:", err)
		// 	break
		// }

		log.Println("[broadcast] from ", r.RemoteAddr)
		broadcastMessage(mt, message, c)
	}

	log.Println("[con_info] [discon] to", r.RemoteAddr)
}

// 广播消息到所有客户端
func broadcastMessage(mt int, message []byte, sender *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		// 不将消息发送回发送者
		if client != sender {
			err := client.WriteMessage(mt, message)
			if err != nil {
				log.Println("Error broadcasting to client:", err)
				client.Close()
				// 移除连接错误的客户端
				delete(clients, client)
			}
		}
	}
}

func main() {
	addrress := "0.0.0.0:58080"
	log.Println("Starting server at", addrress)
	http.HandleFunc("/", home)
	http.ListenAndServe(addrress, nil)
}
