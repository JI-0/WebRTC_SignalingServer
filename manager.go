package main

import (
	"net/http"
	"sync"
)

type Manager struct {
	clients   ClientList
	streamers StreamerList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:   make(ClientList),
		streamers: make(StreamerList),
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {

}
