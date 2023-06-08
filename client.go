package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// Clients map with Client: identifier
type ClientList map[*Client]string

// Streamer map with Client: info
type StreamerList map[*Client]string

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	// Egress to prevent concurrent writes to the websocket conn
	egress chan []byte
	timer  time.Timer
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
		timer:      *time.NewTimer(time.Second * 20),
	}
}

func (c *Client) readMessages() {

}

func (c *Client) writeMessages() {

}
