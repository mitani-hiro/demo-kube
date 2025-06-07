package router

import (
	"api/internal/interface/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	rg := r.Group("/api")

	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	rg.POST("/hoge", handler.Hoge)

	// Policy routes
	rg.POST("/policies", handler.CreatePolicy)
	rg.GET("/policies/:id", handler.GetPolicy)
	rg.PUT("/policies/:id", handler.UpdatePolicy)
	rg.DELETE("/policies/:id", handler.DeletePolicy)
	rg.GET("/policies", handler.ListPolicies)

	return r
}
