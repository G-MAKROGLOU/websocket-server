package server

import (
	"fmt"
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

// SocketServer represents the socket server
type SocketServer struct {
	Path string
	Port string
	events SocketServerEvents
}

// SocketServerEvents holds the events that can be emitted on different server actions
type SocketServerEvents interface {
	OnStartError(err error)
	OnSend(data map[string]interface{})
	OnSendError(err error)
}

// NewSocketServer creates a new SocketServer instance and returns a pointer to it
func (s *SocketServer) NewSocketServer(path string, port string, events SocketServerEvents) *SocketServer {
	return &SocketServer {
		Path: path,
		Port: port,
		events: events,
	}
}

// Start starts the socket server
func (s *SocketServer) Start() {
	
	http.Handle(s.Path, websocket.Handler(s.handler))
	
	fmt.Println("Starting socket server on port ", s.Port, ", path ", s.Path)
	
	if err := http.ListenAndServe(s.Port, nil); err != nil {
		s.events.OnStartError(err)
	}
}

func (s *SocketServer) handler(ws *websocket.Conn) {
	sessID := strings.Split(ws.Request().Header.Get("Cookie"), "=")[1]
	
	allConsMutex.Lock()
	allCons[sessID] = ws
	allConsMutex.Unlock()
	
	for {
		var msg map[string]interface{}
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break
		}
		
		msgType := msg["Gm_Ws_Type"].(string)
		delete(msg, "Gm_Ws_Type")

		if msgType != "" && msgType == "gm_ws_join" {
			roomName := msg["Gm_Ws_Room"].(string)
			addToRoom(roomName, ws)
		}

		if msgType != "" && msgType == "gm_ws_leave" {
			roomName := msg["Gm_Ws_Room"].(string)
			removeFromRoom(roomName, ws)
		}

		if msgType != "" && msgType == "gm_ws_disconnect" {
			disconnect(sessID, ws)
		}

		if msgType != "" && msgType == "gm_ws_multicast" {
			roomName := msg["Gm_Ws_Room"].(string)
			delete(msg, "Gm_Ws_Room")
			s.sendJSONTo(ws, sessID, roomName, msg)
		}

		if msgType != "" && msgType == "gm_ws_broadcast" {
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
				s.events.OnSendError(err)
				return
			}
			s.events.OnSend(data)
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
				s.events.OnSendError(err)
				return
			}
			s.events.OnSend(data)
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

