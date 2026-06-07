package cgrpc

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient は呼び出し毎に新しいgRPCコネクションを作成する。
// K8s ClusterIP (L4 LB) がTCPコネクション確立時に振り分けるため、
// 毎回新規接続すればPodが分散する。
func NewClient(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

var (
	sharedConn *grpc.ClientConn
	once       sync.Once
	connErr    error
)

// GetSharedClient はシングルトンのgRPCコネクションを返す。
// gRPC公式推奨のコネクション再利用パターン。HTTP/2多重化により性能は良いが、
// K8s ClusterIP (L4 LB) では全リクエストが同一Podに張り付く。
func GetSharedClient(target string) (*grpc.ClientConn, error) {
	once.Do(func() {
		sharedConn, connErr = grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	return sharedConn, connErr
}
