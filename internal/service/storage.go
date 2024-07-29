//go:generate mockgen -source ./storage.go -destination=./mocks/storage.go -package=mocks
package service

import (
	"context"
	"database/sql"
	"homework-1/internal/model"
)

type storage interface {
	WriteOrderWithUniqueId(ctx context.Context, m *model.Order) error

	GetOrdersByIdList(ctx context.Context, idList []uint64, tx *sql.Tx) ([]*model.Order, error)
	GetFilteredOrdersByClientId(ctx context.Context, id uint64, limit int) ([]*model.Order, error)

	GetOrder(ctx context.Context, orderId uint64) (*model.Order, error)
	UpdateOrdersList(ctx context.Context, list []*model.Order, tx *sql.Tx) error

	GetTakebackList(ctx context.Context, pageNumber int) ([]*model.Order, error)
	BeginTransaction() (*sql.Tx, error)
	MakeTransaction(fn func(tx *sql.Tx) error) error
}
