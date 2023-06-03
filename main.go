package main

import (
	"log"
	"net/http"
)

func main() {
	println("Starting signaling server")
	manager := NewManager()

	mux := http.NewServeMux()
	mux.HandleFunc("/", manager.serveWS)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())

	println("Shutting down signaling server")
}
