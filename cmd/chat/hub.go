package main

import (
	"fmt"
	"log"
	"sync"
)

type BroadcastMessage struct {
	Data   []byte
	Sender *Client
}

func (b BroadcastMessage) String() string {
	return fmt.Sprintf("%s: %s", b.Sender.Username, string(b.Data))
}

type RegisterRequest struct {
	Client   *Client
	Response chan bool
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan RegisterRequest
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
	Mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan RegisterRequest),
		Unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMessage),
		Mutex:      sync.Mutex{},
	}
}

func (h *Hub) GetClientsNumber() int {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return len(h.Clients)
}

func (h *Hub) Run() {
	for {
		select {
		case req := <-h.Register:
			h.Mutex.Lock()
			h.Clients[req.Client.ID] = req.Client
			h.Mutex.Unlock()

			// Send confirmation
			req.Response <- true

			// Notify all other clients that this user connected
			connectionMsg := BroadcastMessage{
				Data:   []byte("Connected"),
				Sender: req.Client,
			}
			h.sendMessageToClients(connectionMsg)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				disconnectionMessage := BroadcastMessage{
					Data:   []byte("Disconnected"),
					Sender: client,
				}
				h.sendMessageToClients(disconnectionMessage)
			}
		case broadcastMsg := <-h.Broadcast:
			log.Printf("Broadcasting message: %s\n", broadcastMsg)
			h.sendMessageToClients(broadcastMsg)
		}
	}
}

func (h *Hub) sendMessageToClients(msg BroadcastMessage) {
	for _, client := range h.Clients {
		if client.ID == msg.Sender.ID {
			continue
		}
		select {
		case client.Send <- []byte(msg.String()):
		default:
			close(client.Send)
			delete(h.Clients, client.ID)
		}
	}
}
