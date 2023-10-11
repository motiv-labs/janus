package models

import "sync"

type RoleManager struct {
	Roles map[string]*Role
	sync.Mutex
}

type Role struct {
	ID       uint64 `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Features []Feature
}

type Feature struct {
	ID          uint64 `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Path        string `json:"path" db:"path"`
	Method      string `json:"method" db:"method"`
}
