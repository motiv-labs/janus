package models

import (
	"encoding/json"
	"github.com/hellofresh/janus/pkg/jwt"
)

type Fact struct {
	ID         uint64           `json:"id"`
	ObjectType string           `json:"objectType"`
	ActionType string           `json:"actionType"`
	Object     *json.RawMessage `json:"object"`
	Claims     *jwt.Claims      `json:"claims"`
}
