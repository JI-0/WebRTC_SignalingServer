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

	srv0 := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	srv1 := &http.Server{
		Addr:    ":3001",
		Handler: mux,
	}

	go log.Fatal(srv0.ListenAndServe())
	log.Fatal(srv1.ListenAndServe())

	println("Shutting down signaling server")
}
