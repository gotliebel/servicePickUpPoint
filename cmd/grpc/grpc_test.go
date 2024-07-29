package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework-1/internal/api"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"homework-1/pkg/api/proto/order/v1/order/v1"
	"log"
	"testing"
)

const bufSize = 1024 * 1024

func TestAcceptFromCourier(t *testing.T) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	st, _ := storage.New()
	serv := service.New(st)
	order.RegisterOrderServer(s, &api.Server{Service: serv})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	conn, err := grpc.NewClient("localhost"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	client := order.NewOrderClient(conn)

	_, err = client.AcceptFromCourier(context.Background(), &order.OrderRequest{OrderId: 1, ClientId: 2, StoredUntil: "12-12-", Pack: nil})
	require.Error(t, err)
}
