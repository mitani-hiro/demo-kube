package handler

import (
	"net/http"

	"api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) Hoge(c *gin.Context) {
	res, err := h.usecase.CallProducer(c.Request.Context(), 999)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       res.Id,
		"name":     res.Name,
		"hostname": res.Hostname,
	})
}
