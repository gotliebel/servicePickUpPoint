package infrastructure

import (
	"bufio"
	"context"
	"fmt"
	"homework-1/internal/constant"
	"homework-1/kafka"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

func Read(in *bufio.Reader) string {
	fmt.Println(constant.WaitForCmdTxt)
	line, isPrefix, err := in.ReadLine()
	if err != nil {
		log.Printf(err.Error())
	}
	if isPrefix {
		log.Printf("input is too long")
	}
	return string(line)
}

func Process(ctx context.Context, in *bufio.Reader, signals chan os.Signal, s *Synchronization) {
	s.NumRoutinesChan <- baseNumRoutines
	producer := kafka.ProduceEventToKafka()
	defer producer.Close()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("process exit")
			return
		default:
			line := Read(in)
			commands := strings.Split(line, " ")
			command := commands[0]
			source := commands[1:]
			event := kafka.Event{
				CreatedAt: time.Now(),
				Method:    command,
				RawQuery:  line,
			}
			if kafka.WriteToKafka {
				kafka.SendEvent(producer, event)
			} else {
				fmt.Println("Created at: ", event.CreatedAt, "with method: ", event.Method, "with query: ", event.RawQuery)
			}
			if command == constant.CommandExit {
				signals <- syscall.SIGTERM
				return
			}
			s.CommandChan <- append(source, command)
		}
	}
}
