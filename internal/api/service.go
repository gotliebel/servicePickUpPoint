package api

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework-1/internal/service"
	"homework-1/pkg/api/proto/order/v1/order/v1"
)

type Server struct {
	*service.Service
	order.UnimplementedOrderServer
}

func (s *Server) AcceptFromCourier(ctx context.Context, req *order.OrderRequest) (*emptypb.Empty, error) {
	if errValidate := req.ValidateAll(); errValidate != nil {
		return nil, status.Error(codes.InvalidArgument, errValidate.Error())
	}
	if err := s.Service.AcceptFromCourier(ctx, req.OrderId, req.ClientId, req.StoredUntil, req.Pack.String(), float64(req.Weight), float64(req.Price)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ReturnToCourier(ctx context.Context, req *order.CourierRequest) (*emptypb.Empty, error) {
	if err := s.Service.ReturnToCourier(ctx, req.OrderId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GiveToClient(ctx context.Context, req *order.ClientRequest) (*emptypb.Empty, error) {
	if err := s.Service.GiveToClient(ctx, &req.OrdersId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListForClient(ctx context.Context, req *order.ClientIdRequest) (*order.ListResponse, error) {
	if err := s.Service.ListForClient(ctx, req.ClientId, int(req.Limit)); err != nil {
		return &order.ListResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &order.ListResponse{List: make([]*order.OrderRequest, 0)}, nil
}

func (s *Server) ClientReturnsOrder(ctx context.Context, req *order.ShortOrderRequest) (*emptypb.Empty, error) {
	if err := s.Service.ClientReturnsOrder(ctx, req.OrderId, req.ClientId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListOfTakeBacks(ctx context.Context, req *order.TakebacksRequest) (*order.ListResponse, error) {
	if err := s.Service.ListOfTakeBacks(ctx, int(req.Page)); err != nil {
		return &order.ListResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &order.ListResponse{List: make([]*order.OrderRequest, 0)}, nil
}
