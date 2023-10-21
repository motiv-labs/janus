package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	kafkaWriter *kafka.Writer
}

func NewKafkaProducer(conf *Config) *kafkaProducer {
	return &kafkaProducer{kafkaWriter: &kafka.Writer{
		Addr:         kafka.TCP(conf.KafkaBrokers...),
		Topic:        conf.Producer.Topic,
		BatchTimeout: conf.Producer.BatchTimeout,
		BatchSize:    conf.Producer.BatchSize,
		BatchBytes:   conf.Producer.BatchBytes,

		Balancer:    &kafka.LeastBytes{},
		Compression: kafka.Snappy,
	}}
}

func (producer *kafkaProducer) ProduceMessage(key string, messageBytes []byte) error {
	err := producer.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: messageBytes,
		},
	)
	return err
}

func (producer *kafkaProducer) ProduceMadeMessage(msg kafka.Message) error {
	err := producer.kafkaWriter.WriteMessages(context.Background(), msg)
	return err
}

func (producer *kafkaProducer) Close() error {
	err := producer.kafkaWriter.Close()
	return err
}
