package relaylib

import (
	"encoding/json"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			clientId := client.ID
			for client := range h.Clients {
				msg := []byte("some one join room (ID: " + clientId + ")")
				client.Send <- msg
			}

			h.Clients[client] = true

		case client := <-h.Unregister:
			clientId := client.ID
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			for client := range h.Clients {
				msg := []byte("some one leave room (ID: " + clientId + ")")
				client.Send <- msg
			}
		case userMessage := <-h.Broadcast:
			var data map[string][]byte
			json.Unmarshal(userMessage, &data)

			for client := range h.Clients {
				//prevent self receive the message
				if client.ID == string(data["id"]) {
					continue
				}
				select {
				case client.Send <- data["message"]:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}