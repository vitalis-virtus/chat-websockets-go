package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vitalis-virtus/go-chat-websockets/internal/ws"
	wshandlers "github.com/vitalis-virtus/go-chat-websockets/internal/ws/handlers"

	"github.com/vitalis-virtus/go-chat-websockets/pkg/array"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	clients ws.WebSocketClientsPool
	mu      sync.Mutex
)

// WebSocketConnect handlers responsible for the WebSocket connection
func WebSocketConnect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go startClient(c, ws)
}

func startClient(ctx context.Context, wsConn *websocket.Conn) {
	client := ws.NewWebSocketClient(wsConn)

	mu.Lock()
	clients = append(clients, client)
	mu.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("error: %v", err)
		}

		mu.Lock()
		defer mu.Unlock()
		clients = array.Except(clients, func(item ws.Client) bool { return item.ID() == client.ID() })
		client.Close()
	}()

	client.Launch(ctx)

	for {
		select {
		case msg, ok := <-client.Listen():
			if !ok {
				return
			} else {
				switch msg.Type {
				case wshandlers.MESSAGE:
					wshandlers.NewMessage(clients, client, msg.Content["message"])
				default:
					log.Printf("unknown message type: %v", msg.Type)
					return
				}

			}
		case err := <-client.Error():
			log.Printf("websocket error: %v", err)
		case <-client.Done():
			wshandlers.MemberLeave(clients, client)
			return
		}
	}
}
