package packages

import (
	"fmt"
	"github.com/bojanz/currency"
	"homework-1/internal/model"
	"strconv"
)

type PackagingStrategy interface {
	ApplyPackaging(order *model.Order) error
}

func FillPackage(order *model.Order, name string, maxWeight float64, price float64) error {
	order.Pack = name
	base, _ := currency.NewAmount(strconv.FormatFloat(price, 'f', 6, 64), Currency)
	if order.Weight < maxWeight || maxWeight == 0 {
		order.Price, _ = order.Price.Add(base)
		return nil
	}
	return fmt.Errorf(model.ErrOverWeight, name, maxWeight)
}
