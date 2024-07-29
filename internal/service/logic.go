package service

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"homework-1/internal/constant"
	"homework-1/internal/model"
	"math"
	"strconv"
	"strings"
	"time"
)

func (s *Service) AcceptFromCourier(ctx context.Context, orderId, clientId uint64, dateString string, pack string, weight, price float64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service AcceptFromCourier")
	defer span.Finish()
	storedUntil, err := CheckDate(dateString)
	if err != nil {
		return err
	}

	strategy, err := CheckPackage(pack)
	if err != nil {
		return err
	}

	order, err := CreateOrder(orderId, clientId, storedUntil, strategy, weight, price)
	if err != nil {
		return err
	}

	err = s.Storage.WriteOrderWithUniqueId(ctx, order)
	return err
}

func (s *Service) ReturnToCourier(ctx context.Context, orderId uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service ReturnToCourier")
	defer span.Finish()
	err := CheckOrderId(orderId)
	if err != nil {
		return err
	}

	order, err := s.Storage.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if !order.StoredUntil.Before(time.Now()) {
		return fmt.Errorf(model.ErrTooSoonToTakeBackTemplate, order.StoredUntil.Format(constant.TimeFormat))
	}
	if !order.TakenAt.IsZero() {
		return model.ErrOrderAlreadyTaken
	}
	order.ReturnedToCourier = true
	tx, err := s.Storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = s.Storage.UpdateOrdersList(ctx, []*model.Order{order}, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) GiveToClient(ctx context.Context, orderIdList *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service GiveToClient")
	defer span.Finish()
	idStringList := strings.Split(*orderIdList, ",")
	idList := make([]uint64, 0, len(idStringList))
	for _, id := range idStringList {
		num, _ := strconv.ParseUint(id, 10, 64)
		idList = append(idList, num)
	}
	return s.GetOrdersAndUpdateWithSafety(ctx, idList)

}

func (s *Service) ListForClient(ctx context.Context, clientId uint64, limit int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service ListForClient")
	defer span.Finish()
	err := CheckClientId(clientId)
	if err != nil {
		return err
	}
	err = CheckLimit(limit)
	if err != nil {
		return err
	}
	if limit == 0 {
		limit = math.MaxUint32
	}
	orders, err := s.Storage.GetFilteredOrdersByClientId(ctx, clientId, limit)
	if err != nil {
		return err
	}
	for _, order := range orders {
		fmt.Println(order)
	}
	return nil
}

func (s *Service) ClientReturnsOrder(ctx context.Context, orderId, clientId uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service ClientReturnsOrder")
	defer span.Finish()
	err := CheckOrderId(orderId)
	if err != nil {
		return err
	}
	err = CheckClientId(clientId)
	if err != nil {
		return err
	}

	order, err := s.Storage.GetOrder(ctx, orderId)
	if err != nil {
		return err
	}

	if order.ClientID != clientId {
		return model.ErrOrderMadeByOtherClient
	}
	if order.TakenAt.Add(TimeToReturn).Before(time.Now()) {
		return model.ErrTimeToReturnHasExpired
	}
	if order.TakenAt.IsZero() {
		return model.ErrOrderWasAlreadyIssued
	}
	order.TakenAt = time.Time{}
	order.TakeBackAt = time.Now()
	tx, err := s.Storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = s.Storage.UpdateOrdersList(ctx, []*model.Order{order}, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) ListOfTakeBacks(ctx context.Context, pageNumber int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service ListOfTakeBacks")
	defer span.Finish()
	err := CheckPageNumber(pageNumber)
	if err != nil {
		return err
	}

	takebacks, err := s.Storage.GetTakebackList(ctx, pageNumber)
	if err != nil {
		return err
	}

	for _, orders := range takebacks {
		fmt.Println(orders.AsTakeback())
	}
	return nil
}

func (s *Service) GetOrdersAndUpdateWithSafety(ctx context.Context, idList []uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service GetOrdersAndUpdateWithSafety")
	defer span.Finish()
	tx, err := s.Storage.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orders, err := s.Storage.GetOrdersByIdList(ctx, idList, tx)
	if err != nil {
		return err
	}
	if err := CheckLen(len(orders), len(idList)); err != nil {
		return err
	}
	clientId := orders[0].ClientID
	curTime := time.Now()
	for _, order := range orders {
		err = CheckOrder(order, clientId, curTime)
		if err != nil {
			return err
		}
		order.TakenAt = curTime
	}

	err = s.Storage.UpdateOrdersList(ctx, orders, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}
