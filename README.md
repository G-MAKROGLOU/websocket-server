# INSTALLATION

```go
go get github.com/G-MAKROGLOU/websocket-server@latest
```


# SERVER (serverevents.go)

You can have access to the data being sent/received to/from the server by implementing the following interface.
In case you want the server to just handle traffic and nothing more, you can use the provided NOOPSocketServerEvents{}

```go
package main

// Events implements SocketServerEvents interface
type Events struct{}

// OnSent is used when you want to perform a server action after a message is sent
func (e Events) OnSent(data map[string]interface{}) {
	// TODO: do stuff here...
	// logging, update db, internal states etc.
}

// OnSendError is used when you want to handle a send error. When the server fails to send a message
// to a connection, it disconnects the socket, clears all collections, and propagates the error to you through this
// function. You can then perform any further action, like logging, updating states etc. 
func (e Events) OnSendError(ws *websocket.Conn, err error) {
    // TODO: do stuff here...
    // logging, update db, internal states etc.
}

//OnReceiveError used when the server fails to receive from a socket for any other reason other
// a closed connection. The error is propagated to you through this function and you can perform any further
// action, like logging, updating states etc.
func (e Events) OnReceiveError(ws *websocket.Conn, err error) {
    // TODO: do stuff here...
    // logging, update db, internal states etc.
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
    "log/slog"
    server "github.com/G-MAKROGLOU/websocket-server"
)

func main() {
    s := server.New(Events{}, withPort, withPath)

    err := s.Start();
		
    if err == http.ErrServerClosed {
	slog.Info("server stopped")
    } else {
	msg := fmt.Sprintf("unexpected server error: %s", err.Error())
	slog.Error(msg)
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
    s := server.New(Events{})

    err := s.Start();
		
    if err == http.ErrServerClosed {
	slog.Info("server stopped")
    } else {
	msg := fmt.Sprintf("unexpected server error: %s", err.Error())
	slog.Error(msg)
    }
}
```

To use the server with any other client other than
```github.com/G-MAKROGLOU/websocket-client``` that supports those functionalities out of the box, all your messages will have to
include an extra property named: ```GmWsType``` with a value of: ```join```, ```leave```, ```disconnect```, ```multicast```, or ```broadcast```.
When ```GmWsType``` is ```join```, ```leave```, or ```multicast```, you need to include an extra property named ```GmWsRoom``` with the name of the room
that you want to perform each action. Those properties are removed before the exchange and the payload arrives cleaned up to the clients.
