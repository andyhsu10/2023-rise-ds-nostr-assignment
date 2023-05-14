package controllers

import (
	"fmt"
	"log"
	"net/http"

	"distrise/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	ws services.WsService
}

func NewWsController() (WsController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}
	return &wsController{ws: srv.Ws}, nil
}

func (controller *wsController) Home(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	for {
		// Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		// If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}

		// Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
