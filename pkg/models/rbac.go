package models

type Role struct {
	Name     string
	Features []*Feature
}

type Feature struct {
	Name      string
	Endpoints []*Endpoint
}

type Endpoint struct {
	Name   string
	Path   string
	Method string
}
