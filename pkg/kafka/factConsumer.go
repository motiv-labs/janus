package kafka

import (
	"context"
	"encoding/json"
	"github.com/hellofresh/janus/pkg/cache"
	"github.com/hellofresh/janus/pkg/models"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func StartFactConsumer(kafkaAddr, topic, dlqtopic, consumerGroup string, rolesCache *cache.RolesCache) {
	ConsumeMessages(kafkaAddr, topic, consumerGroup,
		func(msg Message) error {
			var fact models.Fact
			err := json.Unmarshal(msg.Value, &fact)
			if err != nil {
				return err
			}

			var role *models.Role

			err = json.Unmarshal(*fact.Object, &role)
			if err != nil {
				log.Println(err)
			}

			switch fact.Method {
			case http.MethodPost:
				err = rolesCache.Set(role)
				if err != nil {
					log.Printf("Role %s not created, err: %s", role.Name, err.Error())
				}
			case http.MethodPut:
				err = rolesCache.Update(role, fact.PathRole)
				if err != nil {
					log.Printf("Role %s not changed, err:", fact.PathRole, err)
				}
			case http.MethodDelete:
				err = rolesCache.Delete(fact.PathRole)
				if err != nil {
					log.Printf("Role %s not deleted, err:", fact.PathRole, err)
				}
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
