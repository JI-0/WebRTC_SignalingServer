package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	println("Starting signaling server")
	manager := NewManager()

	mux := http.NewServeMux()
	mux.HandleFunc("/", manager.serveWS)
	tlsConf := &tls.Config{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/the.testingwebrtc.com/fullchain.pem",
				"/etc/letsencrypt/live/the.testingwebrtc.com/privkey.pem")
			if err != nil {
				log.Println("Failed to load TLS certificate!")
				return nil, err
			}
			return &cert, nil
		},
	}

	srv0 := &http.Server{
		Addr:      ":3000",
		Handler:   mux,
		TLSConfig: tlsConf,
	}
	srv1 := &http.Server{
		Addr:      ":3001",
		Handler:   mux,
		TLSConfig: tlsConf,
	}

	go srv0.ListenAndServeTLS("", "")
	srv1.ListenAndServeTLS("", "")

	println("Shutting down signaling server")
}
