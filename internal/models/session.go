package models

import (
	"context"
	"time"
)

type Session struct {
	ID      string
	Player  *Player
	Started time.Time
}

type SessionManager interface {
	LoadSessionByID(ctx context.Context, id string) (*Session, error)
	LoadAllSessions(ctx context.Context) ([]*Session, error)
	SaveSession(ctx context.Context, session *Session) (*Session, error)
	DeleteSession(ctx context.Context, session *Session) error
}
