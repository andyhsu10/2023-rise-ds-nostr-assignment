package relaylib

import (
	"encoding/json"
)

// Room maintains the set of active clients and broadcasts messages to the clients.
type Room struct {
	Name string

	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
		case userMessage := <-r.Broadcast:
			var data map[string][]byte
			json.Unmarshal(userMessage, &data)

			for client := range r.Clients {
				// prevent sending message to self
				if client.ID == string(data["id"]) {
					continue
				}

				select {
				case client.Send <- data["message"]:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}
