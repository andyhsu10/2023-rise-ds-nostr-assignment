package routers

import (
	"github.com/gin-gonic/gin"

	"distrise/internal/middlewares"
)

var (
	routerInstance *gin.Engine
)

func GetRouter() (*gin.Engine, error) {
	if routerInstance == nil {
		instance, err := newRouter()
		if err != nil {
			return nil, err
		}
		routerInstance = instance
	}
	return routerInstance, nil
}

func newRouter() (*gin.Engine, error) {
	engine := gin.Default()
	middleware, err := middlewares.GetMiddleware()
	if err != nil {
		return nil, err
	}

	engine.Use(middleware.Cors.Cors())

	err = InitWsRouter(engine, "/")
	if err != nil {
		return nil, err
	}

	return engine, nil
}
