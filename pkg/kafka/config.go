package kafka

import "time"

type producer struct {
	ProducersAmount       int           `yaml:"producersAmount"`
	Topic                 string        `yaml:"topic"`
	DLQTopic              string        `yaml:"DLQTopic"`
	BatchSize             int           `yaml:"batchSize"`
	BatchBytes            int64         `yaml:"batchBytes"`
	BatchTimeout          time.Duration `yaml:"batchTimeout"`
	FactProduceRetryDelay time.Duration `yaml:"factProduceRetryDelay"`
}

type consumer struct {
	ConsumersAmount   int           `yaml:"consumersAmount"`
	Topics            []string      `yaml:"topics"`
	Group             string        `yaml:"group"`
	QueueCapacity     int           `yaml:"queueCapacity"`
	MinBytes          int           `yaml:"minBytes"`
	MaxBytes          int           `yaml:"maxBytes"`
	MaxWait           time.Duration `yaml:"maxWait"`
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval"`
}

type Config struct {
	KafkaAddr string   `yaml:"kafkaAddr"`
	Producer  producer `yaml:"producer"`
	Consumer  consumer `yaml:"consumer"`
}

func (c *Config) Normalize() {
	if c.Producer.ProducersAmount == 0 {
		c.Producer.ProducersAmount = 1
	}
	if c.Producer.BatchSize == 0 {
		c.Producer.BatchSize = 64
	}
	if c.Producer.BatchBytes == 0 {
		c.Producer.BatchBytes = 32000
	}
	if c.Producer.BatchTimeout == 0 {
		c.Producer.BatchTimeout = 500 * time.Millisecond
	}
	if c.Producer.FactProduceRetryDelay == 0 {
		c.Producer.FactProduceRetryDelay = 2 * time.Second
	}

	if c.Consumer.ConsumersAmount == 0 {
		c.Consumer.ConsumersAmount = 1
	}
	if c.Consumer.QueueCapacity == 0 {
		c.Consumer.QueueCapacity = 128
	}
	if c.Consumer.MinBytes == 0 {
		c.Consumer.MinBytes = 1
	}
	if c.Consumer.MaxBytes == 0 {
		c.Consumer.MaxBytes = 1e6
	}
	if c.Consumer.MaxWait == 0 {
		c.Consumer.MaxWait = 10 * time.Second
	}
	if c.Consumer.HeartbeatInterval == 0 {
		c.Consumer.HeartbeatInterval = 500 * time.Millisecond
	}
}
