package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/interface/handler"
	"api/internal/usecase"
	"api/router"

	"common/cgrpc"
	"proto/pb"

	"google.golang.org/grpc"
)

const (
	defaultProducerAddr = "producer-service:50051"
	defaultPort         = "8080"
)

func main() {
	producerAddr := os.Getenv("PRODUCER_ADDR")
	if producerAddr == "" {
		producerAddr = defaultProducerAddr
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	useSharedConn := os.Getenv("USE_SHARED_CONN") != "false"

	userUC, err := newUserUsecase(producerAddr, useSharedConn)
	if err != nil {
		log.Fatalf("init usecase: %v", err)
	}

	userHandler := handler.NewUserHandler(userUC)
	healthHandler := handler.NewHealthHandler()
	r := router.NewRouter(userHandler, healthHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	log.Println("server shutdown complete")
}

// newUserUsecase はコネクション戦略を切り替えて usecase を組み立てる。
// 共有コネクション再利用（全リクエストが同一 Pod に張り付く）と
// リクエスト毎の新規接続（L4 LB により Pod 分散）を環境変数で切り替える。
func newUserUsecase(producerAddr string, useSharedConn bool) (usecase.UserUsecase, error) {
	if useSharedConn {
		conn, err := cgrpc.GetSharedClient(producerAddr)
		if err != nil {
			return nil, fmt.Errorf("get shared grpc client: %w", err)
		}
		return usecase.NewUserUsecase(pb.NewUserServiceClient(conn), nil), nil
	}

	factory := func() (*grpc.ClientConn, error) {
		return cgrpc.NewClient(producerAddr)
	}
	return usecase.NewUserUsecase(nil, factory), nil
}
