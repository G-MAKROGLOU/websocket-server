# INSTALLATION

```go
go get github.com/G-MAKROGLOU/websocket-server
```


# SERVER (serverevents.go)

You can have access to the data being sent/received to/from the server by implementing the following interface.
In case you don't want the server to just handle traffic and nothing more, you can use the provided NOOPSocketServerEvents{}

```go
package main

import "fmt"

type CustomEvents struct {}

func (c CustomEvents) onSend(data map[string]interface{}) {
    b, _ := json.MarshalIndent(data, "", " ")

    fmt.Println("[server] sent: ", string(b))
}

func (c CustomEvents) onSendError(ws *websocket.Conn, err error) {
    fmt.Println("[server] failed to send ", err)
}

func (c CustomEvents) onReceive(data map[string]interface{}) {
    b, _ := json.MarshalIndent(data, "", " ")

    fmt.Println("[server] received: ", string(b))
}

func (c CustomEvents) onReceiveError(ws *websocket.Conn, err error) {
    fmt.Println("[server] failed to receive ", err)
}

```

# SERVER (server.go)

The default configuration for the websocket server is the following:

```go
&SocketServer {
    Path: "/ws",
    Port: ":3000"
}

```

Which you can ovveride by providing any number of functions that receive a pointer to *SocketServer. For example:

To override the port:

```go
func withPort(s *SocketServer) {
    s.Port = ":6000"
}

```

To override the path:

```go
func withPath(s *SocketServer) {
    s.Path = "/wss"
}

```


Then you can start a server with your custom configuration:

```go
package main

import (
    "log"
    server "github.com/G-MAKROGLOU/websocket-server"
)

func main() {
    s := server.New(CustomEvents{}, withPort, withPath);

    if err := s.Start(); err != nil {
        log.Fatalln("Failed to start socket server: ", err)
    }
}

```

or start a server with the default configuration:

```go
package main

import (
    "log"
    server "github.com/G-MAKROGLOU/websocket-server"
)

func main() {
    s := server.New(CustomEvents{});

    if err := s.Start(); err != nil {
        log.Fatalln("Failed to start socket server: ", err)
    }
}
```

To use the server with any other client other than
```github.com/G-MAKROGLOU/websocket-client``` that supports those functionalities out of the box, all your messages will have to
include an extra property named: ```GmWsType``` with a value of: ```join```, ```leave```, ```disconnect```, ```multicast```, or ```broadcast```.
When ```GmWsType``` is ```join```, ```leave```, or ```multicast```, you need to include an extra property named ```GmWsRoom``` with the name of the room
that you want to perform each action. Those properties are removed before the exchange and the payload arrives cleaned up to the clients.
