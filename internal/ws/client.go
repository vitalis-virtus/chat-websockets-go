package ws

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vitalis-virtus/go-chat-websockets/internal/models"
	"sync"
	"time"
)

const (
	maxMessageSize = 1024
	writeWait      = 5 * time.Second
	pingPeriod     = 5 * time.Second
	pongWait       = 10 * time.Second
)

// Client responsible for simplifying work with WebSocket
type Client interface {
	ID() string
	Launch(ctx context.Context)
	Write(m models.WebSocketMessage) error
	Close() error
	Listen() <-chan models.WebSocketMessage
	Done() <-chan interface{}
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

type WebSocketClientsPool []Client

// ID returns a unique identifier of the WebSocket connection
func (c *client) ID() string {
	return c.id
}

// Launch launches the client, so it starts listening to new messages
func (c *client) Launch(ctx context.Context) {
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	c.Do(func() { go c.launch(ctx) })
}

func (c *client) launch(ctx context.Context) {
	var wg sync.WaitGroup

	cancellation, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.write(websocket.CloseMessage)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.read(cancellation)
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.ping(cancellation)
		cancel()
	}()

	wg.Wait()
	c.done <- struct{}{}
}

// Write sends a message `m` back to the client
func (c *client) Write(m models.WebSocketMessage) error {
	c.Lock()
	defer c.Unlock()

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(m)
}

// Close closes WebSocket connection
func (c *client) Close() error {
	close(c.messages)
	return c.ws.Close()
}

// Listen returns a channel with incoming messages
func (c *client) Listen() <-chan models.WebSocketMessage {
	return c.messages
}

// Done returns a channel that closes when work is done (WebSocket connection closed or should be closed)
func (c *client) Done() <-chan interface{} {
	return c.done
}

// Error returns a channel with errors that happened during WebSocket listening
func (c *client) Error() <-chan error {
	return c.errors
}

// write private method sends service messages (ping, close connection)
func (c *client) write(messageType int) {
	c.Lock()
	defer c.Unlock()

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.ws.WriteMessage(messageType, nil); err != nil {
		c.handleError(err)
	}
}

// read is responsible for listening to incoming messages. It publishes them to the channel (the channel is returned by Listen method). The goroutine is finished when the context is done or when the read operation returns an error
func (c *client) read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg models.WebSocketMessage
			err := c.ws.ReadJSON(&msg)
			if err != nil {
				c.handleError(err)
				return
			}

			c.messages <- msg
		}
	}
}

// ping is responsible for sending periodical ping. The goroutine is finished when the context is done
func (c *client) ping(ctx context.Context) {
	pingTicker := time.NewTicker(pingPeriod)

	for {
		select {
		case <-ctx.Done():
			return

		case <-pingTicker.C:
			c.write(websocket.PingMessage)
		}
	}
}

func (c *client) handleError(err error) {
	var closeError *websocket.CloseError
	if errors.As(err, &closeError) {
		return
	}

	if errors.Is(err, websocket.ErrCloseSent) {
		return
	}

	c.errors <- err
}

func NewWebSocketClient(ws *websocket.Conn) Client {
	return &client{
		id:       uuid.NewString(),
		ws:       ws,
		messages: make(chan models.WebSocketMessage),
		errors:   make(chan error),
		done:     make(chan interface{}),
	}
}
