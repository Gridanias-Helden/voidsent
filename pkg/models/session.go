package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	ID       string          `json:"id"`
	PlayerID string          `json:"player_id"`
	Avatar   string          `json:"avatar"`
	Username string          `json:"username"`
	Updated  time.Time       `json:"updated"`
	WS       *websocket.Conn `json:"-"`
}
