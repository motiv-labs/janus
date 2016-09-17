package main

// Specification for basic configurations
type Specification struct {
	DatabaseDSN string `envconfig:"DATABASE_DSN"`
	Port        int    `envconfig:"PORT"`
	Debug       bool   `envconfig:"DEBUG"`
	StorageDSN  string `envconfig:"REDIS_DSN"`
}
