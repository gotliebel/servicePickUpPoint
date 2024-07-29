package packages

import (
	"homework-1/internal/model"
)

const (
	PackageWrapping = "wrapping"
	PriceWrapping   = 1
)

type Wrapping struct {
	name  string
	price float64
}

func NewWrapping() *Wrapping {
	return &Wrapping{
		name:  PackageWrapping,
		price: PriceWrapping,
	}
}

func (w *Wrapping) ApplyPackaging(order *model.Order) error {
	return FillPackage(order, w.name, 0, w.price)
}
