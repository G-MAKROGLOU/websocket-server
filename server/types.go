package server

import (
	"net/http"

	"golang.org/x/net/websocket"
)

// SocketServer represents the socket server
type SocketServer struct {
	Path string
	Port string
	events SocketServerEvents
	server *http.Server
}

// SocketServerEvents holds the events that can be emitted on different server actions
type SocketServerEvents interface {
	OnSent(data map[string]interface{})
	OnSendError(ws *websocket.Conn, err error)
	OnReceiveError(ws *websocket.Conn, err error)
}

// NOOPSocketServerEvents is a default struct that has no implementation for the Server events
type NOOPSocketServerEvents struct{}

// OnSent is used when you want to perform a server action after a message is sent
func (n NOOPSocketServerEvents) OnSent(data map[string]interface{}) {} 

// OnSendError is used when you want to handle a send error. When the server fails to send a message
// to a connection, it disconnects the socket, clears all collections, and propagates the error to you through this
// function. You can then perform any further action, like logging, updating states etc. 
func (n NOOPSocketServerEvents) OnSendError(ws *websocket.Conn, err error) {} 

//OnReceiveError used when the server fails to receive from a socket for any other reason other
// a closed connection. The error is propagated to you through this function and you can perform any further
// action, like logging, updating states etc.
func (n NOOPSocketServerEvents) OnReceiveError(ws *websocket.Conn, err error) {} 

// ConfigFunc is a generic function that can be passed to New() to configure the Server
type ConfigFunc func(*SocketServer)
