package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"strings"

	"github.com/coder/websocket"
)

type Server struct {
	Hub *Hub
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Received websocket request")
	conn, error := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"chat"},
	})
	if error != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.CloseNow()

	log.Println("Websocket connection established")
	client := NewClient(
		rand.Text(), // Generate a random ID for the client
		r.URL.Query().Get("user"),
		conn,
	)

	// Synchronous registration
	response := make(chan bool)
	s.Hub.Register <- RegisterRequest{
		Client:   client,
		Response: response,
	}
	<-response // Wait for registration to complete

	log.Printf("Client registered, id: %s, username: %s\n", client.ID, client.Username)
	client.SendWelcomeMessage()
	client.InformHowManyUsersAreInTheRoom(s.Hub.GetClientsNumber())

	defer func() {
		log.Printf("Unregistering client, id: %s, username: %s\n", client.ID, client.Username)
		s.Hub.Unregister <- client
	}()
	for {
		ctx := context.Background()
		_, msg, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusNoStatusRcvd {
				log.Println("Client closed the connection gracefully")
				return
			}
			stc := websocket.CloseStatus(err)
			log.Println("Error reading message:", err)
			log.Println("Closing connection with status:", stc)
			return
		}
		log.Printf("Received message: %s\n", msg)

		msgStr := string(msg)
		if strings.HasPrefix(msgStr, "/") {
			if s.handleCommand(msgStr, client) {
				return // Exit the loop if command requests disconnection
			}
		} else {
			s.Hub.Broadcast <- BroadcastMessage{
				Data:   msg,
				Sender: client,
			}
		}
	}
}
