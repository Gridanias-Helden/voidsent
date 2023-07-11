package memory

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/oklog/ulid"

	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/storage"
)

type sessionManager struct {
	sessions map[string]models.Session
	ttl      time.Duration
}

func NewSessions(ttl time.Duration) storage.Sessions {
	return &sessionManager{
		sessions: make(map[string]models.Session),
		ttl:      ttl,
	}
}

func (mm *sessionManager) SessionByID(ctx context.Context, id string) (models.Session, error) {
	session, ok := mm.sessions[id]
	if !ok {
		return models.Session{}, fmt.Errorf("invalid session")
	}

	if session.Updated.Add(mm.ttl).Before(time.Now().UTC()) {
		return models.Session{}, fmt.Errorf("session expired")
	}

	return session, nil
}

func (mm *sessionManager) SaveSession(ctx context.Context, session models.Session) (models.Session, error) {
	if session.ID == "" {
		session.ID = ulid.MustNew(ulid.Now(), rand.Reader).String()
	}

	mm.sessions[session.ID] = session

	return session, nil
}

func (mm *sessionManager) DeleteSession(ctx context.Context, session models.Session) error {
	delete(mm.sessions, session.ID)

	return nil
}
