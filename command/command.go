package main

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/exec"

	_ "embed"

	"github.com/antlabs/greatws"
)

//go:embed index.html
var indexHTML []byte

var m *greatws.MultiEventLoop

func init() {

	m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(256), greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	m.Start()
}
func executeCommand(cmd string) []byte {
	var stdout, stderr bytes.Buffer
	command := exec.Command("sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return []byte(fmt.Sprintf("Error: %s\n%s", err.Error(), stderr.String()))
	}
	return stdout.Bytes()
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(indexHTML)
}
func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := greatws.Upgrade(w, r, greatws.WithServerMultiEventLoop(m), greatws.WithServerOnMessageFunc(func(c *greatws.Conn, op greatws.Opcode, data []byte) {

			log.Printf("Received command: %s", data)
			result := executeCommand(string(data))
			err := c.WriteMessage(greatws.Text, []byte(result))
			if err != nil {
				log.Printf("Write error: %s", err)
				c.Close()
				return
			}
			log.Printf("Sent response: %s", result)
		}))
		if err != nil {
			log.Printf("Upgrade error: %s", err)
			return
		}
		conn.StartReadLoop()
	})

	log.Println("Server started on ws://localhost:8080")
	log.Println("open http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %s", err)
	}
}
