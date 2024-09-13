/*
* 网络服务接口
 */

package net

type Net interface {
	Connect()
	Close()
	BroadcastMsg(string)
	MsgChan() <-chan string // 网络消息管道的只读接口
}
