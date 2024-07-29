//go:build !integration

package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"homework-1/internal/constant"
	"homework-1/internal/model"
	"homework-1/internal/service/mocks"
	"testing"
	"time"
)

func ParseDate(date string) time.Time {
	storedUntil, _ := time.Parse(constant.TimeFormat, date)
	return storedUntil
}
func TestService_AcceptFromCourier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		orderId    uint64
		clientId   uint64
		pack       string
		weight     float64
		date       string
		basicPrice float64
		error      bool
	}{
		{
			name:       "simple",
			date:       "12-12-2024",
			basicPrice: 122,
			orderId:    32,
			clientId:   33,
			pack:       "wrapping",
			weight:     8,
			error:      false,
		},
		{
			name:       "wrong format date",
			date:       "2024-12-12",
			basicPrice: 122,
			orderId:    34,
			clientId:   35,
			pack:       "wrapping",
			weight:     8,
			error:      true,
		},
		{
			name:       "wrong package name",
			date:       "12-12-2024",
			basicPrice: 122,
			orderId:    35,
			clientId:   36,
			pack:       "carton",
			weight:     8,
			error:      true,
		},
		{
			name:       "wrong weight",
			date:       "12-12-2024",
			basicPrice: 122,
			orderId:    35,
			clientId:   36,
			pack:       "box",
			weight:     31,
			error:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockData := mocks.NewMockstorage(ctrl)
			serv := New(mockData)
			ctx := context.Background()

			if tt.error {
				mockData.EXPECT().WriteOrderWithUniqueId(gomock.Any(), gomock.Any()).Times(0)
				err := serv.AcceptFromCourier(ctx, tt.orderId, tt.clientId, tt.date, tt.pack, tt.weight, tt.basicPrice)
				require.Error(t, err)
			} else {
				mockData.EXPECT().WriteOrderWithUniqueId(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				err := serv.AcceptFromCourier(ctx, tt.orderId, tt.clientId, tt.date, tt.pack, tt.weight, tt.basicPrice)
				require.NoError(t, err)
			}
		})
	}
}

func TestService_ReturnToCourier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		orderId uint64
		order   *model.Order
		error   bool
	}{
		{
			name:    "wrong stored until",
			orderId: 37,
			order: &model.Order{
				OrderID:     37,
				ClientID:    38,
				StoredUntil: ParseDate("12-12-2024"),
			},
			error: true,
		},
		{
			name:    "wrong taken at",
			orderId: 38,
			order: &model.Order{
				OrderID:  38,
				ClientID: 39,
				TakenAt:  time.Now(),
			},
			error: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockData := mocks.NewMockstorage(ctrl)
			serv := New(mockData)
			ctx := context.Background()

			if tt.error {
				mockData.EXPECT().GetOrder(gomock.Any(), tt.orderId).Return(tt.order, nil).Times(1)
				mockData.EXPECT().UpdateOrdersList(gomock.Any(), tt.order, gomock.Any()).Return(nil).Times(0)
				err := serv.ReturnToCourier(ctx, tt.orderId)
				require.Error(t, err)
			} else {
				mockData.EXPECT().GetOrder(gomock.Any(), tt.orderId).Return(tt.order, nil).Times(1)
				mockData.EXPECT().BeginTransaction().Times(1)
				mockData.EXPECT().UpdateOrdersList(ctx, &model.Order{
					OrderID:           tt.order.OrderID,
					ClientID:          tt.order.ClientID,
					StoredUntil:       tt.order.StoredUntil,
					TakenAt:           tt.order.CreatedAt,
					ReturnedToCourier: true,
				}, gomock.Any()).Times(1).Return(nil)
				err := serv.ReturnToCourier(ctx, tt.orderId)
				require.NoError(t, err)
			}
		})
	}
}

func TestService_ListForClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		clientId uint64
		limit    int
		error    bool
	}{
		{
			name:     "wrong client id",
			clientId: 0,
			error:    true,
		},
		{
			name:     "wrong limit",
			clientId: 1,
			limit:    -1,
			error:    true,
		},
		{
			name:     "simple",
			clientId: 1,
			limit:    1,
			error:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockData := mocks.NewMockstorage(ctrl)
			serv := New(mockData)
			ctx := context.Background()

			if tt.error {
				mockData.EXPECT().GetFilteredOrdersByClientId(gomock.Any(), tt.clientId, tt.limit).Times(0)
				err := serv.ListForClient(ctx, tt.clientId, tt.limit)
				require.Error(t, err)
			} else {
				mockData.EXPECT().GetFilteredOrdersByClientId(gomock.Any(), tt.clientId, tt.limit).Return([]*model.Order{}, nil).Times(1)
				err := serv.ListForClient(ctx, tt.clientId, tt.limit)
				require.NoError(t, err)
			}
		})
	}
}

func TestService_ListOfTakeBacks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		pageNumber int
		error      bool
	}{
		{
			name:       "simple",
			pageNumber: 1,
		},
		{
			name:       "wrong page number",
			pageNumber: 0,
			error:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockData := mocks.NewMockstorage(ctrl)
			serv := New(mockData)
			ctx := context.Background()

			if tt.error {
				mockData.EXPECT().GetTakebackList(gomock.Any(), tt.pageNumber).Times(0)
				err := serv.ListOfTakeBacks(ctx, tt.pageNumber)
				require.Error(t, err)
			} else {
				mockData.EXPECT().GetTakebackList(gomock.Any(), tt.pageNumber).Return([]*model.Order{}, nil).Times(1)
				err := serv.ListOfTakeBacks(ctx, tt.pageNumber)
				require.NoError(t, err)
			}
		})
	}

}

func TestService_ClientReturnsOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		clientId uint64
		orderId  uint64
		order    *model.Order
		error    bool
	}{
		{
			name:     "wrong client id",
			clientId: 0,
			error:    true,
		},
		{
			name:     "different client id",
			clientId: 40,
			orderId:  40,
			order: &model.Order{
				OrderID:  40,
				ClientID: 41,
			},
		},
		{
			name:     "wrong taken at",
			clientId: 42,
			orderId:  43,
			order: &model.Order{
				OrderID:  43,
				ClientID: 42,
				TakenAt:  time.Now().Add(-TimeToReturn),
			},
		},
		{
			name:     "wrong taken at",
			clientId: 44,
			orderId:  45,
			order: &model.Order{
				OrderID:  45,
				ClientID: 44,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockData := mocks.NewMockstorage(ctrl)
			serv := New(mockData)
			ctx := context.Background()

			if tt.error {
				mockData.EXPECT().GetOrder(gomock.Any(), tt.orderId).Return(&model.Order{}, nil).Times(0)
				err := serv.ClientReturnsOrder(ctx, tt.orderId, tt.clientId)
				require.Error(t, err)
			} else {
				mockData.EXPECT().GetOrder(gomock.Any(), tt.orderId).Return(tt.order, nil).Times(1)
				err := serv.ClientReturnsOrder(ctx, tt.orderId, tt.clientId)
				require.Error(t, err)
			}
		})
	}
}
