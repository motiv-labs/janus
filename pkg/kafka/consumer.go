package kafka

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"

	log "github.com/sirupsen/logrus"
)

type Message = kafka.Message
type Header = kafka.Header
type KafkaMessageHandler func(Message) error
type FaultyKafkaMessageHandler func(Message, error)

func ConsumeMessages(conf *Config, topic string, msgHandler KafkaMessageHandler, DLQHandler FaultyKafkaMessageHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           conf.KafkaBrokers,
		Topic:             topic,
		GroupID:           conf.Consumer.Group,
		QueueCapacity:     conf.Consumer.QueueCapacity,
		MinBytes:          conf.Consumer.MinBytes,
		MaxBytes:          conf.Consumer.MaxBytes,
		MaxWait:           conf.Consumer.MaxWait,
		HeartbeatInterval: conf.Consumer.HeartbeatInterval,

		GroupBalancers: []kafka.GroupBalancer{&kafka.RoundRobinGroupBalancer{}},
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		if err := r.Close(); err != nil {
			log.Printf("Cound not close kafka reader (Err: %s)\n", err)
		}
	}()

	for {
		msg, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Cound not read kafka message (Err: %s)\n", err.Error())
			DLQHandler(msg, err)
			continue
		}
		err = msgHandler(msg)
		if err != nil {
			log.Printf("Cound not handle kafka message (Err: %s)\n", err.Error())
			DLQHandler(msg, err)
		}
		err = r.CommitMessages(context.Background(), msg)
		if err != nil {
			log.Printf("Cound not commit kafka message (Err: %s)\n", err.Error())
			DLQHandler(msg, err)
		}
	}
}
