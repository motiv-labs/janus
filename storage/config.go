package storage


// Database Specifies database configurations
type Database struct {
	DSN  string `envconfig:"DATABASE_DSN"`
	Name string `envconfig:"DATABASE_NAME"`
}
