package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

func ProduceEventToKafka() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(Brokers, config)
	if err != nil {
		log.Printf("Error creating Kafka producer: %s", err)
		return nil
	}
	return producer
}

func SendEvent(producer sarama.AsyncProducer, event Event) {
	eventJSON, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.StringEncoder(eventJSON),
	}
	producer.Input() <- msg
}
