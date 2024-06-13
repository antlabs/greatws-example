package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/antlabs/greatws"
	reuseport "github.com/libp2p/go-reuseport"
)

var m *greatws.MultiEventLoop

func init() {

	m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(256), greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	m.Start()
}

type echoHandler struct{}

func (e *echoHandler) OnOpen(c *greatws.Conn) {
	// fmt.Printf("OnOpen: %p\n", c)
}

func (e *echoHandler) OnMessage(c *greatws.Conn, op greatws.Opcode, msg []byte) {
	if err := c.WriteTimeout(op, msg, 3*time.Second); err != nil {
		fmt.Println("write fail:", err)
	}
	// if err := c.WriteMessage(op, msg); err != nil {
	//  slog.Error("write fail:", err)
	// }
}

func (e *echoHandler) OnClose(c *greatws.Conn, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	slog.Error("OnClose:", errMsg)
}

type handler struct {
	m *greatws.MultiEventLoop
}

func (h *handler) echo(w http.ResponseWriter, r *http.Request) {
	c, err := greatws.Upgrade(w, r,
		greatws.WithServerReplyPing(),
		// greatws.WithServerDecompression(),
		greatws.WithServerIgnorePong(),
		greatws.WithServerCallback(&echoHandler{}),
		// greatws.WithServerEnableUTF8Check(),
		greatws.WithServerReadTimeout(5*time.Second),
		greatws.WithServerMultiEventLoop(h.m),
	)
	if err != nil {
		slog.Error("Upgrade fail:", "err", err.Error())
	}
	_ = c
}

func main() {
	var h handler

	h.m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(256), greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	h.m.Start()
	fmt.Printf("apiname:%s\n", h.m.GetApiName())

	mux := &http.ServeMux{}
	mux.HandleFunc("/autobahn", h.echo)

	rawTCP, err := reuseport.Listen("tcp", ":9001")
	// rawTCP, err := net.Listen("tcp", ":9001")
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}
	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))
}
