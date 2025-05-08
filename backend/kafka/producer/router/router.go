package router

import (
	"producer/interface/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	rg := r.Group("/api")

	rg.POST("/hoge", handler.Hoge)

	return r
}
