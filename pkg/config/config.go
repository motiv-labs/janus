package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	UserManagementURL string `yaml:"userManagementURL"`
	RbacURL           string `yaml:"rbacURL"`
	ApiVersion        string `yaml:"apiVersion"`
}

func UnmarshalYAML(path string, dest *Config) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, dest)
}
