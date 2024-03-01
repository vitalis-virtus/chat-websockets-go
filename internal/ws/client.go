package ws

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/vitalis-virtus/chat-websockets-go/internal/models"
)

// Responsible for simplifying work with WebSocket
type Client interface {
	// Returns a unique identifier of the WebSocket connection
	ID() string
	// Launches the client, so it starts listening to new messages
	Launch(ctx context.Context)
	// Sends a message `m` back to the client
	Write(m models.WebSocketMessage) error
	// Closes WebSocket connection
	Close()
	// Returns a channel with incoming messages
	Listen() <-chan models.WebSocketMessage
	// Returns a channel that closes when work is done (WebSocket connection closed or should be closed)
	Done() <-chan interface{}
	// Returns a channel with errors that happened during WebSocket listening
	Error() <-chan error
}

type client struct {
	id       string
	ws       *websocket.Conn
	messages chan models.WebSocketMessage
	errors   chan error
	done     chan interface{}
	sync.Mutex
	sync.Once
}
