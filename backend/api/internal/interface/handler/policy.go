package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"common/cgrpc"
	"proto/pb"

	"github.com/gin-gonic/gin"
)

type CreatePolicyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Rules       string `json:"rules" binding:"required"`
}

type UpdatePolicyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Rules       string `json:"rules" binding:"required"`
}

func CreatePolicy(c *gin.Context) {
	var req CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := cgrpc.NewClient("policy:50052")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewPolicyServiceClient(conn)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := client.CreatePolicy(ctx, &pb.CreatePolicyRequest{
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Status.Code != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status.Message})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func GetPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
		return
	}

	conn, err := cgrpc.NewClient("policy:50052")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewPolicyServiceClient(conn)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := client.GetPolicy(ctx, &pb.GetPolicyRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Status.Code == 5 { // NOT_FOUND
		c.JSON(http.StatusNotFound, gin.H{"error": res.Status.Message})
		return
	}

	if res.Status.Code != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status.Message})
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdatePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
		return
	}

	var req UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := cgrpc.NewClient("policy:50052")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewPolicyServiceClient(conn)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := client.UpdatePolicy(ctx, &pb.UpdatePolicyRequest{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Status.Code == 5 { // NOT_FOUND
		c.JSON(http.StatusNotFound, gin.H{"error": res.Status.Message})
		return
	}

	if res.Status.Code != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status.Message})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeletePolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid policy ID"})
		return
	}

	conn, err := cgrpc.NewClient("policy:50052")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewPolicyServiceClient(conn)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := client.DeletePolicy(ctx, &pb.DeletePolicyRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Status.Code == 5 { // NOT_FOUND
		c.JSON(http.StatusNotFound, gin.H{"error": res.Status.Message})
		return
	}

	if res.Status.Code != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status.Message})
		return
	}

	c.JSON(http.StatusOK, res)
}

func ListPolicies(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	conn, err := cgrpc.NewClient("policy:50052")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewPolicyServiceClient(conn)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := client.ListPolicies(ctx, &pb.ListPoliciesRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res.Status.Code != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status.Message})
		return
	}

	c.JSON(http.StatusOK, res)
}