package router

import (
	"api/internal/interface/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *handler.UserHandler, healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.Default()

	rg := r.Group("/api")
	rg.GET("/health", healthHandler.Health)
	rg.POST("/hoge", userHandler.Hoge)

	return r
}
