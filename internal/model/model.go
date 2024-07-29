package model

import (
	"fmt"
	"github.com/bojanz/currency"
	"time"
)

type Order struct {
	OrderID           uint64          `json:"order_id" db:"order_id"`
	ClientID          uint64          `json:"client_id" db:"client_id"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	StoredUntil       time.Time       `json:"stored_until" db:"stored_until"`
	TakenAt           time.Time       `json:"taken_at" db:"taken_at"`
	TakeBackAt        time.Time       `json:"takeback_time" db:"takeback_time"`
	ReturnedToCourier bool            `json:"returned" db:"returned"`
	Pack              string          `json:"package" db:"package"`
	Weight            float64         `json:"weight" db:"weight"`
	Price             currency.Amount `json:"price" db:"price"`
}

func (o Order) String() string {
	return fmt.Sprintf("Order id: %v\t delivered at: %v\t stored until: %v", o.OrderID, o.CreatedAt, o.StoredUntil)
}

func (o Order) AsTakeback() string {
	return fmt.Sprintf("Order id: %v\t taken back at: %v", o.OrderID, o.TakeBackAt)
}

type ByCreated []*Order

func (b ByCreated) Len() int {
	return len(b)
}
func (b ByCreated) Less(i, j int) bool {
	return b[i].CreatedAt.Before(b[j].CreatedAt)
}

func (b ByCreated) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
