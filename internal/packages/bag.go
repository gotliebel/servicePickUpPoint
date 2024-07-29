package packages

import (
	"homework-1/internal/model"
)

const (
	PackageBag = "bag"
	PriceBag   = 5
	WeightBag  = 10
	Currency   = "USD"
)

type Bag struct {
	name      string
	maxWeight float64
	price     float64
}

func NewBag() *Bag {
	return &Bag{
		name:      PackageBag,
		maxWeight: WeightBag,
		price:     PriceBag,
	}
}

func (b *Bag) ApplyPackaging(order *model.Order) error {
	return FillPackage(order, b.name, b.maxWeight, b.price)
}
