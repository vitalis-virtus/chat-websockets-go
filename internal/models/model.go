package models

type WebSocketMessage struct {
	Type    string            `json:"type"`
	Content map[string]string `json:"content"`
}
