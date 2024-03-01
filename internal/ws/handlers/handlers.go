package handlers

import (
	"github.com/vitalis-virtus/go-chat-websockets/internal/models"
	"github.com/vitalis-virtus/go-chat-websockets/internal/ws"
	"github.com/vitalis-virtus/go-chat-websockets/pkg/array"
)

func NewMessage(clients ws.WebSocketClientsPool, client ws.Client, message string) {
	broadcast(clients, client, models.WebSocketMessage{
		Type: MESSAGE,
		Content: map[string]string{
			"id":      client.ID(),
			"message": message,
		}})
}

func MemberJoin(clients ws.WebSocketClientsPool, client ws.Client) {
	broadcast(clients, client, models.WebSocketMessage{
		Type: MEMBER_JOIN,
		Content: map[string]string{
			"id":      client.ID(),
			"message": MEMBER_JOIN,
		},
	})
}

func MemberLeave(clients ws.WebSocketClientsPool, client ws.Client) {
	broadcast(clients, client, models.WebSocketMessage{
		Type: MEMBER_LEAVE,
		Content: map[string]string{
			"id": client.ID(),
		},
	})
}

func broadcast(clients ws.WebSocketClientsPool, client ws.Client, message models.WebSocketMessage) {
	// the message is broadcasted to every WebSocket client except the sender
	array.ForEach(
		array.Except(clients, func(item ws.Client) bool { return item.ID() == client.ID() }),
		func(item ws.Client) { item.Write(message) },
	)
}
