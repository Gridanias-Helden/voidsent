package models

import (
	"context"
	"time"
)

type Session struct {
	ID      string    `json:"id"`
	Player  *Player   `json:"player"`
	Started time.Time `json:"started"`
	Game    string    `json:"game"`
}

type SessionManager interface {
	LoadSessionByID(ctx context.Context, id string) (*Session, error)
	LoadAllSessions(ctx context.Context) ([]*Session, error)
	SaveSession(ctx context.Context, session *Session) (*Session, error)
	DeleteSession(ctx context.Context, session *Session) error
}
