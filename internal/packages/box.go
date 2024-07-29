package packages

import (
	"homework-1/internal/model"
)

const (
	PackageBox = "box"
	PriceBox   = 20
	WeightBox  = 30
)

type Box struct {
	name      string
	maxWeight float64
	price     float64
}

func NewBox() *Box {
	return &Box{
		name:      PackageBox,
		maxWeight: WeightBox,
		price:     PriceBox,
	}
}

func (b *Box) ApplyPackaging(order *model.Order) error {
	return FillPackage(order, b.name, b.maxWeight, b.price)
}
