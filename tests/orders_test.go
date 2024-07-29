//go:build !integration

package tests

import (
	"context"
	"github.com/bojanz/currency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-1/internal/model"
	"homework-1/internal/packages"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"log"
	"testing"
	"time"
)

func Init() (*service.Service, *storage.Storage) {
	st, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	serv := service.New(st)
	return serv, st
}

func DeleteOrder(st *storage.Storage, id uint64) {
	_ = st.Db.QueryRow(`DELETE FROM orders WHERE order_id = $1`, id).Err()
	_ = st.Db.QueryRow(`DELETE FROM characteristic WHERE order_id = $1`, id).Err()
}

func TestWriteOrderWithUniqueId(t *testing.T) {
	t.Parallel()
	serv, st := Init()
	defer st.CloseStorage()
	ctx := context.Background()
	amount, _ := currency.NewAmount("123", packages.Currency)
	order := &model.Order{
		OrderID:           1,
		ClientID:          2,
		CreatedAt:         time.Now().UTC(),
		StoredUntil:       time.Now().UTC().AddDate(1, 0, 0),
		TakenAt:           time.Time{},
		TakeBackAt:        time.Time{},
		ReturnedToCourier: false,
		Pack:              "bag",
		Weight:            8,
		Price:             amount,
	}
	defer DeleteOrder(st, 1)

	err := serv.Storage.WriteOrderWithUniqueId(ctx, order)
	require.NoError(t, err)

	insertedOrder, err := serv.Storage.GetOrder(ctx, order.OrderID)
	require.NoError(t, err)
	assert.Equal(t, order.ClientID, insertedOrder.ClientID)
	assert.Equal(t, order.CreatedAt, insertedOrder.CreatedAt)
	assert.Equal(t, order.StoredUntil, insertedOrder.StoredUntil)
	assert.Equal(t, order.TakenAt, insertedOrder.TakenAt)
	assert.Equal(t, order.TakeBackAt, insertedOrder.TakeBackAt)
	assert.Equal(t, order.ReturnedToCourier, insertedOrder.ReturnedToCourier)
}

func TestGetOrdersByIdList(t *testing.T) {
	t.Parallel()
	serv, st := Init()
	defer st.CloseStorage()
	ctx := context.Background()
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 3, Pack: "box"})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 4, Pack: "box"})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 5, Pack: "box"})

	defer DeleteOrder(st, 3)
	defer DeleteOrder(st, 4)
	defer DeleteOrder(st, 5)
	idList := []uint64{3, 4, 5}
	tx, _ := st.Db.Begin()
	defer tx.Rollback()

	res, err := serv.Storage.GetOrdersByIdList(ctx, idList, tx)
	require.NoError(t, err)

	assert.Equal(t, 3, len(res))
	for _, order := range res {
		assert.Contains(t, idList, order.OrderID)
	}
}

func TestGetFilteredOrdersByClientId(t *testing.T) {
	t.Parallel()
	serv, st := Init()
	defer st.CloseStorage()
	ctx := context.Background()
	limit := 10
	id := 5
	lst := []*model.Order{{OrderID: 6, ClientID: 5, Pack: "box"},
		{OrderID: 7, ClientID: 6, Pack: "box"},
		{OrderID: 8, ClientID: 7, Pack: "box"},
		{OrderID: 9, ClientID: 5, Pack: "box", ReturnedToCourier: true}}
	for _, order := range lst {
		_ = serv.Storage.WriteOrderWithUniqueId(ctx, order)
		defer DeleteOrder(st, order.OrderID)
	}
	res, err := serv.Storage.GetFilteredOrdersByClientId(ctx, uint64(id), limit)
	require.NoError(t, err)

	assert.Equal(t, 1, len(res))
	for i, order := range res {
		assert.Equal(t, lst[i].OrderID, order.OrderID)
		assert.Equal(t, lst[i].ClientID, order.ClientID)
		assert.Equal(t, lst[i].CreatedAt, order.CreatedAt)
		assert.Equal(t, lst[i].TakeBackAt, order.TakeBackAt)
		assert.Equal(t, lst[i].TakeBackAt, order.TakeBackAt)
	}
}

func TestUpdateOrdersList(t *testing.T) {
	t.Parallel()
	serv, st := Init()
	defer st.CloseStorage()
	ctx := context.Background()
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 10, ClientID: 11, Pack: "box"})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 11, ClientID: 24, Pack: "box"})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 12, ClientID: 35, Pack: "box"})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 13, ClientID: 40, Pack: "box"})
	tx, _ := st.Db.Begin()
	defer DeleteOrder(st, 10)
	defer DeleteOrder(st, 11)
	defer DeleteOrder(st, 12)
	defer DeleteOrder(st, 13)
	defer tx.Rollback()
	lst := []*model.Order{{OrderID: 10, ClientID: 5, Pack: "box"},
		{OrderID: 11, ClientID: 6, Pack: "box"},
		{OrderID: 12, ClientID: 7, Pack: "box"},
		{OrderID: 13, ClientID: 5, ReturnedToCourier: true, Pack: "box"}}

	err := serv.Storage.UpdateOrdersList(ctx, lst, tx)
	require.NoError(t, err)

	for _, order := range lst {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM orders WHERE order_id = $1 and client_id=$2", order.OrderID, order.ClientID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	}
}

func TestGetTakebackList(t *testing.T) {
	t.Parallel()
	serv, st := Init()
	defer st.CloseStorage()
	ctx := context.Background()
	idList := []uint64{46, 50}
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 46, ClientID: 47, Pack: "box", TakenAt: time.Now().UTC()})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 48, ClientID: 49, Pack: "box"})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 50, ClientID: 51, Pack: "box", TakenAt: time.Now().UTC()})
	_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: 51, ClientID: 52, Pack: "box"})
	defer DeleteOrder(st, 46)
	defer DeleteOrder(st, 48)
	defer DeleteOrder(st, 50)
	defer DeleteOrder(st, 51)

	res, err := serv.Storage.GetTakebackList(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, len(res))
	for _, order := range res {
		assert.Contains(t, idList, order.OrderID)
	}
}
