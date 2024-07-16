# SERVER (server.go)

```go
package main

import (
	"fmt"
	"log"
	server "github.com/G-MAKROGLOU/websocket-server"

	"golang.org/x/net/websocket"
)

func main() {
	s := server.SocketServer{
		Path: "/ws",
		Port: ":5000",
	}

	if err := s.Start(); err != nil {
		fmt.Println("Failed to start socket server: ", err)
	}
}

```
