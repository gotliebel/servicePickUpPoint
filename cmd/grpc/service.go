package grpc

import (
	"google.golang.org/grpc"
	"homework-1/internal/api"
	"homework-1/internal/service"
	"homework-1/pkg/api/proto/order/v1/order/v1"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	grpcPort = ":50051"
	httpPort = ":63342"
)

func NewGrpcServer(serv *service.Service) {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	order.RegisterOrderServer(s, &api.Server{
		Service: serv,
	})
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		gwServer := &http.Server{
			Addr: httpPort,
		}

		errHttp := gwServer.ListenAndServe()
		if errHttp != nil {
			log.Println(errHttp)
			return
		}

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	s.Stop()
}
