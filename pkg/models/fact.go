package models

import (
	"encoding/json"
	"github.com/hellofresh/janus/pkg/jwt"
)

type Fact struct {
	ID       uint64           `json:"id"`
	PathRole string           `json:"objectType"`
	Method   string           `json:"actionType"`
	Object   *json.RawMessage `json:"object"`
	Claims   *jwt.Claims      `json:"claims"`
}
