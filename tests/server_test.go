package tests

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/G-MAKROGLOU/websocket-server/server"
	"github.com/stretchr/testify/assert"
)

func withFalsePort(config *server.SocketServer) {
	config.Port = "5000"
}

func withPort(config *server.SocketServer) {
	config.Port = ":5500"
}

func withPath(config *server.SocketServer) {
	config.Path = "/wss"
}

func withAnotherPath(config *server.SocketServer) {
	config.Path = "/wsss"
}

func TestWrongPort(t *testing.T){
	
	s := server.New(server.NOOPSocketServerEvents{}, withPath, withFalsePort)
	
	hasError := make(chan bool)
	
	go func() {
		err := s.Start()
		
		if assert.Error(t, err) {
			hasError <- assert.Equal(t, err.(*net.OpError), err)
			close(hasError)
		}
	}()
	
	time.Sleep(5 * time.Second)

	s.Stop()

	assert.Equal(t, true, <-hasError)
}

func TestWithDefaultConfig(t *testing.T) {

	s := server.New(server.NOOPSocketServerEvents{})
	
	go func() {
		time.Sleep(20 * time.Second)
		err := s.Stop()
		assert.Equal(t, nil, err)
	}()
		
	err := s.Start()
	
	assert.Equal(t, errors.New("http: Server closed"), err)
}

func TestWithCorrectCustomConfig(t *testing.T) {

	s := server.New(server.NOOPSocketServerEvents{}, withPort, withAnotherPath)
	
	go func() {
		time.Sleep(20 * time.Second)
		err := s.Stop()
		assert.Equal(t, nil, err)
	}()
		
	err := s.Start()
	
	assert.Equal(t, errors.New("http: Server closed"), err)
}
