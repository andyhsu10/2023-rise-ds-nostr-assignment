package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"distrise/internal/databases"
	"distrise/internal/libs/relaylib"
	"distrise/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsController interface {
	Home(ctx *gin.Context)
}

type wsController struct {
	ws    services.WsService
	rooms map[*relaylib.Room]bool
}

type RequestFilter struct {
	Ids      *[]string `json:"ids,omitempty"`
	Authors  *[]string `json:"authors,omitempty"`
	Kinds    *[]int    `json:"kinds,omitempty"`
	EventIds *[]string `json:"#e,omitempty"`
	Pubkeys  *[]string `json:"#p,omitempty"`
	Since    *int      `json:"since,omitempty"`
	Until    *int      `json:"until,omitempty"`
	Limit    *int      `json:"limit,omitempty"`
}

type RequestMessage struct {
	Action       string `json:"action"`
	Subscription string `json:"subscription"`
	Filters      *RequestFilter
}

func (r *RequestMessage) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&r.Action, &r.Subscription, &r.Filters}

	desiredLength := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}

	if g, e := len(tmp), desiredLength; g != e {
		return fmt.Errorf("wrong number of fields in RequestMessage: %d != %d", g, e)
	}

	return nil
}

func NewWsController() (WsController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}

	return &wsController{
		ws:    srv.Ws,
		rooms: make(map[*relaylib.Room]bool),
	}, nil
}

func (c *wsController) Home(ctx *gin.Context) {
	newCtx, span := otel.Tracer("wsController").Start(ctx, "Home")
	defer span.End()

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}

	db, err := databases.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	_, message, _ := conn.ReadMessage()
	var msg RequestMessage
	json.Unmarshal(message, &msg)

	if msg.Action == "REQ" {
		roomName := msg.Subscription

		room := c.findRoomByName(roomName)
		if room == nil {
			room = c.createRoom(roomName)
		}

		client := &relaylib.Client{
			Room: room,
			Conn: conn,
			Send: make(chan []byte, 256),
			DB:   db,
		}

		client.Room.Register <- client
		client.ID = GenUserId()
		client.Addr = conn.RemoteAddr().String()

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.WritePump(newCtx)
		go client.ReadPump(newCtx)

		msg := []byte("[\"EOSE\", \"" + roomName + "\"]")
		client.Send <- msg
	}
}

func (c *wsController) findRoomByName(name string) *relaylib.Room {
	var foundRoom *relaylib.Room
	for room := range c.rooms {
		if room.Name == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (c *wsController) createRoom(name string) *relaylib.Room {
	room := relaylib.NewRoom(name)
	go room.Run()
	c.rooms[room] = true

	return room
}

func GenUserId() string {
	uid := uuid.NewString()
	return uid
}
