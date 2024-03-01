package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitalis-virtus/go-chat-websockets/internal/handlers"
)

func New() *http.Server {
	r := gin.Default()

	r.GET("/ws", handlers.WebSocketConnect)

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
