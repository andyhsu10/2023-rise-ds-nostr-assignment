package routers

import (
	"github.com/gin-gonic/gin"

	"distrise/internal/controllers"
)

func InitWsRouter(engine *gin.Engine, path string) error {
	ctl, err := controllers.GetController()
	if err != nil {
		return nil
	}
	group := engine.Group(path)

	group.GET("/", ctl.Ws.Home)

	return nil
}