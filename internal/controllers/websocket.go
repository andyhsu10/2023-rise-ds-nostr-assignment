package controllers

import (
	"log"
	"net/http"
	"time"

	"distrise/internal/libs/relaylib"
	"distrise/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsController interface {
	Home(hub *relaylib.Hub, ctx *gin.Context)
}

type wsController struct {
	ws services.WsService
}

func NewWsController() (WsController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}
	return &wsController{ws: srv.Ws}, nil
}

func (controller *wsController) Home(hub *relaylib.Hub, ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &relaylib.Client{
		Hub: hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	client.Hub.Register <- client
	client.ID = GenUserId()
	client.Addr = conn.RemoteAddr().String()
	client.EnterAt = time.Now()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

	client.Send <- []byte("Welcome")
}

func GenUserId() string {
	uid := uuid.NewString()
	return uid
}
