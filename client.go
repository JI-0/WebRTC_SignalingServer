package main

import (
	"log"
	"strings"
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
		messageType, payload, err := c.connection.ReadMessage()
		log.Println("New:")
		log.Println(string(payload))

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// log.Println("Read error ", err)
			}
			break
		}

		// Check if text
		if messageType == 1 || messageType == 2 {
			message := string(payload)
			// //Prevent string with '$'
			// if strings.Contains(message, "$") {
			// 	log.Println("Message contains $")
			// 	break
			// }
			// payloadParts := strings.Split(message, "\n")
			// if len(payloadParts) < 2 {
			// 	// log.Println("Message not in the right format")
			// 	break
			// }
			// if !parseMessage(c, payloadParts[0], payloadParts) {
			// 	break
			// }
			if !parseMessage(c, message) {
				break
			}
		}
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
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(3000, "Connection lost unexpectedly")); err != nil {
					// log.Println("Closed connection: ", err)
				}
				return
			}

			if string(message) == "<CK>OK" {
				if err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "")); err != nil {
				}
				return
			} else if strings.HasPrefix(string(message), "<CK>") {
				if err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(3000, strings.TrimPrefix(string(message), "<CK>"))); err != nil {
				}
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("Message sent: ", string(message))

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
