package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	KafkaAddr          string `yaml:"kafkaAddr"`
	KafkaFactTopic     string `yaml:"kafkaFactTopic"`
	KafkaDLQTopic      string `yaml:"kafkaDLQTopic"`
	KafkaConsumerGroup string `yaml:"kafkaConsumerGroup"`
	RBACUrl            string `yaml:"RBACUrl"`
}

func UnmarshalYAML(path string, dest *Config) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, dest)
}
