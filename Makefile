all:
	go build -o reuseport.server ./reuseport/server/server.go
	go build -o gin-ws.server ./gin-upgrade/server/server.go
	go build -o gin-ws.client ./gin-upgrade/client/client.go
	go build -o net-http-ws.server ./net-http-upgrade/server/server.go
	go build -o net-http-ws.client ./net-http-upgrade/client/client.go
	go build -o with-html.server ./with-html/server/server.go
	go build -o dirwatch.out ./dirwatch/dirwatch.go
	go build -o command.out ./command/command.go
	go build -o chatroom.out ./chatroom/chatroom.go