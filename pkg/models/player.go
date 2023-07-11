package models

import (
	"time"
)

type Player struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	SessionID string    `json:"session_id"`
	Updated   time.Time `json:"updated"`
}
