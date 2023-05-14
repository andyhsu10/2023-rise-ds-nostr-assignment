package routers

import (
	"github.com/gin-gonic/gin"

	"distrise/internal/controllers"
	"distrise/internal/libs/relaylib"
)

func InitWsRouter(engine *gin.Engine, path string) error {
	ctl, err := controllers.GetController()
	if err != nil {
		return nil
	}
	group := engine.Group(path)

	hub := relaylib.NewHub()
	go hub.Run()

	group.GET("/", func(c *gin.Context) {
		ctl.Ws.Home(hub, c)
	})

	return nil
}