package memory

import (
	"context"
	"fmt"

	"github.com/gridanias-helden/voidsent/internal/models"
)

type SessionManager struct {
	Sessions map[string]*models.Session
}

func NewSessionManager() models.SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*models.Session, 0),
	}
}

func (mm *SessionManager) LoadSessionByID(ctx context.Context, id string) (*models.Session, error) {
	session, ok := mm.Sessions[id]
	if !ok {
		return nil, fmt.Errorf("session %q not found", id)
	}

	return session, nil
}

func (mm *SessionManager) LoadAllSessions(ctx context.Context) ([]*models.Session, error) {
	var Sessions []*models.Session
	for _, Session := range mm.Sessions {
		Sessions = append(Sessions, Session)
	}

	return Sessions, nil
}

func (mm *SessionManager) SaveSession(ctx context.Context, Session *models.Session) (*models.Session, error) {
	mm.Sessions[Session.ID] = Session

	return Session, nil
}

func (mm *SessionManager) DeleteSession(ctx context.Context, Session *models.Session) error {
	delete(mm.Sessions, Session.ID)
	return nil
}
