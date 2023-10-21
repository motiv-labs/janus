package models

import "sync"

type RoleManager struct {
	Roles map[string]*Role
	sync.Mutex
}

type Role struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Features []Feature
}

type Feature struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
}
