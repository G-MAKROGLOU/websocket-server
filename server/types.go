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
	onSend(data map[string]interface{})
	onSendError(ws *websocket.Conn, err error)
	onReceive(data map[string]interface{})
	onReceiveError(ws *websocket.Conn, err error)
}

// NOOPSocketServerEvents is a default struct that has no implementation for the Server events
type NOOPSocketServerEvents struct{}

func (n NOOPSocketServerEvents) onSend(data map[string]interface{}) {} 
func (n NOOPSocketServerEvents) onSendError(ws *websocket.Conn, err error) {} 
func (n NOOPSocketServerEvents) onReceive(data map[string]interface{}) {} 
func (n NOOPSocketServerEvents) onReceiveError(ws *websocket.Conn, err error) {} 

// ConfigFunc is a generic function that can be passed to New() to configure the Server
type ConfigFunc func(*SocketServer)
