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

func ConsumeMessages(conf *Config, topic string, msgHandler KafkaMessageHandler, DLQHandler FaultyKafkaMessageHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           []string{conf.KafkaAddr},
		Topic:             topic,
		GroupID:           conf.Consumer.Group,
		QueueCapacity:     conf.Consumer.QueueCapacity,
		MinBytes:          conf.Consumer.MinBytes,
		MaxBytes:          conf.Consumer.MaxBytes,
		MaxWait:           conf.Consumer.MaxWait,
		HeartbeatInterval: conf.Consumer.HeartbeatInterval,

		GroupBalancers: []kafka.GroupBalancer{&kafka.RangeGroupBalancer{}},
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
