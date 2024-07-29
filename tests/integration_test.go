//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-1/internal/constant"
	"homework-1/internal/model"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"log"
	"testing"
	"time"
)

func Init() (*service.Service, *storage.Storage, context.Context) {
	st, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	serv := service.New(st)
	return serv, st, context.Background()
}

func DeleteOrder(st *storage.Storage, id uint64) {
	_ = st.Db.QueryRow(`DELETE FROM orders WHERE order_id = $1`, id).Err()
	_ = st.Db.QueryRow(`DELETE FROM characteristic WHERE order_id = $1`, id).Err()
}

func TestAcceptFromCourier(t *testing.T) {
	t.Parallel()
	serv, st, ctx := Init()
	defer st.CloseStorage()
	testCases := []struct {
		name     string
		orderId  uint64
		clientId uint64
		date     string
		pack     string
		weight   float64
		price    float64
		error    bool
	}{
		{name: "simple",
			orderId:  14,
			clientId: 1,
			date:     "12-12-2024",
			pack:     "bag",
			weight:   1.0,
			price:    1.0,
		},
		{
			name:     "with missing client id",
			orderId:  15,
			clientId: 0,
			date:     "12-12-2024",
			pack:     "bag",
			error:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := serv.AcceptFromCourier(ctx, tc.orderId, tc.clientId, tc.date, tc.pack, tc.weight, tc.price)
			if tc.error {
				require.Error(t, err)
			} else {
				var count int
				require.NoError(t, err)
				err = st.Db.QueryRow("SELECT COUNT(*) FROM orders WHERE order_id = $1", tc.orderId).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 1, count)

				err = st.Db.QueryRow("SELECT COUNT(*) FROM characteristic WHERE order_id = $1", tc.orderId).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 1, count)
				DeleteOrder(st, tc.orderId)
			}
		})
	}
}

func TestReturnToCourier(t *testing.T) {
	t.Parallel()
	serv, st, ctx := Init()
	defer st.CloseStorage()

	testCases := []struct {
		name    string
		orderId uint64
		date    time.Time
		error   bool
	}{
		{
			name:    "further end of storage",
			orderId: 16,
			date:    time.Now().Add(time.Hour * 24 * 30),
			error:   true,
		},
		{
			name:    "simple",
			orderId: 17,
			date:    time.Now().Add(-time.Hour * 24),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: tc.orderId, StoredUntil: tc.date, Pack: "box"})
			err := serv.ReturnToCourier(ctx, tc.orderId)
			if tc.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				var count int
				err = st.Db.QueryRow("SELECT COUNT(*) FROM orders WHERE order_id = $1 and returned = true", tc.orderId).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 1, count)
			}
			DeleteOrder(st, tc.orderId)
		})
	}
}

func TestGiveToClient(t *testing.T) {
	t.Parallel()
	serv, st, ctx := Init()
	defer st.CloseStorage()

	testCases := []struct {
		name      string
		orderId   []uint64
		clientId  uint64
		orderList string
		date      string
	}{
		{
			name:      "simple",
			orderId:   []uint64{18, 19, 20},
			clientId:  18,
			orderList: "18,19,20",
			date:      "12-12-2024",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			date, _ := time.Parse(constant.TimeFormat, tc.date)
			for _, id := range tc.orderId {
				_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: id, StoredUntil: date, Pack: "box"})
				defer DeleteOrder(st, id)
			}
			err := serv.GiveToClient(ctx, &tc.orderList)
			require.NoError(t, err)
		})
	}
}

func TestClientReturnsOrder(t *testing.T) {
	t.Parallel()
	serv, st, ctx := Init()
	defer st.CloseStorage()

	testCases := []struct {
		name     string
		orderId  uint64
		clientId uint64
		date     time.Time
		error    bool
	}{
		{
			name:     "cannot be returned",
			orderId:  21,
			clientId: 22,
			date:     time.Time{},
			error:    true,
		},
		{
			name:     "simple",
			orderId:  23,
			clientId: 24,
			date:     time.Now().Add(time.Hour * 24 * 30),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = serv.Storage.WriteOrderWithUniqueId(ctx, &model.Order{OrderID: tc.orderId, ClientID: tc.clientId, TakenAt: tc.date, Pack: "box"})
			defer DeleteOrder(st, tc.orderId)
			err := serv.ClientReturnsOrder(ctx, tc.orderId, tc.clientId)
			if tc.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
