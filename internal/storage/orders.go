package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"homework-1/internal/constant"
	"homework-1/internal/metrics"
	"homework-1/internal/model"
	"log"
)

func CheckTypeError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return model.ErrOrderAlreadyExists
		}
	}
	return err
}

func (s *Storage) WriteOrderWithUniqueId(ctx context.Context, newOrder *model.Order) error {
	err := s.MakeTransaction(func(tx *sql.Tx) error {
		span, _ := opentracing.StartSpanFromContext(ctx, "WriteOrderWithUniqueId")
		defer span.Finish()
		query := `INSERT INTO orders (order_id, client_id, created_at, stored_until, taken_at, takeback_time, returned)
           VALUES ($1, $2, $3, $4, $5, $6, $7)`
		err := tx.QueryRow(query, newOrder.OrderID, newOrder.ClientID, newOrder.CreatedAt, newOrder.StoredUntil, newOrder.TakenAt, newOrder.TakeBackAt, newOrder.ReturnedToCourier).Err()
		if err != nil {
			return CheckTypeError(err)
		}
		query = `INSERT INTO characteristic (order_id, package, weight, price)
           VALUES ($1, $2, $3, $4)`
		err = tx.QueryRow(query, newOrder.OrderID, newOrder.Pack, newOrder.Weight, newOrder.Price).Err()
		if err != nil {
			return err
		}
		metrics.ChangesOrdersInCache.Inc()
		metrics.OrdersIssued.Inc()
		return nil
	})
	if err != nil {
		return err
	}
	s.Cache.Del(newOrder.OrderID)
	return nil
}

func (s *Storage) GetOrdersByIdList(ctx context.Context, idList []uint64, tx *sql.Tx) ([]*model.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetOrdersByIdList")
	defer span.Finish()
	res := make([]*model.Order, 0, len(idList))
	partialRes := make([]*model.Order, 0, len(idList))
	idsInCache := make([]uint64, 0)
	for _, id := range idList {
		if value, ok := s.Cache.Get(id); ok {
			if !value.ReturnedToCourier {
				partialRes = append(partialRes, value)
				idsInCache = append(idsInCache, id)
				metrics.CachedOrdersIssued.Inc()
				metrics.OrdersIssued.Inc()
			}
		}
	}
	if len(partialRes) == len(idList) {
		return partialRes, nil
	}
	query := `SELECT order_id, client_id, created_at, stored_until, taken_at, takeback_time, returned
	           FROM orders WHERE not returned and order_id = ANY($1) and order_id <> ALL($2)`
	rows, err := tx.Query(query, pq.Array(idList), pq.Array(idsInCache))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderID, &order.ClientID, &order.CreatedAt, &order.StoredUntil, &order.TakenAt, &order.TakeBackAt, &order.ReturnedToCourier); err != nil {
			log.Printf(err.Error())
		}
		res = append(res, &order)
		metrics.OrdersIssued.Inc()
		s.Cache.Set(order.OrderID, &order)
	}
	return res, nil
}

func (s *Storage) GetFilteredOrdersByClientId(ctx context.Context, id uint64, limit int) (res []*model.Order, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetFilteredOrdersByClientId")
	defer span.Finish()
	res = make([]*model.Order, 0)
	query := `SELECT order_id, client_id, created_at, stored_until, taken_at, takeback_time, returned
	           FROM orders WHERE client_id = $1 and not returned and taken_at = '0001-01-01 00:00:00' order by created_at limit $2`
	rows, err := s.Db.Query(query, id, limit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderID, &order.ClientID, &order.CreatedAt, &order.StoredUntil, &order.TakenAt, &order.TakeBackAt, &order.ReturnedToCourier); err != nil {
			log.Printf(err.Error())
		}
		res = append(res, &order)
		metrics.OrdersIssued.Inc()
		s.Cache.Set(order.OrderID, &order)
	}
	return res, nil
}
func (s *Storage) UpdateOrdersList(ctx context.Context, list []*model.Order, tx *sql.Tx) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UpdateOrdersList")
	defer span.Finish()
	for _, order := range list {
		query := `UPDATE orders SET client_id = $2, created_at = $3, stored_until = $4, taken_at = $5, takeback_time = $6, returned = $7 WHERE order_id = $1`
		result, err := tx.Exec(query, order.OrderID, order.ClientID, order.CreatedAt, order.StoredUntil, order.TakenAt, order.TakeBackAt, order.ReturnedToCourier)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return model.ErrNoRowsAffected
		}
		metrics.ChangesOrdersInCache.Inc()
		metrics.OrdersIssued.Inc()
		s.Cache.Del(order.OrderID)
	}
	return nil
}

func (s *Storage) GetTakebackList(ctx context.Context, pageNumber int) (res []*model.Order, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetTakebackList")
	defer span.Finish()
	query := `SELECT order_id, client_id, created_at, stored_until, taken_at, takeback_time, returned 
           FROM orders where taken_at != '0001-01-01 00:00:00' offset $1 limit $2`
	rows, err := s.Db.Query(query, (pageNumber-1)*constant.TakebackPageSize, constant.TakebackPageSize)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	res = make([]*model.Order, 0)
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderID, &order.ClientID, &order.CreatedAt, &order.StoredUntil, &order.TakenAt, &order.TakeBackAt, &order.ReturnedToCourier); err != nil {
			log.Printf(err.Error())
		}
		res = append(res, &order)
		metrics.OrdersIssued.Inc()
		s.Cache.Set(order.OrderID, &order)
	}
	return res, nil
}

func (s *Storage) GetOrder(ctx context.Context, orderId uint64) (*model.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetOrder")
	defer span.Finish()
	if value, ok := s.Cache.Get(orderId); ok {
		metrics.CachedOrdersIssued.Inc()
		metrics.OrdersIssued.Inc()
		return value, nil
	}
	query := `SELECT order_id, client_id, created_at, stored_until, taken_at, takeback_time, returned
	           FROM orders WHERE order_id = $1`
	rows, err := s.Db.Query(query, orderId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderID, &order.ClientID, &order.CreatedAt, &order.StoredUntil, &order.TakenAt, &order.TakeBackAt, &order.ReturnedToCourier); err != nil {
			log.Printf(err.Error())
		}
		metrics.OrdersIssued.Inc()
		s.Cache.Set(order.OrderID, &order)
		return &order, nil
	}
	return nil, model.ErrDoesntExists
}
