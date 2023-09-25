package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

	log "github.com/sirupsen/logrus"
)

type Message = kafka.Message
type Header = kafka.Header
type KafkaMessageHandler func(Message) error
type FaultyKafkaMessageHandler func(Message, error)

func ConsumeMessages(kafkaAddr, topic, consumerGroup string, msgHandler KafkaMessageHandler, DLQHandler FaultyKafkaMessageHandler) {
	// Minimum amount of bytes in a message batch
	const minBytes = 10e3
	// Maximum amount of bytes in a message batch
	const maxBytes = 10e6

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaAddr},
		Topic:    topic,
		GroupID:  consumerGroup,
		MinBytes: minBytes,
		MaxBytes: maxBytes,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Cound not read kafka message (Err: %s)\n", err.Error())
			DLQHandler(m, err)
			continue
		}
		err = msgHandler(m)
		if err != nil {
			log.Printf("Cound not handle kafka message (Err: %s)\n", err.Error())
			DLQHandler(m, err)
		}
	}
}
