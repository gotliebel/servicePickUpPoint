package main

import (
	"bufio"
	"context"
	"flag"
	"homework-1/cmd/grpc"
	"homework-1/internal/infrastructure"
	"homework-1/internal/metrics"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"homework-1/internal/trace"
	"homework-1/kafka"
	"log"
	"os"
)

func main() {
	flag.BoolVar(&kafka.WriteToKafka, "kafka", false, "Write events to Kafka")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	defer st.CloseStorage()

	serv := service.New(st)
	in := bufio.NewReader(os.Stdin)
	signals := make(chan os.Signal)
	if kafka.WriteToKafka {
		go kafka.ConsumeEventsFromKafka()
	}
	s := infrastructure.NewSync()
	go infrastructure.Process(ctx, in, signals, s)
	s.Wg.Add(1)
	go infrastructure.ControlRoutines(serv, s)
	go grpc.NewGrpcServer(serv)
	go metrics.CountMetrics()
	c := trace.MustSetup(ctx, "orders_service_jaeger")
	defer c.Close()
	infrastructure.WaitForExitSignal(signals)
	infrastructure.WaitForExit(s)
}
