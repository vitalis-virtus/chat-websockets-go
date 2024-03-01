package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitalis-virtus/chat-websockets-go/internal/handler"
)

func New() *http.Server {
	r := gin.Default()

	r.GET("/ws", handler.WebSocketConnect)

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
