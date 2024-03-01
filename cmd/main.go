package main

import (
	"log"

	"github.com/vitalis-virtus/go-chat-websockets/internal/server"
)

func main() {
	log.Print("running server on port 8080")
	log.Fatal(server.New().ListenAndServe())
}
