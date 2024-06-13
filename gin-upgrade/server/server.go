package main

import (
	"fmt"
	"log/slog"

	"github.com/antlabs/greatws"
	"github.com/gin-gonic/gin"
)

type handler struct{}

var m *greatws.MultiEventLoop

func init() {

	m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(256), greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	m.Start()
}

func (h *handler) OnOpen(c *greatws.Conn) {
	fmt.Printf("服务端收到一个新的连接")
}

func (h *handler) OnMessage(c *greatws.Conn, op greatws.Opcode, msg []byte) {
	// 如果msg的生命周期不是在OnMessage中结束，需要拷贝一份
	// newMsg := makc([]byte, len(msg))
	// copy(newMsg, msg)

	fmt.Printf("收到客户端消息:%s\n", msg)
	c.WriteMessage(op, msg)
	// os.Stdout.Write(msg)
}

func (h *handler) OnClose(c *greatws.Conn, err error) {
	fmt.Printf("服务端连接关闭:%v\n", err)
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		con, err := greatws.Upgrade(c.Writer, c.Request, greatws.WithServerMultiEventLoop(m), greatws.WithServerCallback(&handler{}))
		if err != nil {
			return
		}
		con.StartReadLoop()
	})
	r.Run()
}
