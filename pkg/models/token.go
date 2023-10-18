package models

import (
	"sync"
	"time"
)

type TokenManager struct {
	Tokens map[string]*JWTToken
	sync.RWMutex
}

type JWTToken struct {
	ID             uint64    `json:"id"`
	UserID         uint64    `json:"user_id"`
	Token          string    `json:"token"`
	ExpirationDate time.Time `json:"expiration_date"`
}
