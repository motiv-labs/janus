package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	kafkaWriter *kafka.Writer
}

func NewKafkaProducer(conf *Config) *KafkaProducer {
	return &KafkaProducer{kafkaWriter: &kafka.Writer{
		Addr:         kafka.TCP(conf.KafkaAddr),
		Topic:        conf.Producer.Topic,
		BatchTimeout: conf.Producer.BatchTimeout,
		BatchSize:    conf.Producer.BatchSize,
		BatchBytes:   conf.Producer.BatchBytes,

		Balancer:    &kafka.LeastBytes{},
		Compression: kafka.Snappy,
	}}
}

const BatchSize = 1000

func (producer *KafkaProducer) ProduceMessage(key string, messageBytes []byte) error {
	err := producer.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: messageBytes,
		},
	)
	return err
}

func (producer *KafkaProducer) ProduceMadeMessage(msg kafka.Message) error {
	err := producer.kafkaWriter.WriteMessages(context.Background(), msg)
	return err
}

func (producer *KafkaProducer) Close() error {
	err := producer.kafkaWriter.Close()
	return err
}
