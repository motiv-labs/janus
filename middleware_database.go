package main

import (
	"github.com/kataras/iris"
	"github.com/hellofresh/api-gateway/storage"
)

// a silly example
type Database struct {
	dba storage.DatabaseAccessor
}

func NewDatabase(server storage.DatabaseAccessor) *Database {
	return &Database{server}
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (d Database) Serve(c *iris.Context) {
	reqSession := d.dba.Clone()
	defer reqSession.Close()
	d.dba.Set(c, reqSession)
	c.Next()
}
