package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
)

const (
	letterBytes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
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

func getNewToken(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Got something")

	// Upgrade from http
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create new client
	client := NewClient(conn, m)
	m.addClient(client)

	// Start client process goroutines
	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = "@" + getNewToken(64)
}

func (m *Manager) upgradeUserToStreamer(client *Client, username string, info string) {
	m.Lock()
	defer m.Unlock()

	client.timer.Stop()
	m.clients[client] = username
	m.streamers[client] = info
}

func (m *Manager) getClientUsername(client *Client) string {
	m.Lock()
	defer m.Unlock()

	return m.clients[client]
}

func (m *Manager) getClientFromUsername(username string) (*Client, error) {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.clients {
		if v == username {
			return k, nil
		}
	}
	return nil, errors.New("No user found")
}

func (m *Manager) isStreamer(c *Client) bool {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.streamers[c]; ok {
		return true
	}
	return false
}

func (m *Manager) getRandomStreamer() (*Client, error) {
	m.Lock()
	defer m.Unlock()

	if len(m.streamers) == 0 {
		return nil, errors.New("No streamers")
	}
	k := rand.Intn(len(m.streamers))
	for usr := range m.streamers {
		if k == 0 {
			return usr, nil
		}
		k--
	}
	return nil, errors.New("Func error")
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.streamers[client]; ok {
		delete(m.streamers, client)
	}
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "test.com":
		return false
	default:
		return true
	}
}
