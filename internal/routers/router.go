package routers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

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
	engine.Use(otelgin.Middleware("distrise-relay"))

	err = InitWsRouter(engine, "/")
	if err != nil {
		return nil, err
	}

	return engine, nil
}
