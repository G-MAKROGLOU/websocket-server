package server

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

var allCons = make(map[string]*websocket.Conn)
var allConsMutex = sync.Mutex{}

var rooms = make(map[string][]*websocket.Conn)
var roomsMutex = sync.Mutex{}

var connToRoom = make(map[*websocket.Conn]string)
var connToRoomMutex = sync.Mutex{}

// New creates a new SocketServer instance and returns a pointer to it
func New(events SocketServerEvents, configs ...ConfigFunc) *SocketServer {
	config := defaultConfig()
	config.events = events
	
	for _, fn := range configs {
		fn(&config)
	}
	
	return &config
}

// Start starts the socket server
func (s *SocketServer) Start() error {
	server := &http.Server{Addr: s.Port}

	s.server = server

	http.Handle(s.Path, websocket.Handler(s.jsonHandler))

	return server.ListenAndServe()
}

// Stop stops the server
func (s *SocketServer) Stop() error {
	return s.server.Shutdown(context.Background())
}

func defaultConfig() SocketServer {
	return SocketServer {
		Path: "/ws",
		Port: ":3000",
	}
}

func (s *SocketServer) jsonHandler(ws *websocket.Conn) {
	sessID := strings.Split(ws.Request().Header.Get("Cookie"), "=")[1]
	
	allConsMutex.Lock()
	allCons[sessID] = ws
	allConsMutex.Unlock()
	
	for {
		var msg map[string]interface{}
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			s.events.onReceiveError(ws, err)
		}
		
		msgType := msg["GmWsType"].(string)
		delete(msg, "GmWsType")

		if msgType != "" && msgType == "join" {
			roomName := msg["GmWsRoom"].(string)
			addToRoom(roomName, ws)
		}

		if msgType != "" && msgType == "leave" {
			roomName := msg["GmWsRoom"].(string)
			removeFromRoom(roomName, ws)
		}

		if msgType != "" && msgType == "disconnect" {
			disconnect(sessID, ws)
		}

		if msgType != "" && msgType == "multicast" {
			roomName := msg["GmWsRoom"].(string)
			delete(msg, "GmWsRoom")
			s.sendJSONTo(ws, sessID, roomName, msg)
		}

		if msgType != "" && msgType == "broadcast" {
			s.sendJSON(ws, sessID, msg)
		}
	}
}

// Send sends a broadcast message to all connected sockets on the server
func (s *SocketServer) sendJSON(ws *websocket.Conn, sessID string, data map[string]interface{}) {
	for _, socket := range allCons {
		if  socket != ws {
			err := websocket.JSON.Send(socket, data)
			if err != nil {
				disconnect(sessID, socket)
				s.events.onSendError(ws, err)
				return
			}
			s.events.onSend(data)
		}
	}
}

// SendTo sends a unitcast/multicast message to all sockets in a room
func (s *SocketServer) sendJSONTo(ws *websocket.Conn, sessID string, roomName string, data map[string]interface{}) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	sockets := rooms[roomName]

	for _, socket := range sockets {
		if socket != ws {
			err := websocket.JSON.Send(socket, data)
			if err != nil {
				disconnect(sessID, ws)
				s.events.onSendError(ws, err)
				return
			}
			s.events.onSend(data)
		}
	}
}

func addToRoom(roomName string, ws *websocket.Conn) {
	connToRoomMutex.Lock()
	defer connToRoomMutex.Unlock()

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	connToRoom[ws] = roomName

	_, exists := rooms[roomName]

	if !exists {
		rooms[roomName] = []*websocket.Conn{ ws }
	}

	if exists {
		rooms[roomName] = append(rooms[roomName], ws)
	}
}

// RemoveClient removes a client from a room
func removeFromRoom(roomName string, ws *websocket.Conn) {
	connToRoomMutex.Lock()
	defer connToRoomMutex.Unlock()

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	sockets := rooms[roomName]

	newSockets := []*websocket.Conn{}

	for _, s := range sockets {
		if s != ws {
			newSockets = append(newSockets, s)
		}
	}

	rooms[roomName] = newSockets

	delete(connToRoom, ws)
}

// disconnects a client form the server and removes the client form any possible rooms
func disconnect(sessID string, ws *websocket.Conn) {
	allConsMutex.Lock()
	defer allConsMutex.Unlock()

	connToRoomMutex.Lock()
	defer connToRoomMutex.Unlock()

	room, exists := connToRoom[ws]
	
	if exists {
		removeFromRoom(room, ws)
		delete(connToRoom, ws)
	}

	delete(allCons, sessID)

	ws.Close()
}

