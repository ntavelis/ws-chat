package main

import (
	"context"
	"log"
	"strings"

	"github.com/coder/websocket"
)

func (s *Server) handleCommand(command string, client *Client) bool {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}

	cmd := parts[0]
	switch cmd {
	case "/exit":
		log.Printf("Client %s requested to exit", client.Username)
		s.Hub.Unregister <- client
		return true // Signal to disconnect
	default:
		// Send error message to client
		errorMsg := "Unknown command: " + cmd
		if err := client.Conn.Write(context.Background(), websocket.MessageText, []byte(errorMsg)); err != nil {
			log.Printf("Error sending error message to client %s: %v", client.ID, err)
		}
		return false
	}
}
