package main

import (
	"github.com/hellofresh/api-gateway/storage"
)

// Specification for basic configurations
type Specification struct {
	storage.Database
	Port  int  `envconfig:"PORT"`
	Debug bool `envconfig:"DEBUG"`
	Storage
}

// Storage holds the configuration for a data storage
type Storage struct {
	DSN      string `envconfig:"REDIS_DSN"`
	Password string `envconfig:"REDIS_PASSWORD"`
	Database int64  `envconfig:"REDIS_DB"`
}
