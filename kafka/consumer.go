package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

func ConsumeEventsFromKafka() {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(Brokers, Group, config)
	if err != nil {
		log.Printf("Error creating Kafka consumer: %s", err)
		return
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing Kafka consumer: %s", err)
		}
	}()

	handler := ConsumerHandler{stopChan: make(chan struct{})}
	for {
		select {
		case <-handler.stopChan:
			consumer.Close()
			return
		default:

			err := consumer.Consume(context.Background(), []string{Topic}, &handler)
			if err != nil {
				log.Printf("Error consuming events from Kafka: %s", err)
				return
			}
		}
	}
}

type Event struct {
	CreatedAt time.Time `json:"created_at"`
	Method    string    `json:"method"`
	RawQuery  string    `json:"raw_query"`
}

type ConsumerHandler struct {
	stopChan chan struct{}
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event Event
		err := json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Error decoding event: %s", err)
			continue
		}
		fmt.Printf("Event from Kafka: %+v\n", event)
		session.MarkMessage(message, "")
		if event.Method == "exit" {
			close(h.stopChan)
			return nil
		}
	}
	return nil
}
