package relaylib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"distrise/internal/models"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Room *Room

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	User

	DB *gorm.DB
}

type User struct {
	ID   string
	Addr string
}

type Request struct {
	Action string
	Data   string
}

type Event struct {
	ID        string   `json:"id"`
	Pubkey    string   `json:"pubkey"`
	CreatedAt int      `json:"created_at"`
	Kind      int      `json:"kind"`
	Tags      []string `json:"tags"`
	Content   string   `json:"content"`
	Sig       string   `json:"sig"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var data []json.RawMessage
		err = json.Unmarshal([]byte(message), &data)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if len(data) < 2 {
			fmt.Println("Invalid JSON structure")
			continue
		}

		var req Request
		err = json.Unmarshal(data[0], &req.Action)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if req.Action == "EVENT" && len(data) == 3 {
			err = json.Unmarshal(data[1], &req.Data)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// Save to DB
			if err := crdbgorm.ExecuteTx(context.Background(), c.DB, nil,
				func(tx *gorm.DB) error {
					if err := c.DB.Create(&models.CoreEvent{Name: req.Data, Data: string(data[2])}).Error; err != nil {
						fmt.Println("Error:", err)
						return err
					}
					return nil
				},
			); err != nil {
				// For information and reference documentation, see:
				//   https://www.cockroachlabs.com/docs/stable/error-handling-and-troubleshooting.html
				fmt.Println("Error:", err)
			}

			msg := map[string][]byte{
				"message": []byte("[\"EVENT\", \"" + req.Data + "\", " + string(data[2]) + "]"),
				"id":      []byte(c.ID),
			}
			userMessage, _ := json.Marshal(msg)
			c.Room.Broadcast <- userMessage
		} else if req.Action == "REQ" {
			err = json.Unmarshal(data[1], &req.Data)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			msg := []byte("[\"EOSE\", \"" + req.Data + "\"]")
			c.Send <- msg
		} else if req.Action == "CLOSE" {
			c.Room.Unregister <- c
			c.Conn.Close()
			break
		}

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
