package kafka

import (
	"context"
	"encoding/json"
	"github.com/hellofresh/janus/pkg/models"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

func StartFactConsumer(kafkaAddr, topic, dlqtopic, consumerGroup string) {
	ConsumeMessages(kafkaAddr, topic, consumerGroup,
		func(msg Message) error {
			var fact models.Fact
			err := json.Unmarshal(msg.Value, &fact)
			if err != nil {
				return err
			}
			var role models.Role

			err = json.Unmarshal(*fact.Object, &role)
			if err != nil {
				log.Println(err)
			}

			return nil
		},
		func(msg Message, inerr error) {
			msg.Headers = []kafka.Header{{Key: "Error", Value: []byte(inerr.Error())}}
			producer := NewKafkaProducer(kafkaAddr, dlqtopic)
			defer producer.Close()
			err := producer.ProduceMadeMessage(msg)
			if err != nil {
				log.Println(err)
			}
			return
		},
	)
	return
}

type KafkaProducer struct {
	kafkaWriter *kafka.Writer
}

func NewKafkaProducer(kafkaAddr string, topic string) *KafkaProducer {
	return &KafkaProducer{kafkaWriter: &kafka.Writer{
		Addr:     kafka.TCP(kafkaAddr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}}
}

func (producer *KafkaProducer) ProduceMadeMessage(msg kafka.Message) error {
	err := producer.kafkaWriter.WriteMessages(context.Background(), msg)
	return err
}

func (producer *KafkaProducer) Close() error {
	err := producer.kafkaWriter.Close()
	return err
}
