# SERVER (serverevents.go)

```go
package main

import "fmt"

type CustomEvents struct {}

func (c CustomEvents) OnStartError(err error) {
    fmt.Println("[SERVER] Failed to start ", err)
}

func (c CustomEvents) OnSend(data map[string]interface{}) {
    b, _ := json.MarshalIndent(data, "", " ")

    fmt.Println("[SERVER] SENT: ", string(b))
}

func (c CustomEvents) OnSendError(err error) {
    fmt.Println("[SERVER] Failed to send ", err)
}
```

# SERVER (server.go)

```go
package main

import (
    "fmt"
    server "github.com/G-MAKROGLOU/websocket-server"
)

func main() {
    port := ":5000"
    path := "/ws"
    s := server.NewSocketServer(path, port, CustomEvents{});

    if err := s.Start(); err != nil {
        fmt.Println("Failed to start socket server: ", err)
    }
}
```
