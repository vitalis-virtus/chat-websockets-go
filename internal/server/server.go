package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *http.Server {
	r := gin.Default()

	r.GET("/ws", handler.webSocketConnect)

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
