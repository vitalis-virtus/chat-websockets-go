package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "HELLO WORLD")
	})

	log.Fatal(http.ListenAndServe(":3030", nil))
}
