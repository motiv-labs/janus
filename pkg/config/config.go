package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/hellofresh/janus/pkg/kafka"
)

type Config struct {
	UserManagementURL string        `yaml:"UserManagementURL"`
	RbacURL           string        `yaml:"RbacURL"`
	ApiVersion        string        `yaml:"ApiVersion"`
	KafkaConfig       *kafka.Config `yaml:"kafkaConfig"`
}

func UnmarshalYAML(path string, dest *Config) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, dest)
}
