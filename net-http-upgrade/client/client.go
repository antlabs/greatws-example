package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/antlabs/greatws"
	"github.com/antlabs/wsutil/opcode"
)

type handler struct{}

func (h *handler) OnOpen(c *greatws.Conn) {
	fmt.Printf("客户端连接成功\n")
}

var m *greatws.MultiEventLoop

func init() {

	m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(256), greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	m.Start()
}

func (h *handler) OnMessage(c *greatws.Conn, op greatws.Opcode, msg []byte) {
	// 如果msg的生命周期不是在OnMessage中结束，需要拷贝一份
	// newMsg := makc([]byte, len(msg))
	// copy(newMsg, msg)

	fmt.Printf("收到服务端消息:%s\n", msg)
	c.WriteMessage(op, msg)
	time.Sleep(time.Second)
}

func (h *handler) OnClose(c *greatws.Conn, err error) {
	fmt.Printf("客户端端连接关闭:%v\n", err)
}

func main() {
	c, err := greatws.Dial("ws://127.0.0.1:8080/", greatws.WithClientMultiEventLoop(m), greatws.WithClientCallback(&handler{}))
	if err != nil {
		fmt.Printf("连接失败:%v\n", err)
		return
	}

	c.WriteMessage(opcode.Text, []byte("hello"))
	c.ReadLoop()
}
