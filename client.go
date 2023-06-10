package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 8) / 10
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
	defer func() {
		// Cleanup connection
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		// log.Println(err)
		return
	}

	c.connection.SetReadLimit(4096)

	c.connection.SetPongHandler(c.pongHandler)

	for {

	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		// Cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			return

		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// log.Panicf("Failed to send message: %v", err)
				return
			}

		case <-c.timer.C:
			// Timeout the peer if not a streamer
			if !c.manager.isStreamer(c) {
				if err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "")); err != nil {
				}
				return
			}
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
