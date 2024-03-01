package ws

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vitalis-virtus/chat-websockets-go/internal/models"
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
	// ID returns a unique identifier of the WebSocket connection
	ID() string
	// Launch launches the client, so it starts listening to new messages
	Launch(ctx context.Context)
	// Write sends a message `m` back to the client
	Write(m models.WebSocketMessage) error
	// Close closes WebSocket connection
	Close()
	// Listen returns a channel with incoming messages
	Listen() <-chan models.WebSocketMessage
	// Done returns a channel that closes when work is done (WebSocket connection closed or should be closed)
	Done() <-chan interface{}
	// Error returns a channel with errors that happened during WebSocket listening
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

func (c *client) ID() string {
	return c.id
}

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

func (c *client) Write(m models.WebSocketMessage) error {
	//TODO implement me
	panic("implement me")
}

func (c *client) Close() {
	//TODO implement me
	panic("implement me")
}

func (c *client) Listen() <-chan models.WebSocketMessage {
	//TODO implement me
	panic("implement me")
}

func (c *client) Done() <-chan interface{} {
	//TODO implement me
	panic("implement me")
}

func (c *client) Error() <-chan error {
	//TODO implement me
	panic("implement me")
}

func (c *client) write(message int) {
	//TODO implement me
	panic("implement me")
}

// read is responsible for listening to incoming messages. It publishes them to the channel (the channel is returned by Listen method). The goroutine is finished when the context is done or when the read operation returns an erro
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
			c.send(websocket.PingMessage)
		}
	}
}

func (c *client) handleError(err error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) send(message int) {
	//TODO implement me
	panic("implement me")
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
