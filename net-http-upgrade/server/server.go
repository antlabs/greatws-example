package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/antlabs/greatws"
)

var m *greatws.MultiEventLoop

func init() {

	m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(256), greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	m.Start()
}

type echoHandler struct{}

func (e *echoHandler) OnOpen(c *greatws.Conn) {
	fmt.Println("OnOpen:\n")
}

func (e *echoHandler) OnMessage(c *greatws.Conn, op greatws.Opcode, msg []byte) {
	fmt.Printf("OnMessage: %s, %v\n", msg, op)
	if err := c.WriteTimeout(op, msg, 3*time.Second); err != nil {
		fmt.Println("write fail:", err)
	}
}

func (e *echoHandler) OnClose(c *greatws.Conn, err error) {
	fmt.Println("OnClose: %v", err)
}

// echo测试服务
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := greatws.Upgrade(w, r, greatws.WithServerReplyPing(),
		// greatws.WithServerDecompression(),
		// greatws.WithServerIgnorePong(),
		greatws.WithServerCallback(&echoHandler{}),
		greatws.WithServerReadTimeout(5*time.Second),
		greatws.WithServerMultiEventLoop(m),
	)
	if err != nil {
		fmt.Println("Upgrade fail:", err)
		return
	}

	c.StartReadLoop()
}

func main() {
	http.HandleFunc("/", echo)

	http.ListenAndServe(":8080", nil)
}
