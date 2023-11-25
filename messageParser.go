package main

import (
	"log"
	"strings"
)

func parseMessage(c *Client, msg string) bool {
	parts := strings.Split(msg, "\n")
	switch parts[0] {
	case "U":
		if len(parts) == 3 {
			return upgrade(c, parts[1], parts[2])
		}
	case "R":
		if len(parts) == 2 {
			return requestStreamers(c, parts[1])
		}
	case "O":
		if len(parts) == 3 {
			return relayData(c, parts[1], parts[0], parts[2])
		}
	case "A":
		if len(parts) == 3 {
			return relayData(c, parts[1], parts[0], parts[2])
		}
	case "C":
		if len(parts) == 3 {
			return relayData(c, parts[1], parts[0], parts[2])
		}
	}

	// No function match
	return false
}

func upgrade(c *Client, username string, info string) bool {
	// Upgrade
	c.manager.upgradeUserToStreamer(c, username, info)
	return true
}

func requestStreamers(c *Client, info string) bool {
	username := c.manager.getClientUsername(c)
	if info == "0" {
		if streamer, err := c.manager.getRandomStreamer(); err == nil {
			// c.egress <- []byte("R\n" + streamer)
			streamer.egress <- []byte("R\n" + username)
		} else {
			c.egress <- []byte("<CK>Error no streamers")
			// return false
		}
	} else {
		if streamer, err := c.manager.getClientFromUsername(info); err == nil {
			streamer.egress <- []byte("R\n" + username)
		} else {
			c.egress <- []byte("<CK>Error streamer does not exist")
		}
	}
	return true
}

func relayData(c *Client, peer string, dataType string, data string) bool {
	log.Println("RELAY", dataType)
	username := c.manager.getClientUsername(c)
	peerClient, err := c.manager.getClientFromUsername(peer)
	if err != nil {
		// Could not find user, but keep the connection alive
		return true
	}
	// Relay data
	peerClient.egress <- []byte(dataType + "\n" + username + "\n" + data)
	return true
}
