package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"distrise/internal/configs"
)

type CorsMiddleware interface {
	Cors() gin.HandlerFunc
}

type corsMiddleware struct {
	Config cors.Config
}

func NewCorsMiddleware() (CorsMiddleware, error) {
	serverConfig := configs.GetConfig().Server

	config := cors.DefaultConfig()
	config.AllowOrigins = serverConfig.Cors
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host"}
	config.AllowWildcard = true
	return &corsMiddleware{Config: config}, nil
}

func (middleware *corsMiddleware) Cors() gin.HandlerFunc {
	return cors.New(middleware.Config)
}
