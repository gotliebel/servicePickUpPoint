package service

import (
	"fmt"
	"github.com/bojanz/currency"
	"homework-1/internal/constant"
	"homework-1/internal/model"
	"homework-1/internal/packages"
	"strconv"
	"time"
)

func CheckDate(dateString string) (time.Time, error) {
	if dateString == "" {
		return time.Time{}, model.ErrDateMissing
	}
	storedUntil, err := time.Parse(constant.TimeFormat, dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf(model.ErrDateUnparsableTemplate, err.Error())
	}
	if storedUntil.Before(time.Now()) {
		return time.Time{}, model.ErrStorageTimeIsIncorrect
	}
	return storedUntil, nil
}

func CreateOrder(orderId, clientId uint64, storedUntil time.Time, strategy packages.PackagingStrategy, weight, price float64) (*model.Order, error) {
	if orderId == 0 {
		return nil, model.ErrOrderIdMissing
	}
	if clientId == 0 {
		return nil, model.ErrClientIdMissing
	}
	amount, _ := currency.NewAmount(strconv.FormatFloat(price, 'f', 6, 64), packages.Currency)
	order := &model.Order{
		OrderID:     orderId,
		ClientID:    clientId,
		StoredUntil: storedUntil,
		CreatedAt:   time.Time{},
		Weight:      weight,
		Price:       amount,
	}
	err := strategy.ApplyPackaging(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func CheckOrderId(orderId uint64) error {
	if orderId == 0 {
		return model.ErrOrderIdMissing
	}
	return nil
}

func CheckClientId(clientId uint64) error {
	if clientId == 0 {
		return model.ErrClientIdMissing
	}
	return nil
}

func CheckOrder(order *model.Order, clientId uint64, curTime time.Time) error {
	if clientId != order.ClientID {
		return model.ErrOrdersByDifferentClients
	}
	if order.StoredUntil.Before(curTime) {
		return model.ErrTimeHasExpired
	}
	if !order.TakenAt.IsZero() {
		return model.ErrOrderAlreadyTaken
	}
	return nil
}

func CheckPageNumber(pageNumber int) error {
	if pageNumber <= 0 {
		return model.ErrPageMissing
	}
	return nil
}

func CheckLen(expect, real int) error {
	if expect != real {
		return model.ErrSomeOrdersWereNotFound
	}
	return nil
}

func CheckPackage(pack string) (packages.PackagingStrategy, error) {
	if pack == "" {
		return nil, model.ErrPackageMissing
	}
	if _, ok := packages.PackageMap[pack]; !ok {
		return nil, fmt.Errorf(model.ErrPackageDoesntExist, pack)
	}
	var strategy packages.PackagingStrategy
	strategy = packages.PackageMap[pack]
	return strategy, nil
}

func CheckLimit(limit int) error {
	if limit < 0 {
		return model.ErrNegativeLimit
	}
	return nil
}
