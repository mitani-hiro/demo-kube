package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"common/cgrpc"
	"proto/pb"

	"github.com/gin-gonic/gin"
)

func Hoge(c *gin.Context) {
	fmt.Println("Hoge processing...")

	conn, err := cgrpc.NewClient("producer-service:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 999})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Hoge end...")
	c.JSON(http.StatusOK, res)
}
