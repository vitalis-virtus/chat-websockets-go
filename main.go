package main

import (
	"log"

	"github.com/vitalis-virtus/chat-websockets-go/internal/server"
)

func main() {
	log.Fatal(server.New().ListenAndServe())
}
