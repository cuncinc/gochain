/*
* 服务器，用于接收终端传来的交易请求
 */

package net

import (
	"fmt"
	"log"
	"net"
)

type ServerNet struct {
	listener net.Listener
	msgChan  chan string
}

func (n *ServerNet) Listen(port int) {
	addr := "0.0.0.0:" + fmt.Sprintf("%d", port)
	var e error
	n.listener, e = net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("[net] [server] [listen_error]", e)
	}
	log.Println("[net] [server] [listen] on", addr)

	n.msgChan = make(chan string, 10)

	go func() {
		for {
			conn, err := n.listener.Accept()
			if err != nil {
				log.Println("[net] [server] [accepting]", err)
				continue
			}
			go handleConnection(conn, n.msgChan) // 启动一个 goroutine 处理每个连接
		}
	}()
}

func (n *ServerNet) Close() {
	close(n.msgChan)
	n.listener.Close()
	log.Println("[net] [server] [close]")
}

func (n *ServerNet) MsgChan() <-chan string {
	return n.msgChan
}

func handleConnection(conn net.Conn, msgChan chan<- string) {
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf) //todo：tcp是流式读取，数据可能会超出缓冲区，这里并没有做处理
		if err != nil {
			log.Println("[net] [server] [receive] Error reading:", err)
			return
		}

		msg := string(buf[:n])
		log.Printf("[net] [server] [receive] %d \n", len(msg))
		msgChan <- msg

		// Echo back the received data
		rsp := fmt.Sprintf("recived %d", len(msg))
		_, err = conn.Write([]byte(rsp))
		if err != nil {
			log.Println("[net] [server] [response] Error writing:", err)
			return
		}
	}
}
