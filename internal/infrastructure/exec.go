package infrastructure

import (
	"context"
	"flag"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"homework-1/internal/constant"
	"homework-1/internal/model"
	"homework-1/internal/service"
	"log"
)

func execChange(source []string, s *Synchronization) {
	span, _ := opentracing.StartSpanFromContext(context.Background(), "execute change")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for accept", flag.ContinueOnError)
	num := flags.Int(constant.FlagNumRoutines, 1, "number of go-routines")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}

	s.NumRoutinesChan <- *num
}

func execAccept(serv *service.Service, source []string) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "execute accept")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for accept", flag.ContinueOnError)
	orderId := flags.Uint64(constant.FlagOrderId, 0, "id of order")
	clientId := flags.Uint64(constant.FlagClientId, 0, "id of client")
	storedUntil := flags.String(constant.FlagStoredUntil, "", "date in dd-mm-yyyy format")
	pack := flags.String(constant.FlagPack, "", "type of package")
	weight := flags.Float64(constant.FlagWeight, 0, "weight of order")
	price := flags.Float64(constant.FlagPrice, 0, "price of order")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}

	if err := serv.AcceptFromCourier(ctx, *orderId, *clientId, *storedUntil, *pack, *weight, *price); err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(constant.Success)
}

func execBack(serv *service.Service, source []string) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "execute back")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for back", flag.ContinueOnError)
	orderId := flags.Uint64(constant.FlagOrderId, 0, "id of order")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}

	if err := serv.ReturnToCourier(ctx, *orderId); err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(constant.Success)
}

func execPickup(serv *service.Service, source []string) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "execute pickup")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for pickup", flag.ContinueOnError)
	orderIdList := flags.String("list", "", "a string to collect all order IDs")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}
	if *orderIdList == "" {
		log.Printf(model.ErrIdListMissing.Error())
		return
	}

	if err := serv.GiveToClient(ctx, orderIdList); err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(constant.Success)
}

func execList(serv *service.Service, source []string) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "execute list")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for list of orders", flag.ContinueOnError)
	clientId := flags.Uint64(constant.FlagClientId, 0, "id of client")
	limit := flags.Int("limit", 0, "to limit list of orders by client")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}

	if err := serv.ListForClient(ctx, *clientId, *limit); err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(constant.Success)
}

func execReturn(serv *service.Service, source []string) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "execute return")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for return", flag.ContinueOnError)
	orderId := flags.Uint64(constant.FlagOrderId, 0, "id of order")
	clientId := flags.Uint64(constant.FlagClientId, 0, "id of client")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}

	if err := serv.ClientReturnsOrder(ctx, *orderId, *clientId); err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(constant.Success)
}

func execTakebacks(serv *service.Service, source []string) {
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "execute takeback")
	defer span.Finish()
	flags := flag.NewFlagSet("flags for takeback", flag.ContinueOnError)
	takebackPage := flags.Int("page", 0, "to specify a page")
	if err := flags.Parse(source); err != nil {
		log.Printf(err.Error())
		return
	}

	if err := serv.ListOfTakeBacks(ctx, *takebackPage); err != nil {
		log.Printf(err.Error())
		return
	}
	fmt.Println(constant.Success)
}
