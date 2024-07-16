# SERVER (server.go)

```go
package main

import (
	"fmt"
	"log"
	sockets "github.com/G-MAKROGLOU/websocket-server"

	"golang.org/x/net/websocket"
)

func main() {
	server := server.SocketServer{
		Path: "/ws",
		Port: ":5000",
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start socket server: ", err)
	}
}

```
